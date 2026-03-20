package handler

import sharedhttp "github.com/Manuelda10/pcshop-backend/shared/http"

const (
	CodeUserCreated        sharedhttp.ResponseCode = "USER_CREATED"
	CodeEmailAlreadyExists sharedhttp.ResponseCode = "EMAIL_ALREADY_EXISTS"
	CodeInvalidCredentials sharedhttp.ResponseCode = "INVALID_CREDENTIALS"
)
