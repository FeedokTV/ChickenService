package middlewares

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type TokenValidationResponse struct {
	UserID int    `json:"user_id"`
	Status string `json:"valid"`
}

// TokenValidationMiddleware checks token from header by sendind request to auth service
func TokenValidationMiddleware(validationServiceURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Split from "Bearer"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}
		token := tokenParts[1]

		// Create context with timeout for request
		ctx, contextCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer contextCancel()

		// Prepare HTTP request

		client := &http.Client{}
		req, err := http.NewRequestWithContext(ctx, "GET", validationServiceURL+"/auth/validate", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create validation request"})
			c.Abort()
			return
		}

		// "Forward" all headers from request to authorization service
		req.Header.Set("Authorization", "Bearer "+token)

		req.Header.Set("User-Agent", c.GetHeader("User-Agent"))
		req.Header.Set("Accept-Language", c.GetHeader("Accept-Language"))
		req.Header.Set("X-Forwarded-For", c.ClientIP())

		// Make request to auth service
		resp, err := client.Do(req)

		// Is service unavailable?
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				c.JSON(http.StatusGatewayTimeout, gin.H{"error": "The service is taking too long to respond"})
			} else if errors.Is(err, &net.OpError{}) {
				c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Service unavailable, please try again later"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to communicate with the authentitication service"})
			}

			c.Abort()
			return
		}
		defer resp.Body.Close()

		// Token is invalid
		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token validation failed"})
			c.Abort()
			return
		}

		// Decoding answer
		var validationResponse TokenValidationResponse
		if err := json.NewDecoder(resp.Body).Decode(&validationResponse); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid response from validation service"})
			c.Abort()
			return
		}

		// Check token status
		if validationResponse.Status != "yes" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set user id in context, so then we can use it in handlers
		c.Set("user_id", validationResponse.UserID)
		c.Next()
	}
}
