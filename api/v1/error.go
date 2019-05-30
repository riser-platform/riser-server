package v1

// APIError gets formatted and returned via the API
type APIError struct {
	// HTTPCode is the HTTP status code to return to the client
	HTTPCode int
	// Message is a client facing error message
	Message string
	// Internal error is used for fatal errors. This will get logged as an error. The client never sees this error.
	Internal error
}

func NewAPIError(httpCode int, message string) *APIError {
	return &APIError{HTTPCode: httpCode, Message: message}
}

func (e *APIError) Error() string {
	return e.Message
}
