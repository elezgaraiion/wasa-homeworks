package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/aritz/wasa-homeworks/service/models"
	"github.com/gofrs/uuid"
)

func (db *appdbimpl) GetOrCreateOneOnOneConversation(myID, targetUserID string) (models.Conversation, error) {

	var existingConvID string
	err := db.c.QueryRow(`
        SELECT c.id
        FROM conversations c
        JOIN conversation_participants p1 ON p1.conversation_id = c.id
        JOIN conversation_participants p2 ON p2.conversation_id = c.id
        WHERE c.type = 'direct'
          AND p1.user_id = ?
          AND p2.user_id = ?
        LIMIT 1
    `, myID, targetUserID).Scan(&existingConvID)

	if err == nil {
		parts, err := db.getParticipantsByConversation(existingConvID)
		if err != nil {
			return models.Conversation{}, err
		}

		var targetName string
		var targetPhoto sql.NullString
		err = db.c.QueryRow("SELECT name, photo FROM users WHERE id = ?", targetUserID).Scan(&targetName, &targetPhoto)
		if err != nil {
			return models.Conversation{}, err
		}

		return models.Conversation{
			ID:           existingConvID,
			Type:         "direct",
			Name:         targetName,
			Photo:        targetPhoto.String,
			Participants: parts,
		}, nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return models.Conversation{}, err
	}

	tx, err := db.c.Begin()
	if err != nil {
		return models.Conversation{}, err
	}
	defer func() { _ = tx.Rollback() }()

	newUUID, err := uuid.NewV4()
	if err != nil {
		return models.Conversation{}, err
	}
	newID := newUUID.String()
	now := time.Now().UTC().Format(time.RFC3339Nano)

	_, err = tx.Exec(`INSERT INTO conversations (id, type) VALUES (?, 'direct')`, newID)
	if err != nil {
		return models.Conversation{}, err
	}

	_, err = tx.Exec(`INSERT INTO conversation_participants (conversation_id, user_id) VALUES (?, ?)`, newID, myID)
	if err != nil {
		return models.Conversation{}, err
	}
	_, err = tx.Exec(`INSERT INTO conversation_user_meta (conversation_id, user_id, joined_at) VALUES (?, ?, ?)`, newID, myID, now)
	if err != nil {
		return models.Conversation{}, err
	}

	_, err = tx.Exec(`INSERT INTO conversation_participants (conversation_id, user_id) VALUES (?, ?)`, newID, targetUserID)
	if err != nil {
		return models.Conversation{}, err
	}
	_, err = tx.Exec(`INSERT INTO conversation_user_meta (conversation_id, user_id, joined_at) VALUES (?, ?, ?)`, newID, targetUserID, now)
	if err != nil {
		return models.Conversation{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.Conversation{}, err
	}

	parts, err := db.getParticipantsByConversation(newID)
	if err != nil {
		return models.Conversation{}, err
	}

	var targetName string
	var targetPhoto sql.NullString
	err = db.c.QueryRow("SELECT name, photo FROM users WHERE id = ?", targetUserID).Scan(&targetName, &targetPhoto)
	if err != nil {
		return models.Conversation{}, err
	}

	return models.Conversation{
		ID:           newID,
		Type:         "direct",
		Name:         targetName,
		Photo:        targetPhoto.String,
		Participants: parts,
	}, nil
}