package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"go_di_architecture/internal/domain/models/module"
	"go_di_architecture/internal/domain/models/response"
	moduleService "go_di_architecture/internal/domain/service/module"
	moduleRepo "go_di_architecture/internal/infra/db/module"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ModuleHandler handles HTTP requests for module entities.
//
// This handler implements the same pattern as the .NET example, using:
//   - A response mapper to create consistent API responses
//   - Clear separation of success/error handling
//   - Proper HTTP status code semantics
//   - Standardized response structure
//
// The handler follows the requirement to use "controllers as handlers" while
// maintaining the generic response structure that supports Swagger documentation.
//
// All responses follow the APIResponse structure with:
//   - success: Boolean indicating operation success
//   - message: Brief operation result message
//   - data: Actual payload (on success, with specific type per endpoint)
//   - error: Error details (on failure)
//   - meta: Additional metadata (request ID, timestamp)
type ModuleHandler struct {
	service *moduleService.ModuleService
}

// NewModuleHandler creates a new instance of ModuleHandler.
//
// Returns:
//   - *ModuleHandler: A new handler instance
func NewModuleHandler() *ModuleHandler {
	repo := moduleRepo.NewModuleRepository()
	service := moduleService.NewModuleService(repo)
	return &ModuleHandler{service: service}
}

// CreateModule godoc
// @Summary Create a new module
// @Description Creates a new module entity with the provided details
// @Tags modules
// @Accept json
// @Produce json
// @Param request body module.ModuleRequest true "Module creation payload"
// @Success 201 {object} response.APIResponse{data=module.ModuleResponse} "Module created successfully"
// @Failure 400 {object} response.APIResponse "Validation error"
// @Failure 409 {object} response.APIResponse "Module name already exists"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /modules [post]
//
// Sample Request:
//
//	POST /api/v1/modules
//	{
//	  "name": "Inventory",
//	  "description": "Handles product stock management",
//	  "isActive": true
//	}
//
// Sample Success Response (201):
//
//	{
//	  "success": true,
//	  "message": "Resource created successfully",
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
// Sample Error Response (400):
//
//	{
//	  "success": false,
//	  "message": "Invalid request parameters",
//	  "error": {
//	    "code": "VALIDATION_ERROR",
//	    "message": "Name must be 3-50 characters",
//	    "details": {
//	      "name": ["Name must be 3-50 characters"]
//	    }
//	  },
//	  "meta": {
//	    "requestId": "a1b2c3d4",
//	    "timestamp": "2023-08-15T14:30:00Z"
//	  }
//	}
func (h *ModuleHandler) CreateModule(ctx *gin.Context) {
	// Step 1: Get request ID from context
	requestID := ctx.GetString("request_id")

	// Step 2: Create response mapper
	mapper := response.NewResponseMapper(requestID)

	// Step 3: Validate request payload
	var request module.ModuleRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		// Map validation errors to our format
		details := extractValidationErrors(err)

		// Use mapper to create error response
		response, statusCode := mapper.Error(
			"VALIDATION_ERROR",
			response.StatusToMessage(http.StatusBadRequest),
			details,
			http.StatusBadRequest,
		)
		ctx.JSON(statusCode, response)
		return
	}

	// Step 4: Execute business logic
	responseData, err := h.service.CreateModule(request)
	if err != nil {
		// Map service errors to appropriate responses
		handleServiceError(ctx, err, mapper)
		return
	}

	// Step 5: Use mapper to create success response
	response, statusCode := mapper.Success(
		responseData,
		response.StatusToMessage(http.StatusCreated),
		http.StatusCreated,
	)

	// Step 6: Return standardized response
	ctx.Header("Location", "/api/v1/modules/"+strconv.Itoa(responseData.ID))
	ctx.JSON(statusCode, response)
}

// GetModuleById godoc
// @Summary Get a module by ID
// @Description Retrieves a specific module by its unique identifier
// @Tags modules
// @Produce json
// @Param id path int true "Module ID"
// @Success 200 {object} response.APIResponse{data=module.ModuleResponse} "Module retrieved successfully"
// @Failure 404 {object} response.APIResponse "Module not found"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /modules/{id} [get]
func (h *ModuleHandler) GetModuleById(ctx *gin.Context) {
	requestID := ctx.GetString("request_id")
	mapper := response.NewResponseMapper(requestID)

	id := ctx.Param("id")
	module, err := h.service.GetModuleById(id)
	if err != nil {
		handleServiceError(ctx, err, mapper)
		return
	}

	// Use mapper to create success response
	response, statusCode := mapper.Success(
		module,
		response.StatusToMessage(http.StatusOK),
		http.StatusOK,
	)
	ctx.JSON(statusCode, response)
}

// handleServiceError processes errors from the business layer into standardized responses.
//
// This function maps business layer errors to appropriate HTTP status codes
// and creates consistent error responses following the APIResponse structure.
//
// Parameters:
//   - ctx: Gin context for the request
//   - err: The error returned from the business layer
//   - mapper: The response mapper to use for creating responses
func handleServiceError(ctx *gin.Context, err error, mapper *response.ResponseMapper) {
	statusCode := http.StatusInternalServerError
	code := "INTERNAL_ERROR"
	message := response.StatusToMessage(statusCode)

	switch {
	case errors.Is(err, moduleService.ErrNameRequired),
		errors.Is(err, moduleService.ErrNameLength),
		errors.Is(err, moduleService.ErrDescriptionLength):
		statusCode = http.StatusBadRequest
		code = "VALIDATION_ERROR"
		message = response.StatusToMessage(statusCode)

	case errors.Is(err, moduleService.ErrNameExists):
		statusCode = http.StatusConflict
		code = "RESOURCE_CONFLICT"
		message = response.StatusToMessage(statusCode)

	case errors.Is(err, moduleService.ErrNotFound):
		statusCode = http.StatusNotFound
		code = "NOT_FOUND"
		message = response.StatusToMessage(statusCode)
	}

	// For validation errors, extract field details
	var details map[string][]string
	if statusCode == http.StatusBadRequest {
		details = map[string][]string{
			"name": {err.Error()},
		}
	}

	// Use mapper to create error response
	response, statusCode := mapper.Error(
		code,
		message,
		details,
		statusCode,
	)
	ctx.JSON(statusCode, response)
}

// extractValidationErrors converts Gin validation errors to our format.
//
// Parameters:
//   - err: The validation error
//
// Returns:
//   - map[string][]string: Field-specific error messages
func extractValidationErrors(err error) map[string][]string {
	errors := make(map[string][]string)

	if verr, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range verr {
			field := fieldErr.Field()
			tag := fieldErr.Tag()

			message := "Validation failed"
			switch tag {
			case "required":
				message = "This field is required"
			case "min":
				message = "Value is too short"
			case "max":
				message = "Value exceeds maximum length"
			}

			errors[field] = append(errors[field], message)
		}
	}

	return errors
}
