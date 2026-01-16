package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/aritz/wasa-homeworks/service/models"
)

func (db *appdbimpl) GetMessagesInConversation(
	userID string,
	convID string,
	limit int,
	before string,
) ([]models.Message, error) {

	var joinedAt string
	err := db.c.QueryRow(`
        SELECT joined_at
        FROM conversation_user_meta
        WHERE conversation_id = ? AND user_id = ?
    `, convID, userID).Scan(&joinedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrForbidden
	}
	if err != nil {
		return nil, err
	}

	query := `
        SELECT 
            m.id, m.sender_id, u.name, u.photo, m.conversation_id,
            m.text, m.photo, m.reply_to_message_id, m.created_at, m.status,
            r.text, r.photo, ru.name
        FROM messages m
        JOIN users u ON u.id = m.sender_id
        LEFT JOIN messages r ON m.reply_to_message_id = r.id
        LEFT JOIN users ru ON r.sender_id = ru.id
        WHERE m.conversation_id = ?
          AND m.created_at >= ?
    `
	args := []any{convID, joinedAt}

	if before != "" {
		query += " AND m.created_at < ?"
		args = append(args, before)
	}

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
		var senderPhoto, txt, photo, replyToID sql.NullString
		var createdAtStr, status string
		var rText, rPhoto, rSenderName sql.NullString

		err = rows.Scan(
			&m.ID, &m.Sender.ID, &m.Sender.Name, &senderPhoto, &m.ConversationID,
			&txt, &photo, &replyToID, &createdAtStr, &status,
			&rText, &rPhoto, &rSenderName,
		)
		if err != nil {
			return nil, err
		}

		if senderPhoto.Valid {
			m.Sender.Photo = senderPhoto.String
		}
		if txt.Valid {
			m.Text = txt.String
		}
		if photo.Valid {
			m.Photo = photo.String
		}
		m.Status = status

		t, err := time.Parse(layoutSQLite, createdAtStr)
		if err != nil {
			t, _ = time.Parse(time.RFC3339, createdAtStr)
		}
		m.CreatedAt = t

		if replyToID.Valid && replyToID.String != "" {
			m.ReplyToMessageID = replyToID.String
			m.ReplyTo = &models.Message{
				ID:     replyToID.String,
				Text:   rText.String,
				Photo:  rPhoto.String,
				Sender: models.User{Name: rSenderName.String},
			}
		}

		reactions, err := db.GetReactions(m.ID)
		if err != nil {
			return nil, err
		}

		if reactions == nil {
			m.Reactions = []models.Reaction{}
		} else {
			m.Reactions = reactions
		}

		msgs = append(msgs, m)
	}
	if err = rows.Err(); err != nil {
		return nil, err
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

	var reactions []models.Reaction

	for rows.Next() {
		var r models.Reaction
		var createdAtStr string

		var userPhoto sql.NullString

		err = rows.Scan(
			&r.ID,
			&r.User.ID,
			&r.User.Name,
			&userPhoto,
			&r.Emoji,
			&createdAtStr,
		)
		if err != nil {
			return nil, err
		}

		if userPhoto.Valid {
			r.User.Photo = userPhoto.String
		}

		t, err := time.Parse(time.RFC3339, createdAtStr)
		if err == nil {
			r.CreatedAt = t
		}

		reactions = append(reactions, r)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	if reactions == nil {
		reactions = []models.Reaction{}
	}

	return reactions, nil
}
