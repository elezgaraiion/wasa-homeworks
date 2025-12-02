package database

import(
	"github.com/aritz/wasa-homeworks/service/models"

)
func (db *appdbimpl) RemoveReaction(userID, convID, msgID, reactionID string) error {
    // Validar existencia de conversación y mensaje
    var exists int
    err := db.c.QueryRow(`SELECT COUNT(*) FROM conversations WHERE id = ?`, convID).Scan(&exists)
    if err != nil || exists == 0 {
        return models.ErrConversationNotFound
    }

    err = db.c.QueryRow(`SELECT COUNT(*) FROM messages WHERE id = ? AND conversation_id = ?`, msgID, convID).Scan(&exists)
    if err != nil || exists == 0 {
        return models.ErrMessageNotFound
    }

    // Eliminar sólo si la reacción pertenece al usuario
    res, err := db.c.Exec(`DELETE FROM reactions WHERE id = ? AND user_id = ? AND message_id = ?`, reactionID, userID, msgID)
    if err != nil {
        return err
    }

    affected, _ := res.RowsAffected()
    if affected == 0 {
        return models.ErrReactionNotFound
    }

    return nil
}
