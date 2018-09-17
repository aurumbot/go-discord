package moderation

import (
	"github.com/boltdb/bolt"
	f "github.com/whitman-colm/go-discord"
	"github.com/whitman-colm/go-discord/dat"
	"time"
	"github.com/whitman-colm/go-discord/flags"
)

func init() {
	err := reloadBuckets()
	if err != nil {
		// Er... yeah if something goes wrong here then you're ******.
		panic(err)
	}
}



/* Reloads the buckets from the primary "logs"
*
* TODO: Kill off active goroutines and relaunch them here.
 */
func reloadBuckets() err {
	now := time.Now()
	defer db.Close()

	db, err := bolt.Open(moderation.db, 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	// Checks active infractions to see if any have expired.
	if err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("active"))
		if err != nil {
			dat.Log.Println(err)
			return err
		}

		if err := b.ForEach(func(k, v, []byte) error {
			until, _ := time.Parse("2006-01-02@15:04:05", v.Until)
			if now.After(until) {
				reverseAction((time.Nanosecond), b, v)
				return nil
			} else {
				go reverseAction(time.Until(until), b, v)
				return nil
			}
		}); err != nil {
			dat.Log.Println(err)
			return err
		}
	}); err != nil {
		dat.Log.Println(err)
		return
	}
}

/* The thing that logs things to the database
* log is a simple wrapper to save an infraction to the appropriate
* bolt bucket(s).
*
* Parameters:
* - inf incident{} : the incedent to be logged.
*
* Returns:
* - _ error : an error, if it came up. The error has already
*		been logged. This is just to pass to an AlertDiscord()
*
* TODO: Yeah just another reminder that when working on decay, a third bucket
*	will be needed.
 */
func log(inf incident) error {
	db, err := bolt.Open(moderation.db, 0666, &bolt.Options{Timeout: 1 * time.Second})
	defer db.Close()

	if err = db.Update(func(tx *bolt.Tx) error {
		// Figures out if there is a "logs" bucket and
		// stores values into it
		logs, err := tx.CreateBucketIfNotExists([]byte("logs"))
		if err != nil {
			return err
		}
		if err = logs.Put([]byte(inf.ID), []byte(inf)); err != nil {
			return err
		}

		// does the same thing above but for a bucket named after the user
		// who had action taken on them.
		logs, err = tx.CreateBucketIfNotExists([]byte(inf.User))
		if err != nil {
			return err
		}
		if err = logs.Put([]byte(inf.ID), []byte(inf)); err != nil {
			return err
		}

		// does the same thing above above but for a bucket named after the guild
		// this allows for "multi-server" logs.
		logs, err = tx.CreateBucketIfNotExists([]byte(inf.Guild))
		if err != nil {
			return err
		}
		if err = logs.Put([]byte(inf.ID), []byte(inf)); err != nil {
			return err
		}

		if inf.Action != 2|4 { // This should be a check to see if a timed mute/ban was enacted
			return nil
		}

		// Store values into the active punishment list (timed mutes, tempbans)
		logs, err = tx.CreateBucketIfNotExists([]byte("active"))
		if err != nil {
			return err
		}
		if err = logs.Put([]byte(inf.ID), []byte(inf)); err != nil {
			return err
		}
		return nil

	}); err != nil {
		dat.Log.Println(err)
		return err
	}
	return nil
}

/* Searches the logs
* This is a simple search tool that finds infractions with components that fit
* the flagged terms. Iterates over the bucket of the guild that triggered it.
* Arguments:
* - guild (string) : the guild of the server that searched, this filters
*		     for infractions by that server only.
* - f []flags.Flag : the flag arguments as passed by the sercher. As listed:
*	- <-t|--type> <and|or> : the logic that the infractions must fulfill to
*		be added to the list. AND being all values, OR being one value.
*		Default is "and"
*	- [-u|--user] <string> : the id of the user who had action taken upon
*		them.
*	- <-a|--after> <time.Time> : the earliest an infraction can occur to be
*		counted. default is 1 week ago from the current time.
*	- [-b|--before] <time.Time> : the latest an infraction can occur to be
*		counted.
*	- [-i|--action] <int> : the enum code for an action as defined by the
*		backend.go.
*	- [-m|--moderator] <string> : the id moderator who initiated the action
*	TODO: Make this more intuitive through regex.
*
* Returns:
* - []infraction : all the infractions that fit the parameters
* - error : an error. Note it has already been logged, just alert discord.
*/
func search(guild string, flags []flags.Flag) ([]infraction, error) {
	if len(f) == 0 {
		return [], nil
	}
	infs := make([]infraction)
	db, err := bolt.Open(moderation.db, 0444, nil)
	defer db.Close()
	if err != nil {
		dat.Log.Println(err)
		return [], err
	}
	if err = db.Update(func(tx *bolt.Tx) error {
		logs, err = tx.Bucket([]byte(guild))
		if err != nil {
			return err
		}
		logs.ForEach(func(_, v []byte) error {
			
			for _, flag := range flags {
				switch inf := inicident{v} {
					default:
						
				}
			}
		}); err != nil {
			return err
		}
	}); err != nil {
		dat.Log.Println(err)
		return [], err
	}
	return infs, nil
}


