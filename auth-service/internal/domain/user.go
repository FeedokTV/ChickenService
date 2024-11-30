package domain

import (
	"auth-service/internal/utils"
	"time"
)

type (
	UserRepository interface {
		CreateUser(user *User) (int, *utils.APIError)
		GetUserByID(id int) (*User, *utils.APIError)
		GetUserByUsername(username string) (*User, *utils.APIError)
	}

	User struct {
		ID        int       `json:"id"`
		Username  string    `json:"username"`
		Password  string    `json:"password"`
		CreatedAt time.Time `json:"created_at"`
	}

	// DTO for user data
	UserResponse struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
	}
)

func (user *User) ToUserResponse() *UserResponse {
	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
	}
}
