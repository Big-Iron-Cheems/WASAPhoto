package database

import (
	"database/sql"
	"errors"
	. "wasaphoto.uniroma1.it/wasaphoto/service/model"
)

// GetPhoto Get a photo by its id.
func (db *appdbimpl) GetPhoto(photoId uint) (Photo, error) {
	var photo Photo
	if err := db.c.QueryRow("SELECT * FROM Photos WHERE photoId = ?", photoId).Scan(&photo.PhotoId, &photo.Owner, &photo.File, &photo.Description, &photo.UploadTime, &photo.LikeCount); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Photo{}, errors.New("photo does not exist")
		}
	}
	return photo, nil
}

// GetPhotoCount Get the number of photos uploaded by a user.
func (db *appdbimpl) GetPhotoCount(userId uint) (uint, error) {
	var count uint
	if err := db.c.QueryRow("SELECT COUNT(*) FROM Photos WHERE photoId = ?", userId).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("user does not exist")
		}
	}
	return count, nil
}

// UploadPhoto Upload a photo.
func (db *appdbimpl) UploadPhoto(photo Photo) (Photo, error) {
	res, err := db.c.Exec(
		"INSERT INTO Photos (owner, image, description, uploadTime, likeCount, commentsCount) VALUES (?, ?, ?, strftime('%Y-%m-%dT%H:%M:%SZ', 'now'), ?, ?)",
		photo.Owner, photo.File, photo.Description, photo.LikeCount, photo.CommentsCount,
	)
	if err != nil {
		return Photo{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return Photo{}, err
	}

	photo.PhotoId = uint(id)
	return photo, nil
}

// DeletePhoto Delete a photo.
func (db *appdbimpl) DeletePhoto(photo Photo) error {
	// Delete the photo
	_, err := db.c.Exec("DELETE FROM Photos WHERE PhotoId = ? AND Owner = ?", photo.PhotoId, photo.Owner)
	if err != nil {
		return err
	}

	return nil
}
