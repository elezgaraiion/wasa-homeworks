package database

import (
	"database/sql"
	"errors"

	"github.com/aritz/wasa-homeworks/service/models"
)

func (db *appdbimpl) GetConversationProfile(userID, convID string) (models.Conversation, error) {
	var conv models.Conversation

	var name sql.NullString
	var photo sql.NullString

	err := db.c.QueryRow(`
        SELECT c.id, c.type, c.name, c.photo
        FROM conversations c
        JOIN conversation_participants p ON p.conversation_id = c.id
        WHERE c.id = ? AND p.user_id = ?
    `, convID, userID).Scan(&conv.ID, &conv.Type, &name, &photo)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Conversation{}, models.ErrConversationNotFound
	}
	if err != nil {
		return models.Conversation{}, err
	}

	if name.Valid {
		conv.Name = name.String
	}
	if photo.Valid {
		conv.Photo = photo.String
	}

	if conv.Type == "direct" {
		var otherName, otherPhoto sql.NullString
		err = db.c.QueryRow(`
            SELECT u.name, u.photo
            FROM conversation_participants cp
            JOIN users u ON u.id = cp.user_id
            WHERE cp.conversation_id = ? AND cp.user_id != ?
        `, convID, userID).Scan(&otherName, &otherPhoto)

		if err == nil {
			if otherName.Valid {
				conv.Name = otherName.String
			}
			if otherPhoto.Valid {
				conv.Photo = otherPhoto.String
			}
		}
	}

	conv.Participants, err = db.getParticipantsByConversation(conv.ID)
	if err != nil {
		return models.Conversation{}, err
	}

	return conv, nil
}