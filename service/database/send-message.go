package database

import(
	"github.com/aritz/wasa-homeworks/service/models"
	"time"
	"database/sql"
	"github.com/google/uuid"


)
func (db *appdbimpl) SendMessage(
	senderID, convID, text, photoURL, replyToMessageID string,
) (models.Message, error) {

	// 1️⃣ Validar pertenencia
	var count int
	err := db.c.QueryRow(`
		SELECT COUNT(*) 
		FROM conversation_participants
		WHERE conversation_id = ? AND user_id = ?
	`, convID, senderID).Scan(&count)
	if err != nil {
		return models.Message{}, err
	}
	if count == 0 {
		return models.Message{}, models.ErrForbidden
	}

	// 2️⃣ Validar existencia de la conversación
	var convType string
	err = db.c.QueryRow(`SELECT type FROM conversations WHERE id = ?`, convID).Scan(&convType)
	if err == sql.ErrNoRows {
		return models.Message{}, models.ErrConversationNotFound
	}
	if err != nil {
		return models.Message{}, err
	}

	// 3️⃣ Insertar mensaje
	msgID := uuid.New().String()
	createdAt := time.Now().UTC()

	_, err = db.c.Exec(`
		INSERT INTO messages(id, sender_id, conversation_id, text, photo, reply_to_message_id, created_at, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, msgID, senderID, convID, text, photoURL, replyToMessageID, createdAt.Format(time.RFC3339), "delivered")
	if err != nil {
		return models.Message{}, err
	}

	// 4️⃣ Devolver mensaje completo
	msg := models.Message{
		ID:               msgID,
		Sender:           models.User{ID: senderID}, // puedes cargar nombre/foto si quieres
		ConversationID:   convID,
		Text:             text,
		Photo:            photoURL,
		ReplyToMessageID: replyToMessageID,
		CreatedAt:        createdAt,
		Status:           "delivered",
		Reactions:        []models.Reaction{},
	}

	return msg, nil
}