package services

import (
	"crypto/rsa"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateToken(userID int, email string) (string, error)
	ValidateToken(tokenString string) (int, error)
}

type jwtService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
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
	return &jwtService{
		privateKey: privKey,
		publicKey:  &privKey.PublicKey,
	}, nil
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

// ValidateToken validates JWT, returns userID if ok
func (j *jwtService) ValidateToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что используется правильный signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.publicKey, nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		uid, ok := claims["user_id"].(float64)
		if !ok {
			return 0, errors.New("user_id not found or invalid")
		}
		return int(uid), nil
	}
	return 0, errors.New("invalid token")
}
