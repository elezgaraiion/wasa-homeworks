package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	rt.router.POST("/session", rt.doLogin)
	rt.router.GET("/me", rt.doGetCurrentUser)
	rt.router.PUT("/me/username", rt.updateMyUserName)
	rt.router.PUT("/me/photo", rt.UpdatePhoto)
	rt.router.GET("/users", rt.listOrSearchUsers)
	rt.router.GET("/conversations", rt.GetMyConversations)
	return rt.router
}