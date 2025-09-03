package email

import (
	"context"
	"crypto/rand"
	"fmt"
	"html/template"
	"io"
	"math/big"
	"mime"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/madhiyono/base-api-nosql/internal/cache"
	"github.com/madhiyono/base-api-nosql/internal/models"
	"github.com/madhiyono/base-api-nosql/internal/repository"
	"github.com/madhiyono/base-api-nosql/pkg/logger"
)

type EmailService struct {
	cache      cache.Cache
	verifyRepo repository.VerificationRepository
	logger     *logger.Logger
	config     EmailConfig
	ctx        context.Context
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	FromName     string
	TemplatesDir string
	WorkerCount  int
}

func NewEmailService(
	verifyRepo repository.VerificationRepository,
	cache cache.Cache,
	logger *logger.Logger,
	config EmailConfig,
) *EmailService {
	service := &EmailService{
		cache:      cache,
		verifyRepo: verifyRepo,
		logger:     logger,
		config:     config,
		ctx:        context.Background(),
	}

	// Start email processing workers
	if config.WorkerCount <= 0 {
		config.WorkerCount = 3
	}
	go service.startWorkers(config.WorkerCount)

	return service
}

func (s *EmailService) enqueueEmail(email *models.EmailMessage) error {
	// Use your cache package to store email in a queue-like structure
	queueKey := "email_queue"

	// Get existing queue
	var queue []models.EmailMessage
	if err := s.cache.Get(queueKey, &queue); err != nil {
		// If queue doesn't exist, start with empty slice
		queue = []models.EmailMessage{}
	}

	// Add new email to queue
	queue = append(queue, *email)

	// Store updated queue (expire in 1 day)
	return s.cache.Set(queueKey, queue, 24*time.Hour)
}

func (s *EmailService) dequeueEmail() (*models.EmailMessage, error) {
	queueKey := "email_queue"

	// Get existing queue
	var queue []models.EmailMessage
	if err := s.cache.Get(queueKey, &queue); err != nil {
		return nil, err
	}

	if len(queue) == 0 {
		return nil, fmt.Errorf("no emails in queue")
	}

	// Get first email (FIFO)
	email := queue[0]

	// Update queue without first email
	if len(queue) > 1 {
		queue = queue[1:]
		s.cache.Set(queueKey, queue, 24*time.Hour)
	} else {
		// If queue is empty, delete the key
		s.cache.Delete(queueKey)
	}

	return &email, nil
}

func (s *EmailService) startWorkers(workerCount int) {
	for i := 0; i < workerCount; i++ {
		go s.emailWorker(i)
	}
	s.logger.Info("Started %d email workers", workerCount)
}

func (s *EmailService) emailWorker(workerID int) {
	s.logger.Info("Email worker %d started", workerID)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Process pending emails
		email, err := s.dequeueEmail()
		if err != nil {
			// No emails in queue, continue
			continue
		}

		// Process email
		if err := s.processEmail(workerID, email); err != nil {
			s.logger.Error("Worker %d: Failed to process email to %s: %v", workerID, email.To, err)
			s.handleEmailFailure(email)
		} else {
			s.logger.Info("Worker %d: Email sent successfully to %s", workerID, email.To)
			s.handleEmailSuccess(email)
		}
	}
}

func (s *EmailService) processEmail(workerID int, email *models.EmailMessage) error {
	s.logger.Info("Worker %d: Processing email to %s", workerID, email.To)

	// Update status to sending
	email.Status = models.EmailStatusSending
	email.UpdatedAt = time.Now()
	if err := s.updateEmailStatus(email); err != nil {
		s.logger.Error("Worker %d: Failed to update email status: %v", workerID, err)
	}

	// Create SMTP auth
	auth := smtp.PlainAuth("", s.config.SMTPUser, s.config.SMTPPassword, s.config.SMTPHost)

	// Create email message
	from := fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
	to := []string{email.To}

	subject := mime.QEncoding.Encode("UTF-8", email.Subject)
	msg := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: multipart/alternative; boundary=frontier\r\n"+
			"\r\n"+
			"--frontier\r\n"+
			"Content-Type: text/plain; charset=UTF-8\r\n"+
			"Content-Transfer-Encoding: quoted-printable\r\n"+
			"\r\n"+
			"%s\r\n"+
			"--frontier\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"Content-Transfer-Encoding: quoted-printable\r\n"+
			"\r\n"+
			"%s\r\n"+
			"--frontier--\r\n",
		from, email.To, subject, email.BodyText, email.BodyHTML,
	)

	// Send email
	addr := s.config.SMTPHost + ":" + s.config.SMTPPort
	return smtp.SendMail(addr, auth, s.config.FromEmail, to, []byte(msg))
}

func (s *EmailService) handleEmailSuccess(email *models.EmailMessage) {
	email.Status = models.EmailStatusSent
	now := time.Now()
	email.SentAt = &now
	email.UpdatedAt = now
	s.updateEmailStatus(email)
}

func (s *EmailService) handleEmailFailure(email *models.EmailMessage) {
	email.RetryCount++
	email.UpdatedAt = time.Now()

	if email.RetryCount < email.MaxRetries {
		// Retry the email
		email.Status = models.EmailStatusRetry
		s.updateEmailStatus(email)
		s.enqueueEmail(email) // Re-queue for retry
		s.logger.Info("Email to %s queued for retry (%d/%d)", email.To, email.RetryCount, email.MaxRetries)
	} else {
		// Mark as failed
		email.Status = models.EmailStatusFailed
		s.updateEmailStatus(email)
		s.logger.Error("Email to %s failed after %d retries", email.To, email.RetryCount)
	}
}

func (s *EmailService) updateEmailStatus(email *models.EmailMessage) error {
	// For tracking purposes, we could store in a separate key
	// Here we'll just log the status change
	s.logger.Info("Email %s status updated to %s", email.ID, email.Status)
	return nil
}

func (s *EmailService) VerifyEmail(token string) error {
	verification, err := s.verifyRepo.GetByToken(token)
	if err != nil {
		return fmt.Errorf("invalid or expired verification token")
	}

	// Mark as used
	if err := s.verifyRepo.MarkAsUsed(verification.ID); err != nil {
		return err
	}

	return nil
}

func (s *EmailService) generateToken(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	token := make([]byte, length)

	for i := range token {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		token[i] = charset[n.Int64()]
	}

	return string(token), nil
}

func (s *EmailService) generateMessageID() string {
	// Generate unique message ID
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return fmt.Sprintf("msg_%d_%d", time.Now().Unix(), n.Int64())
}

func (s *EmailService) loadTemplate(filename string) (string, error) {
	path := fmt.Sprintf("%s/%s", s.config.TemplatesDir, filename)
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (s *EmailService) renderTemplate(templateStr string, data map[string]string) (string, error) {
	tmpl, err := template.New("email").Parse(templateStr)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// GetQueueStats returns statistics about the email queue
func (s *EmailService) GetQueueStats() (map[string]interface{}, error) {
	queueKey := "email_queue"

	var queue []models.EmailMessage
	if err := s.cache.Get(queueKey, &queue); err != nil {
		queue = []models.EmailMessage{}
	}

	return map[string]interface{}{
		"pending_emails": len(queue),
		"workers":        s.config.WorkerCount,
	}, nil
}
