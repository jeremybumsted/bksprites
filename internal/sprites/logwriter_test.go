package sprites

import (
	"bytes"
	"strings"
	"sync"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
)

// TestLogWriter_SingleCompleteLine tests writing a single line with newline
func TestLogWriter_SingleCompleteLine(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf)
	logger.SetReportTimestamp(false) // Disable timestamp for easier testing

	writer := newLogWriter(logger, log.InfoLevel)

	// Write a complete line
	n, err := writer.Write([]byte("hello world\n"))
	assert.NoError(t, err)
	assert.Equal(t, 12, n)

	// Verify it was logged
	output := buf.String()
	assert.Contains(t, output, "hello world")

	// Verify buffer is empty (nothing left to flush)
	assert.Equal(t, 0, len(writer.buf))
}

// TestLogWriter_MultipleLines tests writing multiple complete lines in one Write() call
func TestLogWriter_MultipleLines(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf)
	logger.SetReportTimestamp(false)

	writer := newLogWriter(logger, log.InfoLevel)

	// Write multiple complete lines at once
	n, err := writer.Write([]byte("line1\nline2\nline3\n"))
	assert.NoError(t, err)
	assert.Equal(t, 18, n)

	// Verify all three lines are logged
	output := buf.String()
	assert.Contains(t, output, "line1")
	assert.Contains(t, output, "line2")
	assert.Contains(t, output, "line3")

	// Count occurrences - should have 3 INFO entries
	infoCount := strings.Count(output, "INFO")
	assert.Equal(t, 3, infoCount, "expected 3 INFO log entries")

	// Verify buffer is empty after
	assert.Equal(t, 0, len(writer.buf))
}

// TestLogWriter_PartialLine tests writing data without newline
func TestLogWriter_PartialLine(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf)
	logger.SetReportTimestamp(false)

	writer := newLogWriter(logger, log.InfoLevel)

	// Write partial line (no newline)
	n, err := writer.Write([]byte("hello"))
	assert.NoError(t, err)
	assert.Equal(t, 5, n)

	// Verify nothing is logged yet
	output := buf.String()
	assert.Empty(t, output)

	// Verify data remains in buffer
	assert.Equal(t, 5, len(writer.buf))
	assert.Equal(t, "hello", string(writer.buf))
}

// TestLogWriter_PartialThenComplete tests writing partial data then completing the line
func TestLogWriter_PartialThenComplete(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf)
	logger.SetReportTimestamp(false)

	writer := newLogWriter(logger, log.InfoLevel)

	// Write "hello" (no newline)
	n, err := writer.Write([]byte("hello"))
	assert.NoError(t, err)
	assert.Equal(t, 5, n)

	// Verify nothing logged yet
	assert.Empty(t, buf.String())

	// Write " world\n" to complete the line
	n, err = writer.Write([]byte(" world\n"))
	assert.NoError(t, err)
	assert.Equal(t, 7, n)

	// Verify "hello world" is logged as one line
	output := buf.String()
	assert.Contains(t, output, "hello world")

	// Verify buffer is empty
	assert.Equal(t, 0, len(writer.buf))
}

// TestLogWriter_FlushPartialLine tests Flush() with incomplete line
func TestLogWriter_FlushPartialLine(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf)
	logger.SetReportTimestamp(false)

	writer := newLogWriter(logger, log.InfoLevel)

	// Write "incomplete" (no newline)
	n, err := writer.Write([]byte("incomplete"))
	assert.NoError(t, err)
	assert.Equal(t, 10, n)

	// Verify nothing logged yet
	assert.Empty(t, buf.String())

	// Call Flush()
	writer.Flush()

	// Verify "incomplete" gets logged
	output := buf.String()
	assert.Contains(t, output, "incomplete")

	// Verify buffer is empty
	assert.Equal(t, 0, len(writer.buf))
}

