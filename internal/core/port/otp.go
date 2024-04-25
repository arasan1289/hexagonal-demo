package port

import (
	"context"

	"github.com/arasan1289/hexagonal-demo/internal/core/domain"
)

type OtpService interface {
	GenerateOTP(ctx context.Context, length uint) (*domain.OTP, error)
	VerifyOTP(ctx context.Context, otp *domain.OTP) (bool, error)
}
