package wraperror

import (
	"net/http"
)

type ValidationError struct {
	messages map[string]interface{}
	Err      error
}

func NewValidationError(
	messages map[string]interface{},
	err error,
) *ValidationError {
	return &ValidationError{
		messages: messages,
		Err: NewApiDisplayableError(
			http.StatusBadRequest,
			messages,
			err,
		),
	}
}

func (err *ValidationError) Error() string {
	return err.Err.Error()
}

func (err *ValidationError) Unwrap() error {
	return err.Err
}
