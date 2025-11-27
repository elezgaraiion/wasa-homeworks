package models

import "time"

type User struct {
	ID    string `db:"id"`             
	Name  string `db:"name"`           
	Photo string `db:"photo,omitempty"` 
}

type Conversation struct {
	ID                string    `db:"id"`                            
	Type              string    `db:"type"`                          
	Name              string    `db:"name,omitempty"`                
	Photo             string    `db:"photo,omitempty"`               
	LastMessagePreview string   `db:"last_message_preview,omitempty"`
	LastMessageAt     time.Time `db:"last_message_at,omitempty"`     
	Participants      []User    `db:"participants"`  
	TempOrderAt time.Time `json:"-"`                
}

type Reaction struct {
	ID        string    `db:"id"`         
	User      User      `db:"user"`       
	Emoji     string    `db:"emoji"`      
	CreatedAt time.Time `db:"created_at"` 
}

type Message struct {
	ID               string     `db:"id"`                 
	Sender           User       `db:"sender"`             
	ConversationID   string     `db:"conversation_id"`    
	Text             string     `db:"text,omitempty"`     
	Photo            string     `db:"photo,omitempty"`    
	ReplyToMessageID string     `db:"reply_to_message_id,omitempty"` 
	CreatedAt        time.Time  `db:"created_at"`         
	Reactions        []Reaction `db:"reactions,omitempty"`
	Status           string     `db:"status"`             
}