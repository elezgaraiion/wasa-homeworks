package database

import (
    "time"
	"errors"
)

func (db *appdbimpl) MarkConversationSeen(userID, convID string) error {
    // 1️⃣ Validar que pertenece a la conversación
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
        return errors.New("forbidden")
    }

    // 2️⃣ Guardar última lectura
    now := time.Now().UTC()

    _, err = db.c.Exec(`
        INSERT INTO conversation_user_meta (conversation_id, user_id, last_seen_message_at)
        VALUES (?, ?, ?)
        ON CONFLICT(conversation_id, user_id)
        DO UPDATE SET last_seen_message_at = excluded.last_seen_message_at
    `, convID, userID, now)

    return err
}
