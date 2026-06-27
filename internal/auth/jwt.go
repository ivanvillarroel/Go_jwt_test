package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	Generate(subject string, permissions []string) (string, error)
	Validate(tokenValue string) (*Claims, error)
}

type Claims struct {
	Subject     string
	Permissions []string
}

type JWTService struct {
	secret []byte
}

func NewJWTService(secret string) JWTService {
	return JWTService{secret: []byte(secret)}
}

type jwtClaims struct {
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

func (s JWTService) Generate(subject string, permissions []string) (string, error) {
	claims := jwtClaims{
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			Issuer:    "go-jwt-gin",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s JWTService) Validate(tokenValue string) (*Claims, error) {
	claims := &jwtClaims{}
	token, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}

		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid || claims.Subject == "" {
		return nil, errors.New("invalid token")
	}

	return &Claims{
		Subject:     claims.Subject,
		Permissions: claims.Permissions,
	}, nil
}

func (c Claims) HasPermission(permission string) bool {
	for _, current := range c.Permissions {
		if current == permission {
			return true
		}
	}

	return false
}
