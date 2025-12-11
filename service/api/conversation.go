package api

import (
	"encoding/json"
	"net/http"
	"errors"
	"fmt"
	"strconv"
	"github.com/julienschmidt/httprouter"
	"github.com/aritz/wasa-homeworks/service/models"


)

func (rt *_router) getConversationProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := r.Header.Get("Authorization")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	convID := ps.ByName("conversationId")
	conv, err := rt.db.GetConversationProfile(userID, convID)
	if err != nil {
		http.Error(w, "Not found or forbidden", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conv)
}
func (rt *_router) markConversationSeen(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    userID := r.Header.Get("Authorization")
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
            fmt.Println("MARK SEEN ERROR:", err)
            http.Error(w, "Internal server error", 500)
        }
        return
    }

    w.WriteHeader(204)
}
func (rt *_router) listConversationMessages(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    userID := r.Header.Get("Authorization")
    if userID == "" {
        http.Error(w, "Unauthorized", 401)
        return
    }

    convID := ps.ByName("conversationId")
    if convID == "" {
        http.Error(w, "Missing conversationId", 400)
        return
    }

    // Query params
    limit := 50
    before := r.URL.Query().Get("before")

    // OPTIONAL limit override
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
            fmt.Println("GET MESSAGES ERROR:", err)
            http.Error(w, "Internal server error", 500)
        }
        return
    }

    w.WriteHeader(200)
    json.NewEncoder(w).Encode(msgs)
}
func (rt *_router) createPrivateConversation(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    userID := r.Header.Get("Authorization")
    
    // Leer el JSON { "targetUserId": "..." }
    var body struct {
        TargetUserID string `json:"targetUserId"`
    }
    json.NewDecoder(r.Body).Decode(&body)

    // LLAMADA MÁGICA: Si existe lo devuelve, si no lo crea
    conv, err := rt.db.GetOrCreateOneOnOneConversation(userID, body.TargetUserID)
    if err != nil {
        http.Error(w, "Error", 500)
        return
    }

    json.NewEncoder(w).Encode(conv)
}
func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := r.Header.Get("Authorization")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	convID := ps.ByName("conversationId")
	if convID == "" {
		http.Error(w, "Missing conversationId", http.StatusBadRequest)
		return
	}

	var payload struct {
		Text            string `json:"text"`
		PhotoURL        string `json:"photoUrl"`
		ReplyToMessageID string `json:"replyToMessageId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validación: al menos text o photo
	if payload.Text == "" && payload.PhotoURL == "" {
		http.Error(w, "text or photoUrl required", http.StatusBadRequest)
		return
	}

	msg, err := rt.db.SendMessage(userID, convID, payload.Text, payload.PhotoURL, payload.ReplyToMessageID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrConversationNotFound):
			http.Error(w, "Conversation not found", http.StatusNotFound)
		case errors.Is(err, models.ErrForbidden):
			http.Error(w, "Forbidden", http.StatusForbidden)
		default:
			fmt.Println("SEND MESSAGE ERROR:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)
}
func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := r.Header.Get("Authorization")
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
			fmt.Println("DELETE MESSAGE ERROR:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204
}