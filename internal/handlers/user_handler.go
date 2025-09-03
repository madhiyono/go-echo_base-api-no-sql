package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/madhiyono/base-api-nosql/internal/cache"
	"github.com/madhiyono/base-api-nosql/internal/models"
	"github.com/madhiyono/base-api-nosql/pkg/response"
	"github.com/madhiyono/base-api-nosql/pkg/validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// In CreateUser method, you might want to check if the authenticated user has permission
// to create other users (admin-only functionality)
func (h *UserHandler) CreateUser(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		h.logger.Error("Failed to Bind User: %v", err)
		return response.BadRequest(c, "Failed to Create User: Invalid Request Format", nil)
	}

	// Validate user data (log detailed errors but return generic message)
	if err := validation.ValidateStruct(user); err != nil {
		// Log detailed validation errors for debugging
		validationErrors := validation.ValidateStructDetailed(user)
		for _, vErr := range validationErrors {
			h.logger.Error("Validation Error for User: %s", vErr)
		}
		return response.BadRequest(c, "Failed to Create User: Validation Error", nil)
	}

	if err := h.userRepo.Create(user); err != nil {
		h.logger.Error("Failed to Create User: %v", err)
		return response.InternalServerError(c, "Failed to Create User: Internal Server Error", nil)
	}

	// Invalidate user list cache using tags
	if err := h.cache.InvalidateTag(cache.UsersListTag); err != nil {
		h.logger.Error("Failed to Invalidate Users List Cache: %v", err)
	}

	return response.Created(c, "User Created Successfully", user)
}

// GetUser: Retrieves a user by ID
func (h *UserHandler) GetUser(c echo.Context) error {
	id := c.Param("id")

	// Try to get user from cache first
	cacheKey := fmt.Sprintf("%s%s", cache.UserCachePrefix, id)

	var cachedUser models.User
	if err := h.cache.Get(cacheKey, &cachedUser); err == nil {
		h.logger.Info("User Retrieved from Cache: %s", id)
		return response.Success(c, "User Retrieved Successfully", cachedUser)
	}

	// If not in cache, get from database
	user, err := h.userRepo.GetByID(id)
	if err != nil {
		h.logger.Error("Failed to Get User: %v", err)
		return response.NotFound(c, "User Not Found")
	}

	// Check if user is trying to access their own profile or has permission
	authUserID, _ := c.Get("user_id").(primitive.ObjectID)
	roleID := c.Get("role_id").(primitive.ObjectID)

	hasAdminPermission, _ := h.authService.HasPermission(roleID, "users", "read")
	if user.ID != authUserID && !hasAdminPermission {
		return response.Error(c, http.StatusForbidden, "Access Denied to This User Record", nil)
	}

	// Cache the user data with tags for easy invalidation
	tags := []string{cache.UsersTag}
	if err := h.cache.SetWithTags(cacheKey, *user, tags, cache.DefaultExpiration); err != nil {
		h.logger.Error("Failed to Cache User Data: %v", err)
		// Don't return error, just continue without caching
	}

	h.logger.Info("User Retrieved from Database and Cached: %s", id)
	return response.Success(c, "User Retrieved Successfully", user)
}

// In other methods, you might want to add authorization checks
// For example, users can only update their own profile
func (h *UserHandler) UpdateUser(c echo.Context) error {
	id := c.Param("id")

	// Check authorization - users can only update their own profile unless they have admin permissions
	authUserID, _ := c.Get("user_id").(primitive.ObjectID)
	roleID, _ := c.Get("role_id").(primitive.ObjectID)

	// Get existing user to check permissions
	existingUser, err := h.userRepo.GetByID(id)
	if err != nil {
		return response.NotFound(c, "User Not Found!")
	}

	// Non-admin users can only update their own profile
	hasAdminPermission, _ := h.authService.HasPermission(roleID, "users", "update")
	if existingUser.ID != authUserID && !hasAdminPermission {
		return response.Error(c, http.StatusForbidden, "Cannot Update This User Record", nil)
	}

	user := new(models.User)
	if err := c.Bind(user); err != nil {
		h.logger.Error("Failed to Bind User: %v", err)
		return response.BadRequest(c, "Failed to Update User: Invalid Request Format", nil)
	}

	// Validate user data (log detailed errors but return generic message)
	if err := validation.ValidateStruct(user); err != nil {
		// Log detailed validation errors for debugging
		validationErrors := validation.ValidateStructDetailed(user)
		for _, vErr := range validationErrors {
			h.logger.Error("Validation Error for User: %s", vErr)
		}
		return response.BadRequest(c, "Failed to Update User: Validation Error", nil)
	}

	if err := h.userRepo.Update(id, user); err != nil {
		h.logger.Error("Failed to Update User: %v", err)
		return response.InternalServerError(c, "Failed to Update User: Internal Server Error", nil)
	}

	// Invalidate cache for this specific user
	cacheKey := fmt.Sprintf("%s%s", cache.UserCachePrefix, id)
	if err := h.cache.Delete(cacheKey); err != nil {
		h.logger.Error("Failed to invalidate user cache: %v", err)
	}

	// Invalidate cache for this specific user using tags
	userCacheKey := fmt.Sprintf("%s%s", cache.UserCachePrefix, id)
	if err := h.cache.InvalidateTag(cache.UsersTag); err != nil {
		h.logger.Error("Failed to Invalidate User Cache by Tag: %v", err)
	}

	// Also explicitly delete the specific user key
	if err := h.cache.Delete(userCacheKey); err != nil {
		h.logger.Error("Failed to Delete Specific User Cache: %v", err)
	}

	// Invalidate user list cache
	if err := h.cache.InvalidateTag(cache.UsersListTag); err != nil {
		h.logger.Error("Failed to Invalidate Users List Cache: %v", err)
	}

	return response.Success(c, "User Updated Successfully", user)
}

