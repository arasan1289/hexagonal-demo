package port

import (
	"context"

	"github.com/arasan1289/hexagonal-demo/internal/adapters/config"
	"github.com/arasan1289/hexagonal-demo/internal/core/domain"
)

// IUserRepository interface defines the methods for interacting with the user repository
type IUserRepository interface {
	// UpsertUser inserts or updates a user in the repository
	UpsertUser(ctx context.Context, user *domain.User) (*domain.User, error)

	// GetUser retrieves a user from the repository by ID
	GetUser(ctx context.Context, id string) (*domain.User, error)

	// GetUserByPhoneNumber retrieves a user from the repository by phone number hash
	GetUserByPhoneNumber(ctx context.Context, hash string) (*domain.User, error)
}

// UserService interface defines the methods for interacting with the user service
type IUserService interface {
	// Register creates a new user and saves it to the repository
	Register(ctx context.Context, user *domain.User, conf *config.App) (*domain.User, error)

	// GetUser retrieves a user from the repository by ID
	GetUser(ctx context.Context, id string, conf *config.App) (*domain.User, error)

	// GetUserByPhoneNumber retrieves a user from the repository by phone number
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error)
}
