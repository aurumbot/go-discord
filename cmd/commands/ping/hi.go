package ping

import (
	dsg "github.com/bwmarrin/discordgo"
	f "github.com/whitman-colm/go-discord"
)

var Commands = make(map[string]*f.Command)

func init() {
	Commands["hi"] = &f.Command{
		Name:    "Ping The Bot",
		Help:    "Pings the bot to see if its online.",
		Perms:   -1,
		Version: "v1.1",
		Action:  ping,
	}
}

func ping(session *dsg.Session, message *dsg.MessageCreate) {
	session.ChannelMessageSend(message.ChannelID, "Hello, World!")
}
