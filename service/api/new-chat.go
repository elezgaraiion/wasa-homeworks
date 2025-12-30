package api

import (
	"encoding/json"
	"net/http"
	"fmt"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) listOrSearchUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID := r.Header.Get("Authorization")
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	q := r.URL.Query().Get("q") 
users, err := rt.db.ListUsers(q, userID)
if err != nil {
    http.Error(w, "internal server error", http.StatusInternalServerError)
    return
}
json.NewEncoder(w).Encode(users)
}
func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {
	userID := r.Header.Get("Authorization")
	if userID == "" {
		http.Error(w, "Unauthorized", 401)
		return
	}

	var body struct {
		Name  string   `json:"name"`
		Users []string `json:"users"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}
	if len(body.Name) == 0 {
		http.Error(w, "Group name required", 400)
		return
	}

	conv, err := rt.db.CreateGroup(userID, body.Name, body.Users)
	if err != nil {
		fmt.Println("CREATE GROUP ERROR:", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(conv)
}