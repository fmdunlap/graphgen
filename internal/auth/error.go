package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

type KeyParsingError struct {
	KeyType string
	Err     error
}

func (e *KeyParsingError) Error() string {
	return fmt.Sprintf("error parsing %s key: %v", e.KeyType, e.Err)
}

type InvalidTokenError struct {
	Token *jwt.Token
}

func (e *InvalidTokenError) Error() string {
	return fmt.Sprintf("invalid token: %v", e.Token)
}
