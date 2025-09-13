package response

import (
	"net/http"
	"time"
)

// APIResponse represents the standardized response structure for all API endpoints.
//
// This generic response wrapper follows the pattern used in the .NET documentation example,
// with a consistent structure that allows for different data types in the Data field.
//
// The structure supports Swagger documentation through:
//   - Clear separation of success/error states
//   - Proper documentation of the generic Data field
//   - Example responses for different scenarios
//
// Example Success Response (200):
//
//	{
//	  "success": true,
//	  "message": "Module retrieved successfully",
//	  "data": {
//	    "id": 123,
//	    "name": "Inventory",
//	    "description": "Handles product stock management",
//	    "isActive": true,
//	    "createdAt": "2023-08-15T14:30:00Z"
//	  },
//	  "meta": {
//	    "requestId": "a1b2c3d4",
//	    "timestamp": "2023-08-15T14:30:00Z"
//	  }
//	}
//
// Example Error Response (400):
//
//	{
//	  "success": false,
//	  "message": "Validation error",
//	  "error": {
//	    "code": "VALIDATION_ERROR",
//	    "message": "Name must be 3-50 characters"
//	  },
//	  "meta": {
//	    "requestId": "a1b2c3d4",
//	    "timestamp": "2023-08-15T14:30:00Z"
//	  }
//	}
type APIResponse struct {
	// Indicates if the request was successful
	Success bool `json:"success"`

	// Brief message about the result of the operation
	Message string `json:"message"`

	// The actual data payload (only present on success)
	// swagger:allOf
	Data interface{} `json:"data,omitempty"`

	// Error details (only present on failure)
	Error *APIError `json:"error,omitempty"`

	// Additional metadata about the response
	Meta ResponseMeta `json:"meta"`
}

// APIError represents standardized error information.
type APIError struct {
	// Machine-readable error code
	Code string `json:"code"`

	// Human-readable error message
	Message string `json:"message"`

	// Field-specific validation errors
	Details map[string][]string `json:"details,omitempty"`
}

// ResponseMeta contains additional metadata about the response.
type ResponseMeta struct {
	// Unique identifier for the request (for tracing)
	RequestId string `json:"requestId"`

	// Timestamp when the request was processed
	Timestamp string `json:"timestamp"`
}

// ResponseMapper provides methods to create standardized API responses.
//
// This mapper implements the same pattern as the .NET example, where:
//   - Controllers use the mapper to create consistent responses
//   - The Data field can contain any specific response type
//   - Error responses follow the same structure as success responses
type ResponseMapper struct {
	requestID string
}

// NewResponseMapper creates a new response mapper with the request ID.
//
// Parameters:
//   - requestID: The unique identifier for the current request
//
// Returns:
//   - *ResponseMapper: A configured response mapper
func NewResponseMapper(requestID string) *ResponseMapper {
	return &ResponseMapper{requestID: requestID}
}

// Success creates a standardized success response.
//
// Parameters:
//   - data: The actual data payload to return
//   - message: Brief success message
//   - statusCode: HTTP status code for the response
//
// Returns:
//   - *APIResponse: A properly formatted success response
//   - int: The HTTP status code
func (m *ResponseMapper) Success(data interface{}, message string, statusCode int) (*APIResponse, int) {
	return &APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta: ResponseMeta{
			RequestId: m.requestID,
			Timestamp: time.Now().Format(time.RFC3339),
		},
	}, statusCode
}

// Error creates a standardized error response.
//
// Parameters:
//   - code: Machine-readable error code
//   - message: Human-readable error message
//   - details: Field-specific validation errors
//   - statusCode: HTTP status code for the response
//
// Returns:
//   - *APIResponse: A properly formatted error response
//   - int: The HTTP status code
func (m *ResponseMapper) Error(code, message string, details map[string][]string, statusCode int) (*APIResponse, int) {
	return &APIResponse{
		Success: false,
		Message: message,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
		Meta: ResponseMeta{
			RequestId: m.requestID,
			Timestamp: time.Now().Format(time.RFC3339),
		},
	}, statusCode
}

// NewSuccessResponse creates a standardized success response.
//
// Parameters:
//   - data: The actual data payload to return
//   - message: Brief success message
//   - requestId: Unique identifier for the request
//
// Returns:
//   - *APIResponse: A properly formatted success response
func NewSuccessResponse(data interface{}, message string, requestId string) *APIResponse {
	return &APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta: ResponseMeta{
			RequestId: requestId,
			Timestamp: time.Now().Format(time.RFC3339),
		},
	}
}

// NewErrorResponse creates a standardized error response.
//
// Parameters:
//   - code: Machine-readable error code
//   - message: Human-readable error message
//   - details: Field-specific validation errors
//   - requestId: Unique identifier for the request
//
// Returns:
//   - *APIResponse: A properly formatted error response
func NewErrorResponse(code, message string, details map[string][]string, requestId string) *APIResponse {
	return &APIResponse{
		Success: false,
		Message: message,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
		Meta: ResponseMeta{
			RequestId: requestId,
			Timestamp: time.Now().Format(time.RFC3339),
		},
	}
}

// StatusToMessage maps HTTP status codes to standard messages.
//
// Parameters:
//   - statusCode: HTTP status code
//
// Returns:
//   - string: Standardized message for the status code
func StatusToMessage(statusCode int) string {
	switch statusCode {
	case http.StatusOK:
		return "Operation completed successfully"
	case http.StatusCreated:
		return "Resource created successfully"
	case http.StatusBadRequest:
		return "Invalid request parameters"
	case http.StatusNotFound:
		return "Resource not found"
	case http.StatusConflict:
		return "Resource already exists"
	default:
		return "An unexpected error occurred"
	}
}