// TestLogWriter_Close tests Close() flushes partial line
func TestLogWriter_Close(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf)
	logger.SetReportTimestamp(false)

	writer := newLogWriter(logger, log.InfoLevel)

	// Write partial line
	n, err := writer.Write([]byte("partial"))
	assert.NoError(t, err)
	assert.Equal(t, 7, n)

	// Verify nothing logged yet
	assert.Empty(t, buf.String())

	// Call Close()
	err = writer.Close()
	assert.NoError(t, err)

	// Verify partial line gets logged (Close should call Flush)
	output := buf.String()
	assert.Contains(t, output, "partial")

	// Verify buffer is empty
	assert.Equal(t, 0, len(writer.buf))
}

// TestLogWriter_DifferentLevels tests logging at different levels
func TestLogWriter_DifferentLevels(t *testing.T) {
	tests := []struct {
		name          string
		level         log.Level
		expectedLevel string
	}{
		{
			name:          "debug level",
			level:         log.DebugLevel,
			expectedLevel: "DEBU",
		},
		{
			name:          "info level",
			level:         log.InfoLevel,
			expectedLevel: "INFO",
		},
		{
			name:          "warn level",
			level:         log.WarnLevel,
			expectedLevel: "WARN",
		},
		{
			name:          "error level",
			level:         log.ErrorLevel,
			expectedLevel: "ERRO",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := log.New(&buf)
			logger.SetReportTimestamp(false)
			logger.SetLevel(log.DebugLevel) // Set to debug to capture all levels

			writer := newLogWriter(logger, tt.level)

			// Write a complete line
			_, err := writer.Write([]byte("test message\n"))
			assert.NoError(t, err)

			// Verify the correct level appears in output
			output := buf.String()
			assert.Contains(t, output, tt.expectedLevel)
			assert.Contains(t, output, "test message")
		})
	}
}

// TestLogWriter_EmptyWrites tests writing empty byte slices
func TestLogWriter_EmptyWrites(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf)
	logger.SetReportTimestamp(false)

	writer := newLogWriter(logger, log.InfoLevel)

	// Write empty byte slice
	n, err := writer.Write([]byte{})
	assert.NoError(t, err)
	assert.Equal(t, 0, n)

	// Verify no panic, no log output
	assert.Empty(t, buf.String())

	// Write another empty slice
	n, err = writer.Write(nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, n)

	// Still no output
	assert.Empty(t, buf.String())
}

// TestLogWriter_MultipleNewlines tests lines with consecutive newlines
func TestLogWriter_MultipleNewlines(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf)
	logger.SetReportTimestamp(false)

	writer := newLogWriter(logger, log.InfoLevel)

	// Write with multiple consecutive newlines
	_, err := writer.Write([]byte("line1\n\nline2\n"))
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "line1")
	assert.Contains(t, output, "line2")

	// Count INFO entries - should have 3 (line1, empty line, line2)
	infoCount := strings.Count(output, "INFO")
	assert.Equal(t, 3, infoCount, "expected 3 INFO entries including empty line")

	// Verify buffer is empty
	assert.Equal(t, 0, len(writer.buf))
}

// TestLogWriter_ConcurrentWrites tests thread safety with concurrent writes
func TestLogWriter_ConcurrentWrites(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf)
	logger.SetReportTimestamp(false)

	writer := newLogWriter(logger, log.InfoLevel)

	const numGoroutines = 10
	const numWrites = 100
	var wg sync.WaitGroup

	// Launch multiple goroutines writing to same logWriter
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numWrites; j++ {
				line := []byte("goroutine message\n")
				_, err := writer.Write(line)
				assert.NoError(t, err)
			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Flush any remaining data
	writer.Flush()

	// Verify all data was written (should have numGoroutines * numWrites log entries)
	output := buf.String()
	infoCount := strings.Count(output, "INFO")
	assert.Equal(t, numGoroutines*numWrites, infoCount,
		"expected %d INFO entries", numGoroutines*numWrites)

	// Verify no data corruption (all messages should be "goroutine message")
	assert.Equal(t, numGoroutines*numWrites, strings.Count(output, "goroutine message"))

	// Verify buffer is empty after flush
	assert.Equal(t, 0, len(writer.buf))
}

