package api

import (
	"net/http"
)

func (rt *_router) Handler() http.Handler {
	rt.router.POST("/session", rt.wrap(rt.doLogin))
	rt.router.GET("/me", rt.wrap(rt.doGetCurrentUser))
	rt.router.PUT("/me/username", rt.wrap(rt.updateMyUserName))
	rt.router.PUT("/me/photo", rt.wrap(rt.updatePhoto))
	rt.router.GET("/users", rt.wrap(rt.listOrSearchUsers))
	rt.router.GET("/conversations", rt.wrap(rt.getMyConversations))
	rt.router.POST("/conversations", rt.wrap(rt.createGroup))
	rt.router.GET("/conversations/:conversationId", rt.wrap(rt.getConversationProfile))
	rt.router.POST("/conversations/:conversationId/seen", rt.wrap(rt.markConversationSeen))
	rt.router.GET("/conversations/:conversationId/messages", rt.wrap(rt.listConversationMessages))
	rt.router.POST("/conversations/:conversationId/messages", rt.wrap(rt.sendMessage))
	rt.router.GET("/conversations/:conversationId/messages/:messageId", rt.wrap(rt.getMessageById))
	rt.router.POST("/conversations/:conversationId/messages/:messageId", rt.wrap(rt.forwardMessage))
	rt.router.DELETE("/conversations/:conversationId/messages/:messageId", rt.wrap(rt.deleteMessage))
	rt.router.GET("/conversations/:conversationId/messages/:messageId/reactions", rt.wrap(rt.listReactionsForMessage))
	rt.router.POST("/conversations/:conversationId/messages/:messageId/reactions", rt.wrap(rt.addReactionToMessage))
	rt.router.DELETE("/conversations/:conversationId/messages/:messageId/reactions/:reactionId", rt.wrap(rt.removeReactionFromMessage))
	rt.router.POST("/conversations/:conversationId/users", rt.wrap(rt.addUserToGroup))
	rt.router.DELETE("/conversations/:conversationId/users/me", rt.wrap(rt.leaveGroup))
	rt.router.PUT("/conversations/:conversationId/name", rt.wrap(rt.setGroupName))
	rt.router.PUT("/conversations/:conversationId/photo", rt.wrap(rt.setGroupPhoto))
	rt.router.POST("/chats", rt.wrap(rt.createPrivateConversation))
	return rt.router
}
