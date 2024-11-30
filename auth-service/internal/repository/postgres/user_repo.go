package repositories

import (
	"auth-service/internal/domain"
	"auth-service/internal/utils"
	"context"

	logger "auth-service/internal"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type PostgresUserRepo struct {
	db *pgx.Conn
}

func NewPostgresUserRepo(db *pgx.Conn) *PostgresUserRepo {
	return &PostgresUserRepo{db: db}
}

// Returns id of new user and service must insert its in already existing user structure
// Or you can make it returns full user
func (repo *PostgresUserRepo) CreateUser(user *domain.User) (int, *utils.APIError) {
	query := "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING ID"

	// Thanks pgx for doing escaping of special characters for us <3
	var id int
	err := repo.db.QueryRow(context.Background(), query, user.Username, user.Password).Scan(&id)
	if err != nil {
		logger.Error("Cannot create user",
			zap.Error(err))
		return 0, ClassifyDBerror(err)
	}

	return id, nil
}

func (repo *PostgresUserRepo) GetUserByID(id int) (*domain.User, *utils.APIError) {
	query := "SELECT * FROM users WHERE id = $1"
	row := repo.db.QueryRow(context.Background(), query, id)

	var user domain.User
	err := row.Scan(&user.ID, &user.Username, &user.CreatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		logger.Error("Cannot get user by ID",
			zap.Error(err))
		return nil, ClassifyDBerror(err)
	}
	return &user, nil
}

func (repo *PostgresUserRepo) GetUserByUsername(username string) (*domain.User, *utils.APIError) {
	query := "SELECT * FROM users WHERE username = $1"
	row := repo.db.QueryRow(context.Background(), query, username)

	var user domain.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		logger.Error("Cannot get user by Username",
			zap.Error(err))
		return nil, ClassifyDBerror(err)
	}
	return &user, nil
}
