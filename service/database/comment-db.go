package database

import . "wasaphoto.uniroma1.it/wasaphoto/service/model"

// CommentPhoto Add a comment under a photo.
func (db *appdbimpl) CommentPhoto(photoId uint, comment Comment) (Comment, error) {
	res, err := db.c.Exec("INSERT INTO Comments (owner, photoId, content) VALUES (?, ?, ?)",
		comment.Owner, photoId, comment.Content,
	)
	if err != nil {
		return Comment{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return Comment{}, err
	}

	// Increment the CommentsCount field of the photo
	_, err = db.c.Exec("UPDATE Photos SET CommentsCount = CommentsCount + 1 WHERE PhotoId = ?", photoId)
	if err != nil {
		return Comment{}, err
	}

	comment.CommentId = uint(id)
	return comment, nil
}

func (db *appdbimpl) UncommentPhoto(photoId uint, comment Comment) error {
	// Delete the comment
	_, err := db.c.Exec("DELETE FROM Comments WHERE CommentId = ? AND Owner = ? AND PhotoId = ?",
		comment.CommentId, comment.Owner, photoId,
	)
	if err != nil {
		return err
	}

	// Decrement the CommentsCount field of the photo
	_, err = db.c.Exec("UPDATE Photos SET CommentsCount = CommentsCount - 1 WHERE PhotoId = ?", photoId)
	if err != nil {
		return err
	}

	return nil
}
