package port

import (
	"context"

	"github.com/arasan1289/hexagonal-demo/internal/core/domain"
)

// IAuthService defines the interface for authentication and authorizations related operations
type IAuthService interface {
	GenerateJWT(ctx context.Context, user *domain.User) (*domain.JWTToken, error)
	VerifyJWT(ctx context.Context, accessToken string) (*domain.UserClaims, error)
	RefreshJWT(ctx context.Context, refreshToken string, user *domain.User) (*domain.JWTToken, error)
}
