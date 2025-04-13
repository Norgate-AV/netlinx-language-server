package logger

import (
	"os"

	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"
	"github.com/sirupsen/logrus"
)

type Logger interface {
	Printf(format string, v ...any)
	Println(v ...any)

	Debug(msg string, fields logrus.Fields)
	Info(msg string, fields logrus.Fields)
	Warn(msg string, fields logrus.Fields)
	Error(msg string, fields logrus.Fields)

	LogMessage(method string, id any)
	LogRequest(method string, id any)
	LogResponse(method string, id any)
	LogNotification(method string)
	LogDocumentEvent(event string, uri string)
	LogServerEvent(event string)

	WithComponent(Component string) Logger
}

type StructuredLogger struct {
	Log       *logrus.Logger
	Component string
}

func NewFileLogger(fileName string) (*StructuredLogger, error) {
	log := logrus.New()

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o666)
	if err != nil {
		return nil, err
	}

	log.SetOutput(file)
	log.SetFormatter(GetFormatter())

	return &StructuredLogger{
		Log:       log,
		Component: "server",
	}, nil
}

func NewStdLogger() *StructuredLogger {
	log := logrus.New()

	log.SetOutput(os.Stderr)
	log.SetFormatter(GetFormatter())

	return &StructuredLogger{
		Log:       log,
		Component: "server",
	}
}

func GetLogrusLogger(log Logger) *logrus.Logger {
	if sl, ok := log.(*StructuredLogger); ok {
		return sl.Log
	}

	return nil
}

func (l *StructuredLogger) WithComponent(Component string) Logger {
	return &StructuredLogger{
		Log:       l.Log,
		Component: Component,
	}
}

func (l *StructuredLogger) Info(msg string, fields logrus.Fields) {
	if fields == nil {
		fields = logrus.Fields{}
	}

	fields["component"] = l.Component
	l.Log.WithFields(logrus.Fields(fields)).Info(msg)
}

func (l *StructuredLogger) LogMessage(method string, id any) {
	if lsp.IsNotification(method) {
		l.LogNotification(method)
	} else {
		l.LogRequest(method, id)
	}
}

func (l *StructuredLogger) LogRequest(method string, id any) {
	l.Info("Request Received", logrus.Fields{
		"type":   "request",
		"method": method,
		"id":     id,
	})
}

func (l *StructuredLogger) LogResponse(method string, id any) {
	l.Info("Response Sent", logrus.Fields{
		"type":   "response",
		"method": method,
		"id":     id,
	})
}

func (l *StructuredLogger) LogNotification(method string) {
	l.Info("Notification Received", logrus.Fields{
		"type":   "notification",
		"method": method,
	})
}

func (l *StructuredLogger) LogDocumentEvent(action string, uri string) {
	l.Info("Document Event", logrus.Fields{
		"action": action,
		"uri":    uri,
	})
}

func (l *StructuredLogger) LogServerEvent(event string) {
	l.Info(event, logrus.Fields{
		"event_type": "server_lifecycle",
	})
}

func (l *StructuredLogger) Printf(format string, v ...any) {
	l.Log.Printf(format, v...)
}

func (l *StructuredLogger) Println(v ...any) {
	l.Log.Println(v...)
}

func (l *StructuredLogger) Debug(msg string, fields logrus.Fields) {
	if fields == nil {
		fields = logrus.Fields{}
	}

	fields["component"] = l.Component
	l.Log.WithFields(fields).Debug(msg)
}

func (l *StructuredLogger) Warn(msg string, fields logrus.Fields) {
	if fields == nil {
		fields = logrus.Fields{}
	}

	fields["component"] = l.Component
	l.Log.WithFields(fields).Warn(msg)
}

func (l *StructuredLogger) Error(msg string, fields logrus.Fields) {
	if fields == nil {
		fields = logrus.Fields{}
	}

	fields["component"] = l.Component
	l.Log.WithFields(fields).Error(msg)
}
