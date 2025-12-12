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

	// 1. Validar pertenencia (Lectura rápida antes de la transacción)
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

	// 2. INICIAR TRANSACCIÓN (Para hacer todo o nada)
	tx, err := db.c.Begin()
	if err != nil { return models.Message{}, err }
	defer tx.Rollback() // Si algo falla, deshacemos todo

	// Validar existencia conversación
	var convType string
	err = tx.QueryRow(`SELECT type FROM conversations WHERE id = ?`, convID).Scan(&convType)
	if err == sql.ErrNoRows {
		return models.Message{}, models.ErrConversationNotFound
	} else if err != nil {
		return models.Message{}, err
	}

	// 3. Preparar datos
	msgID := uuid.New().String()
	createdAt := time.Now().UTC()
	createdAtStr := createdAt.Format(time.RFC3339)

	// Calcular el preview (Si es foto sin texto, ponemos icono)
	preview := text
	if text == "" && photoURL != "" {
		preview = "📷 Foto"
	}
    // Truncar preview si es muy largo (opcional, por seguridad)
    if len(preview) > 50 {
        preview = preview[:47] + "..."
    }

	// 4. INSERTAR EL MENSAJE
	_, err = tx.Exec(`
		INSERT INTO messages(id, sender_id, conversation_id, text, photo, reply_to_message_id, created_at, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, msgID, senderID, convID, text, photoURL, replyToMessageID, createdAtStr, "delivered")
	if err != nil {
		return models.Message{}, err
	}

	// 5. ACTUALIZAR LA CONVERSACIÓN (¡LA CLAVE PARA EL PREVIEW!)
	// Actualizamos la fecha y el resumen para que salga arriba en la lista
	_, err = tx.Exec(`
		UPDATE conversations 
		SET last_message_preview = ?, last_message_at = ?
		WHERE id = ?
	`, preview, createdAtStr, convID)
	if err != nil {
		return models.Message{}, err
	}

    // 6. ACTUALIZAR MI "VISTO" (Para que no me salga bola verde a mí mismo)
    // Upsert: Si no existe la fila meta, la crea. Si existe, actualiza fecha.
    _, err = tx.Exec(`
        INSERT INTO conversation_user_meta (conversation_id, user_id, joined_at, last_seen_message_at)
        VALUES (?, ?, ?, ?)
        ON CONFLICT(conversation_id, user_id) 
        DO UPDATE SET last_seen_message_at = excluded.last_seen_message_at
    `, convID, senderID, createdAtStr, createdAtStr)
    if err != nil {
        return models.Message{}, err
    }

	// CONFIRMAR TRANSACCIÓN
	if err := tx.Commit(); err != nil {
		return models.Message{}, err
	}

	// 7. Devolver mensaje completo (para que el frontend lo pinte)
    // Recuperamos nombre del sender para que quede bonito si el frontend lo necesita
    var senderName string
    db.c.QueryRow("SELECT name FROM users WHERE id = ?", senderID).Scan(&senderName)

	msg := models.Message{
		ID:               msgID,
		Sender:           models.User{ID: senderID, Name: senderName}, 
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