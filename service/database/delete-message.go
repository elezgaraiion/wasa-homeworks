package database

import(
	"database/sql"
	"fmt"
	"github.com/aritz/wasa-homeworks/service/models"
    "errors"

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
    var msgCreatedAt string
    err = db.c.QueryRow(`
        SELECT sender_id, created_at
        FROM messages
        WHERE id = ? AND conversation_id = ?
    `, messageID, convID).Scan(&senderID, &msgCreatedAt)

    if err == sql.ErrNoRows {
        return models.ErrMessageNotFound
    }
    if err != nil {
        return err
    }

    if senderID != userID {
        return models.ErrForbidden
    }

    tx, err := db.c.Begin()
    if err != nil {
        return err
    }
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()

    _, err = tx.Exec(`DELETE FROM reactions WHERE message_id = ?`, messageID)
    if err != nil {
        return fmt.Errorf("error cleaning reactions: %w", err)
    }


    res, err := tx.Exec(`
        DELETE FROM messages
        WHERE id = ? AND conversation_id = ?
    `, messageID, convID)
    if err != nil {
        return fmt.Errorf("delete message sql error: %w", err)
    }

    affected, _ := res.RowsAffected()
    if affected == 0 {
        return models.ErrMessageNotFound
    }

    var newText sql.NullString
    var newPhoto sql.NullString
    var newCreatedAt string 
    err = tx.QueryRow(`
        SELECT text, photo, created_at
        FROM messages
        WHERE conversation_id = ?
        ORDER BY created_at DESC
        LIMIT 1
    `, convID).Scan(&newText, &newPhoto, &newCreatedAt)

    if errors.Is(err, sql.ErrNoRows) {
        _, err = tx.Exec(`
            UPDATE conversations 
            SET last_message_preview = '', last_message_at = ?
            WHERE id = ?
        `, msgCreatedAt, convID)
        if err != nil {
            return fmt.Errorf("update empty conv: %w", err)
        }
    } else if err != nil {
        return err 
    } else {
        preview := newText.String
        
        if preview == "" && newPhoto.String != "" {
            preview = "📷 Foto"
        }
        
        if len(preview) > 50 {
            preview = preview[:47] + "..."
        }

        _, err = tx.Exec(`
            UPDATE conversations 
            SET last_message_preview = ?, last_message_at = ?
            WHERE id = ?
        `, preview, newCreatedAt, convID)

        if err != nil {
            return fmt.Errorf("update conv preview: %w", err)
        }
    }

    err = tx.Commit()
    if err != nil {
        return err
    }

    return nil
}