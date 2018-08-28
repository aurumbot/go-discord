package moderation

import (
	"errors"
	"github.com/boltdb/bolt"
	f "github.com/whitman-colm/go-discord"
	"github.com/whitman-colm/go-discord/dat"
	"time"
)

type config struct {
	MuteRole    string   `json:"muted"`
	BannedWords []string `json:"swears"`
}

/* Moderation logs
* The moderation package comes with built-in logs using json files as a pseudo-
* database. Each "incident" consists of the following fields:
* - ID     : The infraction ID, this is a unique id for each incident
* - User   : The ID of the user who was acted upon.
* - Issuer : The moderator who took action against the user.
* - Time   : The time the action was taken.
* - Until  : When the action ends for a temp-mute or temp-ban
* - Reason : Why the action was issued
* - Action : The ID of the action that was taken. See constants for more info
* - Decay  : Statute of limitations. Prevents the bot from being overly-harsh
*            to repeat offences with large gaps of time between infractions.
 */
type incident struct {
	ID       int    `json:"id"`
	User     string `json:"user"`
	Actor    string `json:"actor"`
	Guild    string `json:"guild"`
	Time     string `json:"time"`
	Duration string `json:"duration"`
	Reason   string `json:"reason"`
	Action   int    `json:"action"`
	Decay    string `json:"decay"`
}

// Defined enum codes for each action a moderator can take.
const (
	warn     = 0 // an official warning, no immediate action.
	delmsg   = 1 // message is deleted (censorship by bot).
	mute     = 2 // user is muted
	kick     = 3 // user is kicked
	tempban  = 4 // user is temporarily banned
	permaban = 5 // user is permanently banned
	restore  = 6 // user has their mute revoked
	pardon   = 7 // user is pardoned from a ban
)

var (
	Commands = make(map[string]*f.Command)
	ps       string
	cfg      *config
	logs     *modlogs
)

func init() {
	ps = dat.OSCheck()
	dat.Load("moderation"+ps+"config.json", &cfg)
}

func reverseAction(after time.Duration, b *bolt.Bucket, inf *incident) {
	time.Sleep(after)
	switch inf.Action {
	case 2: //mute
		//remove mute
	case 4: //tempban
		//unban & re-invite
	default:
		dat.Log.Println(errors.New(fmt.Sprintf("Cannot file Action ID %d (ID %s)",
			inf.Action, inf.ID)))
	}
	if err := b.Delete([]byte(inf.ID)); err != nil {
		dat.Log.Println(err)
	}
}
