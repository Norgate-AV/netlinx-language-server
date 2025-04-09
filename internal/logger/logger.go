package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

func NewFileLogger(fileName string) (Logger, error) {
	log := logrus.New()

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o666)
	if err != nil {
		return nil, err
	}

	log.SetOutput(file)
	log.SetFormatter(getFormatter())

	return log, nil
}

func NewStdLogger() Logger {
	log := logrus.New()

	log.SetOutput(os.Stderr)
	log.SetFormatter(getFormatter())

	return log
}

func GetLogrusLogger(log Logger) *logrus.Logger {
	if l, ok := log.(*logrus.Logger); ok {
		return l
	}

	return nil
}

func getFormatter() *PrefixFormatter {
	return &PrefixFormatter{
		Prefix: "[netlinx-language-server]",
		Formatter: &logrus.TextFormatter{
			FullTimestamp: true,
		},
	}
}
