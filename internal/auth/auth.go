package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/madhiyono/base-api-nosql/internal/models"
	"github.com/madhiyono/base-api-nosql/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepo repository.AuthRepository
	userRepo repository.UserRepository
	jwtKey   []byte
}

func NewAuthService(authRepo repository.AuthRepository, userRepo repository.UserRepository, jwtKey string) *AuthService {
	return &AuthService{
		authRepo: authRepo,
		userRepo: userRepo,
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

func (s *AuthService) GenerateToken(user *models.User, role models.UserRole) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &models.Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   role,
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
			return nil, fmt.Errorf("Unexpected Signing Method: %v", token.Header["alg"])
		}
		return s.jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("Invalid Token!")
	}

	return claims, nil
}

func (s *AuthService) Register(request *models.RegisterRequest) (*models.AuthResponse, error) {
	// Check if user already exists
	if _, err := s.authRepo.GetByEmail(request.Email); err == nil {
		return nil, fmt.Errorf("User Already Exists!")
	}

	// Create user
	user := &models.User{
		Name:  request.Name,
		Email: request.Email,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
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
		Role:     models.RoleUser, // Default role
		IsActive: true,
	}

	if err := s.authRepo.Create(auth); err != nil {
		return nil, err
	}

	// Generate token
	token, err := s.GenerateToken(user, auth.Role)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token: token,
		User:  user,
		Role:  auth.Role,
	}, nil
}

func (s *AuthService) Login(request *models.LoginRequest) (*models.AuthResponse, error) {
	// Get auth record
	auth, err := s.authRepo.GetByEmail(request.Email)
	if err != nil {
		return nil, fmt.Errorf("Invalid Credentials!")
	}

	// Check if account is active
	if !auth.IsActive {
		return nil, fmt.Errorf("Account is Deactivated!")
	}

	// Check password
	if !s.CheckPasswordHash(request.Password, auth.Password) {
		return nil, fmt.Errorf("Invalid Credentials!")
	}

	// Get user details
	user, err := s.userRepo.GetByID(auth.UserID.Hex())
	if err != nil {
		return nil, err
	}

	// Generate token
	token, err := s.GenerateToken(user, auth.Role)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token: token,
		User:  user,
		Role:  auth.Role,
	}, nil
}
