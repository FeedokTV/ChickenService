package repositories

import (
	"context"
	"database/sql"
	"message-service/internal/domain"
	"message-service/internal/utils"
	"time"

	logger "message-service/internal"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

type PostgresMessageRepo struct {
	db *pgx.Conn
}

func NewPostgresMessageRepo(db *pgx.Conn) *PostgresMessageRepo {
	return &PostgresMessageRepo{db: db}
}

func (r *PostgresMessageRepo) SendMessage(ctx context.Context, message *domain.Message) (*domain.Message, *utils.APIError) {
	// Read about formatting here https://www.geeksforgeeks.org/time-formatting-in-golang/
	today := time.Now()
	query := "INSERT INTO messages_" + today.Format("2006_01_02") + ` (recipient_id, sender_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, timestamp, status`

	err := r.db.QueryRow(
		ctx,
		query,
		message.RecipientID,
		message.SenderID,
		message.Content).Scan(&message.MessageID, &message.Timestamp, &message.Status)

	if err != nil {
		logger.Error("Cannot send message",
			zap.Int("From id", message.SenderID),
			zap.Int("To id", message.RecipientID),
			zap.Error(err))
		// This error says that we dont have actual partition in db
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "42P01" {
			apiErr := r.CreatePartition(ctx)

			if apiErr != nil {
				return nil, apiErr
			}

			// Retry query
			retryErr := r.db.QueryRow(
				ctx,
				query,
				message.RecipientID,
				message.SenderID,
				message.Content).Scan(&message.MessageID, &message.Timestamp, &message.Status)

			if retryErr != nil {
				logger.Error("Cannot send message",
					zap.Int("From id", message.SenderID),
					zap.Int("To id", message.RecipientID),
					zap.Error(retryErr))
				return nil, ClassifyDBerror(retryErr)
			}
			return message, nil
			// User not found error
		} else if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23503" {
			return nil, utils.NewAPIError(404, "User does not exist", "")
		}
		return nil, ClassifyDBerror(err)
	}

	return message, nil
}

func (r *PostgresMessageRepo) CreatePartition(ctx context.Context) *utils.APIError {
	today := time.Now()
	partitionQuery := `CREATE TABLE IF NOT EXISTS messages_` + today.Format("2006_01_02") + ` PARTITION OF messages
				FOR VALUES FROM ('` + today.Format("2006-01-02") + `') TO ('` + today.Add(24*time.Hour).Format("2006-01-02") + `');`

	_, createErr := r.db.Exec(ctx, partitionQuery)
	if createErr != nil {
		logger.Error("Cannot create new partition messages_"+today.Format("2006_01_02"),
			zap.Error(createErr))
		return utils.NewAPIError(500, "Internal server error", "Please report about this")
	}

	logger.Info("Created new partition messages_" + today.Format("2006_01_02"))

	return nil
}
func (r *PostgresMessageRepo) GetMessagesBetweenUsers(ctx context.Context, user1ID, user2ID int) (*[]domain.Message, *utils.APIError) {
	today := time.Now()
	query := `SELECT id, sender_id, recipient_id, content, timestamp, status FROM ` + "messages_" + today.Format("2006_01_02") +
		` WHERE 
			(sender_id = $1 AND recipient_id = $2)
			OR (sender_id = $2 AND recipient_id = $1)
		ORDER BY 
			timestamp DESC;`

	rows, err := r.db.Query(ctx, query, user1ID, user2ID)
	if err != nil {
		logger.Error("Cannot get conversation",
			zap.Int("User 1", user1ID),
			zap.Int("User 2", user2ID),
			zap.Error(err))
		return nil, ClassifyDBerror(err)
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var msg domain.Message
		if err := rows.Scan(
			&msg.MessageID,
			&msg.SenderID,
			&msg.RecipientID,
			&msg.Content,
			&msg.Timestamp,
			&msg.Status,
		); err != nil {
			return nil, ClassifyDBerror(err)
		}
		messages = append(messages, msg)
	}

	return &messages, nil
}
func (r *PostgresMessageRepo) UpdateMessageStatus(ctx context.Context, messageID int, status string) *utils.APIError {
	today := time.Now()
	query := "UPDATE messages_" + today.Format("2006_01_02") + ` SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, status, messageID)
	if err != nil {
		logger.Error("Cannot update message status",
			zap.Int("Message ID", messageID),
			zap.String("Status", status),
			zap.Error(err))
		return ClassifyDBerror(err)
	}
	return nil
}

func (r *PostgresMessageRepo) GetMessageByID(ctx context.Context, messageID int) (*domain.Message, *utils.APIError) {
	query := `
		SELECT message_id, recipient_id, sender_id, content, timestamp, status
		FROM messages
		WHERE message_id = $1
	`
	var msg domain.Message
	err := r.db.QueryRow(ctx, query, messageID).Scan(
		&msg.MessageID,
		&msg.RecipientID,
		&msg.SenderID,
		&msg.Content,
		&msg.Timestamp,
		&msg.Status)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error("Cannot get message by ID",
			zap.Int("Message ID", messageID),
			zap.Error(err))
		return nil, ClassifyDBerror(err)
	}

	return &msg, nil
}
