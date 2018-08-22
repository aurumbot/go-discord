package ping

import (
	"fmt"
	dsg "github.com/bwmarrin/discordgo"
	gp "github.com/sparrc/go-ping"
	f "github.com/whitman-colm/go-discord"
)

func init() {
	Commands["ping"] = &f.Command{
		Name:    "Get Connection Info About Bot",
		Help:    "Pings the bot to see if its online, then reads out networking details.",
		Perms:   -1,
		Version: "v1.0-beta.1",
		Action:  pingDetail,
	}
}

func pingDetail(session *dsg.Session, message *dsg.MessageCreate) {
	pinger, err := gp.NewPinger("www.discordapp.com")
	if err != nil {
		panic(err)
	}
	pinger.Count = 4
	pinger.Run()
	stats := pinger.Statistics()
	var send string
	send += fmt.Sprintf("Pong!\n%d packets transmitted, %d received. (%v%% loss)\n",
		stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
	send += fmt.Sprintf("round-trip min/avg/max/stddev = %v/%v/%v/%v",
		stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)

	session.ChannelMessageSend(message.ChannelID, send)
}
