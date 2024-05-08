package database

import (
	"database/sql"
	"errors"
	. "github.com/Big-Iron-Cheems/WASAPhoto/service/model"
)

// GetBansList returns the list of banned users.
func (db *appdbimpl) GetBansList(userId uint) ([]User, error) {
	rows, err := db.c.Query(`
        SELECT bannedUserId FROM Bans
        WHERE userId = ?`,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var bannedIds []uint
	for rows.Next() {
		var bannedId uint
		if err = rows.Scan(&bannedId); err != nil {
			return nil, err
		}
		bannedIds = append(bannedIds, bannedId)
	}

	bannedUsers := make([]User, 0)
	for _, bannedId := range bannedIds {
		row := db.c.QueryRow(`
            SELECT username FROM Users
            WHERE userId = ?`,
			bannedId,
		)
		var username string
		if err = row.Scan(&username); err != nil {
			return nil, err
		}
		bannedUsers = append(bannedUsers, User{UserId: bannedId, Username: username})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bannedUsers, nil
}

// GetBansCount returns the number of banned users.
func (db *appdbimpl) GetBansCount(userId uint) (uint, error) {
	var count uint
	err := db.c.QueryRow(`
        SELECT COUNT(*) FROM Bans
        WHERE userId = ?`,
		userId,
	).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

/*
GetBanStatus returns the ban status between two users.

If the user with id `userId` has banned the user with id `targetUserId`, it returns true.
*/
func (db *appdbimpl) GetBanStatus(userId uint, targetUserId uint) (bool, error) {
	var res bool
	if err := db.c.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM Bans
            WHERE userId = ? AND bannedUserId = ?
        )`,
		userId, targetUserId,
	).Scan(&res); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, &UserNotFoundByIdError{UserId: targetUserId}
		}
	}
	return res, nil
}

/*
BanUser adds a user to the ban list.

Also, breaks the follow relationship between the two users,
deletes all comments and likes of the banned user from the requester's photos.
*/
func (db *appdbimpl) BanUser(userId uint, targetUserId uint) error {
	// Ban the target user
	_, err := db.c.Exec(`
        INSERT INTO Bans (bannedUserId, userId)
        VALUES (?, ?)`,
		targetUserId, userId,
	)
	if err != nil {
		return err
	}

	// Unfollow the target user
	err = db.UnfollowUser(userId, targetUserId)
	if err != nil {
		return err
	}

	// Remove the user from the target's followers
	err = db.UnfollowUser(targetUserId, userId)
	if err != nil {
		return err
	}

	// Remove all comments of the target user from the requester's photos
	err = db.DeleteCommentsByBannedUser(userId, targetUserId)
	if err != nil {
		return err
	}

	// Remove all likes of the target user from the requester's photos
	err = db.DeleteLikesByBannedUser(userId, targetUserId)
	if err != nil {
		return err
	}

	return nil
}

// UnbanUser removes a user from the ban list.
func (db *appdbimpl) UnbanUser(userId uint, targetUserId uint) error {
	_, err := db.c.Exec(`
        DELETE FROM Bans
        WHERE bannedUserId = ? AND userId = ?`,
		targetUserId, userId,
	)
	if err != nil {
		return err
	}

	return nil
}

// DeleteCommentsByBannedUser Remove all comments of a banned user from the photos of the requester.
func (db *appdbimpl) DeleteCommentsByBannedUser(userId uint, bannedUserId uint) error {
	_, err := db.c.Exec(`
		DELETE FROM Comments
		WHERE ownerId = ? AND photoId IN (
			SELECT photoId FROM Photos WHERE ownerId = ?
		)`,
		bannedUserId, userId,
	)
	if err != nil {
		return err
	}

	return nil
}

// DeleteLikesByBannedUser Remove all likes of a banned user from the photos of the requester.
func (db *appdbimpl) DeleteLikesByBannedUser(userId uint, bannedUserId uint) error {
	_, err := db.c.Exec(`
        DELETE FROM Likes
        WHERE userId = ? AND photoId IN (
            SELECT photoId FROM Photos WHERE ownerId = ?
        )`,
		bannedUserId, userId,
	)
	if err != nil {
		return err
	}

	return nil
}
