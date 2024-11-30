package services

import (
	"auth-service/internal/auth"
	"auth-service/internal/domain"
	"auth-service/internal/utils"
)

type AuthService struct {
	userService *UserService
}

func NewAuthService(userService *UserService) *AuthService {
	return &AuthService{userService: userService}
}

func (s *AuthService) Register(user *domain.User) (*domain.User, *utils.APIError) {
	return s.userService.CreateUser(user)
}

func (s *AuthService) Login(username, password string) (string, *utils.APIError) {
	user, apiErr := s.userService.GetUserByUsername(username)
	if apiErr != nil {
		return "", apiErr
	}

	if !auth.CheckPassword(password, user.Password) {
		return "", utils.NewAPIError(401, "Invalid credentials", "")
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return "", utils.NewAPIError(500, "Error generating token", "")
	}

	return token, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*domain.User, *utils.APIError) {
	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		return nil, utils.NewAPIError(401, "Invalid token", "")
	}
	user, apiErr := s.userService.GetUserByID(claims.UserID)
	if apiErr != nil {
		return nil, apiErr
	}

	return user, nil
}
