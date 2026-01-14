package database

import (
	"github.com/aritz/wasa-homeworks/service/models"
)

func (db *appdbimpl) RemoveReaction(userID, convID, msgID, reactionID string) error {
	res, err := db.c.Exec(`
        DELETE FROM reactions 
        WHERE id = ? AND user_id = ? AND message_id = ?
    `, reactionID, userID, msgID)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return models.ErrReactionNotFound
	}
	return nil
}
