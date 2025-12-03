package api

import (
	"encoding/json"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/aritz/wasa-homeworks/service/models"
)

type addUserToGroupReq struct {
	UserID string `json:"userId"`
}

func (rt *_router) addUserToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    authUserID := r.Header.Get("Authorization")
    if authUserID == "" {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    conversationID := ps.ByName("conversationId")
    if conversationID == "" {
        http.Error(w, "Missing conversationId", http.StatusBadRequest)
        return
    }

    // parse body
    var body addUserToGroupReq
    if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }

    if body.UserID == "" {
        http.Error(w, "Missing userId", http.StatusBadRequest)
        return
    }

    // ✅ Llamada a la DB con IDs de usuario, no structs
    conv, err := rt.db.AddUserToGroup(authUserID, conversationID, body.UserID)
    if err != nil {
        switch err {
        case models.ErrForbidden:
            http.Error(w, "Forbidden", http.StatusForbidden)
        case models.ErrConversationNotFound:
            http.Error(w, "Conversation not found", http.StatusNotFound)
        case models.ErrUserNotFound:
            http.Error(w, "User not found", http.StatusNotFound)
        default:
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(conv)
}
func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    authUserID := r.Header.Get("Authorization")
    if authUserID == "" {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    conversationID := ps.ByName("conversationId")
    if conversationID == "" {
        http.Error(w, "Missing conversationId", http.StatusBadRequest)
        return
    }

    err := rt.db.LeaveGroup(authUserID, conversationID)
    if err != nil {
        switch err {
        case models.ErrConversationNotFound:
            http.Error(w, "Conversation not found", http.StatusNotFound)
        case models.ErrForbidden:
            http.Error(w, "Forbidden", http.StatusForbidden)
        default:
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
        return
    }

    w.WriteHeader(http.StatusNoContent) // 204
}