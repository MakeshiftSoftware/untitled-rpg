package service

import "net/http"

// httpError represents a custom http error, defined by an http status
// code and error message.
type httpError struct {
	code          int    // code is the http status code of the error.
	message       string // message is the error message.
	internalError error  // internalError is the original error that occurred.
}

// Error returns the error message satisfying the Error interface.
func (e *httpError) Error() string {
	return e.message
}

// newValidationError creates a custom validation error.
// This error is typically returned to the client when some data is provided
// that does not pass validation, or if the data is missing some required fields.
func newValidationError(err error) *httpError {
	return &httpError{
		code:    http.StatusBadRequest,
		message: err.Error(),
	}
}

// newBadRequestError creates a custom bad request error.
// This error is typically returned to the client if the client provides incorrectly
// formatted data or leaves out a request body.
func newBadRequestError(message string) *httpError {
	return &httpError{
		code:    http.StatusBadRequest,
		message: message,
	}
}

// newNotFoundError creates a custom not found error.
// This error is typically returned to the client when the requested resource
// does not exist.
func newNotFoundError(message string) *httpError {
	return &httpError{
		code:    http.StatusNotFound,
		message: message,
	}
}

// newUnauthorizedError creates a custom conflict error.
// This error is typically returned to the client following an unsuccessful login attempt
// or if an invalid auth token is provided when interacting with a service.
func newUnauthorizedError() *httpError {
	return &httpError{
		code:    http.StatusUnauthorized,
		message: "Unauthorized",
	}
}

// newConflictError creates a custom conflict error.
// This error is typically returned to the client when a unique index is violated.
func newConflictError(message string) *httpError {
	return &httpError{
		code:    http.StatusConflict,
		message: message,
	}
}

// newInternalServerError creates a custom internal server error.
// This error is typically returned to the client when an unexpected or unknown error occurs.
func newInternalServerError(err error) *httpError {
	return &httpError{
		code:          http.StatusInternalServerError,
		message:       "Something went wrong",
		internalError: err,
	}
}
