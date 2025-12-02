package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) GetConversationProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := r.Header.Get("Authorization")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	convID := ps.ByName("conversationId")
	conv, err := rt.db.GetConversationProfile(userID, convID)
	if err != nil {
		http.Error(w, "Not found or forbidden", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conv)
}
