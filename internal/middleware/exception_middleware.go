package middleware

import (
	"fmt"
	"net/http"

	"go_di_architecture/internal/domain/models/response"

	"github.com/gin-gonic/gin"
)

// ExceptionHandler captures and handles unhandled exceptions.
//
// This middleware handler:
//   - Catches panics and unhandled errors
//   - Creates standardized error responses
//   - Logs errors with request context
//   - Prevents stack traces in production
//
// The error response follows the same structure as all other API responses.
//
// Returns:
//   - gin.HandlerFunc: A middleware handler function
func ExceptionHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := ctx.GetString("request_id")

		defer func() {
			if err := recover(); err != nil {
				// Log the error
				fmt.Printf("[ERROR] [%s] Unhandled panic: %v\n", requestID, err)

				// Create standardized error response
				response := response.NewErrorResponse(
					"INTERNAL_ERROR",
					response.StatusToMessage(http.StatusInternalServerError),
					nil,
					requestID,
				)

				// Return error response
				ctx.JSON(http.StatusInternalServerError, response)
				ctx.Abort()
			}
		}()

		// Continue processing the request
		ctx.Next()

		// Handle errors from controllers
		if len(ctx.Errors) > 0 {
			err := ctx.Errors[0]
			handleError(ctx, err.Err, requestID)
		}
	}
}

// handleError processes errors into standardized responses.
func handleError(ctx *gin.Context, err error, requestID string) {
	statusCode := ctx.Writer.Status()
	if statusCode == http.StatusOK {
		statusCode = http.StatusInternalServerError
	}

	code := "INTERNAL_ERROR"
	message := response.StatusToMessage(statusCode)

	ctx.JSON(statusCode, response.NewErrorResponse(
		code,
		message,
		map[string][]string{"error": {err.Error()}},
		requestID,
	))
}
