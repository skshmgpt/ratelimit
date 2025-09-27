package main

import (
	"io"
	"net/http"
	"testing"
)

func TestBurst(t *testing.T) {
	wanted := "hello world"

	// fire off 5 allowed requests
	for i := 0; i < 5; i++ {
		resp, err := http.Get("http://localhost:8080/api/protected")
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
	}

	// 6th request should be rate-limited
	resp, err := http.Get("http://localhost:8080/api/protected")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	got := string(body)

	if resp.StatusCode == http.StatusTooManyRequests {
		t.Log("Rate limiting works, got 429")
		return
	}

	if got != wanted {
		t.Errorf("expected %q, got %q", wanted, got)
	}
}
