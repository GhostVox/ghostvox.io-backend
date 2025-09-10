package utils_test

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/GhostVox/ghostvox.io-backend/internal/utils"
)

// mockWriter implements io.Writer and can be configured to fail
type mockWriter struct {
	buffer     bytes.Buffer
	shouldFail bool
}

func (m *mockWriter) Write(p []byte) (n int, err error) {
	if m.shouldFail {
		return 0, errors.New("mock write failure")
	}
	return m.buffer.Write(p)
}

func (m *mockWriter) String() string {
	return m.buffer.String()
}

func TestNewLogger(t *testing.T) {
	writer := &bytes.Buffer{}
	logger := utils.NewLogger(writer)

	if logger == nil {
		t.Fatal("NewLogger returned nil")
	}
}

func TestLogJob(t *testing.T) {
	tests := []struct {
		name             string
		jobName          string
		message          string
		writerShouldFail bool
		wantError        bool
		wantContains     []string
	}{
		{
			name:             "Basic logging",
			jobName:          "test-job",
			message:          "test message",
			writerShouldFail: false,
			wantError:        false,
			wantContains:     []string{"[test-job]", "test message"},
		},
		{
			name:             "Failed writer",
			jobName:          "failed-job",
			message:          "this should fail",
			writerShouldFail: true,
			wantError:        true,
			wantContains:     []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockW := &mockWriter{shouldFail: tt.writerShouldFail}
			logger := utils.NewLogger(mockW)

			n, err := logger.LogJob(tt.jobName, tt.message)

			// Check error expectation
			if (err != nil) != tt.wantError {
				t.Errorf("LogJob() error = %v, wantError %v", err, tt.wantError)
				return
			}

			// If error was expected, no need to check content
			if tt.wantError {
				if n != 0 {
					t.Errorf("LogJob() returned n = %d, want 0 for error case", n)
				}
				return
			}

			// Check content contains expected strings
			output := mockW.String()
			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("LogJob() output = %q, should contain %q", output, want)
				}
			}

			// Check that timestamp format is reasonable
			if !strings.Contains(output, time.Now().Format("2006-01-02")) {
				t.Errorf("LogJob() output = %q, should contain today's date", output)
			}
		})
	}
}

func TestLogError(t *testing.T) {
	tests := []struct {
		name             string
		err              error
		writerShouldFail bool
		wantError        bool
		wantContains     []string
	}{
		{
			name:             "Basic error logging",
			err:              errors.New("test error"),
			writerShouldFail: false,
			wantError:        false,
			wantContains:     []string{"test error"},
		},
		{
			name:             "Nil error",
			err:              nil,
			writerShouldFail: false,
			wantError:        false,
			wantContains:     []string{},
		},
		{
			name:             "Failed writer",
			err:              errors.New("writer will fail"),
			writerShouldFail: true,
			wantError:        true,
			wantContains:     []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockW := &mockWriter{shouldFail: tt.writerShouldFail}
			logger := utils.NewLogger(mockW)

			// Skip test if error is nil because the current implementation might not handle this correctly
			if tt.err == nil {
				t.Skip("Skipping nil error test - current implementation doesn't handle nil errors")
				return
			}

			n, err := logger.LogError(tt.err)

			// Check error expectation
			if (err != nil) != tt.wantError {
				t.Errorf("LogError() error = %v, wantError %v", err, tt.wantError)
				return
			}

			// If error was expected, no need to check content
			if tt.wantError {
				if n != 0 {
					t.Errorf("LogError() returned n = %d, want 0 for error case", n)
				}
				return
			}

			// Check content contains expected strings
			output := mockW.String()
			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("LogError() output = %q, should contain %q", output, want)
				}
			}

			// Check that timestamp format is reasonable
			if !strings.Contains(output, time.Now().Format("2006-01-02")) {
				t.Errorf("LogError() output = %q, should contain today's date", output)
			}
		})
	}
}

// TestConcurrency tests that the logger can handle concurrent writes
func TestConcurrency(t *testing.T) {
	writer := &bytes.Buffer{}
	logger := utils.NewLogger(writer)

	const numGoroutines = 100
	done := make(chan bool, numGoroutines*2)

	// Launch multiple goroutines writing job logs
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			_, err := logger.LogJob("job", fmt.Sprintf("message-%d", id))
			if err != nil {
				t.Errorf("Concurrent LogJob failed: %v", err)
			}
			done <- true
		}(i)
	}

	// Launch multiple goroutines writing error logs
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			_, err := logger.LogError(fmt.Errorf("error-%d", id))
			if err != nil {
				t.Errorf("Concurrent LogError failed: %v", err)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines*2; i++ {
		<-done
	}

	// Check that the output contains the expected number of log entries
	output := writer.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	expectedLines := numGoroutines * 2 // Both job and error logs

	if len(lines) != expectedLines {
		t.Errorf("Expected %d log lines, got %d", expectedLines, len(lines))
	}
}

// Test that nil errors are handled properly
func TestLogErrorWithNilError(t *testing.T) {
	writer := &bytes.Buffer{}
	logger := utils.NewLogger(writer)

	n, err := logger.LogError(nil)

	// Current implementation should return early without writing
	// and without error if nil is passed
	if n != 0 {
		t.Errorf("LogError(nil) should return 0 bytes written, got %d", n)
	}

	if err != nil {
		t.Errorf("LogError(nil) should not return an error, got %v", err)
	}

	if writer.Len() > 0 {
		t.Errorf("LogError(nil) should not write anything, got %q", writer.String())
	}
}
