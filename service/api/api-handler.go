package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Login operations
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	// User profile operations
	rt.router.GET("/stream", rt.wrap(rt.getMyStream))
	rt.router.GET("/users", rt.wrap(rt.getAllUsers))
	rt.router.PUT("/users/:username", rt.wrap(rt.setMyUsername))
	rt.router.GET("/users/:username/profile", rt.wrap(rt.getUserProfile))

	// Ban operations
	rt.router.GET("/users/:username/bans/list", rt.wrap(rt.getBansList))
	rt.router.GET("/users/:username/bans/list/:targetUsername", rt.wrap(rt.getBanStatus))
	rt.router.POST("/users/:username/bans", rt.wrap(rt.banUser))
	rt.router.DELETE("/users/:username/bans/:targetUsername", rt.wrap(rt.unbanUser))

	// Follow operations
	rt.router.GET("/users/:username/followers/list", rt.wrap(rt.getFollowersList))
	rt.router.GET("/users/:username/followers/list/:targetUsername", rt.wrap(rt.getFollowStatus))
	rt.router.GET("/users/:username/following/list", rt.wrap(rt.getFollowingList))
	rt.router.POST("/users/:username/followers", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/:username/followers/:targetUsername", rt.wrap(rt.unfollowUser))

	// Photo operations
	rt.router.GET("/users/:username/photos", rt.wrap(rt.getPhotoList))
	rt.router.POST("/users/:username/photos", rt.wrap(rt.uploadPhoto))
	rt.router.DELETE("/users/:username/photos/:photoId", rt.wrap(rt.deletePhoto))

	// Like operations
	rt.router.GET("/users/:username/photos/:photoId/likes/list/:targetUsername", rt.wrap(rt.getLikeStatus))
	rt.router.POST("/users/:username/photos/:photoId/likes", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/users/:username/photos/:photoId/likes/:likeId", rt.wrap(rt.unlikePhoto))

	// Comment operations
	rt.router.GET("/users/:username/photos/:photoId/comments", rt.wrap(rt.getPhotoComments))
	rt.router.POST("/users/:username/photos/:photoId/comments", rt.wrap(rt.commentPhoto))
	rt.router.DELETE("/users/:username/photos/:photoId/comments/:commentId", rt.wrap(rt.uncommentPhoto))

	// Special routes
	rt.router.GET("/liveness", rt.wrap(rt.liveness))

	return rt.router
}
