package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (h *writerHook) Fire(entry *logrus.Entry) error {
	s, err := entry.String()
	if err != nil {
		return err
	}

	for _, writer := range h.Writer {
		_, err = writer.Write([]byte(s))
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *writerHook) Levels() []logrus.Level {
	return h.LogLevels
}

func GetLogger() *Logger {
	return &Logger{e}
}

func init() {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
		DisableColors: false,
		FullTimestamp: true,
	}

	err := os.MkdirAll("logs", 0755)
	if err != nil {
		panic(err)
	}

	logFile, err := os.OpenFile("logs/log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	l.SetOutput(io.Discard)
	l.AddHook(&writerHook{
		Writer:    []io.Writer{logFile},
		LogLevels: logrus.AllLevels,
	})

	l.SetLevel(logrus.TraceLevel)

	e = logrus.NewEntry(l)
}
