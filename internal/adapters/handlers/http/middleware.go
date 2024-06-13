package http

import (
	"strings"
	"time"

	"github.com/arasan1289/hexagonal-demo/internal/adapters/logger"
	"github.com/arasan1289/hexagonal-demo/internal/core/domain"
	"github.com/arasan1289/hexagonal-demo/internal/core/port"
	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
)

// NewRateLimiter creates a new rate limiter middleware using the token bucket strategy
func NewRateLimiter(capacity float64, fillInterval time.Duration) gin.HandlerFunc {

	limiter := tollbooth.NewLimiter(capacity, &limiter.ExpirableOptions{
		DefaultExpirationTTL: fillInterval,
	})

	return func(ctx *gin.Context) {
		err := tollbooth.LimitByRequest(limiter, ctx.Writer, ctx.Request)
		if err != nil {
			handleError(ctx, domain.ErrRateLimitExceeded)
			ctx.Abort()
		} else {
			ctx.Next()
		}

	}
}

// StructuredLogger logs a gin HTTP request in JSON format
func GinStructuredLogger(logger *logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := xid.New().String()
		ctx.Set("request_id", requestID)
		logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("request_id", requestID)
		})
		ctx.Header("x-request-id", requestID)

		start := time.Now() // Start timer
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery

		// Process request
		ctx.Next()

		// Fill the params
		param := gin.LogFormatterParams{}

		param.TimeStamp = time.Now() // Stop timer
		param.Latency = param.TimeStamp.Sub(start)
		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}

		param.ClientIP = ctx.ClientIP()
		param.Method = ctx.Request.Method
		param.StatusCode = ctx.Writer.Status()
		param.ErrorMessage = strings.Join(ctx.Errors.Errors(), ",")
		param.BodySize = ctx.Writer.Size()
		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path

		// Log using the params
		var logEvent *zerolog.Event
		if ctx.Writer.Status() >= 400 {
			logEvent = logger.Error().Str("error_message", param.ErrorMessage)
		} else {
			logEvent = logger.Info()
		}

		userClaims, ok := ctx.Value("user").(*domain.UserClaims)
		if !ok || userClaims == nil {
			logEvent.Str("user_id", "")
		} else {
			logEvent.Str("user_id", userClaims.RegisteredClaims.Subject)
		}

		logEvent.Str("client_id", param.ClientIP).
			Str("user_agent", ctx.Request.UserAgent()).
			Str("method", param.Method).
			Int("status_code", param.StatusCode).
			Int("body_size", param.BodySize).
			Str("path", param.Path).
			Dur("latency", param.Latency).
			Dur("elapsed_ms", time.Since(start)).
			Msg("Request completed.")
	}
}

// JWTAuthMiddleware is the middleware for JWT authentication
func NewJWTAuthMiddleware(as port.IAuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			handleError(ctx, domain.ErrEmptyAuthorizationHeader)
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			handleError(ctx, domain.ErrInvalidAuthorizationHeader)
			ctx.Abort()
			return
		}

		tokenStr := parts[1]
		userClaims, err := as.VerifyJWT(ctx, tokenStr)
		if err != nil {
			handleError(ctx, err)
			ctx.Abort()
			return
		}

		// Add user information to the context
		ctx.Set("user", userClaims)
		ctx.Next()
	}
}
