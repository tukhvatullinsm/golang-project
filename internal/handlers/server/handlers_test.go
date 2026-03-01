package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetValues(t *testing.T) {
	handler := &WebApp{}
	//server := httptest.NewServer(handler)
	//defer server.Close()

	tests := []struct {
		name      string
		url       string
		method    string
		status    int
		mediaType string
	}{
		{
			name:      "Incorrect URL Parameters",
			url:       "/update/gauge/Test/Test",
			method:    http.MethodPost,
			status:    http.StatusBadRequest,
			mediaType: "text/plain",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.method, test.url, nil)
			req.Header.Set("Content-Type", test.mediaType)
			w := httptest.NewRecorder()
			handler.SetValues(w, req)
			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != test.status {
				t.Errorf("Response status code is not %v: %v", test.status, resp.StatusCode)
			}
		})
	}

}
