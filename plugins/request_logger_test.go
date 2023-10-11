package plugins

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestRequestLogger_OnRequestEnd(t *testing.T) {
	// Prepare the request and response objects
	req := httptest.NewRequest("GET", "/test", nil)
	res := &http.Response{
		StatusCode: http.StatusOK,
	}

	// Create a buffer to capture the logger output
	var buf bytes.Buffer
	logger := &requestLogger{
		out: &buf,
	}

	// Call the OnRequestEnd method
	logger.OnRequestStart(req)
	logger.OnRequestEnd(req, res)
	timeCurrent := time.Now().Format("02/Jan/2006 03:04:05")

	// Verify the output
	expected := timeCurrent + " GET /test 200 [0ms]\n"
	if buf.String() != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, buf.String())
	}
}

func TestRequestLogger_OnError(t *testing.T) {
	// Prepare the request object
	req := httptest.NewRequest("GET", "/test", nil)

	// Create a buffer to capture the logger error output
	var buf bytes.Buffer
	logger := &requestLogger{
		errOut: &buf,
	}

	// Simulate an error and call the OnError method
	err := http.ErrServerClosed
	logger.OnRequestStart(req)
	logger.OnError(req, err)

	// Format the time according to the desired layout
	timeCurrent := time.Now().Format("02/Jan/2006 03:04:05")

	// Verify the output
	expected := timeCurrent + " GET /test [0ms] ERROR: " + err.Error() + "\n"
	if buf.String() != expected {
		t.Errorf("Expected: %s\nGot: %s", expected, buf.String())
	}
}

func TestGetRequestDuration(t *testing.T) {
	// Simulate a request duration of 500ms
	now := time.Now()
	startTime := now.Add(-500 * time.Millisecond)
	ctx := context.WithValue(context.Background(), reqTime, startTime)
	duration := getRequestDuration(ctx)
	expectedDuration := 500 * time.Millisecond
	tolerance := 1 * time.Millisecond

	if duration < expectedDuration-tolerance || duration > expectedDuration+tolerance {
		t.Errorf("Expected: %v\nGot: %v", expectedDuration, duration)
	}
}

func TestGetRequestDurationStartZero(t *testing.T) {
	now := time.Now()
	startTime := now.Add(-500 * time.Millisecond)
	ctx := context.WithValue(context.Background(), "reqTimeTemp", startTime)
	duration := getRequestDuration(ctx)
	expectedDuration := 0 * time.Millisecond
	tolerance := 1 * time.Millisecond

	if duration < expectedDuration-tolerance || duration > expectedDuration+tolerance {
		t.Errorf("Expected: %v\nGot: %v", expectedDuration, duration)
	}
}

func TestGetRequestDurationZero(t *testing.T) {
	type CustomContext struct {
		context.Context
		value interface{}
	}
	type CustomValue struct {
		Value time.Time
	}

	customValue := &CustomValue{
		Value: time.Now(),
	}
	// Create a CustomContext with the customValue
	ctx := &CustomContext{
		Context: context.Background(),
		value:   customValue,
	}
	ctxU := context.WithValue(ctx, reqTime, "u")
	duration := getRequestDuration(ctxU)

	// Assert that the returned duration is 0 since start.(time.Time) is not ok
	if duration != 0 {
		t.Errorf("Expected duration to be 0, but got %v", duration)
	}
}

func TestNewRequestNoInput(t *testing.T) {
	logger := NewRequestLogger(nil, nil)
	if logger.(*requestLogger).out != os.Stdout {
		t.Errorf("Expected os.Stdout as default output, but got %v", logger.(*requestLogger).out)
	}
	if logger.(*requestLogger).errOut != os.Stderr {
		t.Errorf("Expected os.Stderr as default error output, but got %v", logger.(*requestLogger).errOut)
	}

}
