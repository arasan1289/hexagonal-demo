package http

import (
	"strings"
	"time"

	"github.com/arasan1289/hexagonal-demo/internal/adapters/logger"
	"github.com/arasan1289/hexagonal-demo/internal/core/domain"
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
	return func(c *gin.Context) {
		requestID := xid.New().String()
		c.Set("request_id", requestID)
		logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("request_id", requestID)
		})
		c.Header("x-request-id", requestID)

		start := time.Now() // Start timer
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Fill the params
		param := gin.LogFormatterParams{}

		param.TimeStamp = time.Now() // Stop timer
		param.Latency = param.TimeStamp.Sub(start)
		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}

		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = strings.Join(c.Errors.Errors(), ",")
		param.BodySize = c.Writer.Size()
		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path

		// Log using the params
		var logEvent *zerolog.Event
		if c.Writer.Status() >= 400 {
			logEvent = logger.Error().Str("error_message", param.ErrorMessage)
		} else {
			logEvent = logger.Info()
		}

		logEvent.Str("client_id", param.ClientIP).
			Str("user_agent", c.Request.UserAgent()).
			Str("method", param.Method).
			Int("status_code", param.StatusCode).
			Int("body_size", param.BodySize).
			Str("path", param.Path).
			Dur("latency", param.Latency).
			Dur("elapsed_ms", time.Since(start)).
			Msg("request completed with")
	}
}
