package utils

import (
	dsg "github.com/bwmarrin/discordgo"
	f "github.com/whitman-colm/go-discord"
	"github.com/whitman-colm/go-discord/dat"
	"strconv"
)

var Commands = make(map[string]*f.Command)

func init() {
	Commands["getroles"] = &f.Command{
		Name:   "Get Server Roles",
		Help:   "Goes through all of the server's roles and posts them and their IDs.",
		Action: getRoles,
	}
	Commands["getperms"] = &f.Command{
		Name: "Get Permission for Users",
		Help: `Gets the value of permission integer for users in each of the server's channels.
		The permissions are set as 53-bit integers calculated using bitwise operations.
		For more info see https://discordapp.com/developers/docs/topics/permissions and
		https://discordapi.com/permissions.html.
		User mentions should be passed as arguments. Multiple users at a time are supported.

		Usage: getperms @someuser @otheruser`,
		Action: getPerms,
	}
	Command["getchannels"] = &f.Command{
		Name:   "Get a List of Channels and Their IDs",
		Help:   "",
		Action: getChannels,
	}
}
