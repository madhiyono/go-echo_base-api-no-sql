package models

import (
	"time"
)

type EmailTemplate string

const (
	TemplateVerification  EmailTemplate = "verification"
	TemplateWelcome       EmailTemplate = "welcome"
	TemplateResetPassword EmailTemplate = "reset_password"
)

type EmailMessage struct {
	ID         string            `json:"id" redis:"id"`
	To         string            `json:"to" redis:"to"`
	Subject    string            `json:"subject" redis:"subject"`
	BodyHTML   string            `json:"body_html" redis:"body_html"`
	BodyText   string            `json:"body_text" redis:"body_text"`
	Template   EmailTemplate     `json:"template" redis:"template"`
	Variables  map[string]string `json:"variables" redis:"variables"`
	Status     EmailStatus       `json:"status" redis:"status"`
	RetryCount int               `json:"retry_count" redis:"retry_count"`
	MaxRetries int               `json:"max_retries" redis:"max_retries"`
	Priority   int               `json:"priority" redis:"priority"` // 0 = high, 1 = normal, 2 = low
	CreatedAt  time.Time         `json:"created_at" redis:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at" redis:"updated_at"`
	SentAt     *time.Time        `json:"sent_at,omitempty" redis:"sent_at"`
}

type EmailStatus string

const (
	EmailStatusPending EmailStatus = "pending"
	EmailStatusSending EmailStatus = "sending"
	EmailStatusSent    EmailStatus = "sent"
	EmailStatusFailed  EmailStatus = "failed"
	EmailStatusRetry   EmailStatus = "retry"
)
