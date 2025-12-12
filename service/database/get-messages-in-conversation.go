package database


import (
	"database/sql"
	"github.com/aritz/wasa-homeworks/service/models"
	"time"
)

func (db *appdbimpl) GetMessagesInConversation(
	userID string,
	convID string,
	limit int,
	before string,
) ([]models.Message, error) {

	// 1. Validar acceso
	var count int
	err := db.c.QueryRow(`
		SELECT COUNT(*)
		FROM conversation_participants
		WHERE conversation_id = ? AND user_id = ?
	`, convID, userID).Scan(&count)

	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, models.ErrForbidden
	}

	// 2. Query real (AÑADIDO m.status)
	query := `
		SELECT 
			m.id,
			m.sender_id,
			u.name,
			u.photo,
			m.conversation_id,
			m.text,
			m.photo,
			m.reply_to_message_id,
			m.created_at,
			m.status  -- <--- ¡AQUÍ FALTABA ESTO!
		FROM messages m
		JOIN users u ON u.id = m.sender_id
		WHERE m.conversation_id = ?
	`
	args := []any{convID}

	if before != "" {
		query += " AND m.created_at < ?"
		args = append(args, before)
	}

	// Ordenamos por fecha DESC
	query += " ORDER BY m.created_at DESC LIMIT ?"
	args = append(args, limit)

	rows, err := db.c.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	msgs := []models.Message{}
	
	const layoutSQLite = "2006-01-02 15:04:05"

	for rows.Next() {
		var m models.Message
		
		// VARIABLES TEMPORALES
		var senderPhoto sql.NullString
		var txt sql.NullString
		var photo sql.NullString
		var replyTo sql.NullString
		var createdAtStr string 
		var status string // <--- VARIABLE PARA EL ESTADO

		err = rows.Scan(
			&m.ID,
			&m.Sender.ID,
			&m.Sender.Name,
			&senderPhoto,
			&m.ConversationID,
			&txt,       
			&photo,     
			&replyTo,   
			&createdAtStr,
			&status, // <--- ESCANEAMOS EL ESTADO
		)
		if err != nil {
			return nil, err
		}

		if senderPhoto.Valid { m.Sender.Photo = senderPhoto.String }
		if txt.Valid { m.Text = txt.String }
		if photo.Valid { m.Photo = photo.String }
		if replyTo.Valid { m.ReplyToMessageID = replyTo.String }
		
		// Asignamos el estado al struct
		m.Status = status 

		// PARSEAR FECHA
		t, err := time.Parse(layoutSQLite, createdAtStr)
		if err != nil {
			t, _ = time.Parse(time.RFC3339, createdAtStr)
		}
		m.CreatedAt = t

		// Reacciones (Opcional, si lo usas)
		// m.Reactions, _ = db.GetReactions(m.ID)

		msgs = append(msgs, m)
	}

	if msgs == nil {
		msgs = []models.Message{}
	}

	return msgs, nil
}
func (db *appdbimpl) GetReactions(messageID string) ([]models.Reaction, error) {

    rows, err := db.c.Query(`
        SELECT 
            r.id,
            r.user_id,
            u.name,
            u.photo,
            r.emoji,
            r.created_at
        FROM reactions r
        JOIN users u ON u.id = r.user_id
        WHERE r.message_id = ?
        ORDER BY r.created_at ASC
    `, messageID)

    if err != nil {
        return nil, err
    }
    defer rows.Close()

    reactions := []models.Reaction{}

    for rows.Next() {
        var r models.Reaction
        err = rows.Scan(
            &r.ID,
            &r.User.ID,
            &r.User.Name,
            &r.User.Photo,
            &r.Emoji,
            &r.CreatedAt,
        )
        if err != nil {
            return nil, err
        }
        reactions = append(reactions, r)
    }

    return reactions, nil
}

func (db *appdbimpl) applySeenStatus(
    convID string,
    requestingUser string,
    convType string,
    msgs []models.Message,
) {
    for i := range msgs {
        msg := &msgs[i]

        // mensajes que no son míos -> no calculo
        if msg.Sender.ID != requestingUser {
            continue
        }

        if convType == "private" {
            db.applyPrivateStatus(convID, msg)
        } else {
            db.applyGroupStatus(convID, msg)
        }
    }
}
func (db *appdbimpl) applyPrivateStatus(convID string, msg *models.Message) {

    var lastSeen time.Time
    err := db.c.QueryRow(`
        SELECT last_seen_message_at
        FROM conversation_user_meta
        WHERE conversation_id = ?
        AND user_id != ?
    `, convID, msg.Sender.ID).Scan(&lastSeen)

    if err == nil && !lastSeen.IsZero() && lastSeen.After(msg.CreatedAt) {
        msg.Status = "read"
    } else {
        msg.Status = "delivered"
    }
}
func (db *appdbimpl) applyGroupStatus(convID string, msg *models.Message) {

    var pending int
    err := db.c.QueryRow(`
        SELECT COUNT(*)
        FROM conversation_user_meta
        WHERE conversation_id = ?
        AND user_id != ?
        AND (last_seen_message_at IS NULL OR last_seen_message_at < ?)
    `, convID, msg.Sender.ID, msg.CreatedAt).Scan(&pending)

    if err == nil && pending == 0 {
        msg.Status = "read"
    } else {
        msg.Status = "delivered"
    }
}
