package database

import (
    "database/sql"
    "github.com/aritz/wasa-homeworks/service/models"
)

func (db *appdbimpl) LeaveGroup(userID, convID string) error {
    // 1. Comprobar que la conversación existe y es grupo
    var convType string
    err := db.c.QueryRow(`
        SELECT type FROM conversations WHERE id = ?
    `, convID).Scan(&convType)
    if err == sql.ErrNoRows {
        return models.ErrConversationNotFound
    }
    if err != nil {
        return err
    }
    if convType != "group" {
        return models.ErrForbidden
    }

    // 2. Comprobar que el usuario es miembro
    var exists int
    err = db.c.QueryRow(`
        SELECT COUNT(*) FROM conversation_participants
        WHERE conversation_id = ? AND user_id = ?
    `, convID, userID).Scan(&exists)
    if err != nil {
        return err
    }
    if exists == 0 {
        return models.ErrForbidden
    }

    // 3. Borrar usuario de participantes y meta
    _, err = db.c.Exec(`
        DELETE FROM conversation_participants
        WHERE conversation_id = ? AND user_id = ?
    `, convID, userID)
    if err != nil {
        return err
    }

    _, err = db.c.Exec(`
        DELETE FROM conversation_user_meta
        WHERE conversation_id = ? AND user_id = ?
    `, convID, userID)
    if err != nil {
        return err
    }

    return nil
}
