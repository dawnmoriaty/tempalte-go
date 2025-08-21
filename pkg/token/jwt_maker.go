package token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JwtMaker struct {
	secretKey string
}

func (maker *JwtMaker) CreateToken(username string, roles []string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, roles, duration)
	if err != nil {
		return "", payload, err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, err
}

func (maker *JwtMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInValidToken
		}
		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.Parse(token, keyFunc)
	if err != nil {
		return nil, err
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInValidToken
	}
	return payload, nil
}

const minSecretKeyLength = 32

func NewJwtMaker(secretKey string) (TokenMaker, error) {
	if len(secretKey) < minSecretKeyLength {
		return nil, fmt.Errorf("secret key length must be %d characters", minSecretKeyLength)
	}
	return &JwtMaker{secretKey}, nil
}
