package database

import (
	"database/sql"
	"errors"
	. "github.com/Big-Iron-Cheems/WASAPhoto/service/model"
	"sort"
	"strings"
)

/*
CreateUser creates a new user in the database.

If the user already exists, it returns the user schema.
Otherwise, it creates it and returns the user schema.
*/
func (db *appdbimpl) CreateUser(user User) (User, error) {
	// Try to insert a new user
	_, err := db.c.Exec(`
        INSERT OR IGNORE INTO Users(username)
        VALUES (?)`,
		user.Username,
	)
	if err != nil {
		return User{}, err
	}

	// Fetch the user details using GetUserProfile
	user, err = db.GetUserProfile(user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// GetAllUsers retrieves all users from the database via paginated requests.
func (db *appdbimpl) GetAllUsers(page int, pageSize int) ([]User, error) {
	offset := (page - 1) * pageSize
	rows, err := db.c.Query(`
        SELECT userId, username FROM Users
        LIMIT ? OFFSET ?`,
		pageSize, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var user User
		if err = rows.Scan(&user.UserId, &user.Username); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetUserProfile retrieves a user's profile from the database given their username.
func (db *appdbimpl) GetUserProfile(user User) (User, error) {
	if err := db.c.QueryRow(`
        SELECT userId, username FROM Users
        WHERE username = ?`,
		user.Username,
	).Scan(&user.UserId, &user.Username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, &UserNotFoundByUsernameError{Username: user.Username}
		}
	}
	return user, nil
}

/*
GetMyStream retrieves the content stream of a user from the database given its id.
The stream is composed of the photos uploaded by the users followed by the user.
The photos are sorted by date, from newest to oldest.
*/
func (db *appdbimpl) GetMyStream(user User) ([]Photo, error) {
	// Get the list of users followed by the user
	following, err := db.GetFollowingList(user.UserId)
	if err != nil {
		return nil, err
	}

	// For each user, get the list of photos they have uploaded
	totalPhotos := make([]Photo, 0)
	for _, followingUser := range following {
		photos, err := db.GetPhotoList(followingUser.UserId)
		if err != nil {
			return nil, err
		}
		totalPhotos = append(totalPhotos, photos...)
	}

	// Sort the photos by date, from newest to oldest
	sort.Slice(totalPhotos, func(i, j int) bool {
		return totalPhotos[i].UploadTime > totalPhotos[j].UploadTime
	})

	return totalPhotos, nil
}

// SetMyUsername updates the username of a user in the database given its id.
func (db *appdbimpl) SetMyUsername(user User, currentUsername string) (User, error) {
	_, err := db.c.Exec(`
        UPDATE Users
        SET username = ?
        WHERE username = ?`,
		user.Username, currentUsername,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, &UserNotFoundByUsernameError{Username: currentUsername}
		}
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return user, &UsernameTakenError{Username: user.Username}
		}
		return user, err
	}

	return user, nil
}
