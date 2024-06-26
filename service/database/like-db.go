package database

// LikePhoto Add a like to a photo.
func (db *appdbimpl) LikePhoto(userId uint, photoId uint) error {
	_, err := db.c.Exec(`
        INSERT INTO Likes (userId, photoId)
        VALUES (?, ?)`,
		userId, photoId,
	)
	if err != nil {
		return err
	}

	// Increment the LikeCount field of the photo
	_, err = db.c.Exec(`
        UPDATE Photos
        SET likeCount = likeCount + 1
        WHERE photoId = ?`,
		photoId,
	)
	if err != nil {
		return err
	}

	return nil
}

// UnlikePhoto Remove a like from a photo.
func (db *appdbimpl) UnlikePhoto(userId uint, photoId uint) error {
	_, err := db.c.Exec(`
        DELETE FROM Likes
        WHERE userId = ? AND photoId = ?`,
		userId, photoId,
	)
	if err != nil {
		return err
	}

	// Decrement the LikeCount field of the photo
	_, err = db.c.Exec(`
        UPDATE Photos
        SET likeCount = likeCount - 1
        WHERE photoId = ?`,
		photoId,
	)
	if err != nil {
		return err
	}

	return nil
}

// GetLikeStatus Check if a user has liked a photo.
func (db *appdbimpl) GetLikeStatus(userId uint, photoId uint) (bool, error) {
	var exists bool
	err := db.c.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM Likes
            WHERE userId = ? AND photoId = ?
        )`,
		userId, photoId,
	).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
