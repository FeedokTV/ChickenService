package services

import (
	"context"
	"message-service/internal/domain"
	"message-service/internal/utils"
)

type MessageService struct {
	repo domain.MessageRepository
}

func NewMessageService(repo domain.MessageRepository) *MessageService {
	return &MessageService{repo: repo}
}
func (s *MessageService) CreateMessage(ctx context.Context, recipientID, senderID int, content string) (*domain.Message, *utils.APIError) {

	if recipientID == senderID {
		return nil, utils.NewAPIError(418, "You can't send a message to yourself", "")
	}

	if len(content) == 0 {
		return nil, utils.NewAPIError(400, "Content cannot be empty", "")
	}

	if len(content) > 5000 {
		return nil, utils.NewAPIError(400, "Content exceeds maximum length", "")
	}

	message := &domain.Message{
		RecipientID: recipientID,
		SenderID:    senderID,
		Content:     content,
	}

	message, apiErr := s.repo.SendMessage(ctx, message)

	if apiErr != nil {
		return nil, apiErr
	}

	return message, nil
}

func (s *MessageService) GetConversationMessages(ctx context.Context, user1ID, user2ID int) (*[]domain.Message, *utils.APIError) {

	messages, apiErr := s.repo.GetMessagesBetweenUsers(ctx, user1ID, user2ID)

	if apiErr != nil {
		return nil, utils.NewAPIError(404, "Conversation or messages not found", "")
	}

	return messages, nil
}

func (s *MessageService) UpdateMessageStatus(ctx context.Context, messageID int, status domain.MessageStatus) *utils.APIError {

	return s.repo.UpdateMessageStatus(ctx, messageID, string(status))
}
