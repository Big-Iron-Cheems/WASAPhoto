/*
Package model contains the structs that model the schemas defined in the API.

These structs are used in both the API and the database packages.
*/
package model

// User struct modeling the schema of a user.
//   - UserId is not modifiable, and it is used to identify the user.
//   - Username is modifiable, and it is used to log in, or to display the user.
type User struct {
	UserId   uint   `json:"userId"`
	Username string `json:"username"`
}

// Profile struct modeling the schema of a user's profile.
//   - Photos is a list of unique Photo.PhotoId owned by the user.
//   - PhotoCount is the number of photos owned by the user.
//   - Followers is a list of unique User.UserId that follow the user.
//   - Following is a list of unique User.UserId that the user follows.
//   - BanList is a list of unique User.UserId that the user has banned.
type Profile struct {
	Photos     []uint `json:"photos"`
	PhotoCount uint   `json:"photoCount"`
	Followers  []uint `json:"followers"`
	Following  []uint `json:"following"`
	BanList    []uint `json:"banList"`
}

// Photo struct modeling the schema of a photo.
//   - PhotoId is not modifiable, and it is used to identify the photo.
//   - Owner is not modifiable, and it is used to identify the owner of the photo. It is a User.UserId.
//   - File is the binary content of the photo. TODO: decide if it is needed.
//   - Description is the description of the photo.
//   - UploadTime is the time when the photo was uploaded.
//   - LikeCount is the number of likes of the photo.
//   - CommentsCount is the number of comments of the photo.
//   - LikeStatus is the status of the like of the photo. It is true if the user has liked the photo, false otherwise.
type Photo struct {
	PhotoId uint `json:"photoId"`
	Owner   uint `json:"owner"`
	//File          []byte
	Description   string `json:"description"`
	UploadTime    string `json:"uploadTime"`
	LikeCount     uint   `json:"likeCount"`
	CommentsCount uint   `json:"commentsCount"`
	LikeStatus    bool   `json:"likeStatus"`
}

// Comment struct modeling the schema of a comment.
//   - CommentId is not modifiable, and it is used to identify the comment.
//   - Owner is not modifiable, and it is used to identify the owner of the comment. It is a User.UserId.
//   - Content is the content of the comment.
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

// Mock structs/maps to emulate DB for small scale tests

var ExampleUser = User{
	UserId:   1234,
	Username: "Pollo",
}

var ExampleProfile = Profile{
	Photos:     nil,
	PhotoCount: 0,
	Followers:  nil,
	Following:  nil,
	BanList:    nil,
}

var tokens = map[User]string{
	ExampleUser: "", // token made from random hash of 16 chars
}

var users = map[uint]User{
	1234: ExampleUser,
	4208: {
		UserId:   4208,
		Username: "Gigio",
	},
}

var profiles = map[uint]Profile{
	1234: ExampleProfile,
}
