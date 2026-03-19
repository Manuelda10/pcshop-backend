package http

type ResponseCode string

const (
	CodeInvalidBody     ResponseCode = "INVALID_BODY"
	CodeValidationError ResponseCode = "VALIDATION_ERROR"
	CodeInternalError   ResponseCode = "INTERNAL_ERROR"
	CodeNotFound        ResponseCode = "NOT_FOUND"
	CodeUnauthorized    ResponseCode = "UNAUTHORIZED"
	CodeInvalidToken    ResponseCode = "INVALID_TOKEN"
	CodeForbidden       ResponseCode = "FORBIDDEN"
)
