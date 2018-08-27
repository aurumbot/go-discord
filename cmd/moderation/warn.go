package moderation

import (
	dsg "github.com/bwmarrin/discordgo"
	f "github.com/whitman-colm/go-discord"
	//"github.com/whitman-colm/go-discord/dat"
	"github.com/whitman-colm/go-discord/flags"
	"strings"
	"time"
)

func init() {
	Commands["warn"] = &f.Command{
		Name: "Warns a user",
		Help: `Officially warns a user of a minor infraction.
The warn will be placed on any users tagged.
	--reason	| reason warned
The warn will be considered "decayed" after 24 hours.
The actor must have the PermissionKickMembers perm to run this command.
Usage: mute @user --reason please do not post spoilers in this channel`,
		Perms:   dsg.PermissionKickMembers,
		Version: "v1.0α",
		Action:  warnUser,
	}
}

func warnUser(session *dsg.Session, message *dsg.MessageCreate) {
	s := session
	m := message

	var (
		reason string
	)
	flagSplit := strings.Split(message.Content, " ")
	flagsParsed := flags.Parse(flagSplit)

	for _, flag := range flagsParsed {
		if flag.Type == flags.DoubleDash && flag.Name == "reason" {
			reason = flag.Value
		}
	}

	for _, user := range m.Mentions {
		s.ChannelMessageSend(m.ChannelID, "Warned "+user.Username+"#"+user.Discriminator+" (ID: "+user.ID+") for \""+reason+"\".")
		duration, _ := time.ParseDuration("19.32m")
		createIncident(user.ID, m.Author.ID, time.Now(), duration, reason, warn)
	}
}
