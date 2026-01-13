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

func (rt *_router) doGetCurrentUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	userHeader := r.Header.Get("Authorization")
	userID := extractBearer(userHeader)
	if userID == "" {
		rt.jsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := rt.db.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			rt.jsonError(w, http.StatusNotFound, "user not found")
			return
		}
		ctx.Logger.WithError(err).Error("GetUserByID failed")
		rt.jsonError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		ctx.Logger.WithError(err).Error("JSON encode failed")
	}
}

func (rt *_router) updateMyUserName(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	userHeader := r.Header.Get("Authorization")
	userID := extractBearer(userHeader)
	if userID == "" {
		rt.jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rt.jsonError(w, http.StatusBadRequest, "bad request")
		return
	}

	if len(req.Name) < 1 || len(req.Name) > 50 {
		rt.jsonError(w, http.StatusBadRequest, "bad request")
		return
	}

	updatedUser, err := rt.db.UpdateUserName(userID, req.Name)
	if err != nil {
		if errors.Is(err, models.ErrNameConflict) {
			rt.jsonError(w, http.StatusConflict, "user with that name already exists")
			return
		}
		ctx.Logger.WithError(err).Error("UpdateUserName failed")
		rt.jsonError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(updatedUser); err != nil {
		ctx.Logger.WithError(err).Error("JSON encode failed")
	}
}

func (rt *_router) updatePhoto(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	userHeader := r.Header.Get("Authorization")
	userID := extractBearer(userHeader)
	if userID == "" {
		rt.jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		rt.jsonError(w, http.StatusBadRequest, "bad request")
		return
	}

	file, handler, err := r.FormFile("photoFile")
	if err != nil {
		rt.jsonError(w, http.StatusBadRequest, "bad request: missing photoFile")
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
		ctx.Logger.WithError(err).Error("saveImageFile failed")
		rt.jsonError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	user, err := rt.db.UpdateMyPhoto(userID, photoURL)
	if err != nil {
		ctx.Logger.WithError(err).Error("UpdateMyPhoto failed")
		rt.jsonError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		ctx.Logger.WithError(err).Error("JSON encode failed")
	}
}