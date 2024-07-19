package http

import (
	"github.com/arasan1289/hexagonal-demo/internal/adapters/config"
	"github.com/arasan1289/hexagonal-demo/internal/adapters/logger"
	"github.com/arasan1289/hexagonal-demo/internal/core/domain"
	"github.com/arasan1289/hexagonal-demo/internal/core/port"
	"github.com/gin-gonic/gin"
)

type OtpHandler struct {
	svc     port.IOtpService
	userSvc port.IUserService
	authSvc port.IAuthService
	log     *logger.Logger
	config  *config.App
}

// NewOtpHandler creates a new instance of OtpHandler
func NewOtpHandler(svc port.IOtpService, userSvc port.IUserService, authSvc port.IAuthService, log *logger.Logger, config *config.App) *OtpHandler {
	return &OtpHandler{
		svc, userSvc, authSvc, log, config,
	}
}

// requestOtp is the request body for the request otp endpoint
type requestOtp struct {
	PhoneNumber string `json:"phone_number" binding:"required,min=10" example:"9876543210"`
}

//	@Summary		Request OTP
//	@Description	Sends OTP to the number if its registered
//	@Tags			Auth
//	@Produce		json
//	@Accept			json
//	@Param			sendOTP	body		requestOtp	true	"Request OTP JSON"
//	@Success		200		{object}	response{data=domain.OTP}
//	@Failure		400		{object}	response
//	@Failure		500		{object}	response
//	@Router			/send-otp [post]
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

// verifyOtp is the request body for the verify otp endpoint
type verifyOtp struct {
	Otp         string `json:"otp" binding:"required,min=6" example:"123456"`
	OtpHash     string `json:"otp_hash" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

//	@Summary		Verify OTP
//	@Description	Verify the given OTP
//	@Tags			Auth
//	@Produce		json
//	@Accept			json
//	@Param			verifyOTP	body		verifyOtp	true	"Verify OTP JSON"
//	@Success		200			{object}	response{data=domain.JWTToken}
//	@Failure		400			{object}	response
//	@Failure		500			{object}	response
//	@Router			/verify-otp [post]
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

	if rsp && !user.IsPhoneNumberVerified {
		user.IsPhoneNumberVerified = true
		user.IsActive = true
		user, err = oh.userSvc.Register(ctx, user, oh.config)
		if err != nil {
			handleError(ctx, err)
			return
		}
	}

	tokens, err := oh.authSvc.GenerateJWT(ctx, user)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, tokens)
}
