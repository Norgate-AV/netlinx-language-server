package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})

	Debug(msg string, fields logrus.Fields)
	Info(msg string, fields logrus.Fields)
	Warn(msg string, fields logrus.Fields)
	Error(msg string, fields logrus.Fields)

	LogRequest(method string, id interface{})
	LogResponse(method string, id interface{})
	LogNotification(method string)
	LogDocumentEvent(event string, uri string)
	LogServerEvent(event string)

	WithComponent(component string) Logger
}

type StructuredLogger struct {
	log       *logrus.Logger
	component string
}

func NewFileLogger(fileName string) (*StructuredLogger, error) {
	log := logrus.New()

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o666)
	if err != nil {
		return nil, err
	}

	log.SetOutput(file)
	log.SetFormatter(getFormatter())

	return &StructuredLogger{
		log:       log,
		component: "server",
	}, nil
}

func NewStdLogger() *StructuredLogger {
	log := logrus.New()

	log.SetOutput(os.Stderr)
	log.SetFormatter(getFormatter())

	return &StructuredLogger{
		log:       log,
		component: "server",
	}
}

func GetLogrusLogger(log Logger) *logrus.Logger {
	if sl, ok := log.(*StructuredLogger); ok {
		return sl.log
	}

	return nil
}

func (l *StructuredLogger) WithComponent(component string) Logger {
	return &StructuredLogger{
		log:       l.log,
		component: component,
	}
}

func (l *StructuredLogger) Info(msg string, fields logrus.Fields) {
	if fields == nil {
		fields = logrus.Fields{}
	}

	fields["component"] = l.component
	l.log.WithFields(logrus.Fields(fields)).Info(msg)
}

func (l *StructuredLogger) LogRequest(method string, id interface{}) {
	l.Info("Request received", logrus.Fields{
		"type":   "request",
		"method": method,
		"id":     id,
	})
}

func (l *StructuredLogger) LogResponse(method string, id interface{}) {
	l.Info("Response sent", logrus.Fields{
		"type":   "response",
		"method": method,
		"id":     id,
	})
}

func (l *StructuredLogger) LogNotification(method string) {
	l.Info("Notification received", logrus.Fields{
		"type":   "notification",
		"method": method,
	})
}

func (l *StructuredLogger) LogDocumentEvent(action string, uri string) {
	l.Info("Document event", logrus.Fields{
		"action": action,
		"uri":    uri,
	})
}

func (l *StructuredLogger) LogServerEvent(event string) {
	l.Info(event, logrus.Fields{
		"event_type": "server_lifecycle",
	})
}

func (l *StructuredLogger) Printf(format string, v ...interface{}) {
	l.log.Printf(format, v...)
}

func (l *StructuredLogger) Println(v ...interface{}) {
	l.log.Println(v...)
}

func (l *StructuredLogger) Debug(msg string, fields logrus.Fields) {
	if fields == nil {
		fields = logrus.Fields{}
	}

	fields["component"] = l.component
	l.log.WithFields(fields).Debug(msg)
}

func (l *StructuredLogger) Warn(msg string, fields logrus.Fields) {
	if fields == nil {
		fields = logrus.Fields{}
	}

	fields["component"] = l.component
	l.log.WithFields(fields).Warn(msg)
}

func (l *StructuredLogger) Error(msg string, fields logrus.Fields) {
	if fields == nil {
		fields = logrus.Fields{}
	}

	fields["component"] = l.component
	l.log.WithFields(fields).Error(msg)
}
