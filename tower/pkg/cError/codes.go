package cError

type ErrorCode string

const (
	ErrInternalServer ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrInvalidInput   ErrorCode = "INVALID_INPUT"
	ErrUnauthorized   ErrorCode = "UNAUTHORIZED"

	// 유저 관련
	ErrUserNotFound       ErrorCode = "USER_NOT_FOUND"
	ErrEmailAlreadyExists ErrorCode = "EMAIL_ALREADY_EXISTS"
	ErrPasswordMismatch   ErrorCode = "PASSWORD_MISMATCH"

	// 인증 관련
	ErrTokenExpired ErrorCode = "TOKEN_EXPIRED"
	ErrInvalidToken ErrorCode = "INVALID_TOKEN"
)
