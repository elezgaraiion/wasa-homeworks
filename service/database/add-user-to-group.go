package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/aritz/wasa-homeworks/service/models"
)

func (db *appdbimpl) AddUserToGroup(requestUserID, convID, targetUserID string) (models.Conversation, error) {

	var exists int
	err := db.c.QueryRow(`
		SELECT COUNT(*)
		FROM conversation_participants
		WHERE conversation_id = ? AND user_id = ?
	`, convID, requestUserID).Scan(&exists)

	if err != nil {
		return models.Conversation{}, err
	}
	if exists == 0 {
		return models.Conversation{}, models.ErrForbidden
	}

	var convType string
	err = db.c.QueryRow(`
		SELECT type FROM conversations WHERE id = ?
	`, convID).Scan(&convType)

	if err == sql.ErrNoRows {
		return models.Conversation{}, models.ErrConversationNotFound
	}
	if err != nil {
		return models.Conversation{}, err
	}
	if convType != "group" {
		return models.Conversation{}, fmt.Errorf("cannot add users to private conversations")
	}

	err = db.c.QueryRow(`
		SELECT COUNT(*)
		FROM conversation_participants
		WHERE conversation_id = ? AND user_id = ?
	`, convID, targetUserID).Scan(&exists)

	if err != nil {
		return models.Conversation{}, err
	}
	if exists > 0 {
		return db.GetConversationProfile(requestUserID, convID)
	}

	_, err = db.c.Exec(`
		INSERT INTO conversation_participants(conversation_id, user_id)
		VALUES (?, ?)
	`, convID, targetUserID)
	if err != nil {
		return models.Conversation{}, err
	}

	joinedAt := time.Now().UTC().Format(time.RFC3339)

	_, err = db.c.Exec(`
		INSERT INTO conversation_user_meta(conversation_id, user_id, joined_at)
		VALUES (?, ?, ?)
	`, convID, targetUserID, joinedAt)
	if err != nil {
		return models.Conversation{}, err
	}

	return db.GetConversationProfile(requestUserID, convID)
}
