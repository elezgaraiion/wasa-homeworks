package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) listOrSearchUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 1. Autenticación
	userID := r.Header.Get("Authorization")
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// 2. Leer query param "q"
	q := r.URL.Query().Get("q") // puede estar vacío
users, err := rt.db.ListUsers(q, userID)
if err != nil {
    http.Error(w, "internal server error", http.StatusInternalServerError)
    return
}
json.NewEncoder(w).Encode(users)
}