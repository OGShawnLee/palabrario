package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Canonical domain errors mapped to HTTP statuses across core-service.
var (
	ErrInvalidInput  = errors.New("invalid input")
	ErrValidation    = errors.New("validation error")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrNotFound      = errors.New("resource not found")
	ErrConflict      = errors.New("conflict or duplicate record")
	ErrUnprocessable = errors.New("unprocessable entity")
)

// DomainError wraps a canonical domain error with a machine-readable code and user-facing message.
type DomainError struct {
	Base    error
	Code    string
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

func (e *DomainError) Unwrap() error {
	return e.Base
}

func NewDomainError(base error, message string) error {
	return &DomainError{Base: base, Message: message}
}

func NewDomainErrorWithCode(base error, code, message string) error {
	return &DomainError{Base: base, Code: code, Message: message}
}

// MapDBError translates pgx / PostgreSQL errors into typed domain errors with machine-readable codes.
func MapDBError(err error, defaultMsg string) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return NewDomainErrorWithCode(ErrNotFound, "NOT_FOUND", defaultMsg)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "22001": // string_data_right_truncation
			return NewDomainErrorWithCode(ErrInvalidInput, "INPUT_TOO_LONG", "Input exceeds maximum allowed length.")
		case "23514": // check_violation
			return NewDomainErrorWithCode(ErrValidation, "VALIDATION_ERROR", defaultMsg)
		case "23503": // foreign_key_violation
			return NewDomainErrorWithCode(ErrInvalidInput, "INVALID_REFERENCE", defaultMsg)
		case "23505": // unique_violation
			switch pgErr.ConstraintName {
			case "idx_accounts_email_lower", "accounts_email_key":
				return NewDomainErrorWithCode(ErrConflict, "EMAIL_ALREADY_EXISTS", "An account with that email already exists.")
			case "idx_users_user_name_lower", "users_user_name_key":
				return NewDomainErrorWithCode(ErrConflict, "USERNAME_ALREADY_EXISTS", "That username is already taken.")
			default:
				return NewDomainErrorWithCode(ErrConflict, "CONFLICT", defaultMsg)
			}
		case "P0001": // raise_exception (e.g., trg_prevent_self_approval RN07)
			return NewDomainErrorWithCode(ErrForbidden, "FORBIDDEN", pgErr.Message)
		}
	}

	return fmt.Errorf("%s: %w", defaultMsg, err)
}

// HTTPStatusAndMessage translates a service error into an HTTP status code and message.
func HTTPStatusAndMessage(err error, fallbackMsg string) (int, string) {
	status, _, msg := HTTPStatusCodeAndMessage(err, "", fallbackMsg)
	return status, msg
}

// HTTPStatusCodeAndMessage translates a service error into an HTTP status code, machine-readable code, and message.
func HTTPStatusCodeAndMessage(err error, fallbackCode, fallbackMsg string) (int, string, string) {
	msg := fallbackMsg
	code := ""
	var domErr *DomainError
	if errors.As(err, &domErr) {
		if domErr.Message != "" {
			msg = domErr.Message
		}
		if domErr.Code != "" {
			code = domErr.Code
		}
	}

	switch {
	case errors.Is(err, ErrInvalidInput):
		if code == "" {
			code = "INVALID_INPUT"
		}
		return http.StatusBadRequest, code, msg
	case errors.Is(err, ErrValidation):
		if code == "" {
			code = "VALIDATION_ERROR"
		}
		return http.StatusBadRequest, code, msg
	case errors.Is(err, ErrUnauthorized):
		if code == "" {
			code = "UNAUTHORIZED"
		}
		return http.StatusUnauthorized, code, msg
	case errors.Is(err, ErrForbidden):
		if code == "" {
			code = "FORBIDDEN"
		}
		return http.StatusForbidden, code, msg
	case errors.Is(err, ErrNotFound):
		if code == "" {
			code = "NOT_FOUND"
		}
		return http.StatusNotFound, code, msg
	case errors.Is(err, ErrConflict):
		if code == "" {
			code = "CONFLICT"
		}
		return http.StatusConflict, code, msg
	case errors.Is(err, ErrUnprocessable):
		if code == "" {
			code = "UNPROCESSABLE_ENTITY"
		}
		return http.StatusUnprocessableEntity, code, msg
	default:
		if fallbackCode != "" {
			code = fallbackCode
		} else {
			code = "INTERNAL_ERROR"
		}
		return http.StatusInternalServerError, code, fallbackMsg
	}
}
