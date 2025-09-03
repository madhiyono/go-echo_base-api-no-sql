package handlers

import (
	"github.com/madhiyono/base-api-nosql/internal/auth"
	"github.com/madhiyono/base-api-nosql/internal/cache"
	"github.com/madhiyono/base-api-nosql/internal/repository"
	"github.com/madhiyono/base-api-nosql/internal/storage"
	"github.com/madhiyono/base-api-nosql/pkg/logger"
)

type Handler struct {
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	authService    *auth.AuthService
	storageService *storage.StorageService
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

func NewAuthHandler(authService *auth.AuthService, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{
		Handler: Handler{
			authService: authService,
			logger:      logger,
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
