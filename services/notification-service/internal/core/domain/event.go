package domain

import "time"

type NotificationType string

const (
	NotificationUserRegistered       NotificationType = "user.registered"
	NotificationPasswordResetRequest NotificationType = "password.reset_requested"
)

// NotificationEvent es el evento que llega desde SNS/SQS.
type NotificationEvent struct {
	Type      NotificationType  `json:"type"`
	Payload   map[string]string `json:"payload"`
	Timestamp time.Time         `json:"timestamp"`
}
