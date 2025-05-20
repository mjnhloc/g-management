package wraperror

type ApiDisplayableError struct {
	httpStatus int
	message    interface{}
	err        error
}

func NewApiDisplayableError(
	httpStatus int,
	message interface{},
	err error,
) *ApiDisplayableError {
	return &ApiDisplayableError{
		httpStatus: httpStatus,
		message:    message,
		err:        err,
	}
}

func (err *ApiDisplayableError) Error() string {
	if err.err != nil {
		return err.err.Error()
	}
	if message, messageIsString := err.message.(string); messageIsString {
		return message
	}

	return "Unknown error"
}

func (err *ApiDisplayableError) Unwrap() error {
	return err.err
}

func (err *ApiDisplayableError) Message() interface{} {
	return err.message
}

func (err *ApiDisplayableError) HttpStatus() int {
	return err.httpStatus
}
