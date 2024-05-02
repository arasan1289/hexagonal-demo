package http

import (
	"github.com/arasan1289/hexagonal-demo/internal/adapters/config"
	"github.com/arasan1289/hexagonal-demo/internal/adapters/logger"
	"github.com/arasan1289/hexagonal-demo/internal/core/domain"
	"github.com/arasan1289/hexagonal-demo/internal/core/port"
	"github.com/gin-gonic/gin"
)

// UserHandler handles HTTP requests related to user management
type UserHandler struct {
	svc    port.IUserService // user service
	config *config.App       // app configuration
	log    *logger.Logger    // logger
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(svc port.IUserService, config *config.App, log *logger.Logger) *UserHandler {
	return &UserHandler{
		svc:    svc,
		config: config,
		log:    log,
	}
}

// registerUser represents the request body for the Register endpoint
type registerUser struct {
	PhoneNumber string `json:"phone_number" binding:"required,min=10" example:"9876543210"`
	FirstName   string `json:"first_name" binding:"required,min=5" example:"Qwerty"`
	LastName    string `json:"last_name" binding:"required,min=1" example:"A"`
}

// Register handles the registration of a new user
func (uh *UserHandler) Register(ctx *gin.Context) {
	var req registerUser
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	user := domain.User{
		PhoneNumber: req.PhoneNumber,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Role:        domain.Admin,
	}

	rsp, err := uh.svc.Register(ctx, &user, uh.config)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, rsp)
}

// getUserRequest represents the request parameters for the GetUser endpoint
type getUserRequest struct {
	ID string `uri:"id" binding:"required,ulid"`
}

// GetUser handles the retrieval of a user by their ID
func (uh *UserHandler) GetUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	rsp, err := uh.svc.GetUser(ctx, req.ID, uh.config)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, rsp)
}
