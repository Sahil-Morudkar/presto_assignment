package model

// APIResponse defines a consistent API response format
type APIResponse struct {
	Status  string      `json:"status"`            // "success" or "error"
	Message string      `json:"message"`           // Human-readable message
	Data    interface{} `json:"data,omitempty"`    // Optional payload
}