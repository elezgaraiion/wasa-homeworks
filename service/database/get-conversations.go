package database

import (
	"sort"
	"time"

	"github.com/aritz/wasa-homeworks/service/models"
)

func (db *appdbimpl) GetMyConversations(userID string) ([]models.Conversation, error) {
	// 1. Traer todas las conversaciones
	rows, err := db.c.Query("SELECT id, type, name, photo, last_message_preview, last_message_at FROM conversations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allConvs []models.Conversation
	for rows.Next() {
		var conv models.Conversation
		var lastMsgAt string
		err := rows.Scan(&conv.ID, &conv.Type, &conv.Name, &conv.Photo, &conv.LastMessagePreview, &lastMsgAt)
		if err != nil {
			return nil, err
		}
		if lastMsgAt != "" {
			conv.LastMessageAt, _ = time.Parse(time.RFC3339, lastMsgAt)
		}
		// TODO: Traer participantes de alguna forma (si los tienes en JSON o tabla aparte)
		conv.Participants, _ = db.getParticipantsByConversation(conv.ID)

		allConvs = append(allConvs, conv)
	}

	// 2. Filtrar solo conversaciones donde participa el user
	var myConvs []models.Conversation
	for _, conv := range allConvs {
		for _, user := range conv.Participants {
			if user.ID == userID {
				myConvs = append(myConvs, conv)
				break
			}
		}
	}

	// 3. Ordenar por lastMessageAt descendente
	sort.Slice(myConvs, func(i, j int) bool {
		return myConvs[i].LastMessageAt.After(myConvs[j].LastMessageAt)
	})

	return myConvs, nil
}

// Ejemplo de función para obtener participantes
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
		if err := rows.Scan(&u.ID, &u.Name, &u.Photo); err != nil {
			return nil, err
		}
		participants = append(participants, u)
	}

	return participants, nil
}
