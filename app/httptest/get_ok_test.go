package httptest

import (
	"context"
	"net/http"
	"testing"
)

func Test_GetOK(t *testing.T) {
	testCases := map[string]struct {
		expectStatusCode int
	}{
		"正常系": {
			expectStatusCode: 200,
		},
	}

	ctx := context.Background()
	url := "http://localhost:8080/ok"
	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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
