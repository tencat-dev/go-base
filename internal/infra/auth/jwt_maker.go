package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/tencat-dev/go-base/internal/biz"
	"github.com/tencat-dev/go-base/internal/conf"
)

type JWTClaims struct {
	SessionID string        `json:"sid,omitempty"`
	Role      string        `json:"role,omitempty"`
	Type      biz.TokenType `json:"type,omitempty"`
	jwt.RegisteredClaims
}

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(c *conf.JWT) biz.TokenMaker {
	return &JWTMaker{secretKey: c.Secret}
}

func (j *JWTMaker) CreateAccessToken(payload biz.AccessPayload) (string, error) {
	now := time.Now().UTC()

	claims := JWTClaims{
		SessionID: payload.SessionID.String(),
		Type:      biz.AccessToken,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   payload.UserID.String(),
			ExpiresAt: jwt.NewNumericDate(now.Add(payload.TTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTMaker) CreateRefreshToken(payload biz.RefreshPayload) (string, error) {
	now := time.Now().UTC()

	claims := JWTClaims{
		SessionID: payload.SessionID.String(),
		Type:      biz.RefreshToken,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   payload.UserID.String(),
			ExpiresAt: jwt.NewNumericDate(now.Add(payload.TTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}
