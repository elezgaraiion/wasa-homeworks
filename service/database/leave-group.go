package database

import (
	"database/sql"
	"errors"

	"github.com/aritz/wasa-homeworks/service/models"
)

func (db *appdbimpl) LeaveGroup(userID, convID string) error {
	var convType string
	err := db.c.QueryRow(`
        SELECT type FROM conversations WHERE id = ?
    `, convID).Scan(&convType)
	if errors.Is(err, sql.ErrNoRows) {
		return models.ErrConversationNotFound
	}
	if err != nil {
		return err
	}
	if convType != "group" {
		return models.ErrForbidden
	}

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
