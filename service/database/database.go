package database

import (
	"database/sql"
	"errors"
	"github.com/aritz/wasa-homeworks/service/models"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {

	CreateUser(u models.User) error
    GetUserIdByName(name string) (string, error)
    GetUserByID(id string) (models.User, error)
    UpdateUserName(id,newName string)(models.User, error)
    UpdateMyPhoto(id,photoURL string)(models.User, error)
    ListUsers(query, currentUserID string) ([]models.User, error)	
    GetMyConversations(userID string) ([]models.Conversation, error)
    CreateGroup(creatorID string, name string, users []string) (models.Conversation, error)
    GetConversationProfile(userID, convID string) (models.Conversation, error)
    MarkConversationSeen(userID, convID string) error
    GetMessagesInConversation(userID string,convID string,limit int,before string) ([]models.Message, error)
    SendMessage(senderID, convID, text, photoURL, replyToMessageID string) (models.Message, error)
    GetMessageByID(userID, convID, messageID string) (models.Message, error)
    ForwardMessage(userID, sourceConvID, messageID, targetConvID string) (models.Message, error)
    DeleteMessage(userID, convID, messageID string) error
    GetReactions(messageID string) ([]models.Reaction, error)
    IsUserInConversation(userID, convID string) (bool, error)
    AddReaction(userID, convID, msgID, emoji string) (models.Reaction, error)
    RemoveReaction(userID, convID, msgID, reactionID string) error
    AddUserToGroup(requestUserID, convID, targetUserID string) (models.Conversation, error)
    Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
    if db == nil {
        return nil, errors.New("database connection required")
    }

    // Crear tabla Users
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id TEXT PRIMARY KEY,
            name TEXT NOT NULL,
            photo TEXT
        );
    `)
    if err != nil {
        return nil, err
    }

    // Crear tabla Conversations
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS conversations (
            id TEXT PRIMARY KEY,
            type TEXT NOT NULL,
            name TEXT,
            photo TEXT,
            last_message_preview TEXT,
            last_message_at TEXT
        );

        -- Tabla de participantes de la conversación
        CREATE TABLE IF NOT EXISTS conversation_participants (
            conversation_id TEXT NOT NULL,
            user_id TEXT NOT NULL,
            PRIMARY KEY (conversation_id, user_id),
            FOREIGN KEY (conversation_id) REFERENCES conversations(id),
            FOREIGN KEY (user_id) REFERENCES users(id)
        );
        CREATE TABLE IF NOT EXISTS conversation_user_meta (
            conversation_id TEXT NOT NULL,
            user_id TEXT NOT NULL,
            joined_at TEXT NOT NULL,
            last_seen_message_at TEXT,
            PRIMARY KEY(conversation_id, user_id),
            FOREIGN KEY(conversation_id) REFERENCES conversations(id),
            FOREIGN KEY(user_id) REFERENCES users(id)
        );

    `)
    if err != nil {
        return nil, err
    }

    // Crear tabla Messages
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS messages (
            id TEXT PRIMARY KEY,
            sender_id TEXT NOT NULL,
            conversation_id TEXT NOT NULL,
            text TEXT,
            photo TEXT,
            reply_to_message_id TEXT,
            created_at TEXT NOT NULL,
            status TEXT NOT NULL,
            FOREIGN KEY(sender_id) REFERENCES users(id),
            FOREIGN KEY(conversation_id) REFERENCES conversations(id)
        );
    `)
    if err != nil {
        return nil, err
    }

    // Crear tabla Reactions
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS reactions (
            id TEXT PRIMARY KEY,
            user_id TEXT NOT NULL,
            message_id TEXT NOT NULL,
            emoji TEXT NOT NULL,
            created_at TEXT NOT NULL,
            FOREIGN KEY(user_id) REFERENCES users(id),
            FOREIGN KEY(message_id) REFERENCES messages(id)
        );
    `)
    if err != nil {
        return nil, err
    }

    return &appdbimpl{c: db}, nil
}


func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
