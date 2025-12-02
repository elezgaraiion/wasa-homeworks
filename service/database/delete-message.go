package database

import(
	"database/sql"
	"fmt"
	"github.com/aritz/wasa-homeworks/service/models"


)

func (db *appdbimpl) DeleteMessage(userID, convID, messageID string) error {
	// 1️⃣ Validar pertenencia a la conversación
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

	// 2️⃣ Comprobar que el mensaje existe y pertenece al usuario
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

	// 3️⃣ Borrar el mensaje
	_, err = db.c.Exec(`
		DELETE FROM messages
		WHERE id = ? AND conversation_id = ?
	`, messageID, convID)
	if err != nil {
		return fmt.Errorf("delete message: %w", err)
	}

	return nil
}
