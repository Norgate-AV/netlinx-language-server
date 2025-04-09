package logger

import (
	"bytes"

	"github.com/sirupsen/logrus"
)

type PrefixFormatter struct {
	Prefix    string
	Formatter logrus.Formatter
}

func (f *PrefixFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	formatted, err := f.Formatter.Format(entry)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.Write(formatted[:bytes.IndexByte(formatted, ' ')+1])
	buf.WriteString(f.Prefix + " ")
	buf.Write(formatted[bytes.IndexByte(formatted, ' ')+1:])

	return buf.Bytes(), nil
}
