package service

import "net/http"

// httpError represents a custom http error.
type httpError struct {
	code          int    // code is the http status code of the error.
	message       string // message is the error message.
	internalError error
}

func (e *httpError) Error() string {
	return e.message
}

// newValidationError creates a custom validation error.
func newValidationError(err error) *httpError {
	return &httpError{
		code:    http.StatusBadRequest,
		message: err.Error(),
	}
}

// newBadRequestError creates a custom bad request error.
func newBadRequestError(message string) *httpError {
	return &httpError{
		code:    http.StatusBadRequest,
		message: message,
	}
}

// newNotFoundError creates a custom not found error.
func newNotFoundError(message string) *httpError {
	return &httpError{
		code:    http.StatusNotFound,
		message: message,
	}
}

// newUnauthorizedError creates a custom conflict error.
func newUnauthorizedError() *httpError {
	return &httpError{
		code:    http.StatusUnauthorized,
		message: "Unauthorized",
	}
}

// newConflictError creates a custom conflict error.
func newConflictError(message string) *httpError {
	return &httpError{
		code:    http.StatusConflict,
		message: message,
	}
}

// newInternalServerError creates a custom internal server error.
func newInternalServerError(err error) *httpError {
	return &httpError{
		code:          http.StatusInternalServerError,
		message:       "Something went wrong",
		internalError: err,
	}
}