// DeleteUser: Deletes a user by ID
func (h *UserHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")

	// Check authorization - users can only delete their own profile unless they have admin permissions
	authUserID, _ := c.Get("user_id").(primitive.ObjectID)
	roleID, _ := c.Get("role_id").(primitive.ObjectID)

	existingUser, err := h.userRepo.GetByID(id)
	if err != nil {
		return response.NotFound(c, "User Not Found!")
	}

	// Non-admin users can only delete their own profile
	hasAdminPermission, _ := h.authService.HasPermission(roleID, "users", "delete")
	if existingUser.ID != authUserID && !hasAdminPermission {
		return response.Error(c, http.StatusForbidden, "Cannot Delete This User Record", nil)
	}

	if err := h.userRepo.Delete(id); err != nil {
		h.logger.Error("Failed to Delete User: %v", err)
		return response.InternalServerError(c, "Failed to Delete User: Internal Server Error", nil)
	}

	// Invalidate cache using tags
	if err := h.cache.InvalidateTag(cache.UsersTag); err != nil {
		h.logger.Error("Failed to Invalidate User Cache by Tag: %v", err)
	}

	// Invalidate user list cache
	if err := h.cache.InvalidateTag(cache.UsersListTag); err != nil {
		h.logger.Error("Failed to Invalidate Users List Cache: %v", err)
	}

	return response.Success(c, "User Deleted Successfully", nil)
}

// ListUsers: Returns all users
func (h *UserHandler) ListUsers(c echo.Context) error {
	authUserID := c.Get("user_id").(primitive.ObjectID)
	roleID := c.Get("role_id").(primitive.ObjectID)

	// Try to get users list from cache first
	cacheKey := fmt.Sprintf("users_list:%s:%s", authUserID.Hex(), roleID.Hex())

	var cachedUsers []*models.User
	if err := h.cache.Get(cacheKey, &cachedUsers); err == nil {
		h.logger.Info("Users list retrieved from cache")
		return response.Success(c, "Users retrieved successfully", cachedUsers)
	}

	// Check if user has admin permission to see all users
	hasAdminPermission, _ := h.authService.HasPermission(roleID, "users", "read")

	var users []*models.User
	var err error

	if hasAdminPermission {
		// Admin can see all users
		users, err = h.userRepo.List()
	} else {
		// Regular users can only see their own records or records they own
		return response.Error(c, http.StatusForbidden, "Access Denied to This User Record", nil)
	}

	if err != nil {
		h.logger.Error("Failed to List Users: %v", err)
		return response.InternalServerError(c, "Failed to Retrieve Users: Internal Server Error", nil)
	}

	// Cache the users list with tags for easy invalidation
	tags := []string{cache.UsersListTag, cache.UsersTag}
	if err := h.cache.SetWithTags(cacheKey, users, tags, cache.DefaultExpiration); err != nil {
		h.logger.Error("Failed to Cache Users List: %v", err)
		// Don't return error, just continue without caching
	}

	h.logger.Info("Users List Retrieved from Database and Cached")
	return response.Success(c, "Users Retrieved Successfully", users)
}

