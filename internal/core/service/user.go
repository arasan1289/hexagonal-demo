package service

import (
	"context"

	"github.com/arasan1289/hexagonal-demo/internal/adapters/config"
	"github.com/arasan1289/hexagonal-demo/internal/adapters/logger"
	"github.com/arasan1289/hexagonal-demo/internal/core/domain"
	"github.com/arasan1289/hexagonal-demo/internal/core/port"
	"github.com/arasan1289/hexagonal-demo/internal/core/util"
)

type UserService struct {
	repo port.UserRepository
	log  *logger.Logger
}

func NewUserService(repo port.UserRepository, log *logger.Logger) port.UserService {
	return &UserService{
		repo: repo,
		log:  log,
	}
}

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

func (us *UserService) GetUser(ctx context.Context, id string, config *config.App) (*domain.User, error) {
	usr, err := us.repo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	phoneNumber, err := util.DecryptString(usr.PhoneNumberEncrypted, config.SecretKey)
	if err != nil {
		return nil, err
	}
	usr.PhoneNumber = phoneNumber
	return usr, nil
}

func (us *UserService) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error) {
	us.log.Info().Msg("Get user by phone number")
	phoneNumberHash := util.HashString(phoneNumber)
	usr, err := us.repo.GetUserByPhoneNumber(ctx, phoneNumberHash)
	if err != nil {
		return nil, err
	}
	usr.PhoneNumber = phoneNumber
	return usr, nil
}
