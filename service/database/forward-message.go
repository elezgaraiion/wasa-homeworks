package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/aritz/wasa-homeworks/service/models"
	"github.com/gofrs/uuid"
)

func (db *appdbimpl) ForwardMessage(userID, sourceConvID, messageID, targetConvID string) (models.Message, error) {
	var count int
	err := db.c.QueryRow(`
        SELECT COUNT(*) FROM conversation_participants
        WHERE conversation_id = ? AND user_id = ?
    `, sourceConvID, userID).Scan(&count)
	if err != nil {
		return models.Message{}, err
	}
	if count == 0 {
		return models.Message{}, models.ErrForbidden
	}

	err = db.c.QueryRow(`
        SELECT COUNT(*) FROM conversation_participants
        WHERE conversation_id = ? AND user_id = ?
    `, targetConvID, userID).Scan(&count)
	if err != nil {
		return models.Message{}, err
	}
	if count == 0 {
		return models.Message{}, models.ErrForbidden
	}

	var text sql.NullString
	var photo sql.NullString
	err = db.c.QueryRow(`
        SELECT text, photo 
        FROM messages 
        WHERE id = ? AND conversation_id = ?
    `, messageID, sourceConvID).Scan(&text, &photo)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Message{}, models.ErrMessageNotFound
	} else if err != nil {
		return models.Message{}, err
	}

	tx, err := db.c.Begin()
	if err != nil {
		return models.Message{}, err
	}
	defer func() { _ = tx.Rollback() }()

	newUUID, err := uuid.NewV4()
	if err != nil {
		return models.Message{}, err
	}
	newMessageID := newUUID.String()

	now := time.Now().UTC()
	nowStr := now.Format(time.RFC3339Nano)

	_, err = tx.Exec(`
        INSERT INTO messages(id, sender_id, conversation_id, text, photo, created_at, status)
        VALUES (?, ?, ?, ?, ?, ?, 'delivered')
    `, newMessageID, userID, targetConvID, text, photo, nowStr)
	if err != nil {
		return models.Message{}, fmt.Errorf("insert msg: %w", err)
	}

	preview := ""
	if text.Valid {
		preview = text.String
	}
	if preview == "" && photo.Valid && photo.String != "" {
		preview = "📷 Foto"
	}

	_, err = tx.Exec(`
        UPDATE conversations 
        SET last_message_preview = ?, last_message_at = ?
        WHERE id = ?
    `, preview, nowStr, targetConvID)
	if err != nil {
		return models.Message{}, fmt.Errorf("update conv: %w", err)
	}

	res, err := tx.Exec(`
        UPDATE conversation_user_meta 
        SET last_seen_message_at = ?
        WHERE conversation_id = ? AND user_id = ?
    `, nowStr, targetConvID, userID)
	if err != nil {
		return models.Message{}, err
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {
		return models.Message{}, err
	}
	if rowsAff == 0 {
		_, err = tx.Exec(`
            INSERT INTO conversation_user_meta (conversation_id, user_id, last_seen_message_at, joined_at)
            VALUES (?, ?, ?, ?)
        `, targetConvID, userID, nowStr, nowStr)
		if err != nil {
			return models.Message{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return models.Message{}, err
	}

	var myName string
	err = db.c.QueryRow("SELECT name FROM users WHERE id = ?", userID).Scan(&myName)
	if err != nil {
		return models.Message{}, err
	}

	msg := models.Message{
		ID:             newMessageID,
		ConversationID: targetConvID,
		Sender:         models.User{ID: userID, Name: myName},
		Text:           text.String,
		Photo:          photo.String,
		CreatedAt:      now,
		Status:         "delivered",
	}

	return msg, nil
}