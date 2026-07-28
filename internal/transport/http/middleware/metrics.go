package middleware

import (
	"strings"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func MetricsMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)

		method := strings.Clone(c.Method())
		route := strings.Clone(c.Route().Path)
		status := strings.Clone(strconv.Itoa(c.Response().StatusCode()))

		metrics.HttpRequestsTotal.
			WithLabelValues(method, route, status).
			Inc()

		metrics.HttpRequestDurationSeconds.
			WithLabelValues(method, route).
			Observe(duration.Seconds())

		if c.Response().StatusCode() >= 500 {
			metrics.HttpErrorsTotal.
				WithLabelValues(method, route, status).
				Inc()
		}

		return err
	}

}