// UploadProfilePhoto uploads a profile photo for the user
func (h *UserHandler) UploadProfilePhoto(c echo.Context) error {
	authUserID := c.Get("user_id").(primitive.ObjectID)
	userID := c.Param("id")

	// Check if user is updating their own profile or has admin permission
	if userID != authUserID.Hex() {
		roleID := c.Get("role_id").(primitive.ObjectID)
		hasAdminPermission, _ := h.authService.HasPermission(roleID, "users", "update")
		if !hasAdminPermission {
			return response.Error(c, http.StatusForbidden, "Cannot update other users' profile photos", nil)
		}
	}

	// Parse multipart form with max memory of 32MB
	form, err := c.MultipartForm()
	if err != nil {
		h.logger.Error("Failed to parse multipart form: %v", err)
		return response.BadRequest(c, "Failed to parse upload: Invalid form data", nil)
	}

	// Get files from form
	files := form.File["photo"]
	if len(files) == 0 {
		return response.BadRequest(c, "No photo file provided", nil)
	}

	fileHeader := files[0]

	// Validate file size (max 5MB)
	if fileHeader.Size > 5*1024*1024 {
		return response.BadRequest(c, "File size exceeds 5MB limit", nil)
	}

	// Validate file type
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
	}

	file, err := fileHeader.Open()
	if err != nil {
		h.logger.Error("Failed to open file: %v", err)
		return response.InternalServerError(c, "Failed to process file", nil)
	}
	defer file.Close()

	// Check file type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		h.logger.Error("Failed to read file header: %v", err)
		return response.InternalServerError(c, "Failed to process file", nil)
	}
	_, err = file.Seek(0, 0) // Reset file pointer
	if err != nil {
		h.logger.Error("Failed to reset file pointer: %v", err)
		return response.InternalServerError(c, "Failed to process file", nil)
	}

	contentType := http.DetectContentType(buffer)
	if !allowedTypes[contentType] {
		return response.BadRequest(c, "Invalid file type. Only JPEG, PNG, and GIF are allowed", nil)
	}

	// Upload to storage
	uploadResult, err := h.storageService.UploadProfilePhoto(authUserID, file, fileHeader.Size, fileHeader.Filename)
	if err != nil {
		h.logger.Error("Failed to upload profile photo: %v", err)
		return response.InternalServerError(c, "Failed to upload profile photo", nil)
	}

	// Update user record with photo URL
	err = h.userRepo.UpdateProfilePhoto(userID, uploadResult.URL)
	if err != nil {
		h.logger.Error("Failed to update user with photo URL: %v", err)
		// Try to clean up uploaded file
		h.storageService.DeleteProfilePhoto(uploadResult.Key)
		return response.InternalServerError(c, "Failed to update user profile", nil)
	}

	// Get updated user
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		h.logger.Error("Failed to get updated user: %v", err)
		return response.InternalServerError(c, "Failed to retrieve updated user", nil)
	}

	return response.Success(c, "Profile photo uploaded successfully", user)
}

// DeleteProfilePhoto removes the profile photo
func (h *UserHandler) DeleteProfilePhoto(c echo.Context) error {
	authUserID := c.Get("user_id").(primitive.ObjectID)
	userID := c.Param("id")

	// Check if user is updating their own profile or has admin permission
	if userID != authUserID.Hex() {
		roleID := c.Get("role_id").(primitive.ObjectID)
		hasAdminPermission, _ := h.authService.HasPermission(roleID, "users", "update")
		if !hasAdminPermission {
			return response.Error(c, http.StatusForbidden, "Cannot delete other users' profile photos", nil)
		}
	}

	// Get current user to get photo URL
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		return response.NotFound(c, "User not found")
	}

	// If user has a profile photo, delete it from storage
	if user.ProfilePhoto != "" {
		// Extract key from URL (this is a simplified approach)
		// In production, you might want to store the key separately
		key := h.extractKeyFromURL(user.ProfilePhoto)
		if key != "" {
			err := h.storageService.DeleteProfilePhoto(key)
			if err != nil {
				h.logger.Error("Failed to delete photo from storage: %v", err)
				// Don't return error here, continue with database update
			}
		}
	}

	// Update user record to remove photo URL
	err = h.userRepo.UpdateProfilePhoto(userID, "")
	if err != nil {
		h.logger.Error("Failed to remove photo URL from user: %v", err)
		return response.InternalServerError(c, "Failed to remove profile photo", nil)
	}

	// Get updated user
	updatedUser, err := h.userRepo.GetByID(userID)
	if err != nil {
		h.logger.Error("Failed to get updated user: %v", err)
		return response.InternalServerError(c, "Failed to retrieve updated user", nil)
	}

	return response.Success(c, "Profile photo deleted successfully", updatedUser)
}

// Helper function to extract key from URL
func (h *UserHandler) extractKeyFromURL(url string) string {
	// Simple extraction - in production, you might store the key separately
	// Expected format: http://minio:9000/bucket-name/key
	// This is a basic implementation - you might need to adjust based on your URL structure
	parts := strings.Split(url, "/")
	if len(parts) >= 5 {
		return strings.Join(parts[4:], "/")
	}
	return ""
}
