package repository

import (
	"context"

	postgres "github.com/arasan1289/hexagonal-demo/internal/adapters/storage/db"
	"github.com/arasan1289/hexagonal-demo/internal/core/domain"
	"github.com/arasan1289/hexagonal-demo/internal/core/port"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	db *postgres.Conn
}

func NewUserRepository(conn *postgres.Conn) port.UserRepository {
	return &UserRepository{
		db: conn,
	}
}

func (ur *UserRepository) UpsertUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	data := ur.db.Clauses(clause.OnConflict{UpdateAll: true}, clause.Returning{}).Create(user)
	if data.Error != nil {
		return nil, data.Error
	}
	return user, nil
}

func (ur *UserRepository) GetUser(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	result := ur.db.First(&user, "id=?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (ur *UserRepository) GetUserByPhoneNumber(ctx context.Context, phoneNumberHash string) (*domain.User, error) {
	var user domain.User
	result := ur.db.First(&user, "phone_number_hash=?", phoneNumberHash)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
