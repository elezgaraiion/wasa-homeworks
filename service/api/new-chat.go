package api

import (
	"encoding/json"
	"net/http"

	"github.com/aritz/wasa-homeworks/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) listOrSearchUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	userHeader := r.Header.Get("Authorization")
	userID := extractBearer(userHeader)
	if userID == "" {
		rt.jsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	q := r.URL.Query().Get("q")
	users, err := rt.db.ListUsers(q, userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("ListUsers DB error")
		rt.jsonError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		ctx.Logger.WithError(err).Error("JSON encode failed")
	}
}

func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	userHeader := r.Header.Get("Authorization")
	userID := extractBearer(userHeader)
	if userID == "" {
		rt.jsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var body struct {
		Name  string   `json:"name"`
		Users []string `json:"users"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		ctx.Logger.WithError(err).Error("Invalid JSON in createGroup")
		rt.jsonError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	if len(body.Name) == 0 {
		rt.jsonError(w, http.StatusBadRequest, "Group name required")
		return
	}

	conv, err := rt.db.CreateGroup(userID, body.Name, body.Users)
	if err != nil {
		ctx.Logger.WithError(err).Error("CreateGroup DB error")
		rt.jsonError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(conv); err != nil {
		ctx.Logger.WithError(err).Error("JSON encode failed")
	}
}