package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/aritz/wasa-homeworks/service/api/reqcontext"
	"github.com/aritz/wasa-homeworks/service/models"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getConversationProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userHeader := r.Header.Get("Authorization")
	userID := extractBearer(userHeader)
	if userID == "" {
		rt.jsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	convID := ps.ByName("conversationId")
	conv, err := rt.db.GetConversationProfile(userID, convID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Check conversation profile failed")
		rt.jsonError(w, http.StatusNotFound, "Not found or forbidden")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(conv); err != nil {
		ctx.Logger.WithError(err).Error("JSON encode failed")
	}
}

func (rt *_router) markConversationSeen(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	err := rt.db.MarkConversationSeen(userID, convID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrConversationNotFound):
			rt.jsonError(w, http.StatusNotFound, "Conversation not found")
		case errors.Is(err, models.ErrForbidden):
			rt.jsonError(w, http.StatusForbidden, "Forbidden")
		default:
			ctx.Logger.WithError(err).Error("Mark seen error")
			rt.jsonError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) listConversationMessages(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	limit := 50
	before := r.URL.Query().Get("before")

	if q := r.URL.Query().Get("limit"); q != "" {
		if v, err := strconv.Atoi(q); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}

	msgs, err := rt.db.GetMessagesInConversation(userID, convID, limit, before)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrConversationNotFound):
			rt.jsonError(w, http.StatusNotFound, "Conversation not found")
		case errors.Is(err, models.ErrForbidden):
			rt.jsonError(w, http.StatusForbidden, "Forbidden")
		default:
			ctx.Logger.WithError(err).Error("List messages error")
			rt.jsonError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(msgs); err != nil {
		ctx.Logger.WithError(err).Error("JSON encode failed")
	}
}

func (rt *_router) createPrivateConversation(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	userHeader := r.Header.Get("Authorization")
	userID := extractBearer(userHeader)
	if userID == "" {
		rt.jsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var body struct {
		TargetUserID string `json:"targetUserId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		ctx.Logger.WithError(err).Error("Invalid JSON body")
		rt.jsonError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	conv, err := rt.db.GetOrCreateOneOnOneConversation(userID, body.TargetUserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Create private conversation failed")
		rt.jsonError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(conv); err != nil {
		ctx.Logger.WithError(err).Error("JSON encode failed")
	}
}

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		rt.jsonError(w, http.StatusBadRequest, "Error parsing form data (expecting multipart)")
		return
	}

	text := r.FormValue("text")
	replyToMessageID := r.FormValue("replyToMessageId")

	var photoURL string
	file, header, err := r.FormFile("photoFile")
	if err == nil {
		defer file.Close()
		url, err := saveImageFile(file, header)
		if err != nil {
			ctx.Logger.WithError(err).Error("Save image file failed")
			rt.jsonError(w, http.StatusInternalServerError, "Internal server error saving image")
			return
		}
		photoURL = url
	}

	if text == "" && photoURL == "" {
		rt.jsonError(w, http.StatusBadRequest, "text or photoFile required")
		return
	}

	msg, err := rt.db.SendMessage(userID, convID, text, photoURL, replyToMessageID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrConversationNotFound):
			rt.jsonError(w, http.StatusNotFound, "Conversation not found")
		case errors.Is(err, models.ErrForbidden):
			rt.jsonError(w, http.StatusForbidden, "Forbidden")
		default:
			ctx.Logger.WithError(err).Error("Send message DB error")
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

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	err := rt.db.DeleteMessage(userID, convID, messageID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrConversationNotFound):
			rt.jsonError(w, http.StatusNotFound, "Conversation not found")
		case errors.Is(err, models.ErrForbidden):
			rt.jsonError(w, http.StatusForbidden, "Forbidden")
		case errors.Is(err, models.ErrMessageNotFound):
			rt.jsonError(w, http.StatusNotFound, "Message not found")
		default:
			ctx.Logger.WithError(err).Error("Delete message error")
			rt.jsonError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
