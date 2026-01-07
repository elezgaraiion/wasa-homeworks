package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aritz/wasa-homeworks/service/api/reqcontext" 
	"github.com/aritz/wasa-homeworks/service/models"
	"github.com/julienschmidt/httprouter"
)

type addUserToGroupReq struct {
	UserID string `json:"userId"`
}

func (rt *_router) addUserToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authUserHeader := r.Header.Get("Authorization")
	authUserID := extractBearer(authUserHeader)
	if authUserID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conversationID := ps.ByName("conversationId")
	if conversationID == "" {
		http.Error(w, "Missing conversationId", http.StatusBadRequest)
		return
	}

	var body addUserToGroupReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if body.UserID == "" {
		http.Error(w, "Missing userId", http.StatusBadRequest)
		return
	}

	conv, err := rt.db.AddUserToGroup(authUserID, conversationID, body.UserID)
	if err != nil {
		switch err {
		case models.ErrForbidden:
			http.Error(w, "Forbidden", http.StatusForbidden)
		case models.ErrConversationNotFound:
			http.Error(w, "Conversation not found", http.StatusNotFound)
		case models.ErrUserNotFound:
			http.Error(w, "User not found", http.StatusNotFound)
		default:
			ctx.Logger.WithError(err).Error("AddUserToGroup DB error")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(conv)
}

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authUserHeader := r.Header.Get("Authorization")
	authUserID := extractBearer(authUserHeader)
	if authUserID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conversationID := ps.ByName("conversationId")
	if conversationID == "" {
		http.Error(w, "Missing conversationId", http.StatusBadRequest)
		return
	}

	err := rt.db.LeaveGroup(authUserID, conversationID)
	if err != nil {
		switch err {
		case models.ErrConversationNotFound:
			http.Error(w, "Conversation not found", http.StatusNotFound)
		case models.ErrForbidden:
			http.Error(w, "Forbidden", http.StatusForbidden)
		default:
			ctx.Logger.WithError(err).Error("LeaveGroup DB error")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type setGroupNameReq struct {
	Name string `json:"name"`
}

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authUserHeader := r.Header.Get("Authorization")
	authUserID := extractBearer(authUserHeader)
	if authUserID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conversationID := ps.ByName("conversationId")
	if conversationID == "" {
		http.Error(w, "Missing conversationId", http.StatusBadRequest)
		return
	}

	var body setGroupNameReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if len(body.Name) == 0 || len(body.Name) > 100 {
		http.Error(w, "Invalid group name", http.StatusBadRequest)
		return
	}

	conv, err := rt.db.SetGroupName(authUserID, conversationID, body.Name)
	if err != nil {
		switch err {
		case models.ErrConversationNotFound:
			http.Error(w, "Conversation not found", http.StatusNotFound)
		case models.ErrForbidden:
			http.Error(w, "Forbidden", http.StatusForbidden)
		default:
			ctx.Logger.WithError(err).Error("SetGroupName DB error")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(conv)
}

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authUserHeader := r.Header.Get("Authorization")
	authUserID := extractBearer(authUserHeader)
	if authUserID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conversationID := ps.ByName("conversationId")
	if conversationID == "" {
		http.Error(w, "Missing conversationId", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("photoFile")
	if err != nil {
		http.Error(w, "Missing photoFile", http.StatusBadRequest)
		return
	}
	defer file.Close()

	contentType := handler.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		http.Error(w, "unsupported media type", http.StatusUnsupportedMediaType)
		return
	}

	photoURL, err := saveImageFile(file, handler)
	if err != nil {
		ctx.Logger.WithError(err).Error("SaveImageFile failed")
		http.Error(w, "Failed to save photo", http.StatusInternalServerError)
		return
	}

	conv, err := rt.db.SetGroupPhoto(authUserID, conversationID, photoURL)
	if err != nil {
		switch err {
		case models.ErrForbidden:
			http.Error(w, "Forbidden", http.StatusForbidden)
		case models.ErrConversationNotFound:
			http.Error(w, "Conversation not found", http.StatusNotFound)
		default:
			ctx.Logger.WithError(err).Error("SetGroupPhoto DB error")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(conv)
}