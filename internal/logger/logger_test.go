package logger

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLoggerOutput(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer

	// Create a custom logrus instance writing to our buffer
	log := logrus.New()
	log.SetOutput(&buf)
	log.SetFormatter(getFormatter())

	// Use our custom logger
	log.Info("Test message")

	// Verify log output contains expected text
	output := buf.String()
	if !strings.Contains(output, "Test message") {
		t.Errorf("Expected log to contain 'Test message', got: %s", output)
	}
	if !strings.Contains(output, "[netlinx-language-server]") {
		t.Errorf("Expected log to contain prefix, got: %s", output)
	}
}
