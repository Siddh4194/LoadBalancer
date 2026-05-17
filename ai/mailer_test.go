package ai

import (
	"strings"
	"testing"
)

func TestBuildSMTPMessage(t *testing.T) {
	from := "sender@example.com"
	to := "recipient@example.com"
	subject := "Test Subject"
	body := "<h1>Test</h1><p>Hello world</p>"

	msg := buildSMTPMessage(from, to, subject, body)

	if !strings.Contains(msg, "From: sender@example.com\r\n") {
		t.Fatalf("expected From header, got %q", msg)
	}

	if !strings.Contains(msg, "To: recipient@example.com\r\n") {
		t.Fatalf("expected To header, got %q", msg)
	}

	if !strings.Contains(msg, "Subject: Test Subject\r\n") {
		t.Fatalf("expected Subject header, got %q", msg)
	}

	if !strings.Contains(msg, "Content-Type: text/html; charset=\"UTF-8\"\r\n") {
		t.Fatalf("expected Content-Type header, got %q", msg)
	}

	if !strings.HasSuffix(msg, body) {
		t.Fatalf("expected body to be appended, got %q", msg)
	}
}
