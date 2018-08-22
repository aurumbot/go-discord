package cmd

import (
	dsg "github.com/bwmarrin/discordgo"
	f "github.com/whitman-colm/go-discord"
	"github.com/whitman-colm/go-discord/dat"
	"github.com/whitman-colm/go-discord/flags"
	"strings"
)

func init() {
	Cmd["help"] = &f.Command{
		Name: "Command Help Page Search",
		Help: `Info  : The built-in helper to get information about all of the bots commands
Flags:
-c --command <command>	: get help for the specific <command>
-ls --list		: get a list of all available commands
Usage : ` + f.MyBot.Prefs.Prefix + `help -c <command>
	` + f.MyBot.Prefs.Prefix + `help -ls`,
		Perms:   -1,
		Version: "v2.0β",
		Action:  help,
	}
}

/* # Get bot help
* Overcomplecated for little good reason
*
* Note that this function handles responding instead of returning a value to
* its parent to be sent out.
*
* Flags:
* -d  | Sends the result via dm (not working). //TODO: Figure this out
* -ls | Gets a list of available commands given the users perm level
* -c  | Gets the result for a specific command
 */
func help(session *dsg.Session, message *dsg.MessageCreate) {
	msg := strings.Split(message.Content, " ")
	if len(msg) <= 1 {
		h := "Help Page Found:\n```" + Cmd["help"].Name + "\n" + Cmd["help"].Help + "```"
		session.ChannelMessageSend(message.ChannelID, h)
		return
	}

	flagsParsed := flags.Parse(msg)

	// These are some cop-out variables so I don't nest to eternity.
	var (
		cmd  *flags.Flag
		list bool
	)

	//dm := false

	for i := range flagsParsed {
		if flagsParsed[i].Type == flags.Dash && flagsParsed[i].Name == "c" {
			list = false
			cmd = flagsParsed[i]
		} else if flagsParsed[i].Type == flags.Dash && flagsParsed[i].Name == "ls" {
			list = true
		} else if flagsParsed[i].Type == flags.DoubleDash && flagsParsed[i].Name == "command" {
			list = false
			cmd = flagsParsed[i]
		}
	}

	if !list {
		for command, action := range Cmd {
			if cmd.Value == command {
				help := "Help Page Found:\n```" + action.Name + "\n" + action.Help + "```"
				session.ChannelMessageSend(message.ChannelID, help)
				return
			}
		}
		session.ChannelMessageSend(message.ChannelID, "Sorry, but I couldn't find a help page for that command.")
		return
	} else {
		msg := "**Available Commands:**"
		for command, action := range Cmd {
			u, err := f.HasPermissions(session, message.Message, message.Author.ID, action.Perms)
			if err != nil {
				dat.Log.Println(err)
				dat.AlertDiscord(session, message, err)
				return
			}
			if u {
				msg += "\n" + f.MyBot.Prefs.Prefix + command + " : " + action.Name
			}
		}
		msg += "\nUse `" + f.MyBot.Prefs.Prefix + "help -c <command>` to get more info on a command"
		session.ChannelMessageSend(message.ChannelID, msg)
	}
}
