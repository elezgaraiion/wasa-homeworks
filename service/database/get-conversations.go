package database

import (
	"database/sql"
	"sort"
	"time"

	"github.com/aritz/wasa-homeworks/service/models"
)

func (db *appdbimpl) GetMyConversations(userID string) ([]models.Conversation, error) {
	query := `
        SELECT 
            c.id, c.type, c.name, c.photo,
            c.last_message_preview, c.last_message_at,
            meta.joined_at,
            COALESCE(msg.sender_id, ''),
            COALESCE(msg.status, ''),
            COALESCE(u_sender.name, ''),
            (SELECT COUNT(*) 
             FROM messages m2 
             WHERE m2.conversation_id = c.id 
               AND m2.created_at > COALESCE(meta.last_seen_message_at, meta.joined_at)
               AND m2.sender_id != ?
            ) AS unread_count
        FROM conversations c
        JOIN conversation_participants p ON p.conversation_id = c.id
        LEFT JOIN conversation_user_meta meta ON meta.conversation_id = c.id AND meta.user_id = ?
        LEFT JOIN messages msg ON msg.conversation_id = c.id AND msg.created_at = c.last_message_at
        LEFT JOIN users u_sender ON u_sender.id = msg.sender_id
        WHERE p.user_id = ?
    `

	rows, err := db.c.Query(query, userID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var convs []models.Conversation

	const layoutSQLite = "2006-01-02 15:04:05"

	for rows.Next() {
		var c models.Conversation

		var lastMsgAtStr sql.NullString
		var joinedAtStr sql.NullString
		var preview sql.NullString
		var name sql.NullString
		var photo sql.NullString
		var senderID sql.NullString
		var msgStatus sql.NullString
		var senderName sql.NullString
		var unreadCount int

		err = rows.Scan(
			&c.ID, &c.Type, &name, &photo, &preview, &lastMsgAtStr, &joinedAtStr,
			&senderID, &msgStatus, &senderName, &unreadCount,
		)
		if err != nil {
			return nil, err
		}

		c.UnreadCount = unreadCount
		if name.Valid {
			c.Name = name.String
		}
		if photo.Valid {
			c.Photo = photo.String
		}
		if preview.Valid {
			c.LastMessagePreview = preview.String
		}
		if senderID.Valid {
			c.LastMessageSenderID = senderID.String
		}
		if msgStatus.Valid {
			c.LastMessageStatus = msgStatus.String
		}
		if senderName.Valid {
			c.LastMessageSenderName = senderName.String
		}

		var tMsg time.Time
		if lastMsgAtStr.Valid && lastMsgAtStr.String != "" {
			t, err := time.Parse(time.RFC3339, lastMsgAtStr.String)
			if err != nil {
				t, _ = time.Parse(layoutSQLite, lastMsgAtStr.String)
			}
			tMsg = t
			c.LastMessageAt = t
		}

		var tJoined time.Time
		if joinedAtStr.Valid && joinedAtStr.String != "" {
			t, err := time.Parse(time.RFC3339, joinedAtStr.String)
			if err != nil {
				t, _ = time.Parse(layoutSQLite, joinedAtStr.String)
			}
			tJoined = t
		}

		if tMsg.After(tJoined) {
			c.TempOrderAt = tMsg
		} else {
			c.TempOrderAt = tJoined
		}

		c.Participants, err = db.getParticipantsByConversation(c.ID)
		if err != nil {
			return nil, err
		}

		if c.Type == "direct" {
			for _, p := range c.Participants {
				if p.ID != userID {
					c.Name = p.Name
					c.Photo = p.Photo
					break
				}
			}
		}

		convs = append(convs, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

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
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return participants, nil
}
