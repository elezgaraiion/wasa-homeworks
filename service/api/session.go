package api

import (
	"encoding/json"
	"net/http"

	"github.com/aritz/wasa-homeworks/service/models"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 1. Parsear body
	var req struct {
		Name string `json:"name"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || len(req.Name) < 3 || len(req.Name) > 16 {
		http.Error(w, "invalid name", http.StatusBadRequest)
		return
	}

	// 2. Buscar si existe usuario
	id, err := rt.db.GetUserIdByName(req.Name)
	if err == nil {
		// Existe → devolver solo el ID
		json.NewEncoder(w).Encode(map[string]string{"identifier": id})
		return
	}

	// 3. Crear usuario nuevo
	newUser := models.User{
		ID:    uuid.New().String(),
		Name:  req.Name,
		Photo: "",
	}

	err = rt.db.CreateUser(newUser)
	if err != nil {
		http.Error(w, "cannot create user", http.StatusInternalServerError)
		return
	}

	// 4. Devolver solo el ID del nuevo usuario
	json.NewEncoder(w).Encode(map[string]string{"identifier": newUser.ID})
}
func (rt *_router) GetMyConversations(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID := r.Header.Get("Authorization")
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	convs, err := rt.db.GetMyConversations(userID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(convs)
}