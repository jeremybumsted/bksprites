package sprites

import (
	"bytes"
	"sync"

	"github.com/charmbracelet/log"
)

// logWriter implements io.Writer that sends output to charmbracelet/log.
// It buffers writes until a newline is encountered, then logs the complete line.
type logWriter struct {
	logger *log.Logger
	level  log.Level
	mu     sync.Mutex
	buf    []byte
}

// newLogWriter creates a new logWriter that logs to the given logger at the specified level.
func newLogWriter(logger *log.Logger, level log.Level) *logWriter {
	return &logWriter{
		logger: logger,
		level:  level,
		buf:    make([]byte, 0, 256),
	}
}

// Write implements io.Writer. It buffers input until complete lines are received,
// then logs each line at the configured level.
func (w *logWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.buf = append(w.buf, p...)

	// Process complete lines
	for {
		idx := bytes.IndexByte(w.buf, '\n')
		if idx == -1 {
			break
		}

		line := string(w.buf[:idx])
		w.logLine(line)

		w.buf = w.buf[idx+1:]
	}

	return len(p), nil
}

// logLine logs a single line at the configured level.
func (w *logWriter) logLine(line string) {
	switch w.level {
	case log.DebugLevel:
		w.logger.Debug(line)
	case log.InfoLevel:
		w.logger.Info(line)
	case log.WarnLevel:
		w.logger.Warn(line)
	case log.ErrorLevel:
		w.logger.Error(line)
	default:
		w.logger.Info(line)
	}
}

// Flush logs any remaining buffered data, even if it doesn't end with a newline.
func (w *logWriter) Flush() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if len(w.buf) > 0 {
		line := string(w.buf)
		w.logLine(line)
		w.buf = w.buf[:0]
	}
}

// Close flushes any remaining buffered data and closes the writer.
func (w *logWriter) Close() error {
	w.Flush()
	return nil
}
