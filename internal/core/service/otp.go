package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/arasan1289/hexagonal-demo/internal/adapters/config"
	"github.com/arasan1289/hexagonal-demo/internal/adapters/logger"
	"github.com/arasan1289/hexagonal-demo/internal/core/domain"
	"github.com/arasan1289/hexagonal-demo/internal/core/port"
	"github.com/arasan1289/hexagonal-demo/internal/core/util"
)

// OtpService struct contains a logger and configuration object
type OtpService struct {
	log    *logger.Logger
	config *config.App
}

// NewOtpService constructor function
func NewOtpService(log *logger.Logger, config *config.App) port.OtpService {
	return &OtpService{
		log:    log,
		config: config,
	}
}

// GenerateOTP generates a new OTP with the given length
func (os *OtpService) GenerateOTP(ctx context.Context, length uint) (*domain.OTP, error) {
	// Create a byte slice of the given length
	b := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return nil, err
	}

	// Define a character set for the OTP
	const charset = "1234567890"
	charsetLen := len(charset)

	// Convert the random bytes to digits using the character set
	for i := range b {
		b[i] = charset[int(b[i])%charsetLen]
	}

	otp := domain.OTP{
		Otp: string(b),
	}
	// Encrypt the OTP and add the expiry time to the hash
	expiry := time.Now().Add(time.Minute * 10).UnixMilli()
	otpEnc, err := util.EncryptString(otp.Otp, os.config.OtpSecretKey)
	if err != nil {
		return nil, err
	}
	otp.OtpHash = fmt.Sprintf("%d.%s", expiry, otpEnc)

	return &otp, nil
}

// VerifyOTP verifies the given OTP
func (os *OtpService) VerifyOTP(ctx context.Context, otp *domain.OTP) (bool, error) {
	// Split the OTP hash into the expiry time and the encrypted OTP
	exp, otpEnc, _ := strings.Cut(otp.OtpHash, ".")
	expiry, err := strconv.ParseInt(exp, 10, 64)
	if err != nil {
		return false, err
	}
	// Check if the OTP has expired
	currTime := time.Now().UnixMilli()
	if currTime >= expiry {
		return false, domain.ErrOTPExpired
	}
	// Decrypt the OTP and compare it to the given OTP
	otpDec, err := util.DecryptString(otpEnc, os.config.OtpSecretKey)
	if err != nil {
		return false, err
	}
	if otpDec != otp.Otp {
		return false, domain.ErrOTPMismatch
	}
	return true, nil
}
