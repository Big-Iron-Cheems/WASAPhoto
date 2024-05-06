package database

import . "github.com/Big-Iron-Cheems/WASAPhoto/service/model"

// GetPhotoComments Get all comments under a photo.
func (db *appdbimpl) GetPhotoComments(photoId uint) ([]Comment, error) {
	rows, err := db.c.Query(`
		SELECT Comments.commentId, Comments.ownerId, Users.username, Comments.content
		FROM Comments
		INNER JOIN Users ON Comments.ownerId = Users.userId
		WHERE Comments.photoId = ?`, photoId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]Comment, 0)
	for rows.Next() {
		var comment Comment
		err = rows.Scan(&comment.CommentId, &comment.OwnerId, &comment.OwnerUsername, &comment.Content)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

// CommentPhoto Add a comment under a photo.
func (db *appdbimpl) CommentPhoto(photoId uint, comment Comment) (Comment, error) {
	res, err := db.c.Exec(`
        INSERT INTO Comments (ownerId, photoId, content)
        VALUES (?, ?, ?)`,
		comment.OwnerId, photoId, comment.Content,
	)
	if err != nil {
		return Comment{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return Comment{}, err
	}

	// Increment the CommentsCount field of the photo
	_, err = db.c.Exec(`
        UPDATE Photos
        SET commentsCount = commentsCount + 1
        WHERE photoId = ?`,
		photoId,
	)
	if err != nil {
		return Comment{}, err
	}

	comment.CommentId = uint(id)
	return comment, nil
}

// UncommentPhoto Remove a comment from a photo.
func (db *appdbimpl) UncommentPhoto(photoId uint, comment Comment) error {
	// Delete the comment
	_, err := db.c.Exec(`
        DELETE FROM Comments
        WHERE commentId = ? AND ownerId = ? AND photoId = ?`,
		comment.CommentId, comment.OwnerId, photoId,
	)
	if err != nil {
		return err
	}

	// Decrement the CommentsCount field of the photo
	_, err = db.c.Exec(`
        UPDATE Photos
        SET commentsCount = commentsCount - 1
        WHERE PhotoId = ?`,
		photoId,
	)
	if err != nil {
		return err
	}

	return nil
}
