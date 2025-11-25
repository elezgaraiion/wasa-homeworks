package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"github.com/julienschmidt/httprouter"

)

func (rt *_router) doGetCurrentUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})		
		return
	}

	parts := strings.Fields(authHeader)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		http.Error(w, "invalid authorization header", http.StatusUnauthorized)
		return
	}
	userID := parts[1]

	user, err := rt.db.GetUserByID(userID)
	if err != nil {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (rt *_router) updateMyUserName(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	
		// 1. Obtener el ID del usuario desde el header Authorization
		userID := r.Header.Get("Authorization")
		if userID == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
	
		// 2. Parsear el body JSON
		var req struct {
			Name string `json:"name"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil || len(req.Name) < 1 || len(req.Name) > 50 {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
	
		// 3. Actualizar el nombre en la base de datos
		updatedUser, err := rt.db.UpdateUserName(userID, req.Name)
		if err != nil {
			if err.Error() == "user with that name already exists" {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	
		// 4. Devolver el usuario actualizado como JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedUser)
	}
	