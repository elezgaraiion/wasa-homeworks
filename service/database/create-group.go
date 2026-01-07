package database

import (
	"fmt"
	"time"
	"github.com/google/uuid"
	"github.com/aritz/wasa-homeworks/service/models"
	"database/sql"
)


func (db *appdbimpl) createSystemMessage(tx *sql.Tx, convID, senderID, text, nowStr string) error {
	msgID := uuid.New().String()
	_, err := tx.Exec(`
		INSERT INTO messages (id, conversation_id, sender_id, text, created_at, status)
		VALUES (?, ?, ?, ?, ?, 'delivered')
	`, msgID, convID, senderID, text, nowStr)
	return err
}

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
	baseTime := time.Now().UTC()
	
	_, err = tx.Exec(`
		INSERT INTO conversations(id, type, name, last_message_at, photo)
		VALUES (?, 'group', ?, ?, '')
	`, convID, name, baseTime.Format(time.RFC3339Nano))
	if err != nil {
		return models.Conversation{}, fmt.Errorf("INSERT conversation: %w", err)
	}

	_, err = tx.Exec(`INSERT INTO conversation_participants(conversation_id, user_id) VALUES (?, ?)`, convID, creatorID)
	if err != nil { return models.Conversation{}, fmt.Errorf("INSERT creator part: %w", err) }
	
	creatorJoinTime := baseTime.Format(time.RFC3339Nano)
	_, err = tx.Exec(`INSERT INTO conversation_user_meta(conversation_id, user_id, joined_at, last_seen_message_at) VALUES (?, ?, ?, ?)`, convID, creatorID, creatorJoinTime, creatorJoinTime)
	if err != nil { return models.Conversation{}, fmt.Errorf("INSERT creator meta: %w", err) }

	msgTime1 := baseTime.Add(time.Millisecond)
	err = db.createSystemMessage(tx, convID, creatorID, "SYS:Group created", msgTime1.Format(time.RFC3339Nano))
	if err != nil {
		return models.Conversation{}, fmt.Errorf("sys msg created: %w", err)
	}

	lastPreviewText := "Group created"
	
	for i, u := range users {
		if u == creatorID { continue }

		loopTime := baseTime.Add(time.Millisecond * time.Duration((i+2)*10))
		loopTimeStr := loopTime.Format(time.RFC3339Nano)

		_, err = tx.Exec(`INSERT INTO conversation_participants(conversation_id, user_id) VALUES (?, ?)`, convID, u)
		if err != nil { return models.Conversation{}, fmt.Errorf("INSERT part %s: %w", u, err) }

		_, err = tx.Exec(`INSERT INTO conversation_user_meta(conversation_id, user_id, joined_at, last_seen_message_at) VALUES (?, ?, ?, ?)`, convID, u, loopTimeStr, loopTimeStr)
		if err != nil { return models.Conversation{}, fmt.Errorf("INSERT meta %s: %w", u, err) }

		var targetName string
		err = tx.QueryRow("SELECT name FROM users WHERE id = ?", u).Scan(&targetName)
		if err != nil { targetName = "User" } 

		sysText := fmt.Sprintf("SYS:added %s to the group", targetName)
		msgTimeLoop := loopTime.Add(time.Millisecond) 
		
		err = db.createSystemMessage(tx, convID, creatorID, sysText, msgTimeLoop.Format(time.RFC3339Nano))
		if err != nil { return models.Conversation{}, fmt.Errorf("sys msg added: %w", err) }
		
		lastPreviewText = fmt.Sprintf("added %s to the group", targetName)
	}

	finalTime := baseTime.Add(time.Second).Format(time.RFC3339Nano)
	_, err = tx.Exec(`
		UPDATE conversations 
		SET last_message_preview = ?, last_message_at = ?
		WHERE id = ?
	`, lastPreviewText, finalTime, convID)
	if err != nil {
		return models.Conversation{}, fmt.Errorf("UPDATE preview: %w", err)
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
		LastMessagePreview: lastPreviewText,
		LastMessageAt:      baseTime,
	}
	return conv, nil
}