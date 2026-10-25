package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Errores canónicos de dominio según la sección 1.2 de .context/INTERFACES.md.
var (
	ErrInvalidInput  = errors.New("datos de entrada inválidos")
	ErrValidation    = errors.New("error de validación")
	ErrUnauthorized  = errors.New("no autorizado")
	ErrForbidden     = errors.New("acceso denegado")
	ErrNotFound      = errors.New("recurso no encontrado")
	ErrConflict      = errors.New("conflicto de estado o registro duplicado")
	ErrUnprocessable = errors.New("entidad no procesable")
)

// DomainError envuelve un error canónico de dominio con un mensaje legible para el usuario final.
type DomainError struct {
	Base    error
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

// MapDBError traduce errores de pgx / PostgreSQL (como unique_violation 23505
// o excepciones de triggers PL/pgSQL P0001 de RN07) a errores tipados de dominio.
func MapDBError(err error, defaultMsg string) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return NewDomainError(ErrNotFound, defaultMsg)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return NewDomainError(ErrConflict, defaultMsg)
		case "P0001": // raise_exception (ej. trg_prevent_self_approval RN07)
			return NewDomainError(ErrForbidden, pgErr.Message)
		case "23514": // check_violation
			return NewDomainError(ErrConflict, defaultMsg)
		}
	}

	return fmt.Errorf("%s: %w", defaultMsg, err)
}

// HTTPStatusAndMessage traduce cualquier error devuelto por un servicio SOA
// al código de estado HTTP y mensaje JSON correspondiente.
func HTTPStatusAndMessage(err error, fallbackMsg string) (int, string) {
	msg := fallbackMsg
	var domErr *DomainError
	if errors.As(err, &domErr) && domErr.Message != "" {
		msg = domErr.Message
	}

	switch {
	case errors.Is(err, ErrInvalidInput), errors.Is(err, ErrValidation):
		return http.StatusBadRequest, msg
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized, msg
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden, msg
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound, msg
	case errors.Is(err, ErrConflict):
		return http.StatusConflict, msg
	case errors.Is(err, ErrUnprocessable):
		return http.StatusUnprocessableEntity, msg
	default:
		return http.StatusInternalServerError, fallbackMsg
	}
}
