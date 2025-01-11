package helpers

import (
	"net/http"
)

type APIResponse struct {
	Error      bool        `json:"error"`
	StatusCode int         `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

// SuccessResponse creates a standardized success response.
func SuccessResponse(data interface{}, code int, message string) APIResponse {
	if message == "" {
		message = getStatusMessage(code)
	}
	return APIResponse{
		Error:      false,
		StatusCode: code,
		Message:    message,
		Data:       data,
	}
}

// ErrorResponse creates a standardized error response.
func ErrorResponse(data interface{}, code int, message string) APIResponse {
	if message == "" {
		message = getStatusMessage(code)
	}
	return APIResponse{
		Error:      true,
		StatusCode: code,
		Message:    message,
		Data:       data,
	}
}

// getStatusMessage maps HTTP status codes to default messages.
func getStatusMessage(code int) string {
	switch code {
	case http.StatusOK:
		return "OK"
	case http.StatusCreated:
		return "Created"
	case http.StatusNoContent:
		return "No Content"
	case http.StatusBadRequest:
		return "Bad Request"
	case http.StatusUnauthorized:
		return "Unauthorized"
	case http.StatusForbidden:
		return "Forbidden"
	case http.StatusNotFound:
		return "Not Found"
	case http.StatusUnprocessableEntity:
		return "Unprocessable Entity"
	case http.StatusInternalServerError:
		return "Internal Server Error"
	default:
		return "Unknown Code"
	}
}
