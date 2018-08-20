package utils

import (
	dsg "github.com/bwmarrin/discordgo"
	f "github.com/whitman-colm/go-discord"
	"github.com/whitman-colm/go-discord/dat"
	"strconv"
)

/* # Get server roles
* A g-d impossibility.
*
* Parameters/return values:
* This function complies with the foundation's action function protocol.
* For documentation on that, please see https://github.com/whitman-colm/discord-public
*
* TODO: Make a godoc for our nonsence.
*
* NOTE: If you print this into a discord chat, it WILL mention @everyone
 */
func getRoles(session *dsg.Session, message *dsg.MessageCreate) {
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

	// Retrieves roles, puts them into a slice
	guild, err := f.GetGuild(f.DG, m)
	if err != nil {
		dat.Log.Println(err.Error())
		return
	}
	roles, err := f.DG.GuildRoles(guild.ID)
	if err != nil {
		dat.Log.Println(err.Error())
		return
	}

	// Sends out the info, this is complex in case the character count is
	// over 2000, which won't send.
	rolemsg := "Role list for server:\n```\n"
	multimsg := false //in case message trails over the 2k char limit
	for _, role := range roles {
		if multimsg {
			rolemsg = "```\n"
			multimsg = false
		}
		rolemsg += "Role: " + role.Name + ", ID: " + channel.ID + "\n"

		if len(rolemsg) > 1900 {
			rolemsg += "```"
			s.ChannelMessageSend(m.ChannelID, rolemsg)
			rolemsg = ""
			multimsg = true
		}
	}
	if len(rolemsg) > 0 {
		rolemsg += "```"
		s.ChannelMessageSend(m.ChannelID, rolemsg)
	}
}
