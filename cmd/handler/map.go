package cmd

import (
	f "github.com/whitman-colm/go-discord"
	"github.com/whitman-colm/go-discord-customcmds/rolepicker"
	"github.com/whitman-colm/go-discord/cmd/commands/info"
	"github.com/whitman-colm/go-discord/cmd/commands/moderation"
	"github.com/whitman-colm/go-discord/cmd/commands/ping"
	"github.com/whitman-colm/go-discord/cmd/commands/utils"
)

var Cmd = map[string]*f.Command{}

/* FOR THE PERSON RUNNING THIS BOT: Adding packages to the command list
* As of now, the bot has no commands set to it so while it may boot up, it
* won't actually do anything. You will need to add the maps of the command
* modules you have imported or made into the main Cmd map. To do this, add
* each of the command's public map[string]*f.Command type into the following
* init statment. 2 commands, `info` and `ping` have already been added to help
* show what you need to do:
 */

func init() {
	for key, value := range ping.Commands {
		Cmd[key] = value
	}
	for key, value := range info.Commands {
		Cmd[key] = value
	}
	for key, value := range utils.Commands {
		Cmd[key] = value
	}
	for key, value := range moderation.Commands {
		Cmd[key] = value
	}
	for key, value := range rolepicker.Commands {
		Cmd[key] = value
	}
	//for key, value := range PACKAGENAME.Commands {
	//        Cmd[key] = value
	//}
}
