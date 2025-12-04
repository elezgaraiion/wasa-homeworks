package database

import (
	"database/sql"
	"fmt"
	
	"github.com/aritz/wasa-homeworks/service/models"
)

func (db *appdbimpl) SetGroupPhoto(userID, convID, photoURL string) (models.Conversation, error) {
	// 1. Verificar conversación y tipo
	var convType string
	err := db.c.QueryRow(`
		SELECT type FROM conversations WHERE id = ?
	`, convID).Scan(&convType)
	if err == sql.ErrNoRows {
		return models.Conversation{}, models.ErrConversationNotFound
	}
	if err != nil {
		return models.Conversation{}, err
	}
	if convType != "group" {
		return models.Conversation{}, fmt.Errorf("cannot set photo for private conversation")
	}

	// 2. Verificar que usuario pertenece al grupo
	var exists int
	err = db.c.QueryRow(`
		SELECT COUNT(*) FROM conversation_participants
		WHERE conversation_id = ? AND user_id = ?
	`, convID, userID).Scan(&exists)
	if err != nil {
		return models.Conversation{}, err
	}
	if exists == 0 {
		return models.Conversation{}, models.ErrForbidden
	}

	// 3. Actualizar photo en DB
	_, err = db.c.Exec(`
		UPDATE conversations
		SET photo = ?
		WHERE id = ?
	`, photoURL, convID)
	if err != nil {
		return models.Conversation{}, err
	}

	// 4. Devolver información actualizada
	return db.GetConversationProfile(userID, convID)
}
