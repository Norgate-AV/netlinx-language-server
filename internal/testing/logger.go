package testing

import (
	"testing"

	"github.com/Norgate-AV/netlinx-language-server/internal/logger"
	"github.com/sirupsen/logrus"
)

type TestLogger struct {
	t *testing.T
}

func NewTestLogger(t *testing.T) *TestLogger {
	return &TestLogger{t: t}
}

func (l *TestLogger) Debug(message string, fields logrus.Fields) {
	l.t.Logf("[DEBUG] %s %v", message, fields)
}

func (l *TestLogger) Info(message string, fields logrus.Fields) {
	l.t.Logf("[INFO] %s %v", message, fields)
}

func (l *TestLogger) Warn(message string, fields logrus.Fields) {
	l.t.Logf("[WARN] %s %v", message, fields)
}

func (l *TestLogger) Error(message string, fields logrus.Fields) {
	l.t.Logf("[ERROR] %s %v", message, fields)
}

func (l *TestLogger) LogServerEvent(event string) {
	l.t.Logf("[SERVER] %s", event)
}

func (l *TestLogger) LogClientEvent(message string, args ...any) {
	l.t.Logf("[CLIENT] %s", message)
}

func (l *TestLogger) GetFilePath() string {
	return "/path/to/logfile.log"
}

// Printf implements the same method from Logger interface
func (l *TestLogger) Printf(format string, v ...any) {
	l.t.Logf(format, v...)
}

// Println implements the same method from Logger interface
func (l *TestLogger) Println(v ...any) {
	l.t.Log(v...)
}

// LogMessage implements the same method from Logger interface
func (l *TestLogger) LogMessage(method string, id any) {
	l.t.Logf("[MESSAGE] Method: %s, ID: %v", method, id)
}

// LogRequest implements the same method from Logger interface
func (l *TestLogger) LogRequest(method string, id any) {
	l.t.Logf("[REQUEST] Method: %s, ID: %v", method, id)
}

// LogResponse implements the same method from Logger interface
func (l *TestLogger) LogResponse(method string, id any) {
	l.t.Logf("[RESPONSE] Method: %s, ID: %v", method, id)
}

// LogNotification implements the same method from Logger interface
func (l *TestLogger) LogNotification(method string) {
	l.t.Logf("[NOTIFICATION] Method: %s", method)
}

// LogDocumentEvent implements the same method from Logger interface
func (l *TestLogger) LogDocumentEvent(event string, uri string) {
	l.t.Logf("[DOCUMENT] Event: %s, URI: %s", event, uri)
}

// WithComponent implements the same method from Logger interface
func (l *TestLogger) WithComponent(component string) logger.Logger {
	return &TestLogger{t: l.t}
}
