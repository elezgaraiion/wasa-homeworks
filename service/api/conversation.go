package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/aritz/wasa-homeworks/service/models"
	"github.com/julienschmidt/httprouter"
    
	"github.com/aritz/wasa-homeworks/service/api/reqcontext" 
)

func (rt *_router) getConversationProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userHeader := r.Header.Get("Authorization")
	userID := extractBearer(userHeader)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	convID := ps.ByName("conversationId")
	conv, err := rt.db.GetConversationProfile(userID, convID)
	if err != nil {
        ctx.Logger.WithError(err).Error("Check conversation profile failed")
		http.Error(w, "Not found or forbidden", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conv)
}

func (rt *_router) markConversationSeen(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	err := rt.db.MarkConversationSeen(userID, convID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrConversationNotFound):
			http.Error(w, "Conversation not found", 404)
		case errors.Is(err, models.ErrForbidden):
			http.Error(w, "Forbidden", 403)
		default:
			ctx.Logger.WithError(err).Error("Mark seen error")
			http.Error(w, "Internal server error", 500)
		}
		return
	}

	w.WriteHeader(204)
}

func (rt *_router) listConversationMessages(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
			http.Error(w, "Conversation not found", 404)
		case errors.Is(err, models.ErrForbidden):
			http.Error(w, "Forbidden", 403)
		default:
			ctx.Logger.WithError(err).Error("List messages error")
			http.Error(w, "Internal server error", 500)
		}
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(msgs)
}

func (rt *_router) createPrivateConversation(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	userHeader := r.Header.Get("Authorization")
	userID := extractBearer(userHeader)
	
    var body struct {
		TargetUserID string `json:"targetUserId"`
	}
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        ctx.Logger.WithError(err).Error("Invalid JSON body")
        http.Error(w, "Invalid JSON", 400)
        return
    }

	conv, err := rt.db.GetOrCreateOneOnOneConversation(userID, body.TargetUserID)
	if err != nil {
        ctx.Logger.WithError(err).Error("Create private conversation failed")
		http.Error(w, "Internal Server Error", 500)
		return
	}

	json.NewEncoder(w).Encode(conv)
}

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing form data (expecting multipart)", http.StatusBadRequest)
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
			http.Error(w, "Internal server error saving image", http.StatusInternalServerError)
			return
		}
		photoURL = url
	}

	if text == "" && photoURL == "" {
		http.Error(w, "text or photoFile required", http.StatusBadRequest)
		return
	}

	msg, err := rt.db.SendMessage(userID, convID, text, photoURL, replyToMessageID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrConversationNotFound):
			http.Error(w, "Conversation not found", http.StatusNotFound)
		case errors.Is(err, models.ErrForbidden):
			http.Error(w, "Forbidden", http.StatusForbidden)
		default:
			ctx.Logger.WithError(err).Error("Send message DB error")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)
}

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	err := rt.db.DeleteMessage(userID, convID, messageID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrConversationNotFound):
			http.Error(w, "Conversation not found", http.StatusNotFound)
		case errors.Is(err, models.ErrForbidden):
			http.Error(w, "Forbidden", http.StatusForbidden)
		case errors.Is(err, models.ErrMessageNotFound):
			http.Error(w, "Message not found", http.StatusNotFound)
		default:
			ctx.Logger.WithError(err).Error("Delete message error")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}