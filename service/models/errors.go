package models

import "errors"

var (
	ErrForbidden                 = errors.New("forbidden")
	ErrConversationNotFound      = errors.New("conversation not found")
	ErrMessageNotFound           = errors.New("message not found")
	ErrReactionNotFound          = errors.New("reaction not found")
	ErrUserAlreadyInConversation = errors.New("user already in conversation")
	ErrUserNotFound              = errors.New("user not found")
	ErrNameConflict              = errors.New("name conflict")
)
