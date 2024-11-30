package services

import (
	"auth-service/internal/auth"
	"auth-service/internal/domain"
	"auth-service/internal/utils"

	logger "auth-service/internal"

	"go.uber.org/zap"
)

type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user *domain.User) (*domain.User, *utils.APIError) {

	foundUser, apiErr := s.repo.GetUserByUsername(user.Username)

	// Unexpectable db error?
	if apiErr != nil {
		logger.Error("Cannot create user. Unexpectable DB error",
			zap.String("error", apiErr.Message),
			zap.String("details", apiErr.Details))
		return nil, utils.NewAPIError(500, "Internal server error", "Database error")
	}

	// User already exists!
	if foundUser != nil {
		return nil, utils.NewAPIError(409, "User already exists", "")
	}

	hashedPassword, err := auth.HashPassword(user.Password)

	if err != nil {
		logger.Error("Cannot hash user password",
			zap.Error(err))
		return nil, utils.NewAPIError(500, "Internal server error while registering your user", "Please try again")
	}

	user.Password = hashedPassword

	newUserId, apiErr := s.repo.CreateUser(user)
	if apiErr != nil {
		logger.Error("Cannot create user",
			zap.String("error", apiErr.Message),
			zap.String("details", apiErr.Details))
		return nil, utils.NewAPIError(500, "Internal server error while registering your user", "Please try again")
	}

	user.ID = newUserId

	return user, nil
}

func (s *UserService) GetUserByUsername(username string) (*domain.User, *utils.APIError) {

	foundUser, err := s.repo.GetUserByUsername(username)

	// Unexpectable db error?
	if err != nil {
		logger.Error("Cannot get user. Unexpectable DB error",
			zap.String("error", err.Message),
			zap.String("details", err.Details))
		return nil, utils.NewAPIError(500, "Internal server error", "Please try again")
	}

	if foundUser == nil {
		return nil, utils.NewAPIError(404, "User not found", "")
	}

	return foundUser, nil
}

func (s *UserService) GetUserByID(id int) (*domain.User, *utils.APIError) {

	foundUser, err := s.repo.GetUserByID(id)

	// Unexpectable db error?
	if err != nil {
		logger.Error("Cannot get user. Unexpectable DB error",
			zap.String("error", err.Message),
			zap.String("details", err.Details))
		return nil, utils.NewAPIError(500, "Internal server error", "Please try again")
	}

	if foundUser == nil {
		return nil, utils.NewAPIError(404, "User not found", "")
	}

	return foundUser, nil
}
