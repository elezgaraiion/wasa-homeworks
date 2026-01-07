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
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	q := r.URL.Query().Get("q")
	users, err := rt.db.ListUsers(q, userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("ListUsers DB error")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	userHeader := r.Header.Get("Authorization")
	userID := extractBearer(userHeader)
	if userID == "" {
		http.Error(w, "Unauthorized", 401)
		return
	}

	var body struct {
		Name  string   `json:"name"`
		Users []string `json:"users"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		ctx.Logger.WithError(err).Error("Invalid JSON in createGroup")
		http.Error(w, "Invalid JSON", 400)
		return
	}
	if len(body.Name) == 0 {
		http.Error(w, "Group name required", 400)
		return
	}

	conv, err := rt.db.CreateGroup(userID, body.Name, body.Users)
	if err != nil {
		ctx.Logger.WithError(err).Error("CreateGroup DB error")
		http.Error(w, "Internal Server Error", 500)
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(conv)
}