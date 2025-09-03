package handlers

import (
	"github.com/madhiyono/base-api-nosql/internal/auth"
	"github.com/madhiyono/base-api-nosql/internal/cache"
	"github.com/madhiyono/base-api-nosql/internal/email"
	"github.com/madhiyono/base-api-nosql/internal/repository"
	"github.com/madhiyono/base-api-nosql/internal/services"
	"github.com/madhiyono/base-api-nosql/internal/storage"
	"github.com/madhiyono/base-api-nosql/pkg/logger"
)

type Handler struct {
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	authRepo       repository.AuthRepository
	verifyRepo     repository.VerificationRepository
	authService    *auth.AuthService
	storageService *storage.StorageService
	emailService   *email.EmailService
	logger         *logger.Logger
}

func NewUserHandler(userRepo repository.UserRepository, authService *auth.AuthService, storageService *storage.StorageService, cache cache.Cache, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		Handler: Handler{
			userRepo:       userRepo,
			authService:    authService,
			storageService: storageService,
			logger:         logger,
		},
		cache: cache,
	}
}

func NewAuthHandler(authRepo repository.AuthRepository, verifyRepo repository.VerificationRepository, authService *auth.AuthService, emailService *email.EmailService, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{
		Handler: Handler{
			authRepo:     authRepo,
			verifyRepo:   verifyRepo,
			authService:  authService,
			emailService: emailService,
			logger:       logger,
		},
	}
}

func NewRoleHandler(roleRepo repository.RoleRepository, authService *auth.AuthService, logger *logger.Logger) *RoleHandler {
	return &RoleHandler{
		Handler: Handler{
			roleRepo:    roleRepo,
			authService: authService,
			logger:      logger,
		},
	}
}

func NewEmailHandler(emailService *email.EmailService, logger *logger.Logger) *EmailHandler {
	return &EmailHandler{
		Handler: Handler{
			emailService: emailService,
			logger:       logger,
		},
	}
}

func NewWebSocketHandler(wsService *services.WebSocketService, logger *logger.Logger) *WebSocketHandler {
	return &WebSocketHandler{
		Handler: Handler{
			logger: logger,
		},
		wsService: wsService,
	}
}

type UserHandler struct {
	Handler
	cache cache.Cache
}

type AuthHandler struct {
	Handler
}

type RoleHandler struct {
	Handler
}

type EmailHandler struct {
	Handler
}

type WebSocketHandler struct {
	Handler
	wsService *services.WebSocketService
}
