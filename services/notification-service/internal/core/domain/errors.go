package domain

import "errors"

var (
	ErrUnknownEventType    = errors.New("unknown notification event type")
	ErrMissingPayloadField = errors.New("missing required payload field")
	ErrTemplateRender      = errors.New("failed to render email template")
	ErrEmailSendFailed     = errors.New("failed to send email")
)
