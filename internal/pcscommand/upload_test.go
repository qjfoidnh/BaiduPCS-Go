package pcscommand

import (
	"strings"
	"testing"
)

func TestUploadFailureError(t *testing.T) {
	if err := uploadFailureError(0); err != nil {
		t.Fatalf("zero failures returned error: %v", err)
	}

	err := uploadFailureError(2)
	if err == nil {
		t.Fatal("expected non-zero failures to return an error")
	}
	if !strings.Contains(err.Error(), "2") {
		t.Fatalf("expected failure count in error, got %v", err)
	}
}
