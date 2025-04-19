package util

import (
	"fmt"
	"io"
	"os"
)

type AppLogWriter struct {
	LogFile *os.File
}

func NewAppLogWriter(logFile *os.File) *AppLogWriter {
	return &AppLogWriter{LogFile: logFile}
}

func (t *AppLogWriter) Write(p []byte) (int, error) {
	os.Stdout.Write(p)
	n, err := t.LogFile.Write(p)
	if err != nil {
		ShowErrorDialog(fmt.Sprintf("failed to write log file: %v", err))
		return n, err
	}
	if n != len(p) {
		err := io.ErrShortWrite
		ShowErrorDialog(fmt.Sprintf("failed to write log file: %v", err))
		return n, err
	}
	if err := t.LogFile.Sync(); err != nil {
		ShowErrorDialog(fmt.Sprintf("failed to sync log file: %v", err))
		return n, nil
	}
	return n, nil
}

func MultiWriter(writers ...io.Writer) io.Writer {
	allWriters := make([]io.Writer, 0, len(writers))
	for _, w := range writers {
		if mw, ok := w.(*multiWriter); ok {
			allWriters = append(allWriters, mw.writers...)
		} else {
			allWriters = append(allWriters, w)
		}
	}
	return &multiWriter{allWriters}
}

type multiWriter struct {
	writers []io.Writer
}

func (t *multiWriter) Write(p []byte) (int, error) {
	for _, w := range t.writers {
		_, err := w.Write(p)
		if err != nil {
			continue
		}
		//if n != len(p) {
		//	continue
		//}
	}
	return len(p), nil
}
