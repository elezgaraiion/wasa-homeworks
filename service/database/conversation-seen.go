package database

import (
	"github.com/aritz/wasa-homeworks/service/models"
	"time"
)

func (db *appdbimpl) MarkConversationSeen(userID, convID string) error {
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

	nowStr := time.Now().UTC().Format(time.RFC3339Nano)

	_, err = db.c.Exec(`
        INSERT OR REPLACE INTO conversation_user_meta (conversation_id, user_id, last_seen_message_at, joined_at)
        VALUES (
            ?, 
            ?, 
            ?,
            COALESCE((SELECT joined_at FROM conversation_user_meta WHERE conversation_id=? AND user_id=?), ?)
        )
    `, convID, userID, nowStr, convID, userID, nowStr)

	if err != nil {
		return err
	}

	query := `
        UPDATE messages
        SET status = 'read'
        WHERE conversation_id = ? 
          AND status = 'delivered'
          AND NOT EXISTS (
              SELECT 1
              FROM conversation_participants cp
              LEFT JOIN conversation_user_meta cum 
                ON cum.conversation_id = cp.conversation_id AND cum.user_id = cp.user_id
              WHERE cp.conversation_id = messages.conversation_id
                AND cp.user_id != messages.sender_id 
                AND (
                    cum.last_seen_message_at IS NULL 
                    OR 
                    cum.last_seen_message_at < messages.created_at
                )
          )
    `

	_, err = db.c.Exec(query, convID)
	return err
}
