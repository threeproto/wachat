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
		return errors.New("message already exists")
	}

	_, err = db.Exec(
		"insert into messages (user_name, content, timestamp, message_hash, is_stored, waku_timestamp) values (?, ?, ?, ?, ?, ?)",
		msg.Name, msg.Content, msg.Timestamp, msg.Hash, false, msg.WakuTimestamp,
	)
	if err != nil {
		return err
	}

	return nil
}

func GetUnstoredMessages(db *sqlx.DB) ([]params.Message, error) {
	messages := []params.Message{}
	err := db.Select(&messages, "select * from messages where is_stored = false")
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func UpdateStoredMessage(db *sqlx.DB, hash string) error {
	_, err := db.Exec("update messages set is_stored = true where message_hash = ?", hash)
	if err != nil {
		return err
	}

	return nil
}
