package cmd

import (
	"errors"
	dsg "github.com/bwmarrin/discordgo"
	f "github.com/whitman-colm/go-discord"
	"github.com/whitman-colm/go-discord/dat"
	"strings"
)

/* # MessageCreate
* The world's bigest switch statment
*
* This is a very big switch statment run commands. It reads all the messages in
* all the servers its in, determines which ones are commands, and then sees
* what in all the commands mean and then takes the appropriate action.
*
* Parameters:
* - s (type *discordgo.Session) | The current running discord session,
*     (discordgo needs that always apparently)
* - m (type *discordgo.Message) | The message thats to be acted upon.
*
* TODO: See if it can be made so it doesn't have to read every single message
*       ever.
*
* TODO: Break this one function up to smaller functions that only run if a user
*       has a certain role
*
* NOTE: Please delegate what the command actually does to a function. This
*       method should only be used to determine what the user is acutally
*       trying to do.
 */
func MessageCreate(s *dsg.Session, m *dsg.MessageCreate) {
	// The message is checked to see if its a command and can be run
	canRunCommand, err := canTriggerBot(s, m.Message)
	if err != nil {
		dat.Log.Println(err.Error())
		dat.AlertDiscord(s, m, err)
		return
	}
	if !canRunCommand {
		return
	}

	// Removing case sensitivity:
	messageSanatized := strings.ToLower(m.Content)

	// The prefix is cut off the message so the commands can be more easily handled.
	var msg []string
	if strings.HasPrefix(m.Content, f.MyBot.Prefs.Prefix) {
		msg = strings.SplitAfterN(messageSanatized, f.MyBot.Prefs.Prefix, 2)
		m.Content = msg[1]
		//TODO: Check if there is a way to use a mention() method of discordgo rather than
		//this string frankenstein
	} else if strings.HasPrefix(m.Content, "<@!"+f.MyBot.Auth.ClientID+">") {
		msg = strings.SplitAfterN(messageSanatized, "<@!"+f.MyBot.Auth.ClientID+">", 2)
		m.Content = strings.TrimSpace(msg[1])
	} else {
		err := errors.New("Message passed 'can run' checks but does not start with prefix:\n" + m.Content)
		dat.Log.Println(err.Error())
		dat.AlertDiscord(s, m, err)
		return
	}

	message := strings.Split(m.Content, " ")

	// Now the message is run to see if its a valid command and acted upon.
	for command, action := range Cmd {
		if message[0] == command {
			if action.Perms != -1 {
				perm, err := f.HasPermissions(s, m.Message, m.Author.ID, action.Perms)
				if err != nil {
					dat.Log.Println(err)
					dat.AlertDiscord(s, m, err)
					return
				}
				if !perm {
					s.ChannelMessageSend(m.ChannelID, "Sorry, you do not have permission to use this command.")
					return
				}
			}
			action.Action(s, m)
			return
		}
	}

	if strings.Contains(m.Message.Content, "@") {
		s.ChannelMessageSend(m.ChannelID, "Sorry <@"+m.Message.Author.ID+">, but I don't understand.")
	} else {
		s.ChannelMessageSend(m.ChannelID, "Sorry <@"+m.Message.Author.ID+">, but I don't understand what you mean by \"`"+m.Message.Content+"`\".")
	}

}

/* # Check if user can run command
* This switch statment makes sure the bot runs when its triggered and the user has the perms to trigger it.
* Prevents:
* - Bot posted something that would trigger itself, possibly creating an infinite loop
* - Message posted doesn't have the bot's prefix
* - Command was posted in a channel where the bot shouldn't respond to commands
* - Bot whitelists channels and the command was run in a channel not on the whitelist.
* - Users with a blacklisted role from running the bot
*
* NOTE: Users who have "admin" roles (according to the bot's json data) or
*       permissions will have the ability to run commands regardless of any
*       other rules
*
* NOTE: IF THESE CONDITIONS ARE MET THEN NO ERROR WILL BE SENT TO EITHER DISCORD OR LOGGED.
* THIS IS BY DESIGN. DON'T CHANGE IT THINKING I WAS JUST LAZY.
 *
func canTriggerBot(s *dsg.Session, m *dsg.Message) (bool, error) {
	if m.Author.Bot {
		return false, nil
	}

	admin, err := f.HasPermissions(s, m, m.Author.ID, dsg.PermissionAdministrator)
	if err != nil {
		dat.Log.Println(err)
		return false, err
	}

	switch true {
	case m.Author.ID == s.State.User.ID:
		return false, nil
	//TODO: look at this stupid line. that seems like it shouldn't work.
	case !strings.HasPrefix(m.Content, f.MyBot.Prefs.Prefix) && !strings.HasPrefix(m.Content, "<@!"+f.MyBot.Auth.ClientID+">"):
		return false, nil
	case admin:
		return true, nil
	case f.Contains(f.MyBot.Perms.BlacklistedChannels, m.ChannelID) == true:
		return false, nil
	case f.MyBot.Perms.WhitelistChannels && !f.Contains(f.MyBot.Perms.WhitelistedChannels, m.ChannelID):
		return false, nil
	}
	for _, b := range f.MyBot.Users.BlacklistedRoles {
		guild, err := f.GetGuild(s, m)
		if err != nil {
			return false, err
		}
		member, err := s.GuildMember(guild.ID, m.Author.ID)
		if err != nil {
			return false, err
		}
		blacklisted := f.Contains(member.Roles, b)
		if blacklisted {
			return false, nil
		}
	}
	return true, nil
}
*/
