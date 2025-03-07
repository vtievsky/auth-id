package services

import (
	"context"

	usersvc "github.com/vtievsky/auth-id/internal/services/users"
)

type SvcLayer struct {
	UserSvc UserService
}

type UserService interface {
	GetUser(ctx context.Context, login string) (*usersvc.User, error)
	GetUsers(ctx context.Context) ([]*usersvc.User, error)
	CreateUser(ctx context.Context, user usersvc.UserCreated) (*usersvc.User, error)
	UpdateUser(ctx context.Context, user usersvc.UserUpdated) (*usersvc.User, error)
	DeleteUser(ctx context.Context, login string) error
}
