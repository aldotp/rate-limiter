package http_test

import (
	"context"
	"testing"

	"github.com/aldotp/rate-limiter/cmd/http"
)

func TestNewHTTPServer(t *testing.T) {
	ctx := context.Background()

	router, addr, err := http.NewHTTPServer(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if router == nil {
		t.Fatal("Expected router to be initialized, got nil")
	}

	if addr == "" {
		t.Fatal("Expected listen address to be set, got empty string")
	}
}
