package mods

import (
	dsg "github.com/bwmarrin/discordgo"
	f "github.com/skilstak/go-discord"
	"github.com/skilstak/go-discord/dat"
	"github.com/skilstak/go-discord/flags"
	"strings"
	"time"
)

type config struct {
	MuteRole    string   `json:"muted"`
	BannedWords []string `json:"swears"`
}

/* Moderation logs
* The moderation package comes with built-in logs using json files as a pseudo-
* database. Each "incident" consists of the following fields:
* - ID     : The infraction ID, this is a unique id for each incident
* - User   : The ID of the user who was acted upon.
* - Issuer : The moderator who took action against the user.
* - Time   : The time the action was taken.
* - Until  : When the action ends for a temp-mute or temp-ban
* - Reason : Why the action was issued
* - Action : The ID of the action that was taken. See constants for more info
* - Decay  : Statute of limitations. Prevents the bot from being overly-harsh
*            to repeat offences with large gaps of time between infractions.
 */
type incident struct {
	ID     int    `json:"case"`
	User   string `json:"user"`
	Issuer string `json:"actor"`
	Time   string `json:"time"`
	Until  string `json:"duration"`
	Reason string `json:"reason"`
	Action int    `json:"action"`
	Decay  string `json:"decay"`
}

type modlogs struct {
	Incidents []incident `json:"incidents"`
}

const (
	warn     = 0 // an official warning, no immediate action.
	delmsg       // message is deleted (censorship by bot).
	mute         // user is muted
	kick         // user is kicked
	tempban      // user is temporarily banned
	permaban     // user is permanently banned
	restore      // user has their mute revoked
	pardon       // user is pardoned from a ban
)

var (
	Commands = make(map[string]*f.Command)
	ps       string
	cfg      *config
	logs     *modlogs
)

func init() {
	ps = dat.OSCheck()
	dat.Load("moderation"+ps+"config.json", &cfg)
	dat.Load("moderation"+ps+"logs.json", &logs)
	Commands["mute"] = &f.Command{
		Name: "Mute a user",
		Help: `Adds the muted role to the mentioned user(s).
The muted role will forbid permission to speak in text and voice channels
unless otherwise and manually changed. A reason and duration (in minutes) can
also be added with the flags:
	--duration	| time muted
	--reason	| reason muted
The muted role will be generated by the bot upon first use.
Usage: mute @user --duration 5 --reason spam`,
		Action: muteUser,
	}
}

/* Logs Moderation Actions
* Create Incident is a private function for use only in this package.
* It handles creating logs and storing them in their JSON proto-"database".
* The function uses the following arguments:
* - u string	: The user affected (ID)
* - a string	: The person who used moderation action (ID)
* - i .Time	: The time that the action was taken
* - e .Duration	: The duration of the action
* - w string	: The reason why the action was taken
* - p int	: The action that was taken (as defined by the constants above)
* The function returns the following values:
* - _ error	: An error. Note that the error has already been logged to the
*		  appropriate place, this return is just in case exterior acts
*		  like alerting the discord-end users that something messed up.
*
* //TODO: I swear I should be using pointers but it breaks the append function.
* //NOTE: The decay time is generated internally by the bot.
*	  //TODO: Make decay times able to be changed in the config to be more
*		  or less harsh.
 */
func createIncident(u string, a string, i time.Time, e time.Duration, w string, p int) error {
	if w == "" {
		w = "No reason was provided"
	}

	eP := i.Add(e)
	//Decay is generated by the function 3√(x^2) * 200 where X is
	//the length of the mute in hopefully minutes.
	decay := eP.Add((((e ^ 2) ^ (1 / 3)) * 200))

	log := incident{
		ID:     len(logs.Incidents),
		User:   u,
		Issuer: a,
		Time:   i.Format("2006-01-02@15:04:05"),
		Until:  eP.Format("2006-01-02@15:04:05"),
		Reason: w,
		Action: p,
		Decay:  decay.Format("2006-01-02@15:04:05"),
	}
	logs.Incidents = append(logs.Incidents, log)
	err := dat.Save("moderation"+ps+"logs.json", logs)
	if err != nil {
		dat.Log.Println(err)
		return err
	}
	return nil
}

func muteUser(session *dsg.Session, message *dsg.MessageCreate) {
	s := session
	m := message

	perm, err := f.HasPermissions(s, m.Message, m.Author.ID, dsg.PermissionKickMembers)
	if err != nil {
		dat.Log.Println(err)
		dat.AlertDiscord(s, m, err)
		return
	}
	if !perm {
		s.ChannelMessageSend(m.ChannelID, "Sorry, you do not have permission to use this command.")
		return
	}

	var (
		reason   string
		duration time.Duration
	)
	flagSplit := strings.Split(message.Content, " ")
	flagsParsed := flags.Parse(flagSplit)

	for _, flag := range flagsParsed {
		if flag.Type == flags.DoubleDash && flag.Name == "reason" {
			reason = flag.Value[0]
		} else if flag.Type == flags.DoubleDash && flag.Name == "reason" {
			duration, err = time.ParseDuration(flag.Value[0])
			if err != nil {
				dat.Log.Println(err)
				dat.AlertDiscord(s, m, err)
				return
			}
		}
	}

	guild, err := f.GetGuild(s, m.Message)
	if err != nil {
		dat.Log.Println(err)
		s.ChannelMessageSend(m.ChannelID, "Error: "+err.Error())
		return
	}

	if cfg.MuteRole == "" {
		createMuteRole(guild, s, m)
	}

	for _, user := range m.Mentions {
		err := s.GuildMemberRoleAdd(guild.ID, user.ID, cfg.MuteRole)
		if err != nil {
			dat.Log.Println(err)
			s.ChannelMessageSend(m.ChannelID, "Error: "+err.Error())
		} else {
			s.ChannelMessageSend(m.ChannelID, user.Username+"#"+user.Discriminator+" (ID: "+user.ID+") has been muted.")
			createIncident(user.ID, m.Author.ID, time.Now(), duration, reason, mute)
		}
	}
}

func createMuteRole(guild *dsg.Guild, s *dsg.Session, m *dsg.MessageCreate) {
	role, err := s.GuildRoleCreate(guild.ID)
	if err != nil {
		dat.Log.Println(err)
		s.ChannelMessageSend(m.ChannelID, "Error: "+err.Error())
		return
	}
	role, err = s.GuildRoleEdit(guild.ID, role.ID, "Muted", 8487814, false, 1049664, false)
	if err != nil {
		dat.Log.Println(err)
		s.ChannelMessageSend(m.ChannelID, "Error: "+err.Error())
		return
	}
	cfg.MuteRole = role.ID
	err = dat.Save("moderation"+ps+"config.json", &cfg)
	if err != nil {
		dat.Log.Println(err)
		s.ChannelMessageSend(m.ChannelID, "Error: "+err.Error())
		return
	}
}
