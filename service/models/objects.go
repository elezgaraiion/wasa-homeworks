package models

import "time"

type User struct {
    ID    string `db:"id" json:"id"`             
    Name  string `db:"name" json:"name"`           
    Photo string `db:"photo,omitempty" json:"photo,omitempty"` 
}

type Conversation struct {
	ID                 string    `db:"id" json:"id"`
	Type               string    `db:"type" json:"type"`
	Name               string    `db:"name,omitempty" json:"name,omitempty"`
	Photo              string    `db:"photo,omitempty" json:"photo,omitempty"`
	LastMessagePreview string    `db:"last_message_preview,omitempty" json:"lastMessagePreview,omitempty"`
	LastMessageAt      time.Time `db:"last_message_at,omitempty" json:"lastMessageAt,omitempty"`
	
	LastMessageSenderID string `json:"lastMessageSenderId"` // Para saber si "fui yo"
	LastMessageStatus   string `json:"lastMessageStatus"`   // Para los ticks (sent, read)
	LastMessageSenderName string `json:"lastMessageSenderName"` // Para poner "Luis: Hola..." en grupos
	
	Participants []User    `db:"participants" json:"participants"`
	TempOrderAt  time.Time `json:"-"`
}

type Reaction struct {
    ID        string    `db:"id" json:"id"`         
    User      User      `db:"user" json:"user"`       
    Emoji     string    `db:"emoji" json:"emoji"`      
    CreatedAt time.Time `db:"created_at" json:"createdAt"` 
}

type Message struct {
    ID               string     `db:"id" json:"id"`                 
    Sender           User       `db:"sender" json:"sender"`             
    ConversationID   string     `db:"conversation_id" json:"conversationId"`    
    Text             string     `db:"text,omitempty" json:"text,omitempty"`     
    Photo            string     `db:"photo,omitempty" json:"photo,omitempty"`    
    ReplyToMessageID string     `db:"reply_to_message_id,omitempty" json:"replyToMessageId,omitempty"` 
    CreatedAt        time.Time  `db:"created_at" json:"createdAt"`         
    Reactions        []Reaction `db:"reactions,omitempty" json:"reactions,omitempty"`
    Status           string     `db:"status" json:"status"`             
}