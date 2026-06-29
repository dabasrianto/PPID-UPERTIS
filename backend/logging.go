package main

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// Global structured logger
var zlog zerolog.Logger

// Request ID context key
const RequestIDKey = "request_id"

func initLogger() {
	// JSON output for production (PM2 captures stdout)
	zerolog.TimeFieldFormat = time.RFC3339
	zlog = zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "upertis-backend").
		Logger()

	zlog.Info().Msg("Structured logger initialized")
}

// requestIDMiddleware generates a unique request_id for every request
// and stores it in Locals for downstream use.
func requestIDMiddleware(c *fiber.Ctx) error {
	rid := c.Get("X-Request-ID")
	if rid == "" {
		rid = uuid.New().String()
	}
	c.Locals(RequestIDKey, rid)
	c.Set("X-Request-ID", rid)
	return c.Next()
}

// structuredLoggerMiddleware logs every request in JSON format.
// Sensitive data (Authorization header, cookies) is never logged.
func structuredLoggerMiddleware(c *fiber.Ctx) error {
	start := time.Now()

	// Process request
	err := c.Next()

	// Calculate latency
	latency := time.Since(start)
	status := c.Response().StatusCode()

	// Get request_id
	rid, _ := c.Locals(RequestIDKey).(string)

	// Get user_id if authenticated
	userID := ""
	if uid := c.Locals("user_id"); uid != nil {
		if s, ok := uid.(string); ok {
			userID = s
		}
	}

	// Build log event
	event := zlog.Info()
	if status >= 500 {
		event = zlog.Error()
	} else if status >= 400 {
		event = zlog.Warn()
	}

	event.
		Str("request_id", rid).
		Str("method", c.Method()).
		Str("path", c.Path()).
		Int("status", status).
		Dur("latency_ms", latency).
		Str("client_ip", c.IP()).
		Str("user_agent", truncate(c.Get("User-Agent"), 100)).
		Int("body_size", len(c.Response().Body()))

	if userID != "" {
		event.Str("user_id", userID)
	}

	if err != nil {
		event.Str("error", err.Error())
	}

	event.Msg("request")

	return err
}

// panicRecoverMiddleware catches panics and logs stack trace internally
// without exposing it to the client.
func panicRecoverMiddleware(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			rid, _ := c.Locals(RequestIDKey).(string)

			// Log panic with stack trace INTERNAL ONLY
			zlog.Error().
				Str("request_id", rid).
				Str("method", c.Method()).
				Str("path", c.Path()).
				Interface("panic", r).
				Msg("PANIC RECOVERED")

			// Generic response to client — no internal details
			_ = c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":      "An internal error occurred",
				"request_id": rid,
			})
		}
	}()
	return c.Next()
}

// truncate limits string length for safe logging
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
