package repository

import (
	"context"
	"errors"

	postgres "github.com/arasan1289/hexagonal-demo/internal/adapters/storage/db"
	"github.com/arasan1289/hexagonal-demo/internal/core/domain"
	"github.com/arasan1289/hexagonal-demo/internal/core/port"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UserRepository is an implementation of the port.UserRepository interface using a PostgreSQL database.
type UserRepository struct {
	db *postgres.Conn
}

// NewUserRepository creates a new instance of UserRepository with the provided database connection.
func NewUserRepository(conn *postgres.Conn) port.IUserRepository {
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
	var result *gorm.DB
	ur.db.Transaction(func(tx *gorm.DB) error {
		c, ok := ur.db.InstanceGet("config")
		if !ok {
			return errors.New("config not found")
		}
		tx = tx.InstanceSet("config", c)
		result = tx.First(&user, "id=?", id)
		return nil
	})
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByPhoneNumber retrieves a user from the database by their phone number hash or email hash.
func (ur *UserRepository) GetUserByPhoneNumberOrEmail(ctx context.Context, hash string) (*domain.User, error) {
	var user domain.User
	var result *gorm.DB
	ur.db.Transaction(func(tx *gorm.DB) error {
		c, ok := ur.db.InstanceGet("config")
		if !ok {
			return errors.New("config not found")
		}
		tx = tx.InstanceSet("config", c)
		result = tx.Where("phone_number_hash=?", hash).Or("email_hash=?", hash).Take(&user)
		return nil
	})
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
