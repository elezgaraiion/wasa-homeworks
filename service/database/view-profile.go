package database

import (
	"fmt"
	"database/sql"
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

	if err != nil {
		return conv, fmt.Errorf("conversation not found or access denied: %w", err)
	}

	if name.Valid { conv.Name = name.String }
	if photo.Valid { conv.Photo = photo.String }

	if conv.Type == "direct" {
		var otherName, otherPhoto sql.NullString
		err = db.c.QueryRow(`
			SELECT u.name, u.photo
			FROM conversation_participants cp
			JOIN users u ON u.id = cp.user_id
			WHERE cp.conversation_id = ? AND cp.user_id != ?
		`, convID, userID).Scan(&otherName, &otherPhoto)
		
		if err == nil {
			if otherName.Valid { conv.Name = otherName.String }
			if otherPhoto.Valid { conv.Photo = otherPhoto.String }
		}
	}

	parts, err := db.getParticipantsByConversation(conv.ID)
	if err == nil {
		conv.Participants = parts
	}

	return conv, nil
}
