package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/madhiyono/base-api-nosql/internal/repository"
	"github.com/madhiyono/base-api-nosql/pkg/logger"
)

type Handler struct {
	userRepo repository.UserRepository
	logger *logger.Logger
}

func NewUserHandler(userRepo repository.UserRepository, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		Handler: Handler{
			userRepo: userRepo,
			logger: logger,
		},
	}
}

type UserHandler struct {
	Handler
}

func (h *UserHandler) RegisterRoutes(e *echo.Echo) {
	users := e.Group("/users")
	{
		users.POST("", h.CreateUser)
		users.GET("/:id", h.GetUser)
		users.PUT("/:id", h.UpdateUser)
		users.DELETE("/:id", h.DeleteUser)
		users.GET("", h.ListUsers)
	}
}