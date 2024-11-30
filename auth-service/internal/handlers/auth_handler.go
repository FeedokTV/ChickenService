package handlers

import (
	"auth-service/internal/auth"
	"auth-service/internal/domain"
	"auth-service/internal/services"
	"auth-service/internal/utils"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	tokenService *services.TokenService
	userService  *services.UserService
}

func NewAuthHandler(tokenService *services.TokenService, userService *services.UserService) *AuthHandler {
	return &AuthHandler{
		tokenService: tokenService,
		userService:  userService,
	}
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var requestForm struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&requestForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	user := &domain.User{
		Username: requestForm.Username,
		Password: requestForm.Password,
	}

	createdUser, apiErr := h.userService.CreateUser(user)
	if apiErr != nil {
		ctx.JSON(apiErr.Code, gin.H{"details": apiErr.Details, "error": apiErr.Message})
		return
	}

	fingerprint := utils.GenerateFingerprint(ctx)
	token, apiErr := h.tokenService.CreateToken(context.Background(), createdUser.ID, fingerprint)

	if apiErr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"details": "Internal server error", "error": "Cannot authorizate"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token.Token})
}

func (h *AuthHandler) Auth(ctx *gin.Context) {

	var authForm struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Check request
	if err := ctx.ShouldBindJSON(&authForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	// Find user
	user, apiErr := h.userService.GetUserByUsername(authForm.Username)

	if apiErr != nil {
		if apiErr.Code == 404 {
			ctx.JSON(http.StatusBadRequest, gin.H{"details": "", "error": "Invalid username or password"})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"details": "Try again", "error": "Internal server error"})
			return
		}
	}

	if !auth.CheckPassword(authForm.Password, user.Password) {
		ctx.JSON(http.StatusBadRequest, gin.H{"details": "", "error": "Invalid username or password"})
		return
	}

	// This thing not merely about creating new token,
	// Its replacing token for current session, so maybe if we use just session or token for redis key
	// There would be a problem here with repeating records
	fingerprint := utils.GenerateFingerprint(ctx)

	token, apiErr := h.tokenService.CreateToken(context.Background(), user.ID, fingerprint)

	if apiErr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"details": "Internal server error", "error": "Cannot authorizate"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token.Token})
}

func (h *AuthHandler) Validate(ctx *gin.Context) {

	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		ctx.Abort()
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	fingerprint := utils.GenerateFingerprint(ctx)
	tokenClaims, apiErr := h.tokenService.ValidateToken(context.Background(), tokenString, fingerprint)

	if tokenClaims == nil || apiErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"valid": "no"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"valid": "yes", "user_id": tokenClaims.UserID})
}
