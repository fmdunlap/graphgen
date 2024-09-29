package auth

import "github.com/golang-jwt/jwt/v5"

type ExampleClaimData struct {
	s string `json:"s"`
}

type JWTClaim struct {
	Data ExampleClaimData
	jwt.RegisteredClaims
}
