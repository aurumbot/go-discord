package utils

import (
	dsg "github.com/bwmarrin/discordgo"
	f "github.com/whitman-colm/go-discord"
	"github.com/whitman-colm/go-discord/dat"
	"strconv"
)

func getPerms(session *dsg.Session, message *dsg.MessageCreate) {
	s := session
	m := message.Message

	// Performs permission check
	perm, err := f.HasPermissions(s, m.Message, m.Author.ID, dsg.PermissionManageServer)
	if err != nil {
		dat.Log.Println(err)
		dat.AlertDiscord(s, m, err)
		return
	}
	if !perm {
		s.ChannelMessageSend(m.ChannelID, "Sorry, you do not have permission to use this command.")
		return
	}

	// Fetches guild info
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

	for _, mention := range m.Mentions {
		permsMSG := "Permissions info for user " + mention.Mention() + ":\n```\n"
		multimsg := false //in case message trails over the 2k char limit

		for _, channel := range channels {
			if multimsg {
				permsMSG = "```\n"
				multimsg = false
			}
			cid := channel.ID
			perms, err := s.UserChannelPermissions(mention.ID, cid)
			if err != nil {
				dat.Log.Println(err.Error())
				// Doesn't quit function. best idea? (fix later if
				// issues arise)
			}
			permsMSG += "Channel: " + channel.Name + " (ID " + channel.ID + ") : " + strconv.Itoa(perms) + ".\n"

			if len(permsMSG) > 1900 {
				permsMSG += "```"
				s.ChannelMessageSend(m.ChannelID, permsMSG)
				permsMSG = ""
				multimsg = true
			}
		}
		if len(permsMSG) > 0 {
			permsMSG += "```"
		}
		s.ChannelMessageSend(m.ChannelID, permsMSG)
	}
}
