package moderation

import (
	"github.com/boltdb/bolt"
	f "github.com/whitman-colm/go-discord"
	"github.com/whitman-colm/go-discord/dat"
	"time"
)

func init() {
	now := time.Now()
	inflog, err := bolt.Open(moderation.db, 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return
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
				//TODO: make actions reverse here.
				return nil
			} else {
				//TODO: make actions reverse here but a goroutine.
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

	// Checks un-decayed infractions to see if any have decayed.
	if err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("decay"))
		if err != nil {
			dat.Log.Println(err)
			return err
		}

		if err := b.ForEach(func(k, v, []byte) error {
			until, _ := time.Parse("2006-01-02@15:04:05", v.Until)
			if now.After(until) {
				//TODO: make actions reverse here.
				return nil
			} else {
				//TODO: make actions reverse here but a goroutine.
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

func Log() error {
	b, err := bolt.Open(infractions.db, 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		dat.Log.Println(err)
		return err
	}
}
