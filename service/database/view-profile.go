package database

import (
	"fmt"

	"github.com/aritz/wasa-homeworks/service/models"
)

// Devuelve info de la conversación para mostrar en perfil
func (db *appdbimpl) GetConversationProfile(userID, convID string) (models.Conversation, error) {
	// Verificar que el usuario pertenece a la conversación
	var conv models.Conversation
	err := db.c.QueryRow(`
		SELECT id, type, name, photo
		FROM conversations c
		JOIN conversation_participants p ON p.conversation_id = c.id
		WHERE c.id = ? AND p.user_id = ?
	`, convID, userID).Scan(&conv.ID, &conv.Type, &conv.Name, &conv.Photo)
	if err != nil {
		return conv, fmt.Errorf("conversation not found or user not participant: %w", err)
	}

	// Solo para grupos, traer participantes
	if conv.Type == "group" {
		parts, err := db.getParticipantsByConversation(conv.ID)
		if err != nil {
			return conv, fmt.Errorf("load participants: %w", err)
		}
		conv.Participants = parts
	}

	return conv, nil
}
