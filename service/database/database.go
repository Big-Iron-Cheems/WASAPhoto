/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	. "wasaphoto.uniroma1.it/wasaphoto/service/model"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	// user-db methods

	CreateUser(user User) (User, error)
	GetUserProfile(user User) (User, error)
	GetMyStream(user User) ([]Photo, error)
	SetMyUsername(user User, currentUsername string) (User, error)
	GetFollowStatus(userId uint, targetUserId uint) (bool, error)
	GetUserBanStatusByRequester(userId uint, targetUserId uint) (bool, error)
	GetRequesterBannedByUser(userId uint, targetUserId uint) (bool, error)

	// ban-db methods

	BanUser(userId uint, targetUserId uint) error
	UnbanUser(userId uint, targetUserId uint) error

	// follow-db methods

	GetFollowers(userId uint) ([]uint, error)
	GetFollowersCount(userId uint) (uint, error)
	GetFollowing(userId uint) ([]uint, error)
	GetFollowingCount(userId uint) (uint, error)
	FollowUser(userId uint, targetUserId uint) error
	UnfollowUser(userId uint, targetUserId uint) error

	// photo-db methods

	GetPhoto(photoId uint) (Photo, error)
	GetPhotoCount(userId uint) (uint, error)
	UploadPhoto(photo Photo) (Photo, error)
	DeletePhoto(photo Photo) error

	// like-db methods

	LikePhoto(userId uint, photoId uint) error
	UnlikePhoto(userId uint, photoId uint) error

	// comment-db methods

	CommentPhoto(photoId uint, comment Comment) (Comment, error)
	UncommentPhoto(photoId uint, comment Comment) error

	Ping() error
}

// appdbimpl: AppDatabaseImpl is the implementation of the AppDatabase interface
type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	tables := map[string]string{
		"Users": `CREATE TABLE Users (
                userId INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
                username TEXT NOT NULL UNIQUE
            );`,
		"Photos": `CREATE TABLE Photos (
			photoId INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
            owner TEXT NOT NULL,
            image BLOB NOT NULL,
            description TEXT,
            uploadTime DATETIME NOT NULL,
            likeCount INTEGER NOT NULL,
            commentsCount INTEGER NOT NULL,
			FOREIGN KEY (owner) REFERENCES Users(userId)
		);`,
		"Likes": `CREATE TABLE Likes (
            userId INTEGER NOT NULL,
            photoId INTEGER NOT NULL,
            PRIMARY KEY (userId, photoId),
            FOREIGN KEY (userId) REFERENCES Users(userId),
            FOREIGN KEY (photoId) REFERENCES Photos(photoId)
		);`,
		"Comments": `CREATE TABLE Comments (
            commentId INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
            owner INTEGER NOT NULL,
            content TEXT NOT NULL,
            photoId INTEGER NOT NULL,
            FOREIGN KEY (owner) REFERENCES Users(userId),
            FOREIGN KEY (photoId) REFERENCES Photos(photoId)
        );`,
		"Bans": `CREATE TABLE Bans (
			bannedUserId INTEGER NOT NULL,
			userId INTEGER NOT NULL,
			PRIMARY KEY (userId, bannedUserId),
            FOREIGN KEY (bannedUserId) REFERENCES Users(userId),
			FOREIGN KEY (userId) REFERENCES Users(userId)
		);`,
		"Followers": `CREATE TABLE Followers (
	        followerUserId INTEGER NOT NULL,
            followingUserId INTEGER NOT NULL,
            PRIMARY KEY (followerUserId, followingUserId),
            FOREIGN KEY (followerUserId) REFERENCES Users(userId),
            FOREIGN KEY (followingUserId) REFERENCES Users(userId)
		);`,
	}

	// Iterate over the tables map
	for tableName, createQuery := range tables {
		// Check if the table exists
		var name string
		err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name=?;`, tableName).Scan(&name)
		if errors.Is(err, sql.ErrNoRows) {
			// If the table does not exist, create it
			if _, err = db.Exec(createQuery); err != nil {
				return nil, fmt.Errorf("error creating `%s` table: %w", tableName, err)
			}
		} else if err != nil {
			// If there was an error checking for the table, return the error
			return nil, fmt.Errorf("error checking for `%s` table: %w", tableName, err)
		}
	}

	return &appdbimpl{c: db}, nil
}

// Ping checks the connection to the database.
func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
