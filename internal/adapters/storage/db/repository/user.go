package repository

import (
	"context"

	postgres "github.com/arasan1289/hexagonal-demo/internal/adapters/storage/db"
	"github.com/arasan1289/hexagonal-demo/internal/core/domain"
	"github.com/arasan1289/hexagonal-demo/internal/core/port"
	"gorm.io/gorm/clause"
)

// UserRepository is an implementation of the port.UserRepository interface using a PostgreSQL database.
type UserRepository struct {
	db *postgres.Conn
}

// NewUserRepository creates a new instance of UserRepository with the provided database connection.
func NewUserRepository(conn *postgres.Conn) port.UserRepository {
	return &UserRepository{
		db: conn,
	}
}

// UpsertUser creates or updates a user in the database.
// It uses the OnConflict clause to update all fields if the user already exists.
func (ur *UserRepository) UpsertUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	data := ur.db.Clauses(clause.OnConflict{UpdateAll: true}, clause.Returning{}).Create(user)
	if data.Error != nil {
		return nil, data.Error
	}
	return user, nil
}

// GetUser retrieves a user from the database by their ID.
func (ur *UserRepository) GetUser(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	result := ur.db.First(&user, "id=?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByPhoneNumber retrieves a user from the database by their phone number hash.
func (ur *UserRepository) GetUserByPhoneNumber(ctx context.Context, phoneNumberHash string) (*domain.User, error) {
	var user domain.User
	result := ur.db.First(&user, "phone_number_hash=?", phoneNumberHash)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
