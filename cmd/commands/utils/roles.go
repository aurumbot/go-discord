package utils

import (
	dsg "github.com/bwmarrin/discordgo"
	f "github.com/whitman-colm/go-discord"
	"github.com/whitman-colm/go-discord/dat"
)

func init() {
	Commands["getroles"] = &f.Command{
		Name:    "Get Server Roles",
		Help:    "Goes through all of the server's roles and posts them and their IDs.",
		Perms:   dsg.PermissionManageChannels,
		Version: "v1.0",
		Action:  getRoles,
	}
}

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
		rolemsg += "Role: " + role.Name + ", ID: " + role.ID + "\n"

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
