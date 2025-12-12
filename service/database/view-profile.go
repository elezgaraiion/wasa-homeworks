package database

import (
	"fmt"
	"database/sql"
	"github.com/aritz/wasa-homeworks/service/models"
)

// Devuelve info de la conversación para mostrar en perfil
func (db *appdbimpl) GetConversationProfile(userID, convID string) (models.Conversation, error) {
	var conv models.Conversation
	
	// Variables auxiliares para manejar NULLs de la base de datos
	var name sql.NullString
	var photo sql.NullString

	// 1. Obtener datos básicos
	err := db.c.QueryRow(`
		SELECT c.id, c.type, c.name, c.photo
		FROM conversations c
		JOIN conversation_participants p ON p.conversation_id = c.id
		WHERE c.id = ? AND p.user_id = ?
	`, convID, userID).Scan(&conv.ID, &conv.Type, &name, &photo)

	if err != nil {
		return conv, fmt.Errorf("conversation not found or access denied: %w", err)
	}

	// Asignar si son válidos
	if name.Valid { conv.Name = name.String }
	if photo.Valid { conv.Photo = photo.String }

	// 2. Si es chat PRIVADO, el nombre/foto deben ser los del OTRO usuario
	if conv.Type == "direct" {
		var otherName, otherPhoto sql.NullString
		// Buscamos al participante que NO soy yo
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

	// 3. Cargar participantes (Importante para ver quién está en el grupo)
	// Lo hacemos siempre para que el modal muestre la lista de gente
	parts, err := db.getParticipantsByConversation(conv.ID)
	if err == nil {
		conv.Participants = parts
	}

	return conv, nil
}
