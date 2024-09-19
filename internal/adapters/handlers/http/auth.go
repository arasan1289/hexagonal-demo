package http

import (
	"github.com/gin-gonic/gin"

	"github.com/arasan1289/hexagonal-demo/internal/adapters/logger"
	"github.com/arasan1289/hexagonal-demo/internal/core/domain"
	"github.com/arasan1289/hexagonal-demo/internal/core/port"
)

type AuthHandler struct {
	authSvc port.IAuthService // auth service
	userSvc port.IUserService // user service
	log     *logger.Logger    // logger
}

func NewAuthHandler(authSvc port.IAuthService, userSvc port.IUserService, log *logger.Logger) *AuthHandler {
	return &AuthHandler{
		authSvc: authSvc,
		userSvc: userSvc,
		log:     log,
	}
}

type loginUser struct {
	PhoneNumber string `json:"phone_number,omitempty" binding:"required_without=Email,omitempty,min=10" example:"9876543210"`
	Email       string `json:"email,omitempty" binding:"required_without=PhoneNumber,omitempty,email" example:"example@example.com"`
	Password    string `json:"password" binding:"required" example:"password"`
}

// @Summary			Login
// @Description		Login user by either of email or phone and password
// @Tags			Auth
// @Produce			json
// @Accept			json
// @Param			login	body		loginUser	true	"Login User JSON"
// @Success			200		{object}	response{data=domain.JWTToken}
// @Failure			400		{object}	response
// @Failure			500		{object}	response
// @Router			/login [post]
func (ah *AuthHandler) Login(ctx *gin.Context) {
	var req loginUser

	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}
	var emailOrPhone string
	if req.PhoneNumber != "" {
		emailOrPhone = req.PhoneNumber
	} else if req.Email != "" {
		emailOrPhone = req.Email
	} else {
		handleError(ctx, domain.ErrInvalidCredentials)
		return
	}

	user, match, err := ah.userSvc.GetUserAndComparePassword(ctx, emailOrPhone, req.Password)
	if err != nil {
		handleError(ctx, err)
		return
	}
	if !match {
		handleError(ctx, domain.ErrInvalidCredentials)
		return
	}
	tokens, err := ah.authSvc.GenerateJWT(ctx, user)
	if err != nil {
		handleError(ctx, err)
		return
	}
	handleSuccess(ctx, tokens)

}
