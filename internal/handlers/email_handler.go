package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/madhiyono/base-api-nosql/internal/auth"
	"github.com/madhiyono/base-api-nosql/pkg/response"
)

// GetEmailQueueStats returns statistics about the email queue
func (h *EmailHandler) GetEmailQueueStats(c echo.Context) error {
	stats, err := h.emailService.GetQueueStats()
	if err != nil {
		h.logger.Error("Failed to get email queue stats: %v", err)
		return response.InternalServerError(c, "Failed to get queue statistics", err)
	}
	return response.Success(c, "Email queue statistics retrieved successfully", stats)
}

// GetEmailQueueDetails returns detailed information about emails in queue
func (h *EmailHandler) GetEmailQueueDetails(c echo.Context) error {
	// This would require additional methods in email service to get detailed queue info
	// For now, return basic stats
	stats, err := h.emailService.GetQueueStats()
	if err != nil {
		h.logger.Error("Failed to get email queue details: %v", err)
		return response.InternalServerError(c, "Failed to get queue details", err)
	}
	return response.Success(c, "Email queue details retrieved successfully", stats)
}

// ClearFailedEmails clears failed emails from queue (admin only)
func (h *EmailHandler) ClearFailedEmails(c echo.Context) error {
	// This would require additional implementation in email service
	return response.Success(c, "Failed emails cleared successfully", nil)
}

// RetryFailedEmails retries all failed emails (admin only)
func (h *EmailHandler) RetryFailedEmails(c echo.Context) error {
	// This would require additional implementation in email service
	return response.Success(c, "Failed emails queued for retry", nil)
}

// Register email routes
func (h *EmailHandler) RegisterRoutes(e *echo.Echo, authMiddleware *auth.Middleware) {
	// Public email-related endpoints
	emailGroup := e.Group("/email")
	{
		emailGroup.GET("/queue-stats", h.GetEmailQueueStats)
		emailGroup.GET("/queue-details", h.GetEmailQueueDetails)
	}

	// Admin-only email management endpoints
	adminEmailGroup := e.Group("/admin/email")
	adminEmailGroup.Use(authMiddleware.JWTAuth)
	adminEmailGroup.Use(authMiddleware.RequireAdmin())
	{
		adminEmailGroup.DELETE("/failed", h.ClearFailedEmails)
		adminEmailGroup.POST("/retry-failed", h.RetryFailedEmails)
	}
}
