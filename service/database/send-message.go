package database

import (
    "database/sql"
    "fmt"
    "time"

    "github.com/aritz/wasa-homeworks/service/models"
    "github.com/google/uuid"
)

func (db *appdbimpl) SendMessage(userID string, conversationID string, text string, photoURL string, replyToMessageID string) (models.Message, error) {

    var count int
    err := db.c.QueryRow(`SELECT COUNT(*) FROM conversation_participants WHERE conversation_id = ? AND user_id = ?`, conversationID, userID).Scan(&count)
    if err != nil {
        return models.Message{}, err
    }
    if count == 0 {
        return models.Message{}, models.ErrForbidden
    }

    var snapText, snapPhoto sql.NullString
    var snapSender string 
    var sqlReplyID sql.NullString

    if replyToMessageID != "" {
        sqlReplyID.String = replyToMessageID
        sqlReplyID.Valid = true

        err := db.c.QueryRow(`
            SELECT m.text, m.photo, u.name 
            FROM messages m
            JOIN users u ON m.sender_id = u.id
            WHERE m.id = ?
        `, replyToMessageID).Scan(&snapText, &snapPhoto, &snapSender)

        if err != nil {
            // Si falla (ej: ID no existe), anulamos la respuesta
            replyToMessageID = ""
            sqlReplyID.Valid = false
        }
    }

    msgID := uuid.New().String()
    now := time.Now().UTC()
    nowStr := now.Format(time.RFC3339)

    var sqlPhoto sql.NullString
    if photoURL != "" {
        sqlPhoto.String = photoURL
        sqlPhoto.Valid = true
    }

    tx, err := db.c.Begin()
    if err != nil {
        return models.Message{}, err
    }
    defer tx.Rollback()

    _, err = tx.Exec(`
        INSERT INTO messages (
            id, conversation_id, sender_id, text, photo, 
            created_at, status, reply_to_message_id,
            reply_snapshot_text, reply_snapshot_sender, reply_snapshot_photo
        )
        VALUES (?, ?, ?, ?, ?, ?, 'delivered', ?, ?, ?, ?)
    `, msgID, conversationID, userID, text, sqlPhoto, nowStr, sqlReplyID, snapText, snapSender, snapPhoto)

    if err != nil {
        return models.Message{}, fmt.Errorf("insert msg: %w", err)
    }

    preview := text
    if preview == "" && photoURL != "" {
        preview = "📷 Foto"
    }
    _, err = tx.Exec(`
        UPDATE conversations 
        SET last_message_preview = ?, last_message_at = ?
        WHERE id = ?
    `, preview, nowStr, conversationID)
    if err != nil {
        return models.Message{}, fmt.Errorf("update conv: %w", err)
    }

    _, err = tx.Exec(`
        INSERT INTO conversation_user_meta (conversation_id, user_id, last_seen_message_at, joined_at)
        VALUES (?, ?, ?, ?)
        ON CONFLICT(conversation_id, user_id) 
        DO UPDATE SET last_seen_message_at = excluded.last_seen_message_at
    `, conversationID, userID, nowStr, nowStr)
    if err != nil {
        return models.Message{}, fmt.Errorf("update meta: %w", err)
    }

    if err := tx.Commit(); err != nil {
        return models.Message{}, err
    }

    var myName string
    db.c.QueryRow("SELECT name FROM users WHERE id = ?", userID).Scan(&myName)

    msg := models.Message{
        ID:               msgID,
        ConversationID:   conversationID,
        Sender:           models.User{ID: userID, Name: myName},
        Text:             text,
        Photo:            photoURL,
        CreatedAt:        now,
        Status:           "delivered",
        ReplyToMessageID: replyToMessageID,
    }

    if replyToMessageID != "" {
        msg.ReplyTo = &models.Message{
            ID:    replyToMessageID,
            Text:  snapText.String,   
            Photo: snapPhoto.String,  
            Sender: models.User{
                Name: snapSender,
            },
        }
    }

    return msg, nil
}