package database

import (
	"database/sql"
	"errors"
	. "wasaphoto.uniroma1.it/wasaphoto/service/model"
)

// GetFollowersList returns the list of followers of a user.
func (db *appdbimpl) GetFollowersList(userId uint) ([]User, error) {
	rows, err := db.c.Query(`
        SELECT followerUserId FROM Followers
        WHERE followingUserId = ?`,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followerIds []uint
	for rows.Next() {
		var followerId uint
		if err = rows.Scan(&followerId); err != nil {
			return nil, err
		}
		followerIds = append(followerIds, followerId)
	}

	followers := make([]User, 0)
	for _, followerId := range followerIds {
		row := db.c.QueryRow(`
            SELECT username FROM Users
            WHERE userId = ?`,
			followerId,
		)
		var username string
		if err = row.Scan(&username); err != nil {
			return nil, err
		}
		followers = append(followers, User{UserId: followerId, Username: username})
	}

	return followers, nil
}

// GetFollowersCount returns the number of followers of a user.
func (db *appdbimpl) GetFollowersCount(userId uint) (uint, error) {
	var count uint
	if err := db.c.QueryRow(`
        SELECT COUNT(*) FROM Followers
        WHERE followingUserId = ?`,
		userId,
	).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, &UserNotFoundByIdError{UserId: userId}
		}
	}
	return count, nil
}

// GetFollowingList returns the list of users followed by a user.
func (db *appdbimpl) GetFollowingList(userId uint) ([]User, error) {
	rows, err := db.c.Query(`
        SELECT followingUserId FROM Followers
        WHERE followerUserId = ?`,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followingIds []uint
	for rows.Next() {
		var followingId uint
		if err = rows.Scan(&followingId); err != nil {
			return nil, err
		}
		followingIds = append(followingIds, followingId)
	}

	following := make([]User, 0)
	for _, followingId := range followingIds {
		row := db.c.QueryRow(`
            SELECT username FROM Users
            WHERE userId = ?`,
			followingId,
		)
		var username string
		if err = row.Scan(&username); err != nil {
			return nil, err
		}
		following = append(following, User{UserId: followingId, Username: username})
	}

	return following, nil
}

// GetFollowingCount returns the number of users followed by a user.
func (db *appdbimpl) GetFollowingCount(userId uint) (uint, error) {
	var count uint
	if err := db.c.QueryRow(`
        SELECT COUNT(*) FROM Followers
        WHERE followerUserId = ?`,
		userId,
	).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, &UserNotFoundByIdError{UserId: userId}
		}
	}
	return count, nil
}

/*
GetFollowStatus returns the follow status between two users.

If the user with id `userId` follows the user with id `targetUserId`, it returns true.
*/
func (db *appdbimpl) GetFollowStatus(userId uint, targetUserId uint) (bool, error) {
	var res bool
	if err := db.c.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM Followers
            WHERE followerUserId = ? AND followingUserId = ?
        )`,
		userId, targetUserId,
	).Scan(&res); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, &UserNotFoundByIdError{UserId: userId}
		}
	}
	return res, nil
}

// FollowUser add the user with id `targetUserId` to the list of following of the user with id `userId`.
func (db *appdbimpl) FollowUser(userId uint, targetUserId uint) error {
	if _, err := db.c.Exec(`
        INSERT INTO Followers (followerUserId, followingUserId)
        VALUES (?, ?)`,
		userId, targetUserId,
	); err != nil {
		return err
	}
	return nil
}

// UnfollowUser remove the user with id `targetUserId` from the list of following of the user with id `userId`.
func (db *appdbimpl) UnfollowUser(userId uint, targetUserId uint) error {
	_, err := db.c.Exec(`
        DELETE FROM Followers
        WHERE followerUserId = ? AND followingUserId = ?`,
		userId, targetUserId,
	)
	if err != nil {
		return err
	}

	return nil
}
