package util

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)


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


type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}


type ErrorDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}


func NewErrorResponse(code int, message string) ErrorResponse {
	return ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
}


func HandleError(c *gin.Context, err error) {
	fmt.Println(err)

	errMsg := strings.ToLower(err.Error())

	switch {
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, err.Error()))
	case errors.Is(err, sql.ErrNoRows):
		c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, "Resource not found"))
	case strings.Contains(errMsg, "not found"):
		c.JSON(http.StatusNotFound, NewErrorResponse(http.StatusNotFound, err.Error()))
	case errors.Is(err, ErrConflict):
		c.JSON(http.StatusConflict, NewErrorResponse(http.StatusConflict, err.Error()))
	case errors.Is(err, ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, err.Error()))
	case errors.Is(err, ErrForbidden):
		c.JSON(http.StatusForbidden, NewErrorResponse(http.StatusForbidden, err.Error()))
	case errors.Is(err, ErrBadRequest):
		c.JSON(http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, err.Error()))
	case errors.Is(err, ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "Invalid email or password"))
	case errors.Is(err, ErrExpiredToken):
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "Token has expired"))
	case strings.Contains(errMsg, "token has expired") || strings.Contains(errMsg, "token expired"):
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "Token has expired"))
	case errors.Is(err, ErrInvalidToken):
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "Invalid token"))
	case strings.Contains(errMsg, "invalid token"):
		c.JSON(http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "Invalid token"))
	default:

		if IsDuplicateKeyError(err) {
			c.JSON(http.StatusConflict, NewErrorResponse(http.StatusConflict, "Resource already exists"))
			return
		}
		c.JSON(http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, "Something went wrong"))
	}
}


func IsDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	
	errMsg := err.Error()
	return contains(errMsg, "duplicate key value") || contains(errMsg, "UNIQUE constraint")
}


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


type SuccessResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}


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


type PaginatedResponse struct {
	Status string         `json:"status"`
	Data   interface{}    `json:"data"`
	Meta   PaginationMeta `json:"meta"`
}


type PaginationMeta struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}


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
