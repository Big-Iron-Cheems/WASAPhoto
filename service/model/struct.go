/*
Package model contains the structs that model the schemas defined in the API.

These structs are used in both the API and the database packages.
*/
package model

import (
	"errors"
	"fmt"
)

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
  - BannedCount is the number of users banned by the user.
*/
type Profile struct {
	UserId         uint   `json:"userId"`
	Username       string `json:"username"`
	PhotoCount     uint   `json:"photoCount"`
	FollowersCount uint   `json:"followersCount"`
	FollowingCount uint   `json:"followingCount"`
	BannedCount    uint   `json:"bannedCount"`
}

/*
Photo struct modeling the schema of a photo.
  - PhotoId is not modifiable, and it is used to identify the photo.
  - OwnerId is not modifiable, and it is used to identify the owner of the photo. It is a User.UserId.
  - OwnerUsername is the User.Username of the owner of the photo.
  - Image is the binary content of the photo.
  - MimeType is the MIME type of the photo.
  - Caption is the text caption of the photo.
  - UploadTime is the time when the photo was uploaded.
  - LikeCount is the number of likes of the photo.
  - CommentsCount is the number of comments of the photo.
*/
type Photo struct {
	PhotoId       uint   `json:"photoId"`
	OwnerId       uint   `json:"ownerId"`
	OwnerUsername string `json:"ownerUsername"` // Calculated via JOIN, not stored in the database
	Image         []byte `json:"image"`
	MimeType      string `json:"mimeType"`
	Caption       string `json:"caption"`
	UploadTime    string `json:"uploadTime"`
	LikeCount     uint   `json:"likeCount"`
	CommentsCount uint   `json:"commentsCount"`
}

/*
Comment struct modeling the schema of a comment.
  - CommentId is not modifiable, and it is used to identify the comment.
  - OwnerId is not modifiable, and it is used to identify the owner of the comment. It is a User.UserId.
  - OwnerUsername is the User.Username of the owner of the comment.
  - Content is the content of the comment.
*/
type Comment struct {
	CommentId     uint   `json:"commentId"`
	OwnerId       uint   `json:"ownerId"`
	OwnerUsername string `json:"ownerUsername"` // Calculated via JOIN, not stored in the database
	Content       string `json:"content"`
}

/*
Error struct
  - Code is the HTTP status code of the error.
  - Message is the error message.

This struct is used to return errors in the API.
Any other error struct here defined should pass its error message to this struct.
*/
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d %s", e.Code, e.Message)
}

func (e *Error) Is(target error) bool {
	var t *Error
	ok := errors.As(target, &t)
	if !ok {
		return false
	}
	return e.Code == t.Code && e.Message == t.Message
}

/*
InvalidPatternError whenever a string does not match a regex pattern.
  - Pattern is the regex pattern that the string should match.
  - Str is the string that does not match the pattern.
*/
type InvalidPatternError struct {
	Pattern string
	Str     string
}

func (e *InvalidPatternError) Error() string {
	return fmt.Sprintf("Input `%s` does not match pattern `%s`", e.Str, e.Pattern)
}

func (e *InvalidPatternError) Is(target error) bool {
	var t *InvalidPatternError
	ok := errors.As(target, &t)
	if !ok {
		return false
	}
	return e.Pattern == t.Pattern && e.Str == t.Str
}

/*
UserNotFoundByUsernameError whenever the db cannot find a user with the given username.
  - Username is the username that the db cannot find.
*/
type UserNotFoundByUsernameError struct {
	Username string
}

func (e *UserNotFoundByUsernameError) Error() string {
	return fmt.Sprintf("User `%s` not found", e.Username)
}

/*
UserNotFoundByIdError whenever the db cannot find a user with the given user ID.
  - UserId is the user ID that the db cannot find.
*/
type UserNotFoundByIdError struct {
	UserId uint
}

func (e *UserNotFoundByIdError) Error() string {
	return fmt.Sprintf("User with ID `%d` not found", e.UserId)
}

/*
UsernameTakenError whenever a request to change username collides with an existing username.
  - Username is the name of the unique constraint that failed.
*/
type UsernameTakenError struct {
	Username string
}

func (e *UsernameTakenError) Error() string {
	return fmt.Sprintf("Username `%s` already taken", e.Username)
}

func (e *UsernameTakenError) Is(target error) bool {
	var t *UsernameTakenError
	ok := errors.As(target, &t)
	if !ok {
		return false
	}
	return e.Username == t.Username
}
