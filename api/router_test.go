package api

import (
	"errors"
	"net/http"
	"testing"

	"github.com/bmcszk/user-service/logic"
)

func Test_getStatusCode(t *testing.T) {
	tests := []struct {
		name         string
		givenErr     error
		expectedCode int
	}{
		{
			name:         "user not found",
			givenErr:     logic.ErrUserNotFound,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "user duplicate",
			givenErr:     logic.ErrUserAlreadyExists,
			expectedCode: http.StatusConflict,
		},
		{
			name:         "user not valid",
			givenErr:     logic.ErrUserNameEmpty,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "any error",
			givenErr:     errors.New("icecream on sidewalk"),
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getStatusCode(tt.givenErr); got != tt.expectedCode {
				t.Errorf("getStatusCode() = %v, want %v", got, tt.expectedCode)
			}
		})
	}
}
