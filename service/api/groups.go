package api

import (
	"encoding/json"
	"errors"
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
		rt.jsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	conversationID := ps.ByName("conversationId")
	if conversationID == "" {
		rt.jsonError(w, http.StatusBadRequest, "Missing conversationId")
		return
	}

	var body addUserToGroupReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		rt.jsonError(w, http.StatusBadRequest, "Bad Request")
		return
	}

	if body.UserID == "" {
		rt.jsonError(w, http.StatusBadRequest, "Missing userId")
		return
	}

	conv, err := rt.db.AddUserToGroup(authUserID, conversationID, body.UserID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrForbidden):
			rt.jsonError(w, http.StatusForbidden, "Forbidden")
		case errors.Is(err, models.ErrConversationNotFound):
			rt.jsonError(w, http.StatusNotFound, "Conversation not found")
		case errors.Is(err, models.ErrUserNotFound):
			rt.jsonError(w, http.StatusNotFound, "User not found")
		default:
			ctx.Logger.WithError(err).Error("AddUserToGroup DB error")
			rt.jsonError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(conv); err != nil {
		ctx.Logger.WithError(err).Error("JSON encode failed")
	}
}

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authUserHeader := r.Header.Get("Authorization")
	authUserID := extractBearer(authUserHeader)
	if authUserID == "" {
		rt.jsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	conversationID := ps.ByName("conversationId")
	if conversationID == "" {
		rt.jsonError(w, http.StatusBadRequest, "Missing conversationId")
		return
	}

	err := rt.db.LeaveGroup(authUserID, conversationID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrConversationNotFound):
			rt.jsonError(w, http.StatusNotFound, "Conversation not found")
		case errors.Is(err, models.ErrForbidden):
			rt.jsonError(w, http.StatusForbidden, "Forbidden")
		default:
			ctx.Logger.WithError(err).Error("LeaveGroup DB error")
			rt.jsonError(w, http.StatusInternalServerError, "Internal server error")
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
		rt.jsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	conversationID := ps.ByName("conversationId")
	if conversationID == "" {
		rt.jsonError(w, http.StatusBadRequest, "Missing conversationId")
		return
	}

	var body setGroupNameReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		rt.jsonError(w, http.StatusBadRequest, "Bad Request")
		return
	}

	if len(body.Name) == 0 || len(body.Name) > 100 {
		rt.jsonError(w, http.StatusBadRequest, "Invalid group name")
		return
	}

	conv, err := rt.db.SetGroupName(authUserID, conversationID, body.Name)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrConversationNotFound):
			rt.jsonError(w, http.StatusNotFound, "Conversation not found")
		case errors.Is(err, models.ErrForbidden):
			rt.jsonError(w, http.StatusForbidden, "Forbidden")
		default:
			ctx.Logger.WithError(err).Error("SetGroupName DB error")
			rt.jsonError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(conv); err != nil {
		ctx.Logger.WithError(err).Error("JSON encode failed")
	}
}

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	authUserHeader := r.Header.Get("Authorization")
	authUserID := extractBearer(authUserHeader)
	if authUserID == "" {
		rt.jsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	conversationID := ps.ByName("conversationId")
	if conversationID == "" {
		rt.jsonError(w, http.StatusBadRequest, "Missing conversationId")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		rt.jsonError(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	file, handler, err := r.FormFile("photoFile")
	if err != nil {
		rt.jsonError(w, http.StatusBadRequest, "Missing photoFile")
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
		rt.jsonError(w, http.StatusInternalServerError, "Failed to save photo")
		return
	}

	conv, err := rt.db.SetGroupPhoto(authUserID, conversationID, photoURL)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrForbidden):
			rt.jsonError(w, http.StatusForbidden, "Forbidden")
		case errors.Is(err, models.ErrConversationNotFound):
			rt.jsonError(w, http.StatusNotFound, "Conversation not found")
		default:
			ctx.Logger.WithError(err).Error("SetGroupPhoto DB error")
			rt.jsonError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(conv); err != nil {
		ctx.Logger.WithError(err).Error("JSON encode failed")
	}
}