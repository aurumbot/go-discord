package moderation

import (
	dsg "github.com/bwmarrin/discordgo"
	f "github.com/whitman-colm/go-discord"
	//"github.com/whitman-colm/go-discord/dat"
	"github.com/whitman-colm/go-discord/flags"
	"strings"
	"time"
)

// The config struct contains a horribly over-complex list of
// channels and guilds where each trigger should be sent.
// TODO: think this out, remember it needs to contain
// guilds, channels, and triggers though not
// necessarily in that order.
type LogConfig struct {
	Delete     map[string][]string `json:"delete"`
	DeleteBulk map[string][]string `json:"deletebulk"`
	Edit       map[string][]string `json:"edit"`
	Ban        map[string][]string `json:"ban"`
	Pardon     map[string][]string `json:"pardon"`
	Join       map[string][]string `json:"join"`
	Leave      map[string][]string `json:"leave"`
}

// establishes config
func init() {

}

// The handlers need to be defined, however its doubtful if the session will
// have already been established by the time I could shove it all in the init
// function, so instead the main.go will shoot this off after the command
// handler is established.
func DefineLogHandlers(s *dsg.Session) {
	s.AddHandler(logMessageDelete)
	s.AddHandler(logMessageDeleteBulk)
	s.AddHandler(logMessageEdit)
	s.AddHandler(logGuildBanAdd)
	s.AddHandler(logGuidBanRemove)
	s.AddHandler(logGuildMemberAdd)
	s.AddHandler(logGuildMemberRemove)
}

// Handles deleted messages
func logMessageDelete(s *dsg.Session, m *dsg.MessageDelete) {
	str := fmt.Sprintf("`[%s] Message deleted`",
		getFormattedTime())
	guild, err := f.GetGuild(s, m.Message)
	if err != nil {
		dat.AlertDiscord(s, m, err)
		dat.Log.Println(err)
		return
	}

	if val, ok := LogConfig.Delete[guild.ID]; ok {
		for _, channel := range val {
			s.ChannelMessageSend(val, str)
		}
	}
}

// Handles deleted messages in bulk. (apparently different handlers)
func logMessageDeleteBulk(s *dsg.Session, m *dsg.MessageDelete) {

}

// Handles edited messages
func logMessageEdit(s *dsg.Session, m *dsg.MessageEdit) {

}

// Handles user bans
func logGuildBanAdd(s *dsg.Session, m *dsg.GuildBanAdd) {

}

// Handles user pardons
func logGuidBanRemove(s *dsg.Session, m *dsg.GuildBanRemove) {

}

// Handles user join
func logGuildMemberAdd(s *dsg.Session, m *dsg.GuildMemberAdd) {

}

// Handles user leaves
func logGuildMemberRemove(s *dsg.Session, m *dsg.GuildMemberRemove) {

}

// formatTime is a small function that gets the current time and
// returns it in the standard format for logs. This is to prevent a million
// character wide lines from repeating myself.
func getFormattedTime() string {
	return time.Now().Format("2006-01-02@15:04:05.000 (-0700 MST)")
}
