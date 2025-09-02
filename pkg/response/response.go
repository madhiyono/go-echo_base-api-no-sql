package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Standard response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Success response
func Success(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created response
func Created(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error response
func Error(c echo.Context, statusCode int, message string, err error) error {
	errorDetail := ""
	if err != nil {
		errorDetail = err.Error()
	}

	return c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error:   errorDetail,
	})
}

// BadRequest response
func BadRequest(c echo.Context, message string, err error) error {
	return Error(c, http.StatusBadRequest, message, err)
}

// NotFound response
func NotFound(c echo.Context, message string) error {
	return Error(c, http.StatusNotFound, message, nil)
}

// InternalServerError response
func InternalServerError(c echo.Context, message string, err error) error {
	return Error(c, http.StatusInternalServerError, message, err)
}
