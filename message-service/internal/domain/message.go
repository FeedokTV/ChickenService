package domain

import (
	"context"
	"message-service/internal/utils"
	"time"
)

type MessageStatus string

const (
	Sent    MessageStatus = "sent"
	Watched MessageStatus = "watched"
)

func (s MessageStatus) IsValid() bool {
	switch s {
	case Sent, Watched:
		return true
	default:
		return false
	}
}

type (
	MessageRepository interface {
		SendMessage(ctx context.Context, message *Message) (*Message, *utils.APIError)
		GetMessagesBetweenUsers(ctx context.Context, user1ID, user2ID int) (*[]Message, *utils.APIError)
		UpdateMessageStatus(ctx context.Context, messageID int, status string) *utils.APIError
		GetMessageByID(ctx context.Context, messageID int) (*Message, *utils.APIError)
	}
	Message struct {
		MessageID   int       `json:"id"`
		RecipientID int       `json:"recipient_id"`
		SenderID    int       `json:"sender_id"`
		Content     string    `json:"content"`
		Timestamp   time.Time `json:"timestamp"`
		Status      string    `json:"status"`
	}
)
