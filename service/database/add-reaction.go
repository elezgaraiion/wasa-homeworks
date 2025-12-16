package database

import(
	"time"
	"github.com/google/uuid"
	"github.com/aritz/wasa-homeworks/service/models"
    "database/sql"
)

func (db *appdbimpl) AddReaction(userID, convID, msgID, emoji string) (models.Reaction, error) {
	// 1. Validar que existen conversación y mensaje
	var exists int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM conversations WHERE id = ?`, convID).Scan(&exists)
	if err != nil || exists == 0 {
		return models.Reaction{}, models.ErrConversationNotFound
	}

	err = db.c.QueryRow(`SELECT COUNT(*) FROM messages WHERE id = ? AND conversation_id = ?`, msgID, convID).Scan(&exists)
	if err != nil || exists == 0 {
		return models.Reaction{}, models.ErrMessageNotFound
	}

	// 2. Comprobar si ya existe reacción de este usuario
	var reactionID string
	err = db.c.QueryRow(`SELECT id FROM reactions WHERE user_id = ? AND message_id = ?`, userID, msgID).Scan(&reactionID)

	now := time.Now().UTC()
	
	if err == sql.ErrNoRows {
		// A) No existe -> INSERTAR
		reactionID = uuid.New().String()
		_, err = db.c.Exec(`
			INSERT INTO reactions(id, user_id, message_id, emoji, created_at)
			VALUES (?, ?, ?, ?, ?)
		`, reactionID, userID, msgID, emoji, now.Format(time.RFC3339))
		if err != nil {
			return models.Reaction{}, err
		}
	} else if err == nil {
		// B) Ya existe -> ACTUALIZAR (UPDATE)
		_, err = db.c.Exec(`
			UPDATE reactions SET emoji = ?, created_at = ?
			WHERE id = ?
		`, emoji, now.Format(time.RFC3339), reactionID)
		if err != nil {
			return models.Reaction{}, err
		}
	} else {
		// Error de base de datos
		return models.Reaction{}, err
	}

	// 3. Devolver la estructura completa
	var userName string
	db.c.QueryRow("SELECT name FROM users WHERE id = ?", userID).Scan(&userName)

	return models.Reaction{
		ID:        reactionID,
		User:      models.User{ID: userID, Name: userName},
		Emoji:     emoji,
		CreatedAt: now,
	}, nil
}