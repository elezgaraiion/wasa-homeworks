package database

import(
	"fmt"
	"time"
	"github.com/google/uuid"
	"github.com/aritz/wasa-homeworks/service/models"

)

func (db *appdbimpl) ForwardMessage(
	userID, sourceConvID, messageID, targetConvID string,
) (models.Message, error) {

	// 1️⃣ Validar pertenencia del usuario a la conversación origen
	var count int
	err := db.c.QueryRow(`
		SELECT COUNT(*) 
		FROM conversation_participants
		WHERE conversation_id = ? AND user_id = ?
	`, sourceConvID, userID).Scan(&count)
	if err != nil {
		return models.Message{}, err
	}
	if count == 0 {
		return models.Message{}, models.ErrForbidden
	}

	// 2️⃣ Validar existencia de mensaje
	var originalMessage models.Message
	originalMessage, err = db.GetMessageByID(userID, sourceConvID, messageID)
	if err != nil {
		return models.Message{}, err
	}

	// 3️⃣ Validar pertenencia a la conversación destino
	err = db.c.QueryRow(`
		SELECT COUNT(*) 
		FROM conversation_participants
		WHERE conversation_id = ? AND user_id = ?
	`, targetConvID, userID).Scan(&count)
	if err != nil {
		return models.Message{}, err
	}
	if count == 0 {
		return models.Message{}, models.ErrForbidden
	}

	// 4️⃣ Insertar mensaje en la conversación destino
	newMessageID := uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)

	_, err = db.c.Exec(`
		INSERT INTO messages(
			id, sender_id, conversation_id, text, photo, reply_to_message_id, created_at, status
		) VALUES (?, ?, ?, ?, ?, ?, ?, 'delivered')
	`, newMessageID, userID, targetConvID, originalMessage.Text, originalMessage.Photo,
		originalMessage.ReplyToMessageID, now)
	if err != nil {
		return models.Message{}, fmt.Errorf("insert forwarded message: %w", err)
	}

	// 5️⃣ Devolver el mensaje insertado
	return db.GetMessageByID(userID, targetConvID, newMessageID)
}
