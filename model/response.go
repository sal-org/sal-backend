package model

// SuccessResponse .
type SuccessResponse struct {
	Meta Meta `json:"meta"`
}

// ErrorResponse .
type ErrorResponse struct {
	Meta Meta `json:"meta"`
}

// Meta .
type Meta struct {
	Message     string `json:"message,omitempty"`
	MessageType string `json:"message_type,omitempty"`
	Status      string `json:"status,omitempty"`
}
