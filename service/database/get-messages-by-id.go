package database

import(
	"database/sql"
	"github.com/aritz/wasa-homeworks/service/models"

)
func (db *appdbimpl) GetMessageByID(userID, convID, messageID string) (models.Message, error) {
	// 1️⃣ Validar pertenencia del usuario
	var count int
	err := db.c.QueryRow(`
		SELECT COUNT(*) 
		FROM conversation_participants
		WHERE conversation_id = ? AND user_id = ?
	`, convID, userID).Scan(&count)
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

	// 3️⃣ Consultar mensaje
	var m models.Message
	var senderID string
	err = db.c.QueryRow(`
		SELECT m.id, m.sender_id, u.name, u.photo, m.conversation_id,
		       m.text, m.photo, m.reply_to_message_id, m.created_at, m.status
		FROM messages m
		JOIN users u ON u.id = m.sender_id
		WHERE m.conversation_id = ? AND m.id = ?
	`, convID, messageID).Scan(
		&m.ID,
		&senderID,
		&m.Sender.Name,
		&m.Sender.Photo,
		&m.ConversationID,
		&m.Text,
		&m.Photo,
		&m.ReplyToMessageID,
		&m.CreatedAt,
		&m.Status,
	)
	if err == sql.ErrNoRows {
		return models.Message{}, models.ErrMessageNotFound
	}
	if err != nil {
		return models.Message{}, err
	}

	m.Sender.ID = senderID

	// 4️⃣ Reacciones
	m.Reactions, _ = db.GetReactions(m.ID)

	// 5️⃣ Aplicar estado leído / entregado si es mío
	if m.Sender.ID == userID {
		if convType == "private" {
			db.applyPrivateStatus(convID, &m)
		} else {
			db.applyGroupStatus(convID, &m)
		}
	}

	return m, nil
}
