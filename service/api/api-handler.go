package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	/*TODO:
	  Verify Authorization header on all routes except /session
	  Implement helper to parse args from path / request body
	  Interface with DB to store/load structs
	*/

	// Login operations
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	// User profile operations
	rt.router.PUT("/users/:username", rt.wrap(rt.setMyUsername))
	rt.router.GET("/users/:username/profile", rt.wrap(rt.getUserProfile))
	rt.router.GET("/users/:username/stream", rt.wrap(rt.getMyStream))

	// Ban operations
	rt.router.POST("/users/:username/bans", rt.wrap(rt.banUser))
	rt.router.DELETE("/users/:username/bans/:targetUsername", rt.wrap(rt.unbanUser))

	// Follow operations
	rt.router.POST("/users/:username/follows", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/:username/follows/:targetUsername", rt.wrap(rt.unfollowUser))

	// Photo operations
	rt.router.POST("/users/:username/photos", rt.wrap(rt.uploadPhoto))
	rt.router.DELETE("/users/:username/photos/:photoId", rt.wrap(rt.deletePhoto))

	// Like operations
	rt.router.POST("/users/:username/photos/:photoId/likes", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/users/:username/photos/:photoId/likes/:likeId", rt.wrap(rt.unlikePhoto))

	// Comment operations
	rt.router.POST("/users/:username/photos/:photoId/comments", rt.wrap(rt.commentPhoto))
	rt.router.DELETE("/users/:username/photos/:photoId/comments/:commentId", rt.wrap(rt.uncommentPhoto))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
