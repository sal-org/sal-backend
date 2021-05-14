package model

// RazorPayTransaction .
type RazorPayTransaction struct {
	ID          string `json:"id"`
	Amount      int    `json:"amount"`
	Status      string `json:"status"`
	Description string `json:"description"`
}
