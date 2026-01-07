package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aritz/wasa-homeworks/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doGetCurrentUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	userHeader := r.Header.Get("Authorization")
	userID := extractBearer(userHeader)
	if userID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	user, err := rt.db.GetUserByID(userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("GetUserByID failed")
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (rt *_router) updateMyUserName(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	userHeader := r.Header.Get("Authorization")
	userID := extractBearer(userHeader)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || len(req.Name) < 1 || len(req.Name) > 50 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	updatedUser, err := rt.db.UpdateUserName(userID, req.Name)
	if err != nil {
		if err.Error() == "user with that name already exists" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		ctx.Logger.WithError(err).Error("UpdateUserName failed")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

func (rt *_router) updatePhoto(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	userHeader := r.Header.Get("Authorization")
	userID := extractBearer(userHeader)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("photoFile")
	if err != nil {
		http.Error(w, "bad request: missing photoFile", http.StatusBadRequest)
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
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	user, err := rt.db.UpdateMyPhoto(userID, photoURL)
	if err != nil {
		ctx.Logger.WithError(err).Error("UpdateMyPhoto failed")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}