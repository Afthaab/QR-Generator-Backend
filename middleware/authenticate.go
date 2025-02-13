package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (m *middleware) Authenticate(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx := c.Request.Context()

		// Get token from Authorization header
		tokenStr, err := getTokenFromHeader(c)
		if err != nil {
			log.Err(err).Msg("could not extract token from Authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: Missing or invalid token",
			})
			return
		}

		// Validate JWT token
		claims, err := m.auth.ValidateToken(tokenStr)
		if err != nil {
			log.Error().Err(err).Msg("could not validate token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: Invalid token",
			})
			return
		}

		// If the token is valid, store the claims in the context
		ctx = context.WithValue(ctx, "claimskey", claims.Uid)

		// Create a new request with the updated context and assign it back to Gin context
		req := c.Request.WithContext(ctx)
		c.Request = req

		next(c)
	}
}

func getTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header is missing")
	}

	// Ensure the header format is "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("Authorization header format must be 'Bearer <token>'")
	}

	return parts[1], nil
}
