package database

import (
	"database/sql"
	"errors"
	. "github.com/Big-Iron-Cheems/WASAPhoto/service/model"
)

// GetPhoto Get a photo by its id.
func (db *appdbimpl) GetPhoto(photoId uint) (Photo, error) {
	var photo Photo
	if err := db.c.QueryRow(
		`SELECT Photos.photoId, Photos.ownerId, Users.username, Photos.image, Photos.mimeType, Photos.caption, Photos.uploadTime, Photos.likeCount, Photos.commentsCount
		FROM Photos
		INNER JOIN Users ON Photos.ownerId = Users.userId
		WHERE photoId = ?`,
		photoId,
	).Scan(
		&photo.PhotoId,
		&photo.OwnerId,
		&photo.OwnerUsername,
		&photo.Image,
		&photo.MimeType,
		&photo.Caption,
		&photo.UploadTime,
		&photo.LikeCount,
		&photo.CommentsCount,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Photo{}, errors.New("photo does not exist")
		}
	}
	return photo, nil
}

/*
GetPhotoList Get a list of photos uploaded by a user.
The photos are sorted by date, from newest to oldest.
*/
func (db *appdbimpl) GetPhotoList(userId uint) ([]Photo, error) {
	rows, err := db.c.Query(
		`SELECT Photos.photoId, Photos.ownerId, Users.username, Photos.image, Photos.mimeType, Photos.caption, Photos.uploadTime, Photos.likeCount, Photos.commentsCount
		FROM Photos
		INNER JOIN Users ON Photos.ownerId = Users.userId
		WHERE ownerId = ?
		ORDER BY uploadTime DESC`,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	photos := make([]Photo, 0)
	for rows.Next() {
		var photo Photo
		if err := rows.Scan(
			&photo.PhotoId,
			&photo.OwnerId,
			&photo.OwnerUsername,
			&photo.Image,
			&photo.MimeType,
			&photo.Caption,
			&photo.UploadTime,
			&photo.LikeCount,
			&photo.CommentsCount,
		); err != nil {
			return nil, err
		}
		photos = append(photos, photo)
	}
	return photos, nil
}

// GetPhotoCount Get the number of photos uploaded by a user.
func (db *appdbimpl) GetPhotoCount(userId uint) (uint, error) {
	var count uint
	if err := db.c.QueryRow(`
        SELECT COUNT(*) FROM Photos
        WHERE ownerId = ?`,
		userId,
	).Scan(&count); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, &UserNotFoundByIdError{UserId: userId}
		}
	}
	return count, nil
}

// UploadPhoto Upload a photo.
func (db *appdbimpl) UploadPhoto(photo Photo) (Photo, error) {
	res, err := db.c.Exec(`
        INSERT INTO Photos (ownerId, image, mimeType, caption, uploadTime, likeCount, commentsCount)
        VALUES (?, ?, ?, ?, strftime('%Y-%m-%dT%H:%M:%SZ', 'now'), ?, ?)`,
		photo.OwnerId, photo.Image, photo.MimeType, photo.Caption, photo.LikeCount, photo.CommentsCount,
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

// DeletePhoto Delete a photo and its associated comments.
func (db *appdbimpl) DeletePhoto(photo Photo) error {
	// Delete the comments associated with the photo
	_, err := db.c.Exec(`
        DELETE FROM Comments
        WHERE photoId = ?`,
		photo.PhotoId,
	)
	if err != nil {
		return err
	}

	// Delete the likes associated with the photo
	_, err = db.c.Exec(`
        DELETE FROM Likes
        WHERE photoId = ?`,
		photo.PhotoId,
	)
	if err != nil {
		return err
	}

	// Delete the photo
	_, err = db.c.Exec(`
        DELETE FROM Photos
        WHERE photoId = ? AND ownerId = ?`,
		photo.PhotoId, photo.OwnerId,
	)
	if err != nil {
		return err
	}

	return nil
}
