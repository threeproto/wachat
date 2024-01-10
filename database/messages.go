package database

import (
	"errors"
	"wachat/params"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func SaveMessage(db *sqlx.DB, msg params.Message) error {
	messages := []params.Message{}
	err := db.Select(&messages, "select * from messages where message_hash = ?", msg.Hash)
	if err != nil {
		return err
	}

	if len(messages) > 0 {
		return errors.New("card name is already used")
	}

	_, err = db.Exec(
		"insert into messages (user_name, content, timestamp, message_hash, is_stored) values (?, ?, ?, ?, ?)",
		msg.Name, msg.Content, msg.Timestamp, msg.Hash, false,
	)
	if err != nil {
		return err
	}

	return nil
}
