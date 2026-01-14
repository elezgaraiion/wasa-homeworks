package api

import (
	"encoding/json"
	"net/http"

	"github.com/aritz/wasa-homeworks/service/api/reqcontext"
	"github.com/aritz/wasa-homeworks/service/models"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.Name) < 3 || len(req.Name) > 16 {
		rt.jsonError(w, http.StatusBadRequest, "invalid name (must be 3-16 chars)")
		return
	}

	id, err := rt.db.GetUserIdByName(req.Name)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"identifier": id}); err != nil {
			ctx.Logger.WithError(err).Error("JSON encode failed")
		}
		return
	}

	newUUID, err := uuid.NewV4()
	if err != nil {
		ctx.Logger.WithError(err).Error("UUID generation failed")
		rt.jsonError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	newUser := models.User{
		ID:    newUUID.String(),
		Name:  req.Name,
		Photo: "",
	}

	err = rt.db.CreateUser(newUser)
	if err != nil {
		ctx.Logger.WithError(err).Error("cannot create user")
		rt.jsonError(w, http.StatusInternalServerError, "cannot create user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"identifier": newUser.ID}); err != nil {
		ctx.Logger.WithError(err).Error("JSON encode failed")
	}
}

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	userHeader := r.Header.Get("Authorization")
	userID := extractBearer(userHeader)
	if userID == "" {
		rt.jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	convs, err := rt.db.GetMyConversations(userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("GetMyConversations failed")
		rt.jsonError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(convs); err != nil {
		ctx.Logger.WithError(err).Error("JSON encode failed")
	}
}
