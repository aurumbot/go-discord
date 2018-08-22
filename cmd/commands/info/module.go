package info

import (
	//"errors"
	dsg "github.com/bwmarrin/discordgo"
	f "github.com/whitman-colm/go-discord"
	"github.com/whitman-colm/go-discord/dat"
	"github.com/whitman-colm/go-discord/flags"
	"strings"
)

type config struct {
	version string `json:"version"`
	embed   bool   `json:"embedDefault"`
	useDMs  bool   `json:"dmDefault"`
}

var (
	Commands = make(map[string]*f.Command)
	cfg      *config
	ps       string
)

func init() {
	ps = dat.OSCheck()
	dat.Load("info"+ps+"config.json", &cfg)
	Commands["info"] = &f.Command{
		Name: "Bot Info",
		Help: `Gets information about the bot, version number, so on.
	Options:
	 -e : Get info as embed
	 -t : Get info as raw text
	 -m : Get info via direct message
	 -c : Post info in channel
Github  : https://github.com/whitman-colm/discord-public/cmd/commands/info`,
		Perms:   -1,
		Version: "v1.0",
		Action:  info,
	}
}

func info(session *dsg.Session, message *dsg.MessageCreate) {
	f := strings.Split(message.Content, " ")
	flg := flags.Parse(f)
	for _, myflags := range flg {
		if myflags.Type == flags.Dash && myflags.Name == "e" {
			cfg.embed = true
		} else if myflags.Type == flags.Dash && myflags.Name == "t" {
			cfg.embed = false
		} else if myflags.Type == flags.Dash && myflags.Name == "m" {
			cfg.useDMs = true
		} else if myflags.Type == flags.Dash && myflags.Name == "c" {
			cfg.useDMs = false
		}
	}

	if cfg.useDMs {
		if cfg.embed {
			session.ChannelMessageSendEmbed(message.Author.ID, getBotInfoAsEmbed(session))
		} else {
			session.ChannelMessageSend(message.Author.ID, getBotInfoAsText())
		}
	} else {
		if cfg.embed {
			session.ChannelMessageSendEmbed(message.ChannelID, getBotInfoAsEmbed(session))
		} else {
			session.ChannelMessageSend(message.ChannelID, getBotInfoAsText())
		}
	}
}

func getBotInfoAsEmbed(session *dsg.Session) *dsg.MessageEmbed {
	user, _ := session.User("@me")
	return &dsg.MessageEmbed{
		Author:      &dsg.MessageEmbedAuthor{},
		Color:       0x073642,
		Title:       "Bot Information",
		Description: "A list of commands can be brought up with `" + f.MyBot.Prefs.Prefix + "help <command>`.",
		Thumbnail: &dsg.MessageEmbedThumbnail{
			URL:    user.AvatarURL("64"),
			Width:  64,
			Height: 64,
		},
		Fields: []*dsg.MessageEmbedField{
			&dsg.MessageEmbedField{
				Name:   "Bot Version",
				Value:  "Version " + f.MyBot.Prefs.Version + ".",
				Inline: true,
			},
			&dsg.MessageEmbedField{
				Name:   "Link to bot framework:",
				Value:  "https://github.com/whitman-colm/discord-public",
				Inline: true,
			},
			&dsg.MessageEmbedField{
				Name:   "Command info:",
				Value:  "Version " + cfg.version + ".",
				Inline: true,
			},
		},
	}
}

func getBotInfoAsText() string {
	return "```" + `Bot information:
A list of commands can be brought up with ` + "`" + f.MyBot.Prefs.Prefix + "help <command>`" + `

Bot github link : https://github.com/whitman-colm/discord-public
Bot version     : ` + f.MyBot.Prefs.Version + `
Command version : ` + cfg.version + "```"
}
