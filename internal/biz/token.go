package biz

import (
	"time"

	"github.com/google/uuid"
)

type TokenMaker interface {
	CreateAccessToken(payload AccessPayload) (string, error)
	CreateRefreshToken(payload RefreshPayload) (string, error)
}

type AccessPayload struct {
	UserID    uuid.UUID
	SessionID uuid.UUID
	TTL       time.Duration
}

type RefreshPayload struct {
	UserID    uuid.UUID
	SessionID uuid.UUID
	TTL       time.Duration
}

const (
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 7 * 24 * time.Hour
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)
