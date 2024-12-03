package logger

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestSetLogLevel(t *testing.T) {
	tests := []struct {
		level          string
		expectedLevel  logrus.Level
		expectedLogMsg string
	}{
		{"debug", logrus.DebugLevel, ""},
		{"info", logrus.InfoLevel, ""},
		{"warning", logrus.WarnLevel, ""},
		{"error", logrus.ErrorLevel, ""},
		{"invalid", logrus.InfoLevel, "Log level incorrect. Set level info"},
	}

	for _, tt := range tests {

		Log.SetLevel(logrus.InfoLevel)

		SetLogLevel(tt.level)

		if Log.Level != tt.expectedLevel {
			t.Errorf("expected level %v, got %v", tt.expectedLevel, Log.Level)
		}
	}
}

func TestLoggingResponseWriter_Write(t *testing.T) {

	data := &responseData{}

	recorder := httptest.NewRecorder()

	loggingWriter := &loggingResponseWriter{
		ResponseWriter: recorder,
		responseData:   data,
	}

	requestBody := []byte("Hello, World!")

	size, err := loggingWriter.Write(requestBody)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if size != len(requestBody) {
		t.Errorf("expected size %d, got %d", len(requestBody), size)
	}

	if data.size != len(requestBody) {
		t.Errorf("expected response data size %d, got %d", len(requestBody), data.size)
	}

	if recorder.Body.String() != string(requestBody) {
		t.Errorf("expected body %q, got %q", string(requestBody), recorder.Body.String())
	}

	if data.status != 0 {
		t.Errorf("expected status %d, got %d", 0, data.status)
	}
}

func TestLoggingResponseWriter_WriteHeader(t *testing.T) {

	data := &responseData{}

	recorder := httptest.NewRecorder()

	loggingWriter := &loggingResponseWriter{
		ResponseWriter: recorder,
		responseData:   data,
	}

	statuses := []int{
		http.StatusOK,
		http.StatusBadRequest,
		http.StatusNotFound,
	}

	for _, status := range statuses {

		loggingWriter.WriteHeader(status)

		if data.status != status {
			t.Errorf("expected response data status %d, got %d", status, data.status)
		}
	}
}

func TestRequestLogger(t *testing.T) {

	var logBuffer logrus.Logger
	logBuffer.SetOutput(httptest.NewRecorder())
	Log = &logBuffer

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	loggedHandler := RequestLogger(testHandler)

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	recorder := httptest.NewRecorder()

	loggedHandler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, recorder.Code)
	}

	if recorder.Body.String() != "Hello, World!" {
		t.Errorf("expected body %q, got %q", "Hello, World!", recorder.Body.String())
	}

}
