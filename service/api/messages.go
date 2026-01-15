package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aritz/wasa-homeworks/service/api/reqcontext"
	"github.com/aritz/wasa-homeworks/service/models"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userHeader := r.Header.Get("Authorization")
	userID := extractBearer(userHeader)
	if userID == "" {
		rt.jsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	convID := ps.ByName("conversationId")
	if convID == "" {
		rt.jsonError(w, http.StatusBadRequest, "Missing conversationId")
		return
	}

	messageID := ps.ByName("messageId")
	if messageID == "" {
		rt.jsonError(w, http.StatusBadRequest, "Missing messageId")
		return
	}

	var payload struct {
		TargetConversationID string `json:"targetConversationId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		rt.jsonError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	if payload.TargetConversationID == "" {
		rt.jsonError(w, http.StatusBadRequest, "targetConversationId required")
		return
	}

	msg, err := rt.db.ForwardMessage(userID, convID, messageID, payload.TargetConversationID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrConversationNotFound):
			rt.jsonError(w, http.StatusNotFound, "Conversation not found")
		case errors.Is(err, models.ErrForbidden):
			rt.jsonError(w, http.StatusForbidden, "Forbidden")
		case errors.Is(err, models.ErrMessageNotFound):
			rt.jsonError(w, http.StatusNotFound, "Message not found")
		default:
			ctx.Logger.WithError(err).Error("Forward message error")
			rt.jsonError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(msg); err != nil {
		ctx.Logger.WithError(err).Error("JSON encode failed")
	}
}

func (rt *_router) addReactionToMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userHeader := r.Header.Get("Authorization")
	userID := extractBearer(userHeader)
	if userID == "" {
		rt.jsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	convID := ps.ByName("conversationId")
	msgID := ps.ByName("messageId")
	if convID == "" || msgID == "" {
		rt.jsonError(w, http.StatusBadRequest, "Missing conversationId or messageId")
		return
	}

	var payload struct {
		Emoji string `json:"emoji"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		ctx.Logger.WithError(err).Error("Invalid JSON body")
		rt.jsonError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	if payload.Emoji == "" {
		rt.jsonError(w, http.StatusBadRequest, "Emoji is required")
		return
	}

	ok, err := rt.db.IsUserInConversation(userID, convID)
	if err != nil {
		ctx.Logger.WithError(err).Error("IsUserInConversation check failed")
		rt.jsonError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if !ok {
		rt.jsonError(w, http.StatusForbidden, "Forbidden")
		return
	}

	reaction, err := rt.db.AddReaction(userID, convID, msgID, payload.Emoji)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrConversationNotFound):
			rt.jsonError(w, http.StatusNotFound, "Conversation not found")
		case errors.Is(err, models.ErrMessageNotFound):
			rt.jsonError(w, http.StatusNotFound, "Message not found")
		default:
			ctx.Logger.WithError(err).Error("AddReaction DB error")
			rt.jsonError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(reaction); err != nil {
		ctx.Logger.WithError(err).Error("JSON encode failed")
	}
}

func (rt *_router) removeReactionFromMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userHeader := r.Header.Get("Authorization")
	userID := extractBearer(userHeader)
	if userID == "" {
		rt.jsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	convID := ps.ByName("conversationId")
	msgID := ps.ByName("messageId")
	reactionID := ps.ByName("reactionId")

	if convID == "" || msgID == "" || reactionID == "" {
		rt.jsonError(w, http.StatusBadRequest, "Missing conversationId, messageId, or reactionId")
		return
	}

	ok, err := rt.db.IsUserInConversation(userID, convID)
	if err != nil {
		ctx.Logger.WithError(err).Error("IsUserInConversation check failed")
		rt.jsonError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if !ok {
		rt.jsonError(w, http.StatusForbidden, "Forbidden")
		return
	}

	err = rt.db.RemoveReaction(userID, convID, msgID, reactionID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrConversationNotFound):
			rt.jsonError(w, http.StatusNotFound, "Conversation not found")
		case errors.Is(err, models.ErrMessageNotFound):
			rt.jsonError(w, http.StatusNotFound, "Message not found")
		case errors.Is(err, models.ErrReactionNotFound):
			rt.jsonError(w, http.StatusNotFound, "Reaction not found")
		default:
			ctx.Logger.WithError(err).Error("RemoveReaction DB error")
			rt.jsonError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
