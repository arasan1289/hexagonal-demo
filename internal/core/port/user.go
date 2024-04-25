package port

import (
	"context"

	"github.com/arasan1289/hexagonal-demo/internal/adapters/config"
	"github.com/arasan1289/hexagonal-demo/internal/core/domain"
)

type UserRepository interface {
	UpsertUser(ctx context.Context, user *domain.User) (*domain.User, error)
	// UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUser(ctx context.Context, id string) (*domain.User, error)
	GetUserByPhoneNumber(ctx context.Context, hash string) (*domain.User, error)
	// ListUsers(ctx context.Context) ([]domain.User, error)
	// DeleteUser(ctx context.Context, id string) (bool, error)
}

type UserService interface {
	Register(ctx context.Context, user *domain.User, conf *config.App) (*domain.User, error)
	GetUser(ctx context.Context, id string, conf *config.App) (*domain.User, error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error)
	// ListUsers(ctx context.Context, skip, limit uint64) ([]domain.User, error)
	// UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	// DeleteUser(ctx context.Context, id string) (bool, error)
}
