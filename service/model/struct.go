/*
Package model contains the structs that model the schemas defined in the API.

These structs are used in both the API and the database packages.
*/
package model

/*
User struct modeling the schema of a user.
  - UserId is not modifiable, and it is used to identify the user.
  - Username is modifiable, and it is used to log in, or to display the user.
*/
type User struct {
	UserId   uint   `json:"userId"`
	Username string `json:"username"`
}

/*
Profile struct modeling the schema of a user's profile.
This struct contains the public data related to a user's profile.
Since this data is public, it is not modifiable.
  - UserId is used to identify the user.
  - Username is used to display the user.
  - PhotoCount is the number of photos of the user.
  - FollowersCount is the number of followers of the user.
  - FollowingCount is the number of users followed by the user.
*/
type Profile struct {
	UserId         uint   `json:"userId"`
	Username       string `json:"username"`
	PhotoCount     uint   `json:"photoCount"`
	FollowersCount uint   `json:"followersCount"`
	FollowingCount uint   `json:"followingCount"`
}

/*
Photo struct modeling the schema of a photo.
  - PhotoId is not modifiable, and it is used to identify the photo.
  - Owner is not modifiable, and it is used to identify the owner of the photo. It is a User.UserId.
  - File is the binary content of the photo. TODO: decide if it is needed.
  - Description is the description of the photo.
  - UploadTime is the time when the photo was uploaded.
  - LikeCount is the number of likes of the photo.
  - CommentsCount is the number of comments of the photo.
*/
type Photo struct {
	PhotoId       uint `json:"photoId"`
	Owner         uint `json:"owner"`
	File          []byte
	Description   string `json:"description"`
	UploadTime    string `json:"uploadTime"`
	LikeCount     uint   `json:"likeCount"`
	CommentsCount uint   `json:"commentsCount"`
}

/*
Comment struct modeling the schema of a comment.
  - CommentId is not modifiable, and it is used to identify the comment.
  - Owner is not modifiable, and it is used to identify the owner of the comment. It is a User.UserId.
  - Content is the content of the comment.
*/
type Comment struct {
	CommentId uint   `json:"commentId"` // A user can make multiple comments, thus we need a CommentId
	Owner     uint   `json:"owner"`
	Content   string `json:"content"`
}

// Error struct
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
