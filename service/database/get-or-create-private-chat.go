package database

import (
	"database/sql"
	"time"
	"github.com/google/uuid"
	"github.com/aritz/wasa-homeworks/service/models"
)

// GetOrCreateOneOnOneConversation busca un chat privado existente o crea uno nuevo
func (db *appdbimpl) GetOrCreateOneOnOneConversation(myID, targetUserID string) (models.Conversation, error) {

	// 1. BUSCAR SI YA EXISTE
	// Hacemos un JOIN doble para encontrar una conversación 'direct' que tenga a LOS DOS usuarios
	var existingConvID string
	err := db.c.QueryRow(`
		SELECT c.id
		FROM conversations c
		JOIN conversation_participants p1 ON p1.conversation_id = c.id
		JOIN conversation_participants p2 ON p2.conversation_id = c.id
		WHERE c.type = 'direct'
		  AND p1.user_id = ?
		  AND p2.user_id = ?
		LIMIT 1
	`, myID, targetUserID).Scan(&existingConvID)

	// CASO A: YA EXISTE
	if err == nil {
		// Recuperamos los participantes para devolver el objeto completo
		parts, _ := db.getParticipantsByConversation(existingConvID)
		
		// Buscamos info del "otro" para rellenar nombre/foto en la estructura
		// (Esto ayuda al frontend a pintarlo rápido)
		var targetName, targetPhoto string
		// Ignoramos error aquí por simplicidad, si falla saldrá vacío
		_ = db.c.QueryRow("SELECT name, photo FROM users WHERE id = ?", targetUserID).Scan(&targetName, &targetPhoto)

		return models.Conversation{
			ID:           existingConvID,
			Type:         "direct",
			Name:         targetName,   
			Photo:        targetPhoto,
			Participants: parts,
		}, nil
	} else if err != sql.ErrNoRows {
		// Si es un error real de la base de datos (no "no encontrado"), fallamos
		return models.Conversation{}, err
	}

	// CASO B: NO EXISTE -> CREAMOS UNO NUEVO
	
	// Iniciamos transacción (para que se guarde todo o nada)
	tx, err := db.c.Begin()
	if err != nil { return models.Conversation{}, err }
	defer tx.Rollback()

	newID := uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)

	// 1. Insertar Conversación
	_, err = tx.Exec(`INSERT INTO conversations (id, type) VALUES (?, 'direct')`, newID)
	if err != nil { return models.Conversation{}, err }

	// 2. Insertarme a mí (Participante + Meta)
	_, err = tx.Exec(`INSERT INTO conversation_participants (conversation_id, user_id) VALUES (?, ?)`, newID, myID)
	if err != nil { return models.Conversation{}, err }
	_, err = tx.Exec(`INSERT INTO conversation_user_meta (conversation_id, user_id, joined_at) VALUES (?, ?, ?)`, newID, myID, now)
	if err != nil { return models.Conversation{}, err }

	// 3. Insertar al otro (Participante + Meta)
	_, err = tx.Exec(`INSERT INTO conversation_participants (conversation_id, user_id) VALUES (?, ?)`, newID, targetUserID)
	if err != nil { return models.Conversation{}, err }
	_, err = tx.Exec(`INSERT INTO conversation_user_meta (conversation_id, user_id, joined_at) VALUES (?, ?, ?)`, newID, targetUserID, now)
	if err != nil { return models.Conversation{}, err }

	// Confirmar transacción
	if err := tx.Commit(); err != nil { return models.Conversation{}, err }

	// Devolver la nueva conversación
	parts, _ := db.getParticipantsByConversation(newID)
	
	// Datos del otro usuario
	var targetName, targetPhoto string
	_ = db.c.QueryRow("SELECT name, photo FROM users WHERE id = ?", targetUserID).Scan(&targetName, &targetPhoto)

	return models.Conversation{
		ID:           newID,
		Type:         "direct",
		Name:         targetName,
		Photo:        targetPhoto,
		Participants: parts,
	}, nil
}