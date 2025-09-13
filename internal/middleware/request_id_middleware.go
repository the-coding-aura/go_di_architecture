package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDHandler generates and manages request IDs for tracing.
//
// This middleware handler:
//   - Generates a unique ID for each request if none is provided
//   - Propagates the request ID through the entire request lifecycle
//   - Includes the ID in all logs and responses
//   - Supports incoming X-Request-Id header for distributed tracing
//
// The request ID is used for:
//   - Correlating logs across services
//   - Debugging specific requests
//   - Providing consistent error responses with traceability
//
// Returns:
//   - gin.HandlerFunc: A middleware handler function
func RequestIDHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get or generate request ID
		requestID := c.GetHeader("X-Request-Id")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Set request ID in context
		c.Set("request_id", requestID)

		// Set request ID in response header
		c.Header("X-Request-Id", requestID)

		// Process request
		c.Next()
	}
}
