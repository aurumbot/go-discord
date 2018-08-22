package utils

import (
	dsg "github.com/bwmarrin/discordgo"
	f "github.com/whitman-colm/go-discord"
	"github.com/whitman-colm/go-discord/dat"
)

func init() {
	Commands["getchannels"] = &f.Command{
		Name:    "Get a List of Channels and Their IDs",
		Help:    "",
		Perms:   dsg.PermissionManageChannels,
		Action:  getChannels,
		Version: "v1.0",
	}
}

func getChannels(session *dsg.Session, message *dsg.MessageCreate) {
	s := session
	m := message.Message

	guild, err := f.GetGuild(s, m)
	if err != nil {
		dat.Log.Println(err.Error())
		return
	}
	channels, err := s.GuildChannels(guild.ID)
	if err != nil {
		dat.Log.Println(err.Error())
		return
	}

	multimsg := false
	msg := "**Channel List:**\n```"
	for _, channel := range channels {
		if multimsg {
			msg = "```\n"
			multimsg = false
		}
		msg += "Channel: " + channel.Name + " (ID :" + channel.ID + ") .\n"

		if len(msg) > 1900 {
			msg += "```"
			s.ChannelMessageSend(m.ChannelID, msg)
			msg = ""
			multimsg = true
		}
	}
	if len(msg) > 0 {
		msg += "```"
	}
	s.ChannelMessageSend(m.ChannelID, msg)
}
