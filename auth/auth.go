package auth

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Uid string `json:"uid"`
	jwt.RegisteredClaims
}

type authentication struct {
	signingKey string
}

type Auth interface {
	GenerateJWT(uid string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

func NewAuth(signingKey string) Auth {
	return &authentication{
		signingKey: signingKey,
	}
}
