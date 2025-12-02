package database

import(
	"time"
	"github.com/google/uuid"
	"github.com/aritz/wasa-homeworks/service/models"
)

func (db *appdbimpl) AddReaction(userID, convID, msgID, emoji string) (models.Reaction, error) {
    // Validar existencia de conversación y mensaje
    var exists int
    err := db.c.QueryRow(`SELECT COUNT(*) FROM conversations WHERE id = ?`, convID).Scan(&exists)
    if err != nil || exists == 0 {
        return models.Reaction{}, models.ErrConversationNotFound
    }

    err = db.c.QueryRow(`SELECT COUNT(*) FROM messages WHERE id = ? AND conversation_id = ?`, msgID, convID).Scan(&exists)
    if err != nil || exists == 0 {
        return models.Reaction{}, models.ErrMessageNotFound
    }

    id := uuid.New().String()
    now := time.Now().UTC()
    _, err = db.c.Exec(`INSERT INTO reactions(id, user_id, message_id, emoji, created_at) VALUES (?, ?, ?, ?, ?)`,
        id, userID, msgID, emoji, now.Format(time.RFC3339))
    if err != nil {
        return models.Reaction{}, err
    }

    reaction := models.Reaction{
        ID:        id,
        User:      models.User{ID: userID},
        Emoji:     emoji,
        CreatedAt: now,
    }
    return reaction, nil
}
