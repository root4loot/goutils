package httputil

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientWithOptionalResolvers(t *testing.T) {
	client, err := ClientWithOptionalResolvers()
	if err != nil {
		t.Fatalf("ClientWithOptionalResolvers failed: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestFindScheme(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	scheme, url, err := FindScheme(server.URL)
	if err != nil {
		t.Fatalf("FindScheme failed: %v", err)
	}

	if scheme != "http" {
		t.Errorf("Expected scheme 'http', got '%s'", scheme)
	}

	if url != server.URL {
		t.Errorf("Expected URL '%s', got '%s'", server.URL, url)
	}
}

func TestRedirectsToHTTPS(t *testing.T) {
	httpsServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer httpsServer.Close()

	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, httpsServer.URL, http.StatusMovedPermanently)
	}))
	defer httpServer.Close()

	redirects, finalURL, err := RedirectsToHTTPS(httpServer.URL)
	if err != nil {
		t.Fatalf("RedirectsToHTTPS failed: %v", err)
	}

	if !redirects {
		t.Errorf("Expected redirect to HTTPS, but it did not occur")
	}

	if finalURL != httpsServer.URL {
		t.Errorf("Expected final URL '%s', got '%s'", httpsServer.URL, finalURL)
	}
}

func TestIsBinaryResponse(t *testing.T) {
	tests := []struct {
		contentType string
		isBinary    bool
	}{
		{"application/octet-stream", true},
		{"text/html", false},
		{"image/jpeg", true},
		{"application/json", false},
	}

	for _, test := range tests {
		resp := &http.Response{
			Header: http.Header{
				"Content-Type": []string{test.contentType},
			},
		}

		if IsBinaryResponse(resp) != test.isBinary {
			t.Errorf("IsBinaryResponse(%s) = %v; want %v", test.contentType, !test.isBinary, test.isBinary)
		}
	}
}
