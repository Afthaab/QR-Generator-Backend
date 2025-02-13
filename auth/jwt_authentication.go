package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (a *authentication) GenerateJWT(uid string) (string, error) {
	// Set the claims
	claims := Claims{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)), // Token expires in 1 hour
			Issuer:    "qrGenerator",
		},
	}

	// Create the token using the claims and the signing key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate and return the signed token
	tokenString, err := token.SignedString([]byte(a.signingKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *authentication) ValidateToken(tokenString string) (*Claims, error) {
	// Parse the token with the key
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Check if the signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(a.signingKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Validate the claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
