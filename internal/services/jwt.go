package services

import (
	"crypto/rsa"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateToken(userID int, email string) (string, error)
}

type jwtService struct {
	privateKey *rsa.PrivateKey
}

func NewJWTService(privateKeyPath string) (JWTService, error) {
	keyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
	if err != nil {
		return nil, err
	}
	return &jwtService{privateKey: privKey}, nil
}

func (j *jwtService) GenerateToken(userID int, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(j.privateKey)
}
