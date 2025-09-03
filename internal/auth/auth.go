package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/madhiyono/base-api-nosql/internal/models"
	"github.com/madhiyono/base-api-nosql/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepo repository.AuthRepository
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
	jwtKey   []byte
}

func NewAuthService(
	authRepo repository.AuthRepository,
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	jwtKey string,
) *AuthService {
	return &AuthService{
		authRepo: authRepo,
		userRepo: userRepo,
		roleRepo: roleRepo,
		jwtKey:   []byte(jwtKey),
	}
}

func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *AuthService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *AuthService) GenerateToken(user *models.User, roleID primitive.ObjectID) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &models.Claims{
		UserID: user.ID,
		Email:  user.Email,
		RoleID: roleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtKey)
}

func (s *AuthService) ValidateToken(tokenString string) (*models.Claims, error) {
	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (s *AuthService) Register(request *models.RegisterRequest) (*models.AuthResponse, error) {
	// Check if user already exists
	if _, err := s.authRepo.GetByEmail(request.Email); err == nil {
		return nil, fmt.Errorf("user already exists")
	}

	// Create user
	user := &models.User{
		Name:  request.Name,
		Email: request.Email,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Get default role (user role)
	defaultRole, err := s.roleRepo.GetByName("user")
	if err != nil {
		return nil, fmt.Errorf("default role not found")
	}

	// Hash password
	hashedPassword, err := s.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	// Create auth record
	auth := &models.UserAuth{
		UserID:   user.ID,
		Email:    request.Email,
		Password: hashedPassword,
		RoleID:   defaultRole.ID,
		IsActive: true,
	}

	if err := s.authRepo.Create(auth); err != nil {
		return nil, err
	}

	// Generate token
	token, err := s.GenerateToken(user, auth.RoleID)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token: token,
		User:  user,
		Role:  defaultRole,
	}, nil
}

func (s *AuthService) Login(request *models.LoginRequest) (*models.AuthResponse, error) {
	// Get auth record
	auth, err := s.authRepo.GetByEmail(request.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if account is active
	if !auth.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	// Check password
	if !s.CheckPasswordHash(request.Password, auth.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Get user details
	user, err := s.userRepo.GetByID(auth.UserID.Hex())
	if err != nil {
		return nil, err
	}

	// Get role details
	role, err := s.roleRepo.GetByID(auth.RoleID)
	if err != nil {
		return nil, err
	}

	// Generate token
	token, err := s.GenerateToken(user, auth.RoleID)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token: token,
		User:  user,
		Role:  role,
	}, nil
}

// Check if user has permission
func (s *AuthService) HasPermission(roleID primitive.ObjectID, resource, action string) (bool, error) {
	return s.roleRepo.HasPermission(roleID, resource, action)
}
