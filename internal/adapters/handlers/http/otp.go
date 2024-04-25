package http

import (
	"github.com/arasan1289/hexagonal-demo/internal/adapters/config"
	"github.com/arasan1289/hexagonal-demo/internal/adapters/logger"
	"github.com/arasan1289/hexagonal-demo/internal/core/domain"
	"github.com/arasan1289/hexagonal-demo/internal/core/port"
	"github.com/gin-gonic/gin"
)

type OtpHandler struct {
	svc     port.OtpService
	userSvc port.UserService
	log     *logger.Logger
	config  *config.App
}

func NewOtpHandler(svc port.OtpService, userSvc port.UserService, log *logger.Logger, config *config.App) *OtpHandler {
	return &OtpHandler{
		svc, userSvc, log, config,
	}
}

type requestOtp struct {
	PhoneNumber string `json:"phone_number" binding:"required,min=10" example:"9876543210"`
}

func (oh *OtpHandler) RequestOtp(ctx *gin.Context) {
	var req requestOtp

	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}
	_, err := oh.userSvc.GetUserByPhoneNumber(ctx, req.PhoneNumber)
	if err != nil {
		handleError(ctx, err)
		return
	}
	rsp, err := oh.svc.GenerateOTP(ctx, oh.config.OtpLength)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, rsp)
}

type verifyOtp struct {
	Otp         string `json:"otp" binding:"required,min=6" example:"123456"`
	OtpHash     string `json:"otp_hash" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

func (oh *OtpHandler) VerifyOtp(ctx *gin.Context) {
	var req verifyOtp

	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	otp := domain.OTP{
		Otp:     req.Otp,
		OtpHash: req.OtpHash,
	}

	rsp, err := oh.svc.VerifyOTP(ctx, &otp)
	if err != nil {
		handleError(ctx, err)
		return
	}

	user, err := oh.userSvc.GetUserByPhoneNumber(ctx, req.PhoneNumber)
	if err != nil {
		handleError(ctx, err)
		return
	}

	if rsp {
		user.IsPhoneNumberVerified = true
		user, err = oh.userSvc.Register(ctx, user, oh.config)
		if err != nil {
			handleError(ctx, err)
			return
		}
	}

	handleSuccess(ctx, user)
}
