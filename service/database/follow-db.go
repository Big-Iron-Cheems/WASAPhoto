package database

import (
	"database/sql"
	"errors"
)

// GetFollowers returns the list of followers of a user.
func (db *appdbimpl) GetFollowers(userId uint) ([]uint, error) {
	rows, err := db.c.Query("SELECT followerUserId FROM Followers WHERE followingUserId = ?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []uint
	for rows.Next() {
		var followerId uint
		if err = rows.Scan(&followerId); err != nil {
			return nil, err
		}
		followers = append(followers, followerId)
	}

	return followers, nil
}

// GetFollowersCount returns the number of followers of a user.
func (db *appdbimpl) GetFollowersCount(userId uint) (uint, error) {
	var count uint
	if err := db.c.QueryRow("SELECT COUNT(*) FROM Followers WHERE followingUserId = ?", userId).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("user does not exist")
		}
	}
	return count, nil
}

// GetFollowing returns the list of users followed by a user.
func (db *appdbimpl) GetFollowing(userId uint) ([]uint, error) {
	rows, err := db.c.Query("SELECT followingUserId FROM Followers WHERE followerUserId = ?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var following []uint
	for rows.Next() {
		var followingId uint
		if err = rows.Scan(&followingId); err != nil {
			return nil, err
		}
		following = append(following, followingId)
	}

	return following, nil
}

// GetFollowingCount returns the number of users followed by a user.
func (db *appdbimpl) GetFollowingCount(userId uint) (uint, error) {
	var count uint
	if err := db.c.QueryRow("SELECT COUNT(*) FROM Followers WHERE followerUserId = ?", userId).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("user does not exist")
		}
	}
	return count, nil
}

// FollowUser add the user with id `targetUserId` to the list of following of the user with id `userId`.
func (db *appdbimpl) FollowUser(userId uint, targetUserId uint) error {
	if _, err := db.c.Exec("INSERT INTO Followers (followerUserId, followingUserId) VALUES (?, ?)", userId, targetUserId); err != nil {
		return err
	}
	return nil
}

// UnfollowUser remove the user with id `targetUserId` from the list of following of the user with id `userId`.
func (db *appdbimpl) UnfollowUser(userId uint, targetUserId uint) error {
	res, err := db.c.Exec("DELETE FROM Followers WHERE followerUserId = ? AND followingUserId = ?", userId, targetUserId)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("user does not exist")
	}

	return nil
}
