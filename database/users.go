package database

import (
	"errors"
	"wachat/params"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func SaveUser(db *sqlx.DB, user params.User) error {
	users := []params.User{}
	err := db.Select(&users, "select * from users where name = ?", user.Name)
	if err != nil {
		return err
	}

	if len(users) > 0 {
		return errors.New("user name already exists")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	tx.Exec(`UPDATE users SET selected = false`)

	_, err = tx.Exec(
		"insert into users (name, selected) values (?, ?)",
		user.Name, user.Selected,
	)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func GetSelectedUser(db *sqlx.DB) (params.User, error) {
	users := []params.User{}
	err := db.Select(&users, "select * from users where selected = true")
	if err != nil {
		return params.User{}, err
	}

	if len(users) == 0 {
		return params.User{}, errors.New("no selected user")
	}
	if len(users) > 1 {
		return params.User{}, errors.New("more than one selected user")
	}

	return users[0], nil
}
