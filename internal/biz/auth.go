package biz

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/matthewhartstonge/argon2"
)

type AuthLogin struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

// Auth is a Auth model.
type Auth struct {
	ID           uuid.UUID `json:"id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Email        string    `json:"email,omitempty"`
	PasswordHash string    `json:"password_hash,omitempty"`
}

// AuthRepo is a Greater repo.
type AuthRepo interface {
	FindByEmail(context.Context, string) (*Auth, error)
}

// AuthBiz is a Auth usecase.
type AuthBiz struct {
	repo  AuthRepo
	authz PermissionChecker
}

// NewAuthBiz new a Auth usecase.
func NewAuthBiz(repo AuthRepo, authz PermissionChecker) *AuthBiz {
	return &AuthBiz{
		repo:  repo,
		authz: authz,
	}
}

// Login creates a Auth, and returns the new Auth.
func (b *AuthBiz) Login(ctx context.Context, u *AuthLogin) (*Auth, error) {
	user, err := b.repo.FindByEmail(ctx, u.Email)
	if err != nil {
		return nil, err
	}

	ok, err := argon2.VerifyEncoded([]byte(u.Password), []byte(user.PasswordHash))
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, fmt.Errorf("invalid password")
	}

	return user, nil
}
