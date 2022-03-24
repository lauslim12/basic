package basic

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Tests the whole codebase.
func TestAuthenticate(t *testing.T) {
	users := map[string]string{"gerysantoso": "gerysantoso"}
	tests := []struct {
		name           string
		username       string
		password       string
		auth           *BasicAuth
		expectedStatus int
	}{
		{
			name:           "test_success",
			username:       "gerysantoso",
			password:       "gerysantoso",
			auth:           NewDefaultBasicAuth(users),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "test_invalid_scheme",
			username:       "",
			password:       "",
			auth:           NewDefaultBasicAuth(nil),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "test_invalid_credentials",
			username:       "gery",
			password:       "gery",
			auth:           NewDefaultBasicAuth(users),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "test_www_authenticate",
			username:       "",
			password:       "",
			auth:           NewCustomBasicAuth(nil, "UTF-8", nil, nil, "Test", users),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "test_empty_charset_users_and_realm",
			username:       "",
			password:       "",
			auth:           NewCustomBasicAuth(nil, "", nil, nil, "", nil),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := tc.auth.Authenticate(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()

			if tc.username != "" && tc.password != "" {
				r.SetBasicAuth(tc.username, tc.password)
			}

			handler(w, r)

			if tc.expectedStatus != w.Code {
				t.Errorf("Expected and actual status code values are different! Expected: %v. Got: %v.", tc.expectedStatus, w.Code)
			}
		})
	}
}