// TestLogWriter_MixedPartialAndComplete tests mixed partial and complete writes
func TestLogWriter_MixedPartialAndComplete(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf)
	logger.SetReportTimestamp(false)

	writer := newLogWriter(logger, log.InfoLevel)

	// Write sequence: partial, complete, partial, complete
	_, err := writer.Write([]byte("part1 "))
	assert.NoError(t, err)
	assert.Empty(t, buf.String()) // Nothing logged yet

	_, err = writer.Write([]byte("part2\nline2\npart3 "))
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "part1 part2") // First complete line
	assert.Contains(t, output, "line2")       // Second complete line
	assert.NotContains(t, output, "part3")    // Still in buffer

	// Complete the last line
	_, err = writer.Write([]byte("part4\n"))
	assert.NoError(t, err)

	output = buf.String()
	assert.Contains(t, output, "part3 part4")

	// Verify buffer is empty
	assert.Equal(t, 0, len(writer.buf))
}

// TestLogWriter_LargeWrites tests writing large amounts of data
func TestLogWriter_LargeWrites(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf)
	logger.SetReportTimestamp(false)

	writer := newLogWriter(logger, log.InfoLevel)

	// Build a large string with multiple lines
	var largeData bytes.Buffer
	const numLines = 1000
	for i := 0; i < numLines; i++ {
		largeData.WriteString("line number ")
		// Use a simple counter instead of rune conversion
		// to avoid multi-byte character issues
		if i < 10 {
			largeData.WriteByte(byte('0' + i))
		} else {
			// Just write a simple marker for larger numbers
			largeData.WriteString("XXX")
		}
		largeData.WriteString("\n")
	}

	// Write all at once
	n, err := writer.Write(largeData.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, largeData.Len(), n)

	// Verify correct number of log entries
	output := buf.String()
	infoCount := strings.Count(output, "INFO")
	assert.Equal(t, numLines, infoCount)

	// Verify buffer is empty
	assert.Equal(t, 0, len(writer.buf))
}

// TestLogWriter_EmptyLineHandling tests how empty lines are handled
func TestLogWriter_EmptyLineHandling(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf)
	logger.SetReportTimestamp(false)

	writer := newLogWriter(logger, log.InfoLevel)

	// Write just a newline (empty line)
	_, err := writer.Write([]byte("\n"))
	assert.NoError(t, err)

	// Should still produce a log entry (for empty line)
	output := buf.String()
	assert.Contains(t, output, "INFO")

	// Verify buffer is empty
	assert.Equal(t, 0, len(writer.buf))
}

// TestLogWriter_DefaultLevel tests that unknown levels default to InfoLevel
func TestLogWriter_DefaultLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf)
	logger.SetReportTimestamp(false)

	// Use an arbitrary level value that's not defined
	writer := newLogWriter(logger, log.Level(999))

	// Write a line
	_, err := writer.Write([]byte("test\n"))
	assert.NoError(t, err)

	// Should default to INFO
	output := buf.String()
	assert.Contains(t, output, "INFO")
	assert.Contains(t, output, "test")
}

// TestLogWriter_FlushEmpty tests flushing when buffer is empty
func TestLogWriter_FlushEmpty(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf)
	logger.SetReportTimestamp(false)

	writer := newLogWriter(logger, log.InfoLevel)

	// Flush with empty buffer
	writer.Flush()

	// Should not produce any output
	assert.Empty(t, buf.String())
}

// TestLogWriter_ConsecutiveFlushes tests multiple consecutive flushes
func TestLogWriter_ConsecutiveFlushes(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf)
	logger.SetReportTimestamp(false)

	writer := newLogWriter(logger, log.InfoLevel)

	// Write partial data
	_, err := writer.Write([]byte("data"))
	assert.NoError(t, err)

	// Flush multiple times
	writer.Flush()
	writer.Flush()
	writer.Flush()

	// Should only log once
	output := buf.String()
	infoCount := strings.Count(output, "INFO")
	assert.Equal(t, 1, infoCount, "should only log once")
	assert.Contains(t, output, "data")
}
