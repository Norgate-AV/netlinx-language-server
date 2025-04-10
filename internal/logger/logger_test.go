package logger

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLoggerOutput(t *testing.T) {
	var buf bytes.Buffer

	log := logrus.New()
	log.SetOutput(&buf)
	log.SetFormatter(getFormatter())

	log.Info("Test message")

	output := buf.String()
	if !strings.Contains(output, "Test message") {
		t.Errorf("Expected log to contain 'Test message', got: %s", output)
	}

	if !strings.Contains(output, "[netlinx-language-server]") {
		t.Errorf("Expected log to contain prefix, got: %s", output)
	}
}

func TestStructuredLoggerOutput(t *testing.T) {
	var buf bytes.Buffer

	log := logrus.New()
	log.SetOutput(&buf)
	log.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: true,
	})

	structLog := &StructuredLogger{
		log:       log,
		component: "test-component",
	}

	structLog.Info("Test info message", logrus.Fields{
		"custom_field": "custom_value",
	})

	var output map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	// Verify fields
	if output["level"] != "info" {
		t.Errorf("Expected level 'info', got %v", output["level"])
	}

	if output["msg"] != "Test info message" {
		t.Errorf("Expected message 'Test info message', got %v", output["msg"])
	}

	if output["component"] != "test-component" {
		t.Errorf("Expected component 'test-component', got %v", output["component"])
	}

	if output["custom_field"] != "custom_value" {
		t.Errorf("Expected custom_field 'custom_value', got %v", output["custom_field"])
	}
}

func TestSpecializedLogMethods(t *testing.T) {
	tests := []struct {
		name           string
		logFunc        func(logger *StructuredLogger)
		expectedFields map[string]interface{}
	}{
		{
			name: "LogRequest",
			logFunc: func(logger *StructuredLogger) {
				logger.LogRequest("initialize", 1)
			},
			expectedFields: map[string]interface{}{
				"msg":       "Request received",
				"type":      "request",
				"method":    "initialize",
				"id":        float64(1), // JSON numbers are floats when unmarshaled
				"level":     "info",
				"component": "test-component",
			},
		},
		{
			name: "LogResponse",
			logFunc: func(logger *StructuredLogger) {
				logger.LogResponse("initialize", 1)
			},
			expectedFields: map[string]interface{}{
				"msg":       "Response sent",
				"type":      "response",
				"method":    "initialize",
				"id":        float64(1),
				"level":     "info",
				"component": "test-component",
			},
		},
		{
			name: "LogNotification",
			logFunc: func(logger *StructuredLogger) {
				logger.LogNotification("exit")
			},
			expectedFields: map[string]interface{}{
				"msg":       "Notification received",
				"type":      "notification",
				"method":    "exit",
				"level":     "info",
				"component": "test-component",
			},
		},
		{
			name: "LogDocumentEvent",
			logFunc: func(logger *StructuredLogger) {
				logger.LogDocumentEvent("open", "file:///test.axs")
			},
			expectedFields: map[string]interface{}{
				"msg":       "Document event",
				"action":    "open",
				"uri":       "file:///test.axs",
				"level":     "info",
				"component": "test-component",
			},
		},
		{
			name: "LogServerEvent",
			logFunc: func(logger *StructuredLogger) {
				logger.LogServerEvent("Starting")
			},
			expectedFields: map[string]interface{}{
				"msg":        "Starting",
				"event_type": "server_lifecycle",
				"level":      "info",
				"component":  "test-component",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer

			log := logrus.New()
			log.SetOutput(&buf)
			log.SetFormatter(&logrus.JSONFormatter{
				DisableTimestamp: true,
			})

			structLog := &StructuredLogger{
				log:       log,
				component: "test-component",
			}

			test.logFunc(structLog)

			var output map[string]interface{}
			if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
				t.Fatalf("Failed to parse log output: %v", err)
			}

			// Verify fields
			for key, expected := range test.expectedFields {
				if output[key] != expected {
					t.Errorf("Expected %s=%v, got %v", key, expected, output[key])
				}
			}
		})
	}
}

func TestWithComponent(t *testing.T) {
	var buf bytes.Buffer

	log := logrus.New()
	log.SetOutput(&buf)
	log.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: true,
	})

	baseLogger := &StructuredLogger{
		log:       log,
		component: "base",
	}

	// Create derived logger with new component
	derivedLogger := baseLogger.WithComponent("derived")
	derivedLogger.Info("Test component message", nil)

	var output map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
		t.Fatalf("Failed to parse log output: %v", err)
	}

	// Check component was changed
	if output["component"] != "derived" {
		t.Errorf("Expected component 'derived', got %v", output["component"])
	}
}

func TestFileLogger(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "log-test-*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	logger, err := NewFileLogger(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create file logger: %v", err)
	}

	logger.Info("Test file logging", nil)

	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	// Verify log contains expected text
	if !strings.Contains(string(content), "Test file logging") {
		t.Errorf("Expected log to contain 'Test file logging', got: %s", content)
	}
}
