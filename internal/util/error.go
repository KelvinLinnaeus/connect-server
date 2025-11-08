package util

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Custom error types
var (
	ErrNotFound           = errors.New("resource not found")
	ErrConflict           = errors.New("resource already exists")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrBadRequest         = errors.New("bad request")
	ErrInternalServer     = errors.New("internal server error")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrExpiredToken       = errors.New("token expired")
	ErrInvalidToken       = errors.New("invalid token")
)

// ErrorResponse represents API error response
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains error code and message
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(code, message string) ErrorResponse {
	return ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
}

// HandleError maps application errors to HTTP responses
func HandleError(c *gin.Context, err error) {
	fmt.Println(err)

	errMsg := strings.ToLower(err.Error())

	switch {
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, NewErrorResponse("not_found", err.Error()))
	case errors.Is(err, sql.ErrNoRows):
		c.JSON(http.StatusNotFound, NewErrorResponse("not_found", "Resource not found"))
	case strings.Contains(errMsg, "not found"):
		c.JSON(http.StatusNotFound, NewErrorResponse("not_found", err.Error()))
	case errors.Is(err, ErrConflict):
		c.JSON(http.StatusConflict, NewErrorResponse("conflict", err.Error()))
	case errors.Is(err, ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, NewErrorResponse("unauthorized", err.Error()))
	case errors.Is(err, ErrForbidden):
		c.JSON(http.StatusForbidden, NewErrorResponse("forbidden", err.Error()))
	case errors.Is(err, ErrBadRequest):
		c.JSON(http.StatusBadRequest, NewErrorResponse("bad_request", err.Error()))
	case errors.Is(err, ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, NewErrorResponse("invalid_credentials", "Invalid email or password"))
	case errors.Is(err, ErrExpiredToken):
		c.JSON(http.StatusUnauthorized, NewErrorResponse("token_expired", "Token has expired"))
	case strings.Contains(errMsg, "token has expired") || strings.Contains(errMsg, "token expired"):
		c.JSON(http.StatusUnauthorized, NewErrorResponse("token_expired", "Token has expired"))
	case errors.Is(err, ErrInvalidToken):
		c.JSON(http.StatusUnauthorized, NewErrorResponse("invalid_token", "Invalid token"))
	case strings.Contains(errMsg, "invalid token"):
		c.JSON(http.StatusUnauthorized, NewErrorResponse("invalid_token", "Invalid token"))
	default:
		// Check for database constraint violations
		if IsDuplicateKeyError(err) {
			c.JSON(http.StatusConflict, NewErrorResponse("duplicate_key", "Resource already exists"))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse("internal_error", "Something went wrong"))
	}
}

// IsDuplicateKeyError checks if error is a duplicate key constraint violation
func IsDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	// PostgreSQL duplicate key error contains "duplicate key value violates unique constraint"
	errMsg := err.Error()
	return contains(errMsg, "duplicate key value") || contains(errMsg, "UNIQUE constraint")
}

// contains checks if string contains substring (case-insensitive helper)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

// NewSuccessResponse creates a success response
func NewSuccessResponse(data interface{}) SuccessResponse {
	return SuccessResponse{
		Status: "success",
		Data:   data,
	}
}

type SuccessErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewSuccessErrorResponse(message string) SuccessErrorResponse {
	return SuccessErrorResponse{
		Status:  "success",
		Message: message,
	}
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Status string         `json:"status"`
	Data   interface{}    `json:"data"`
	Meta   PaginationMeta `json:"meta"`
}

// PaginationMeta contains pagination metadata
type PaginationMeta struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// NewPaginatedResponse creates a paginated response
func NewPaginatedResponse(data interface{}, total, page, limit int) PaginatedResponse {
	return PaginatedResponse{
		Status: "success",
		Data:   data,
		Meta: PaginationMeta{
			Total: total,
			Page:  page,
			Limit: limit,
		},
	}
}
