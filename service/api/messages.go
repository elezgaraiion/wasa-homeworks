package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aritz/wasa-homeworks/service/api/reqcontext" 
	"github.com/aritz/wasa-homeworks/service/models"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getMessageById(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userHeader := r.Header.Get("Authorization")
	userID := extractBearer(userHeader)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	convID := ps.ByName("conversationId")
	if convID == "" {
		http.Error(w, "Missing conversationId", http.StatusBadRequest)
		return
	}

	messageID := ps.ByName("messageId")
	if messageID == "" {
		http.Error(w, "Missing messageId", http.StatusBadRequest)
		return
	}

	msg, err := rt.db.GetMessageByID(userID, convID, messageID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrConversationNotFound):
			http.Error(w, "Conversation not found", http.StatusNotFound)
		case errors.Is(err, models.ErrForbidden):
			http.Error(w, "Forbidden", http.StatusForbidden)
		case errors.Is(err, models.ErrMessageNotFound):
			http.Error(w, "Message not found", http.StatusNotFound)
		default:
			// Log con contexto
			ctx.Logger.WithError(err).Error("Get message error")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msg)
}

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userHeader := r.Header.Get("Authorization")
	userID := extractBearer(userHeader)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	convID := ps.ByName("conversationId")
	if convID == "" {
		http.Error(w, "Missing conversationId", http.StatusBadRequest)
		return
	}

	messageID := ps.ByName("messageId")
	if messageID == "" {
		http.Error(w, "Missing messageId", http.StatusBadRequest)
		return
	}

	var payload struct {
		TargetConversationID string `json:"targetConversationId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	if payload.TargetConversationID == "" {
		http.Error(w, "targetConversationId required", http.StatusBadRequest)
		return
	}

	msg, err := rt.db.ForwardMessage(userID, convID, messageID, payload.TargetConversationID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrConversationNotFound):
			http.Error(w, "Conversation not found", http.StatusNotFound)
		case errors.Is(err, models.ErrForbidden):
			http.Error(w, "Forbidden", http.StatusForbidden)
		case errors.Is(err, models.ErrMessageNotFound):
			http.Error(w, "Message not found", http.StatusNotFound)
		default:
			ctx.Logger.WithError(err).Error("Forward message error")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)
}

func (rt *_router) listReactionsForMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userHeader := r.Header.Get("Authorization")
	userID := extractBearer(userHeader)
	if userID == "" {
		http.Error(w, "Unauthorized", 401)
		return
	}

	convID := ps.ByName("conversationId")
	if convID == "" {
		http.Error(w, "Missing conversationId", 400)
		return
	}

	messageID := ps.ByName("messageId")
	if messageID == "" {
		http.Error(w, "Missing messageId", 400)
		return
	}

	ok, err := rt.db.IsUserInConversation(userID, convID)
	if err != nil {
		ctx.Logger.WithError(err).Error("IsUserInConversation check failed")
		http.Error(w, "Internal server error", 500)
		return
	}
	if !ok {
		http.Error(w, "Forbidden", 403)
		return
	}

	reactions, err := rt.db.GetReactions(messageID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Get reactions error")
		http.Error(w, "Internal server error", 500)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(reactions)
}

func (rt *_router) addReactionToMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userHeader := r.Header.Get("Authorization")
	userID := extractBearer(userHeader)
	if userID == "" {
		http.Error(w, "Unauthorized", 401)
		return
	}

	convID := ps.ByName("conversationId")
	msgID := ps.ByName("messageId")
	if convID == "" || msgID == "" {
		http.Error(w, "Missing conversationId or messageId", 400)
		return
	}

	var payload struct {
		Emoji string `json:"emoji"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		ctx.Logger.WithError(err).Error("Invalid JSON body")
		http.Error(w, "Invalid JSON body", 400)
		return
	}
	if payload.Emoji == "" {
		http.Error(w, "Emoji is required", 400)
		return
	}

	ok, err := rt.db.IsUserInConversation(userID, convID)
	if err != nil {
		ctx.Logger.WithError(err).Error("IsUserInConversation check failed")
		http.Error(w, "Internal server error", 500)
		return
	}
	if !ok {
		http.Error(w, "Forbidden", 403)
		return
	}

	reaction, err := rt.db.AddReaction(userID, convID, msgID, payload.Emoji)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrConversationNotFound):
			http.Error(w, "Conversation not found", 404)
		case errors.Is(err, models.ErrMessageNotFound):
			http.Error(w, "Message not found", 404)
		default:
			ctx.Logger.WithError(err).Error("AddReaction DB error")
			http.Error(w, "Internal server error", 500)
		}
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(reaction)
}

func (rt *_router) removeReactionFromMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userHeader := r.Header.Get("Authorization")
	userID := extractBearer(userHeader)
	if userID == "" {
		http.Error(w, "Unauthorized", 401)
		return
	}

	convID := ps.ByName("conversationId")
	msgID := ps.ByName("messageId")
	reactionID := ps.ByName("reactionId")

	if convID == "" || msgID == "" || reactionID == "" {
		http.Error(w, "Missing conversationId, messageId, or reactionId", 400)
		return
	}

	ok, err := rt.db.IsUserInConversation(userID, convID)
	if err != nil {
		ctx.Logger.WithError(err).Error("IsUserInConversation check failed")
		http.Error(w, "Internal server error", 500)
		return
	}
	if !ok {
		http.Error(w, "Forbidden", 403)
		return
	}

	err = rt.db.RemoveReaction(userID, convID, msgID, reactionID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrConversationNotFound):
			http.Error(w, "Conversation not found", 404)
		case errors.Is(err, models.ErrMessageNotFound):
			http.Error(w, "Message not found", 404)
		case errors.Is(err, models.ErrReactionNotFound):
			http.Error(w, "Reaction not found", 404)
		default:
			ctx.Logger.WithError(err).Error("RemoveReaction DB error")
			http.Error(w, "Internal server error", 500)
		}
		return
	}

	w.WriteHeader(204)
}