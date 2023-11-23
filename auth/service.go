package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strconv"
	"time"
)

type Service interface {
	GenerateToken(userID string) (string, error)
}

type jwtService struct {
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func NewService() *jwtService {
	return &jwtService{}
}

func (s *jwtService) GenerateToken(userID string) (string, error) {
	secretKey := os.Getenv("JWT_AUTH_SECRET")
	jwtExpires := os.Getenv("AUTH_JWT_EXPIRES_IN_MINUTE")
	jwtExpiresInt, err := strconv.ParseInt(jwtExpires, 10, 64)
	if err != nil {
		return "", err
	}

	claim := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(jwtExpiresInt))),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}

func (s *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	secretKey := os.Getenv("JWT_AUTH_SECRET")
	tokens, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return tokens, err
	}

	return tokens, err
}
