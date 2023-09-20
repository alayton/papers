package papers

import (
	"context"
)

type MessageType string

const (
	MessageConfirmation   = MessageType("confirmation")
	MessageForgotPassword = MessageType("forgot_password")
	MessageLocked         = MessageType("locked")
)

type Email struct {
	Address string
	Name    string
}

type Message struct {
	Type    MessageType
	To      []Email
	Cc      []Email
	Bcc     []Email
	ReplyTo []Email
	Data    map[string]interface{}
}

type Mailer interface {
	SendMessage(ctx context.Context, p *Papers, msg Message) error
}
