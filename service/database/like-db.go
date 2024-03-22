package database

import "errors"

// LikePhoto Add a like to a photo.
func (db *appdbimpl) LikePhoto(userId uint, photoId uint) error {
	_, err := db.c.Exec("INSERT INTO Likes (userId, photoId) VALUES (?, ?)", userId, photoId)
	if err != nil {
		return err
	}

	// Increment the LikeCount field of the photo
	_, err = db.c.Exec("UPDATE Photos SET LikeCount = LikeCount + 1 WHERE PhotoId = ?", photoId)
	if err != nil {
		return err
	}

	return nil
}

// UnlikePhoto Remove a like from a photo.
func (db *appdbimpl) UnlikePhoto(userId uint, photoId uint) error {
	res, err := db.c.Exec("DELETE FROM Likes WHERE userId = ? AND photoId = ?", userId, photoId)
	if err != nil {
		return err
	}

	// Decrement the LikeCount field of the photo
	_, err = db.c.Exec("UPDATE Photos SET LikeCount = LikeCount - 1 WHERE PhotoId = ?", photoId)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
