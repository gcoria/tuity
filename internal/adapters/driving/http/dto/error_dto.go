package dto

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func NewErrorResponse(error, message, details string) *ErrorResponse {
	return &ErrorResponse{
		Error:   error,
		Message: message,
		Details: details,
	}
}
