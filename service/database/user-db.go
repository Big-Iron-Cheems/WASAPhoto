package database

import (
	"database/sql"
	"errors"
	"fmt"
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
				return user, fmt.Errorf("user `%s` does not exist", user.Username)
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

// GetUserProfile retrieves a user's profile from the database.
func (db *appdbimpl) GetUserProfile(user User) (User, error) {
	if err := db.c.QueryRow(`SELECT userId, username FROM Users WHERE username = ?`, user.Username).Scan(&user.UserId, &user.Username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, fmt.Errorf("user `%s` does not exist", user.Username)
		}
	}
	return user, nil
}

// GetMyStream retrieves the content stream of a user from the database given its id.
func (db *appdbimpl) GetMyStream(user User) ([]Photo, error) {
	panic("not implemented") // TODO
}

// SetMyUsername updates the username of a user in the database given its id.
func (db *appdbimpl) SetMyUsername(user User, currentUsername string) (User, error) {
	res, err := db.c.Exec("UPDATE Users SET username = ? WHERE username = ?", user.Username, currentUsername)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, fmt.Errorf("user `%s` does not exist", user.Username)
		}
		return user, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return user, err
	}
	if rowsAffected == 0 {
		return user, fmt.Errorf("user `%s` does not exist", user.Username)
	}

	// Fetch the updated user from the database
	updatedUser, err := db.GetUserProfile(User{Username: user.Username}) // TODO: find if this can be sent in the req along the username
	if err != nil {
		return user, fmt.Errorf("unable to fetch updated user profile: %w", err)
	}

	return updatedUser, nil
}

/*
GetFollowStatus returns the follow status between two users.

If the user with id `userId` follows the user with id `targetUserId`, it returns true.
*/
func (db *appdbimpl) GetFollowStatus(userId uint, targetUserId uint) (bool, error) {
	var res bool
	if err := db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM Followers WHERE followerUserId = ? AND followingUserId = ?)`, userId, targetUserId).Scan(&res); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, errors.New("user does not exist")
		}
	}
	return res, nil
}

/*
GetUserBanStatusByRequester returns the ban status between two users.

If the user with id `targetUserId` is banned by the user with id `userId`, it returns true.
*/
func (db *appdbimpl) GetUserBanStatusByRequester(userId uint, targetUserId uint) (bool, error) {
	var res bool
	if err := db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM Bans WHERE userId = ? AND bannedUserId = ?)`, userId, targetUserId).Scan(&res); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, errors.New("user does not exist")
		}
	}
	return res, nil
}

/*
GetRequesterBannedByUser returns the ban status between two users.

If the user with id `userId` is banned by the user with id `targetUserId`, it returns true.
*/
func (db *appdbimpl) GetRequesterBannedByUser(userId uint, targetUserId uint) (bool, error) {
	var res bool
	if err := db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM Bans WHERE userId = ? AND bannedUserId = ?)`, targetUserId, userId).Scan(&res); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, errors.New("user does not exist")
		}
	}
	return res, nil
}
