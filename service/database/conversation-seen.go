package database

import (
    "time"
	"errors"
)

func (db *appdbimpl) MarkConversationSeen(userID, convID string) error {
	// 1. Validar acceso
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

	// 2. Definir fecha (Ahora + 5 seg) para asegurar que supera al mensaje
	nowStr := time.Now().UTC().Add(5 * time.Second).Format(time.RFC3339)

	// 3. ACTUALIZAR METADATOS (Registrar que YO he visto el chat)
	res, err := db.c.Exec(`
		UPDATE conversation_user_meta 
		SET last_seen_message_at = ?
		WHERE conversation_id = ? AND user_id = ?
	`, nowStr, convID, userID)
	if err != nil { return err }

	rows, _ := res.RowsAffected()
	if rows == 0 {
		// Si no existía fila, la creamos
		db.c.Exec(`INSERT INTO conversation_user_meta (conversation_id, user_id, last_seen_message_at, joined_at) VALUES (?, ?, ?, ?)`, convID, userID, nowStr, nowStr)
	}

	// 4. LÓGICA DE GRUPO (TICKS AZULES)
	// Esta Query hace lo siguiente:
	// Actualiza a 'read' los mensajes de este chat QUE ESTÉN EN 'delivered'
	// SIEMPRE Y CUANDO NO EXISTA ningún participante (excluyendo al sender)
	// que tenga un 'last_seen_message_at' anterior a la fecha del mensaje.
	
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
				AND cp.user_id != messages.sender_id -- No contamos al que envió el mensaje
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