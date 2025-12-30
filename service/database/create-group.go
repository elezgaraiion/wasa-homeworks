package database

import (
	"fmt"
	"time"
	"github.com/google/uuid"
	"github.com/aritz/wasa-homeworks/service/models"
)


func (db *appdbimpl) CreateGroup(creatorID string, name string, users []string) (models.Conversation, error) {

	if len(name) == 0 || len(name) > 100 {
		return models.Conversation{}, fmt.Errorf("invalid group name length")
	}

	if len(users) > 511 {
		return models.Conversation{}, fmt.Errorf("max 511 invited users allowed")
	}
	tx, err := db.c.Begin()
	if err != nil {
		return models.Conversation{}, fmt.Errorf("TX BEGIN: %w", err)
	}
	defer tx.Rollback()

	convID := uuid.New().String()
	_, err = tx.Exec(`
		INSERT INTO conversations(id, type, name)
		VALUES (?, 'group', ?)
	`, convID, name)
	if err != nil {
		return models.Conversation{}, fmt.Errorf("INSERT conversation: %w", err)
	}

	joinedAt := time.Now().UTC().Format(time.RFC3339)

	_, err = tx.Exec(`
		INSERT INTO conversation_participants(conversation_id, user_id)
		VALUES (?, ?)
	`, convID, creatorID)
	if err != nil {
		return models.Conversation{}, fmt.Errorf("INSERT creator participant: %w", err)
	}

	_, err = tx.Exec(`
		INSERT INTO conversation_user_meta(conversation_id, user_id, joined_at)
		VALUES (?, ?, ?)
	`, convID, creatorID, joinedAt)
	if err != nil {
		return models.Conversation{}, fmt.Errorf("INSERT creator meta: %w", err)
	}

	for _, u := range users {
		if u == creatorID {
			continue 
		}

		_, err = tx.Exec(`
			INSERT INTO conversation_participants(conversation_id, user_id)
			VALUES (?, ?)
		`, convID, u)
		if err != nil {
			return models.Conversation{}, fmt.Errorf("INSERT participant %s: %w", u, err)
		}

		_, err = tx.Exec(`
			INSERT INTO conversation_user_meta(conversation_id, user_id, joined_at)
			VALUES (?, ?, ?)
		`, convID, u, joinedAt)
		if err != nil {
			return models.Conversation{}, fmt.Errorf("INSERT meta %s: %w", u, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return models.Conversation{}, fmt.Errorf("TX COMMIT: %w", err)
	}

	parts, err := db.getParticipantsByConversation(convID)
	if err != nil {
    	return models.Conversation{}, fmt.Errorf("load participants: %w", err)
	}
	conv := models.Conversation{
    	ID:           convID,
    	Type:         "group",
    	Name:         name,
    	Participants: parts,
	}
	return conv, nil
}
