package database

import (
	"database/sql"
	"fmt"

	"github.com/aritz/wasa-homeworks/service/models"
)

func (db *appdbimpl) SetGroupName(userID, convID, newName string) (models.Conversation, error) {
	// 1. Verificar que la conversación existe y es grupo
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
		return models.Conversation{}, fmt.Errorf("cannot rename private conversation")
	}

	// 2. Verificar que el usuario pertenece al grupo
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

	// 3. Actualizar nombre del grupo
	_, err = db.c.Exec(`
		UPDATE conversations
		SET name = ?
		WHERE id = ?
	`, newName, convID)
	if err != nil {
		return models.Conversation{}, err
	}

	// 4. Devolver información actualizada
	return db.GetConversationProfile(userID, convID)
}
