package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	rt.router.POST("/session", rt.doLogin)
	rt.router.GET("/me", rt.doGetCurrentUser)
	rt.router.PUT("/me/username", rt.updateMyUserName)
	rt.router.PUT("/me/photo", rt.updatePhoto)
	rt.router.GET("/users", rt.listOrSearchUsers)
	rt.router.GET("/conversations", rt.getMyConversations)
	rt.router.POST("/conversations", rt.createGroup)
	rt.router.GET("/conversations/:conversationId", rt.getConversationProfile)
	rt.router.POST("/conversations/:conversationId/seen", rt.markConversationSeen)
	rt.router.GET("/conversations/:conversationId/messages", rt.listConversationMessages)
	rt.router.POST("/conversations/:conversationId/messages", rt.sendMessage)
	rt.router.GET("/conversations/:conversationId/messages/:messageId", rt.getMessageById)
	rt.router.POST("/conversations/:conversationId/messages/:messageId", rt.forwardMessage)
	rt.router.DELETE("/conversations/:conversationId/messages/:messageId", rt.deleteMessage)
	rt.router.GET("/conversations/:conversationId/messages/:messageId/reactions", rt.listReactionsForMessage)
	rt.router.POST("/conversations/:conversationId/messages/:messageId/reactions", rt.addReactionToMessage)
	rt.router.DELETE("/conversations/:conversationId/messages/:messageId/reactions/reactionId", rt.removeReactionFromMessage)
	return rt.router
}