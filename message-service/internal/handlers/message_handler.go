package handlers

import (
	"message-service/internal/domain"
	"message-service/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	messageService *services.MessageService
}

func NewMessageHandler(messageService *services.MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

func (h *MessageHandler) SendMessage(ctx *gin.Context) {

	var requestForm struct {
		RecipientID int    `json:"recipient_id"`
		Content     string `json:"content"`
	}

	if err := ctx.ShouldBindJSON(&requestForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	// Get user from Context
	senderIDString, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	// Convert to int
	senderID, ok := senderIDString.(int)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid User ID type"})
		return
	}

	message, apiErr := h.messageService.CreateMessage(
		ctx.Request.Context(),
		requestForm.RecipientID,
		senderID,
		requestForm.Content,
	)
	if apiErr != nil {
		ctx.JSON(apiErr.Code, gin.H{"details": apiErr.Details, "error": apiErr.Message})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": message})
}

func (h *MessageHandler) GetConversationMessages(ctx *gin.Context) {

	user1IDString, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized request"})
		return
	}

	user1ID, ok := user1IDString.(int)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	recipientIDString := ctx.Query("recipient_id")

	if recipientIDString == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	recipientID, err := strconv.Atoi(recipientIDString)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	messages, apiErr := h.messageService.GetConversationMessages(
		ctx.Request.Context(),
		user1ID,
		recipientID,
	)
	if apiErr != nil {
		ctx.JSON(apiErr.Code, gin.H{"details": apiErr.Details, "error": apiErr.Message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"messages": messages})
}

func (h *MessageHandler) UpdateMessageStatus(ctx *gin.Context) {

	var requestForm struct {
		MessageID int                  `json:"message_id"`
		Status    domain.MessageStatus `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&requestForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	if !requestForm.Status.IsValid() {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	apiErr := h.messageService.UpdateMessageStatus(ctx.Request.Context(), requestForm.MessageID, requestForm.Status)

	if apiErr != nil {
		ctx.JSON(apiErr.Code, gin.H{"details": apiErr.Details, "error": apiErr.Message})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}
