package main

import (
	"io"
	"net/http"
	"testing"
)

func TestBurst(t *testing.T) {
	wanted := "hello world"
	req, _ := http.NewRequest("GET", "http://localhost:8080/api/unprotected", nil)
	// req.Header.Set("X-Api-Key", "68d82305-d680-8325-83c7-971f03a1ce46")
	client := &http.Client{}

	// fire off 5 allowed requests
	for i := 0; i < 20; i++ {
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
	}

	// 6th request should be rate-limited

	resp, err := client.Do(req)

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
