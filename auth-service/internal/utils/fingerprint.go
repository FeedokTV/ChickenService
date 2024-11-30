package utils

import (
	"github.com/gin-gonic/gin"
	murmur3 "github.com/yihleego/murmurhash3"
)

// You can make something like this or more cooler! Check implementations in other languages or projects
func GenerateFingerprint(ctx *gin.Context) string {
	userAgent := ctx.GetHeader("User-Agent")
	clientIP := ctx.ClientIP()
	// For example
	acceptLanguage := ctx.GetHeader("Accept-Language")

	rawData := userAgent + clientIP + acceptLanguage

	hash := murmur3.New128().HashBytes([]byte(rawData)).String()

	return hash
}
