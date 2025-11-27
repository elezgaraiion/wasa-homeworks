package database

import (
	"sort"
	"time"
	"database/sql"
	"github.com/aritz/wasa-homeworks/service/models"
)

func (db *appdbimpl) GetMyConversations(userID string) ([]models.Conversation, error) {
	rows, err := db.c.Query(`
		SELECT c.id, c.type, c.name, c.photo,
		       c.last_message_preview, c.last_message_at,
		       m.joined_at
		FROM conversations c
		JOIN conversation_participants p ON p.conversation_id = c.id
		LEFT JOIN conversation_user_meta m
		  ON m.conversation_id = c.id AND m.user_id = ?
		WHERE p.user_id = ?`, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var convs []models.Conversation

	for rows.Next() {
		var conv models.Conversation

		var lastMsgAt sql.NullString
		var joinedAt sql.NullString
		var preview sql.NullString
		var name sql.NullString
		var photo sql.NullString

		if err := rows.Scan(
			&conv.ID,
			&conv.Type,
			&name,
			&photo,
			&preview,
			&lastMsgAt,
			&joinedAt,
		); err != nil {
			return nil, err
		}

		if name.Valid {
			conv.Name = name.String
		}
		if photo.Valid {
			conv.Photo = photo.String
		}
		if preview.Valid {
			conv.LastMessagePreview = preview.String
		}

		// Parse LastMessageAt si existe, sino default a zero
		conv.LastMessageAt = time.Time{}
		if lastMsgAt.Valid && lastMsgAt.String != "" {
			t, _ := time.Parse(time.RFC3339, lastMsgAt.String)
			conv.LastMessageAt = t
		}

		// Parse JoinedAt si existe
		joinedAtTime := time.Time{}
		if joinedAt.Valid && joinedAt.String != "" {
			t, _ := time.Parse(time.RFC3339, joinedAt.String)
			joinedAtTime = t
		}

		// Orden temporal: el mayor entre last_message_at y joined_at
		conv.TempOrderAt = conv.LastMessageAt
		if joinedAtTime.After(conv.LastMessageAt) {
			conv.TempOrderAt = joinedAtTime
		}

		// Cargar participantes
		conv.Participants, _ = db.getParticipantsByConversation(conv.ID)

		convs = append(convs, conv)
	}

	// Orden descendente por TempOrderAt
	sort.Slice(convs, func(i, j int) bool {
		return convs[i].TempOrderAt.After(convs[j].TempOrderAt)
	})

	return convs, nil
}

func (db *appdbimpl) getParticipantsByConversation(convID string) ([]models.User, error) {
	rows, err := db.c.Query(`
		SELECT u.id, u.name, u.photo 
		FROM users u
		JOIN conversation_participants cp ON cp.user_id = u.id
		WHERE cp.conversation_id = ?`, convID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []models.User

	for rows.Next() {
		var u models.User
		var name sql.NullString
		var photo sql.NullString

		if err := rows.Scan(&u.ID, &name, &photo); err != nil {
			return nil, err
		}

		if name.Valid {
			u.Name = name.String
		}
		if photo.Valid {
			u.Photo = photo.String
		}

		participants = append(participants, u)
	}

	return participants, nil
}
