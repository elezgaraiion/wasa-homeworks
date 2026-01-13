package database

func (db *appdbimpl) IsUserInConversation(userID, convID string) (bool, error) {
	var count int
	err := db.c.QueryRow(`
        SELECT COUNT(*)
        FROM conversation_participants
        WHERE conversation_id = ? AND user_id = ?
    `, convID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}