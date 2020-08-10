package models

// ErrorBody type is used to pass any error information in the response body
type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}
