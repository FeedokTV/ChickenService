package main

import (
	logger "auth-service/internal"
	"auth-service/internal/handlers"
	postgresRepos "auth-service/internal/repository/postgres"
	redisRepos "auth-service/internal/repository/redis"
	"auth-service/internal/services"
	"context"
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	authHandler *handlers.AuthHandler
)

func main() {
	logger.InitLogger()
	defer logger.Sync()

	logger.Info("Starting...")

	// This is actually idea how you can create centralized config manager
	//cfg := config.Load() // Load configuration
	//logger.Info("Config loaded")

	// Open DB connection (Postgres)
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	logger.Info(dbConnString)

	db, err := pgx.Connect(context.Background(), dbConnString)
	if err != nil {
		logger.Fatal("Cannot open DB connection", zap.Error(err))
	}
	defer db.Close(context.Background())

	// Open Redis connection
	redisPort := os.Getenv("REDIS_PORT")
	redisDbId := os.Getenv("REDIS_DB_ID")
	redisConnString := fmt.Sprintf("redis://default:@redis:%s/%s", redisPort, redisDbId)
	opt, _ := redis.ParseURL(redisConnString)

	client := redis.NewClient(opt)

	if err := client.Ping(context.Background()).Err(); err != nil {
		logger.Fatal("Cannot open redis connection", zap.Error(err))
	}

	// Initialize repositories
	tokenRepository := redisRepos.NewRedisTokenRepo(client)
	userRepository := postgresRepos.NewPostgresUserRepo(db)
	logger.Info("Initialized repositories")

	// Initialize services
	tokenService := services.NewTokenService(tokenRepository)
	userService := services.NewUserService(userRepository)
	logger.Info("Initialized services")

	// Initialize handlers
	authHandler = handlers.NewAuthHandler(tokenService, userService)
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

	// Auth router
	authRouter := router.Group("/auth")
	authRouter.POST("/login", authHandler.Auth)
	authRouter.GET("/validate", authHandler.Validate)
	authRouter.POST("/register", authHandler.Register)

	return router
}
