package service

import (
	"context"

	"github.com/arasan1289/hexagonal-demo/internal/adapters/config"
	"github.com/arasan1289/hexagonal-demo/internal/adapters/logger"
	"github.com/arasan1289/hexagonal-demo/internal/core/domain"
	"github.com/arasan1289/hexagonal-demo/internal/core/port"
	"github.com/arasan1289/hexagonal-demo/internal/core/util"
)

// UserService struct represents the user service with its dependencies
type UserService struct {
	repo port.IUserRepository // user repository interface
	log  *logger.Logger       // logger instance
}

// NewUserService constructor function
func NewUserService(repo port.IUserRepository, log *logger.Logger) port.IUserService {
	return &UserService{
		repo: repo,
		log:  log,
	}
}

// Register function: encrypt phone number, generate ID if needed, upsert user
func (us *UserService) Register(ctx context.Context, user *domain.User, config *config.App) (*domain.User, error) {
	phoneNumberEnc, err := util.EncryptString(user.PhoneNumber, config.SecretKey)
	if err != nil {
		return nil, err
	}
	if user.ID == "" {
		user.ID = util.GenerateULID()
	}

	user.PhoneNumberHash = util.HashString(user.PhoneNumber)
	user.PhoneNumberEncrypted = phoneNumberEnc

	usr, err := us.repo.UpsertUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

// GetUser function: retrieve user by ID, decrypt phone number
func (us *UserService) GetUser(ctx context.Context, id string, config *config.App) (*domain.User, error) {
	usr, err := us.repo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

// GetUserByPhoneNumber function: retrieve user by phone number hash
func (us *UserService) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error) {
	us.log.Info().Msg("Get user by phone number")
	phoneNumberHash := util.HashString(phoneNumber)
	usr, err := us.repo.GetUserByPhoneNumber(ctx, phoneNumberHash)
	if err != nil {
		return nil, err
	}
	return usr, nil
}
