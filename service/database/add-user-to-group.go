package database

import (
	"database/sql"
	"fmt"
	"time"
	"github.com/aritz/wasa-homeworks/service/models"
)

func (db *appdbimpl) AddUserToGroup(requestUserID, convID, targetUserID string) (models.Conversation, error) {    var convType string
    err := db.c.QueryRow(`SELECT type FROM conversations WHERE id = ?`, convID).Scan(&convType)
    if err == sql.ErrNoRows { return models.Conversation{}, models.ErrConversationNotFound }
    if err != nil { return models.Conversation{}, err }
    if convType != "group" { return models.Conversation{}, fmt.Errorf("cannot add users to private conversations") }

    var exists int
    err = db.c.QueryRow(`SELECT COUNT(*) FROM conversation_participants WHERE conversation_id = ? AND user_id = ?`, convID, requestUserID).Scan(&exists)
    if err != nil { return models.Conversation{}, err }
    if exists == 0 { return models.Conversation{}, models.ErrForbidden }

    err = db.c.QueryRow(`SELECT COUNT(*) FROM conversation_participants WHERE conversation_id = ? AND user_id = ?`, convID, targetUserID).Scan(&exists)
    if err != nil { return models.Conversation{}, err }
    if exists > 0 { return db.GetConversationProfile(requestUserID, convID) }

    tx, err := db.c.Begin()
    if err != nil { return models.Conversation{}, err }
    defer tx.Rollback()

    _, err = tx.Exec(`INSERT INTO conversation_participants(conversation_id, user_id) VALUES (?, ?)`, convID, targetUserID)
    if err != nil { return models.Conversation{}, err }

    now := time.Now().UTC()
    nowStr := now.Format(time.RFC3339Nano)

    _, err = tx.Exec(`
        INSERT OR REPLACE INTO conversation_user_meta(conversation_id, user_id, joined_at, last_seen_message_at)
        VALUES (?, ?, ?, ?)
    `, convID, targetUserID, nowStr, nowStr)
    if err != nil { return models.Conversation{}, err }

    var targetName string
    err = tx.QueryRow("SELECT name FROM users WHERE id = ?", targetUserID).Scan(&targetName)
    if err != nil { targetName = "User" }

    systemText := fmt.Sprintf("SYS:added %s to the group", targetName)
    
    msgTime := now.Add(time.Millisecond)
    err = db.createSystemMessage(tx, convID, requestUserID, systemText, msgTime.Format(time.RFC3339Nano))
    if err != nil { return models.Conversation{}, fmt.Errorf("error creating system message: %w", err) }

    _, err = tx.Exec(`
        UPDATE conversations 
        SET last_message_preview = ?, last_message_at = ?
        WHERE id = ?
    `, fmt.Sprintf("added %s to the group", targetName), msgTime.Format(time.RFC3339Nano), convID)
    if err != nil { return models.Conversation{}, err }

    if err := tx.Commit(); err != nil { return models.Conversation{}, err }

    return db.GetConversationProfile(requestUserID, convID)
}