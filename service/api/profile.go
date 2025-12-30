package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"github.com/julienschmidt/httprouter"
	"os"
	"image"

)

func (rt *_router) doGetCurrentUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    userID := r.Header.Get("Authorization")
    
    if userID == "" {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})       
        return
    }



    user, err := rt.db.GetUserByID(userID)
    if err != nil {
        http.Error(w, "user not found", http.StatusUnauthorized)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func (rt *_router) updateMyUserName(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	
		userID := r.Header.Get("Authorization")
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
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedUser)
	}
func (rt *_router) updatePhoto(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID := r.Header.Get("Authorization")
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
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	contentType := handler.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		http.Error(w, "unsupported media type", http.StatusUnsupportedMediaType)
		return
	}
	dst := "/tmp/" + handler.Filename
	f, err := os.Create(dst)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer f.Close()
	_, format, err := image.Decode(file)
	if err != nil || format != "jpeg" {
    	http.Error(w, "unsupported media type: only JPEG allowed", http.StatusUnsupportedMediaType)
    	return
	}

	_, err = file.Seek(0, 0)
	if err != nil {
    	http.Error(w, "internal server error", http.StatusInternalServerError)
    	return
	}

	user, err := rt.db.UpdateMyPhoto(userID, dst)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}
