package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Scopes struct{}

type JWTClaims struct {
	Scopes Scopes
	jwt.RegisteredClaims
}

func newJWTClaims(username string) JWTClaims {
	return JWTClaims{
		Scopes: Scopes{},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "graphgen",
			Subject:   username,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationHours * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
}
