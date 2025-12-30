package database

import(
	"time"
	"github.com/google/uuid"
	"github.com/aritz/wasa-homeworks/service/models"
    "database/sql"
)

func (db *appdbimpl) AddReaction(userID, convID, msgID, emoji string) (models.Reaction, error) {
	var exists int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM conversations WHERE id = ?`, convID).Scan(&exists)
	if err != nil || exists == 0 {
		return models.Reaction{}, models.ErrConversationNotFound
	}

	err = db.c.QueryRow(`SELECT COUNT(*) FROM messages WHERE id = ? AND conversation_id = ?`, msgID, convID).Scan(&exists)
	if err != nil || exists == 0 {
		return models.Reaction{}, models.ErrMessageNotFound
	}

	var reactionID string
	err = db.c.QueryRow(`SELECT id FROM reactions WHERE user_id = ? AND message_id = ?`, userID, msgID).Scan(&reactionID)

	now := time.Now().UTC()
	
	if err == sql.ErrNoRows {
		reactionID = uuid.New().String()
		_, err = db.c.Exec(`
			INSERT INTO reactions(id, user_id, message_id, emoji, created_at)
			VALUES (?, ?, ?, ?, ?)
		`, reactionID, userID, msgID, emoji, now.Format(time.RFC3339))
		if err != nil {
			return models.Reaction{}, err
		}
	} else if err == nil {
		_, err = db.c.Exec(`
			UPDATE reactions SET emoji = ?, created_at = ?
			WHERE id = ?
		`, emoji, now.Format(time.RFC3339), reactionID)
		if err != nil {
			return models.Reaction{}, err
		}
	} else {
		return models.Reaction{}, err
	}

	var userName string
	db.c.QueryRow("SELECT name FROM users WHERE id = ?", userID).Scan(&userName)

	return models.Reaction{
		ID:        reactionID,
		User:      models.User{ID: userID, Name: userName},
		Emoji:     emoji,
		CreatedAt: now,
	}, nil
}