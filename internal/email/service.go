package email

import (
	"fmt"
	"time"

	"github.com/madhiyono/base-api-nosql/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *EmailService) SendVerificationEmail(userID primitive.ObjectID, email, name string) error {
	// Generate verification token
	token, err := s.generateToken(32)
	if err != nil {
		return err
	}

	// Create verification record
	verification := &models.EmailVerification{
		UserID:    userID,
		Email:     email,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hours
		IsUsed:    false,
	}

	if err := s.verifyRepo.Create(verification); err != nil {
		return err
	}

	// Load email templates
	htmlTemplate, err := s.loadTemplate("verification.html")
	if err != nil {
		return err
	}

	textTemplate, err := s.loadTemplate("verification.txt")
	if err != nil {
		return err
	}

	// Prepare template data
	verificationURL := fmt.Sprintf("http://localhost:8080/auth/verify/%s", token)
	templateData := map[string]string{
		"Name":            name,
		"Email":           email,
		"VerificationURL": verificationURL,
	}

	// Render templates
	htmlBody, err := s.renderTemplate(htmlTemplate, templateData)
	if err != nil {
		return err
	}

	textBody, err := s.renderTemplate(textTemplate, templateData)
	if err != nil {
		return err
	}

	// Create email message
	emailMsg := &models.EmailMessage{
		ID:         s.generateMessageID(),
		To:         email,
		Subject:    "Verify Your Email Address",
		BodyHTML:   htmlBody,
		BodyText:   textBody,
		Template:   models.TemplateVerification,
		Variables:  templateData,
		Status:     models.EmailStatusPending,
		RetryCount: 0,
		MaxRetries: 3,
		Priority:   1, // Normal priority
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Add to queue
	return s.enqueueEmail(emailMsg)
}
