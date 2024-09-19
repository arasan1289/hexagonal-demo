package http

import (
	"reflect"
	"strings"

	"github.com/arasan1289/hexagonal-demo/internal/adapters/config"
	"github.com/arasan1289/hexagonal-demo/internal/adapters/logger"
	"github.com/arasan1289/hexagonal-demo/internal/core/port"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Router is a wrapper for HTTP router
type Router struct {
	*gin.Engine
}

// NewRouter creates a new Router instance
func NewRouter(config *config.Container, log *logger.Logger, userHandler UserHandler, otpHandler OtpHandler, authService port.IAuthService, authhandler AuthHandler) (*Router, error) {
	// Disable debug mode in production
	if config.App.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	// CORS
	ginConfig := cors.DefaultConfig()
	ginConfig.AllowOrigins = config.HTTP.AllowedOrigins

	router := gin.New()

	// Bind json name as Field() in Validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}

	// Middlewares
	router.Use(GinStructuredLogger(log))
	router.Use(gin.Recovery(), cors.New(ginConfig))

	// Rate limiter
	rateLimit := NewRateLimiter(100, 60)

	// JWT authorization middleware
	authMiddleware := NewJWTAuthMiddleware(authService)

	v1 := router.Group("/api/v1")
	{
		user := v1.Group("/users")
		{
			user.POST("/", userHandler.Register)
			user.GET("/:id", authMiddleware, rateLimit, userHandler.GetUser)
		}
		v1.POST("/send-otp", rateLimit, otpHandler.RequestOtp)
		v1.POST("/verify-otp", rateLimit, otpHandler.VerifyOtp)
		v1.POST("/login", rateLimit, authhandler.Login)
	}
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &Router{
		router,
	}, nil
}

// Serve starts the HTTP server
func (r *Router) Serve(listenAddr string) error {
	return r.Run(listenAddr)
}
