package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// TracingMiddleware adds OpenTelemetry tracing to requests
func TracingMiddleware() gin.HandlerFunc {
	tracer := otel.Tracer("g-management")

	return func(c *gin.Context) {
		// Start a span for this request
		ctx, span := tracer.Start(c.Request.Context(), c.Request.URL.Path,
			trace.WithAttributes(
				attribute.String("http.method", c.Request.Method),
				attribute.String("http.url", c.Request.URL.String()),
				attribute.String("http.client_ip", c.ClientIP()),
			),
		)
		defer span.End()

		// Add the span context to the request context
		c.Request = c.Request.WithContext(ctx)

		// Call the next handler
		c.Next()

		// Add response attributes to the span
		span.SetAttributes(
			attribute.Int("http.status_code", c.Writer.Status()),
		)

		// If there were any errors during handling, add them to the span
		for _, err := range c.Errors {
			span.RecordError(err.Err)
		}
	}
}
