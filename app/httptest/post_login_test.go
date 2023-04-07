package httptest

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func Test_PostLogin(t *testing.T) {
	type arg struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}
	testCases := map[string]struct {
		arg              arg
		expectStatusCode int
	}{
		"正常系": {
			arg:              arg{"user", "password"},
			expectStatusCode: 200,
		},
		"認証エラー": {
			arg:              arg{"user", "Password"},
			expectStatusCode: 401,
		},
	}

	ctx := context.Background()
	url := "http://localhost:8080/login"
	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			b, err := json.Marshal(tt.arg)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(b))
			if err != nil {
				t.Fatal(err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if tt.expectStatusCode != resp.StatusCode {
				t.Errorf("expect status code %d, but got %d", tt.expectStatusCode, resp.StatusCode)
			}
		})
	}
}
