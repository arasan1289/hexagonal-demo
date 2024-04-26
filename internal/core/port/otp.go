package port

import (
	"context"

	"github.com/arasan1289/hexagonal-demo/internal/core/domain"
)

// OtpService defines the interface for OTP related operations
type OtpService interface {
	// GenerateOTP generates a new OTP with the given length
	GenerateOTP(ctx context.Context, length uint) (*domain.OTP, error)

	// VerifyOTP verifies if the given OTP is valid
	VerifyOTP(ctx context.Context, otp *domain.OTP) (bool, error)
}
