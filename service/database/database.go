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
	GetUserProfile(username string) (User, error)
	GetMyStream(userId uint) ([]Photo, error)
	SetMyUsername(userId uint, newUsername string) error

	// ban-db methods

	BanUser(userId int, targetUserId int) error
	UnbanUser(userId int, targetUserId int) error

	// follow-db methods

	FollowUser(userId int, targetUserId int) error
	UnfollowUser(userId int, targetUserId int) error

	// photo-db methods

	UploadPhoto(photo Photo) (Photo, error)
	DeletePhoto(photoId int) error

	// like-db methods

	LikePhoto(userId int, photoId int) error
	UnlikePhoto(userId int, photoId int) error
	// comment-db methods

	CommentPhoto(userId int, photoId int, content string) error
	UncommentPhoto(userId int, photoId int, commentId int) error

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

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='Users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		usersDatabase := `CREATE TABLE Users (
			userId INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE
		);`
		photosDatabase := `CREATE TABLE Photos (
			photoId INTEGER NOT NULL PRIMARY KEY,
            owner TEXT NOT NULL,
            image BLOB NOT NULL,
            description TEXT,
            uploadTime DATETIME NOT NULL,
            likeCount INTEGER NOT NULL,
            commentsCount INTEGER NOT NULL,
            likeStatus BOOLEAN NOT NULL,
			FOREIGN KEY (owner) REFERENCES Users(userId)
		);`
		likesDatabase := `CREATE TABLE Likes (
            likeId INTEGER NOT NULL PRIMARY KEY,
            userId INTEGER NOT NULL,
            photoId INTEGER NOT NULL,
            FOREIGN KEY (userId) REFERENCES Users(userId),
            FOREIGN KEY (photoId) REFERENCES Photos(photoId)
		);`
		commentsDatabase := `CREATE TABLE Comments (
            commentId INTEGER NOT NULL PRIMARY KEY,
            owner INTEGER NOT NULL,
            content TEXT NOT NULL,
            userId INTEGER NOT NULL,
            photoId INTEGER NOT NULL,
            FOREIGN KEY (userId) REFERENCES Users(userId),
            FOREIGN KEY (photoId) REFERENCES Photos(photoId)
        );`
		bansDatabase := `CREATE TABLE Bans (
			bannedUserId INTEGER NOT NULL,
			userId INTEGER NOT NULL,
			PRIMARY KEY (userId, bannedUserId),
            FOREIGN KEY (bannedUserId) REFERENCES Users(userId),
			FOREIGN KEY (userId) REFERENCES Users(userId)
		);`
		followersDatabase := `CREATE TABLE Followers (
	        followerUserId INTEGER NOT NULL,
            followingUserId INTEGER NOT NULL,
            PRIMARY KEY (followerUserId, followingUserId),
            FOREIGN KEY (followerUserId) REFERENCES Users(userId),
            FOREIGN KEY (followingUserId) REFERENCES Users(userId)
		);`

		// Ensure that the tables are created
		tables := []string{usersDatabase, photosDatabase, likesDatabase, commentsDatabase, bansDatabase, followersDatabase}
		for _, table := range tables {
			if _, err = db.Exec(table); err != nil {
				return nil, fmt.Errorf("error creating database structure: %w", err)
			}
		}
	}

	return &appdbimpl{c: db}, nil
}

// Ping checks the connection to the database.
func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
