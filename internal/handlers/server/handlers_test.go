package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeHTTP(t *testing.T) {
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
			name:      "Incorrect HTTP Method",
			url:       "/",
			method:    http.MethodGet,
			status:    http.StatusMethodNotAllowed,
			mediaType: "text/plain",
		},
		{
			name:      "Incorrect Media Type",
			url:       "/update/",
			method:    http.MethodPost,
			status:    http.StatusUnsupportedMediaType,
			mediaType: "text/html",
		},
		{
			name:      "Incomplete URL Path",
			url:       "/update/gauge/10",
			method:    http.MethodPost,
			status:    http.StatusNotFound,
			mediaType: "text/plain",
		},
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
			handler.ServeHTTP(w, req)
			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != test.status {
				t.Errorf("Response status code is not %d: %d", test.status, resp.StatusCode)
			}
		})
	}

}
