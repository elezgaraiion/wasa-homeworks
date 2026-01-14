package database

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/aritz/wasa-homeworks/service/models"
	"github.com/gofrs/uuid"
)

func (db *appdbimpl) createSystemMessage(tx *sql.Tx, convID, senderID, text, nowStr string) error {
	msgUUID, err := uuid.NewV4()
	if err != nil {
		return err
	}
	msgID := msgUUID.String()

	_, err = tx.Exec(`
        INSERT INTO messages (id, conversation_id, sender_id, text, created_at, status)
        VALUES (?, ?, ?, ?, ?, 'delivered')
    `, msgID, convID, senderID, text, nowStr)
	return err
}

func (db *appdbimpl) CreateGroup(creatorID string, name string, users []string) (models.Conversation, error) {
	if len(name) == 0 || len(name) > 100 {
		return models.Conversation{}, errors.New("invalid group name length")
	}
	if len(users) > 511 {
		return models.Conversation{}, errors.New("max 511 invited users allowed")
	}

	tx, err := db.c.Begin()
	if err != nil {
		log.Printf("TX BEGIN error: %v", err)
		return models.Conversation{}, err
	}
	defer func() { _ = tx.Rollback() }()

	convUUID, err := uuid.NewV4()
	if err != nil {
		return models.Conversation{}, err
	}
	convID := convUUID.String()

	baseTime := time.Now().UTC()

	_, err = tx.Exec(`
        INSERT INTO conversations(id, type, name, last_message_at, photo)
        VALUES (?, 'group', ?, ?, '')
    `, convID, name, baseTime.Format(time.RFC3339Nano))
	if err != nil {
		log.Printf("Error inserting conversation: %v", err)
		return models.Conversation{}, err
	}

	_, err = tx.Exec(`INSERT INTO conversation_participants(conversation_id, user_id) VALUES (?, ?)`, convID, creatorID)
	if err != nil {
		log.Printf("Error inserting creator part: %v", err)
		return models.Conversation{}, err
	}

	creatorJoinTime := baseTime.Format(time.RFC3339Nano)
	_, err = tx.Exec(`INSERT INTO conversation_user_meta(conversation_id, user_id, joined_at, last_seen_message_at) VALUES (?, ?, ?, ?)`, convID, creatorID, creatorJoinTime, creatorJoinTime)
	if err != nil {
		log.Printf("Error inserting creator meta: %v", err)
		return models.Conversation{}, err
	}

	msgTime1 := baseTime.Add(time.Millisecond)
	err = db.createSystemMessage(tx, convID, creatorID, "SYS:Group created", msgTime1.Format(time.RFC3339Nano))
	if err != nil {
		log.Printf("Error creating sys msg: %v", err)
		return models.Conversation{}, err
	}

	lastPreviewText := "Group created"

	for i, u := range users {
		if u == creatorID {
			continue
		}

		loopTime := baseTime.Add(time.Millisecond * time.Duration((i+2)*10))
		loopTimeStr := loopTime.Format(time.RFC3339Nano)

		_, err = tx.Exec(`INSERT INTO conversation_participants(conversation_id, user_id) VALUES (?, ?)`, convID, u)
		if err != nil {
			log.Printf("Error inserting part %s: %v", u, err)
			return models.Conversation{}, err
		}

		_, err = tx.Exec(`INSERT INTO conversation_user_meta(conversation_id, user_id, joined_at, last_seen_message_at) VALUES (?, ?, ?, ?)`, convID, u, loopTimeStr, loopTimeStr)
		if err != nil {
			log.Printf("Error inserting meta %s: %v", u, err)
			return models.Conversation{}, err
		}

		var targetName string
		err = tx.QueryRow("SELECT name FROM users WHERE id = ?", u).Scan(&targetName)
		if err != nil {
			targetName = "User"
		}

		sysText := "SYS:added " + targetName + " to the group"
		msgTimeLoop := loopTime.Add(time.Millisecond)

		err = db.createSystemMessage(tx, convID, creatorID, sysText, msgTimeLoop.Format(time.RFC3339Nano))
		if err != nil {
			log.Printf("Error adding sys msg in loop: %v", err)
			return models.Conversation{}, err
		}

		lastPreviewText = "added " + targetName + " to the group"
	}

	finalTime := baseTime.Add(time.Second).Format(time.RFC3339Nano)
	_, err = tx.Exec(`
        UPDATE conversations 
        SET last_message_preview = ?, last_message_at = ?
        WHERE id = ?
    `, lastPreviewText, finalTime, convID)
	if err != nil {
		log.Printf("Error updating preview: %v", err)
		return models.Conversation{}, err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("TX COMMIT error: %v", err)
		return models.Conversation{}, err
	}

	parts, err := db.getParticipantsByConversation(convID)
	if err != nil {
		log.Printf("Error loading participants: %v", err)
		return models.Conversation{}, err
	}

	conv := models.Conversation{
		ID:                 convID,
		Type:               "group",
		Name:               name,
		Participants:       parts,
		LastMessagePreview: lastPreviewText,
		LastMessageAt:      baseTime,
	}
	return conv, nil
}
