package database

import (
	"sort"
	"time"
	"database/sql"
	"github.com/aritz/wasa-homeworks/service/models"
)

func (db *appdbimpl) GetMyConversations(userID string) ([]models.Conversation, error) {
    // MODIFICAMOS LA QUERY:
    // Hacemos JOIN con 'messages' (m) coincidiendo con la fecha del último mensaje
    // Hacemos JOIN con 'users' (u_sender) para saber el nombre del que lo envió
    query := `
        SELECT 
            c.id, c.type, c.name, c.photo,
            c.last_message_preview, c.last_message_at,
            meta.joined_at,
            -- Extraemos datos del mensaje real
            COALESCE(msg.sender_id, ''),
            COALESCE(msg.status, ''),
            COALESCE(u_sender.name, '')
        FROM conversations c
        JOIN conversation_participants p ON p.conversation_id = c.id
        LEFT JOIN conversation_user_meta meta ON meta.conversation_id = c.id AND meta.user_id = ?
        
        -- BUSCAMOS EL ÚLTIMO MENSAJE REAL
        LEFT JOIN messages msg ON msg.conversation_id = c.id AND msg.created_at = c.last_message_at
        
        -- BUSCAMOS EL NOMBRE DEL QUE LO ENVIÓ
        LEFT JOIN users u_sender ON u_sender.id = msg.sender_id
        
        WHERE p.user_id = ?
    `

    rows, err := db.c.Query(query, userID, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var convs []models.Conversation
    const layoutSQLite = "2006-01-02 15:04:05"

    for rows.Next() {
        var conv models.Conversation

        // Variables temporales
        var lastMsgAt sql.NullString
        var joinedAt sql.NullString
        var preview sql.NullString
        var name sql.NullString
        var photo sql.NullString
        
        // Variables nuevas
        var senderID sql.NullString
        var msgStatus sql.NullString
        var senderName sql.NullString

        if err := rows.Scan(
            &conv.ID,
            &conv.Type,
            &name,
            &photo,
            &preview,
            &lastMsgAt,
            &joinedAt,
            &senderID,     // Nuevo
            &msgStatus,    // Nuevo
            &senderName,   // Nuevo
        ); err != nil {
            return nil, err
        }

        if name.Valid { conv.Name = name.String }
        if photo.Valid { conv.Photo = photo.String }
        if preview.Valid { conv.LastMessagePreview = preview.String }
        
        // Asignar nuevos campos al struct
        if senderID.Valid { conv.LastMessageSenderID = senderID.String }
        if msgStatus.Valid { conv.LastMessageStatus = msgStatus.String }
        if senderName.Valid { conv.LastMessageSenderName = senderName.String }

        // Fechas
        conv.LastMessageAt = time.Time{}
        if lastMsgAt.Valid && lastMsgAt.String != "" {
            t, err := time.Parse(layoutSQLite, lastMsgAt.String)
            if err != nil { t, _ = time.Parse(time.RFC3339, lastMsgAt.String) }
            conv.LastMessageAt = t
        }

        joinedAtTime := time.Time{}
        if joinedAt.Valid && joinedAt.String != "" {
            t, err := time.Parse(layoutSQLite, joinedAt.String)
            if err != nil { t, _ = time.Parse(time.RFC3339, joinedAt.String) }
            joinedAtTime = t
        }

        // Lógica de ordenación (Mantenida igual)
        conv.TempOrderAt = conv.LastMessageAt
        if joinedAtTime.After(conv.LastMessageAt) {
            conv.TempOrderAt = joinedAtTime
        }

        // 2. CARGAR PARTICIPANTES (Tu función auxiliar se mantiene igual)
        conv.Participants, _ = db.getParticipantsByConversation(conv.ID)

        // 3. LÓGICA DE NOMBRE (Directo vs Grupo)
        if conv.Type == "direct" {
            for _, p := range conv.Participants {
                if p.ID != userID {
                    conv.Name = p.Name
                    conv.Photo = p.Photo
                    break 
                }
            }
        }

        convs = append(convs, conv)
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

	return participants, nil
}
