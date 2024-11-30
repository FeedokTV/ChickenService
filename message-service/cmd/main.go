package main

import (
	"context"
	"fmt"
	logger "message-service/internal"
	"message-service/internal/handlers"
	"message-service/internal/middlewares"
	repositories "message-service/internal/repository/postgres"
	"message-service/internal/services"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

var (
	messageHandler *handlers.MessageHandler
)

func main() {
	logger.InitLogger()
	defer logger.Sync()
	logger.Info("Starting...")
	//cfg := config.Load() // Load configuration
	//logger.Info("Config loaded")

	// Open DB connection (Postgres)
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := pgx.Connect(context.Background(), dbConnString)
	if err != nil {
		logger.Fatal("Cannot open DB connection", zap.Error(err))
	}
	defer db.Close(context.Background())

	// Initialize repositories
	messageRepository := repositories.NewPostgresMessageRepo(db)
	logger.Info("Initialized repositories")

	// Initialize services
	messageService := services.NewMessageService(messageRepository)
	logger.Info("Initialized services")

	// Initialize middlewares
	logger.Info("Initialized middlewares")

	// Initialize handlers
	messageHandler = handlers.NewMessageHandler(messageService)
	logger.Info("Initialized handlers")

	service_address := fmt.Sprintf("0.0.0.0:%s", os.Getenv("SERVICE_PORT"))
	router := InitRouter() // Setup router
	router.Run(service_address)

}

func InitRouter() *gin.Engine {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	// End-points with auth only
	protected := router.Group("/")
	protected.Use(middlewares.TokenValidationMiddleware(os.Getenv("AUTH_SERVICE_ADDR")))

	protected.POST("/sendMessage", messageHandler.SendMessage)
	protected.GET("/getConversation", messageHandler.GetConversationMessages)
	protected.POST("/updateMessageStatus", messageHandler.UpdateMessageStatus)

	return router
}
