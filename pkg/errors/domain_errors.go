package errors

import "fmt"

type ErrorType string

const (
	ValidationError ErrorType = "validation_error"
	NotFoundError   ErrorType = "not_found_error"
	ConflictError   ErrorType = "conflict_error"
	InternalError   ErrorType = "internal_error"
)

type DomainError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func New(errorType ErrorType, message string, details ...string) *DomainError {
	var detail string
	if len(details) > 0 {
		detail = details[0]
	}
	return &DomainError{
		Type:    errorType,
		Message: message,
		Details: detail,
	}
}

func NewValidationError(message string) *DomainError {
	return New(ValidationError, message)
}

func NewNotFoundError(resource string) *DomainError {
	return New(NotFoundError, fmt.Sprintf("%s not found", resource))
}

func NewConflictError(message string) *DomainError {
	return New(ConflictError, message)
}

func NewInternalError(message string) *DomainError {
	return New(InternalError, message)
}
