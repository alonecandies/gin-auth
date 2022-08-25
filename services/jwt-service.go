package services

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

type JWTService interface {
	GenerateToken(userId uint) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtCustomClaims struct {
	UserId uint `json:"user_id"`
	jwt.RegisteredClaims
}

type jwtService struct {
	secretKey string
	issuer    string
}

func NewJWTService() JWTService {
	return &jwtService{
		secretKey: getSecretKey(),
		issuer:    getIssuer(),
	}
}

func getSecretKey() string {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		secretKey = "secret"
	}
	return secretKey
}

func getIssuer() string {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	issuer := os.Getenv("ISSUER")
	if issuer == "" {
		issuer = "abcxyz"
	}
	return issuer
}

func (j *jwtService) GenerateToken(userId uint) (string, error) {
	claims := jwtCustomClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    j.issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		panic(err)
	}
	return t, nil
}

func (j *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
}
