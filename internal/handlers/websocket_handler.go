package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/madhiyono/base-api-nosql/internal/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *WebSocketHandler) WebSocketConnection(c echo.Context) error {
	// Get user from context (should be authenticated)
	userID, ok := c.Get("user_id").(primitive.ObjectID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	// Get upgrader from WebSocket service
	upgrader := h.wsService.GetUpgrader()

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		h.logger.Error("Failed to upgrade connection: %v", err)
		return err
	}

	// Handle WebSocket connection
	go h.wsService.HandleConnection(conn, userID)

	return nil
}

func (h *WebSocketHandler) RegisterRoutes(e *echo.Echo, authMiddleware *auth.Middleware) {
	// WebSocket endpoint (requires authentication)
	e.GET("/ws", h.WebSocketConnection, authMiddleware.JWTAuth)
}
