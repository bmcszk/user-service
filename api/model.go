package api

import "fmt"

type ApiError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func (e ApiError) Error() string {
	return fmt.Sprintf("status_code: %d, message: %s", e.StatusCode, e.Message)
}
