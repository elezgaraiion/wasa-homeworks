package database

import(
	"database/sql"
	"fmt"
	"github.com/aritz/wasa-homeworks/service/models"
	"time"


)

func (db *appdbimpl) DeleteMessage(userID, convID, messageID string) error {
    var count int
    err := db.c.QueryRow(`
        SELECT COUNT(*) 
        FROM conversation_participants
        WHERE conversation_id = ? AND user_id = ?
    `, convID, userID).Scan(&count)
    if err != nil {
        return err
    }
    if count == 0 {
        return models.ErrForbidden
    }

    var senderID string
    err = db.c.QueryRow(`
        SELECT sender_id
        FROM messages
        WHERE id = ? AND conversation_id = ?
    `, messageID, convID).Scan(&senderID)

    if err == sql.ErrNoRows {
        return models.ErrMessageNotFound
    }
    if err != nil {
        return err
    }

    if senderID != userID {
        return models.ErrForbidden
    }

    _, err = db.c.Exec(`
        DELETE FROM messages
        WHERE id = ? AND conversation_id = ?
    `, messageID, convID)
    if err != nil {
        return fmt.Errorf("delete message: %w", err)
    }


    var newText string 
    var newPhoto string 
    var newCreatedAt time.Time

    err = db.c.QueryRow(`
        SELECT text, photo, created_at
        FROM messages
        WHERE conversation_id = ?
        ORDER BY created_at DESC
        LIMIT 1
    `, convID).Scan(&newText, &newPhoto, &newCreatedAt)

    if err == sql.ErrNoRows {
        _, err = db.c.Exec(`
            UPDATE conversations 
            SET last_message_preview = '', last_message_at = NULL
            WHERE id = ?
        `, convID)
        if err != nil {
            return fmt.Errorf("update empty conv: %w", err)
        }
    } else if err != nil {
        return err
    } else {
        
        preview := newText
        if preview == "" && newPhoto != "" {
            preview = "📷 Photo" 
        }

        _, err = db.c.Exec(`
            UPDATE conversations 
            SET last_message_preview = ?, last_message_at = ?
            WHERE id = ?
        `, preview, newCreatedAt, convID)
        
        if err != nil {
            return fmt.Errorf("update conv preview: %w", err)
        }
    }

    return nil
}