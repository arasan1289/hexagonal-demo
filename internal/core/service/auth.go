package service

import (
	"context"
	"errors"
	"time"

	"github.com/arasan1289/hexagonal-demo/internal/adapters/config"
	"github.com/arasan1289/hexagonal-demo/internal/adapters/logger"
	"github.com/arasan1289/hexagonal-demo/internal/core/domain"
	"github.com/arasan1289/hexagonal-demo/internal/core/port"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	log    *logger.Logger
	config *config.App
}

// NewOtpService constructor function
func NewAuthService(log *logger.Logger, config *config.App) port.IAuthService {
	return &AuthService{
		log:    log,
		config: config,
	}
}

func (as *AuthService) GenerateJWT(ctx context.Context, user *domain.User) (*domain.JWTToken, error) {
	key := []byte(as.config.JWTSecret)
	atClaims := domain.UserClaims{
		Role: string(user.Role),
		Name: user.FirstName + " " + user.LastName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-server",
			Subject:   user.ID,
		},
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)
	rt := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		"iss": "auth-server",
		"sub": user.ID,
	})

	accessToken, err := at.SignedString(key)
	if err != nil {
		return nil, err
	}

	refreshToken, err := rt.SignedString(key)
	if err != nil {
		return nil, err
	}

	return &domain.JWTToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Role:         string(user.Role),
	}, nil
}

func (as *AuthService) VerifyJWT(ctx context.Context, accessToken string) (*domain.UserClaims, error) {
	t, err := jwt.ParseWithClaims(accessToken, &domain.UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(as.config.JWTSecret), nil
	}, jwt.WithLeeway(5*time.Second))

	claims, ok := t.Claims.(*domain.UserClaims)

	switch {
	case t.Valid:
		return claims, nil
	case errors.Is(err, jwt.ErrTokenMalformed), errors.Is(err, jwt.ErrTokenSignatureInvalid), !ok:
		return nil, domain.ErrInvalidToken
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		return nil, domain.ErrExpiredToken
	default:
		return nil, domain.ErrInvalidToken
	}
}

func (as *AuthService) RefreshJWT(ctx context.Context, refreshToken string, user *domain.User) (*domain.JWTToken, error) {
	t, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(as.config.JWTSecret), nil
	}, jwt.WithLeeway(5*time.Second))

	switch {
	case t.Valid:
		return as.GenerateJWT(ctx, user)
	case errors.Is(err, jwt.ErrTokenMalformed), errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return nil, domain.ErrInvalidToken
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		return nil, domain.ErrExpiredToken
	default:
		return nil, domain.ErrInvalidToken
	}
}
