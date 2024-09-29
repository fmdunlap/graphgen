package auth

import (
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"graphgen/internal/config"
)

const (
	expirationHours = 24
)

type Service interface {
	CreateToken(username string) (string, error)
	VerifyToken(tokenString string) (*jwt.Token, *JWTClaims, error)
}

type service struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func New(c *config.AuthConfig) Service {
	privateKey, publicKey, err := parseKeysFromPEM([]byte(c.PrivateKey), []byte(c.PublicKey))
	if err != nil {
		panic(err)
	}

	return &service{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

func (a *service) CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, newJWTClaims(username))

	tokenString, err := token.SignedString(a.privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *service) VerifyToken(tokenString string) (*jwt.Token, *JWTClaims, error) {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.publicKey, nil
	})
	if err != nil {
		return nil, nil, err
	}

	if !token.Valid {
		return nil, nil, &InvalidTokenError{Token: token}
	}

	return token, claims, nil
}

func parseKeysFromPEM(privateKeyBytes, publicKeyBytes []byte) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return nil, nil, &KeyParsingError{
			KeyType: "private",
			Err:     err,
		}
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return nil, nil, &KeyParsingError{
			KeyType: "public",
			Err:     err,
		}
	}

	return privKey, pubKey, nil
}
