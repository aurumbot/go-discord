package util

import (
	dsg "github.com/bwmarrin/discordgo"
	f "github.com/whitman-colm/go-discord"
	"github.com/whitman-colm/go-discord/dat"
	//"strconv"
)

var Commands = make(map[string]*f.Command)

func init() {
	Commands["getroles"] = &f.Command{
		Name:   "Get Server Roles",
		Help:   "Goes through all of the server's roles and posts them and their IDs.",
		Action: getRoles,
	}
	/*	Commands["getperms"] = &f.Command{
				Name: "Get permission for users",
				Help: `Gets the value of permission integer for users in each of the server's channels.
		The permissions are set as 53-bit integers calculated using bitwise operations.
		For more info see https://discordapp.com/developers/docs/topics/permissions and
		https://discordapi.com/permissions.html.
		User mentions should be passed as arguments. Multiple users at a time are supported.

		Usage: getperms @someuser @otheruser`,
				Action: getPerms,
			}*/
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
	roleMSG := "Server role list:\n```\n"
	for _, r := range roles {
		roleMSG += "Name: " + r.Name + "; ID: " + r.ID + ";\n"
	}
	roleMSG += "```"
	s.ChannelMessageSend(m.ChannelID, roleMSG)
}

/*func getPerms(session *dsg.Session, message *dsg.MessageCreate) {
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
}*/
