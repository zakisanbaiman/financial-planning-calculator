package controllers

// ErrorResponse は統一されたエラーレスポンス形式
type ErrorResponse struct {
	Error     string `json:"error"`
	Details   string `json:"details,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}
