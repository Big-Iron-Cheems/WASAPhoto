package database

import (
	"database/sql"
	"errors"
	. "wasaphoto.uniroma1.it/wasaphoto/service/model"
)

// CreateUser creates a new user in the database.
//
// If the user already exists, it returns the user schema.
// Otherwise, it creates it and returns the user schema.
func (db *appdbimpl) CreateUser(user User) (User, error) {
	res, err := db.c.Exec("INSERT INTO Users(username) VALUES (?)", user.Username)
	if err != nil {
		if err = db.c.QueryRow(`SELECT userId, username FROM Users WHERE username = ?`, user.Username).Scan(&user.UserId, &user.Username); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return user, errors.New("user does not exist")
			}
		}
		return user, nil
	}
	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return user, err
	}
	user.UserId = uint(lastInsertID)
	return user, nil
}

// GetUserProfile retrieves a user from the database given its username.
func (db *appdbimpl) GetUserProfile(username string) (User, error) {
	var user User
	if err := db.c.QueryRow(`SELECT userId, username FROM Users WHERE username = ?`, username).Scan(&user.UserId, &user.Username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errors.New("user does not exist")
		}
	}
	return user, nil
}

// GetMyStream retrieves the content stream of a user from the database given its id.
func (db *appdbimpl) GetMyStream(userId uint) ([]Photo, error) {
	panic("not implemented")
}

// SetMyUsername updates the username of a user in the database given its id.
func (db *appdbimpl) SetMyUsername(userId uint, newUsername string) error {
	panic("not implemented")
}
