package biz

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/matthewhartstonge/argon2"
)

// User is a User model.
type User struct {
	ID           uuid.UUID `json:"id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Email        string    `json:"email,omitempty"`
	PasswordHash string    `json:",omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

// UserRepo is a Greater repo.
type UserRepo interface {
	Save(context.Context, *User) (*User, error)
	Update(context.Context, *User) (*User, error)
	FindByID(context.Context, uuid.UUID) (*User, error)
	ListAll(context.Context) ([]*User, error)
	DeleteByID(context.Context, uuid.UUID) error
	ExistByID(context.Context, uuid.UUID) (bool, error)
}

// UserBiz is a User usecase.
type UserBiz struct {
	repo UserRepo
}

// NewUserBiz new a User usecase.
func NewUserBiz(repo UserRepo) *UserBiz {
	return &UserBiz{repo: repo}
}

// CreateUser creates a User, and returns the new User.
func (b *UserBiz) CreateUser(ctx context.Context, u *User) (*User, error) {
	argon := argon2.DefaultConfig()
	passwordHash, err := argon.HashEncoded([]byte(u.PasswordHash))
	if err != nil {
		return nil, err
	}

	u.PasswordHash = string(passwordHash)

	return b.repo.Save(ctx, u)
}

// UpdateUser creates a User, and returns the new User.
func (b *UserBiz) UpdateUser(ctx context.Context, u *User) (*User, error) {
	user, err := b.repo.FindByID(ctx, u.ID)
	if err != nil {
		return nil, err
	}

	if u.Name != "" {
		user.Name = u.Name
	}

	if u.Email != "" {
		user.Email = u.Email
	}

	return b.repo.Update(ctx, user)
}

// FindByID creates a User, and returns the new User.
func (b *UserBiz) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return b.repo.FindByID(ctx, id)
}

func (b *UserBiz) FindAll(ctx context.Context) ([]*User, error) {
	return b.repo.ListAll(ctx)
}

func (b *UserBiz) DeleteByID(ctx context.Context, id uuid.UUID) error {
	exist, err := b.repo.ExistByID(ctx, id)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf("user not found")
	}

	return b.repo.DeleteByID(ctx, id)
}
