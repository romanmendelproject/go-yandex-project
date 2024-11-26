package logger

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func SetLogLevel(level string) {
	switch level {
	case "debug":
		Log.SetLevel(logrus.DebugLevel)
	case "info":
		Log.SetLevel(logrus.InfoLevel)
	case "warning":
		Log.SetLevel(logrus.WarnLevel)
	case "error":
		Log.SetLevel(logrus.ErrorLevel)
	default:
		Log.Warning("Log level incorrect. Set level info")
		Log.SetLevel(logrus.InfoLevel)
	}
}

type ResponseWriter interface {
	Header() http.Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func RequestLogger(h http.Handler) http.Handler {
	logFn := func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()
		uri := req.RequestURI

		method := req.Method

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: res,
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, req)

		duration := time.Since(start)

		Log.WithFields(logrus.Fields{
			"method":   method,
			"path":     uri,
			"duration": duration,
			"status":   responseData.status,
			"size":     responseData.size,
		}).Info("got incoming HTTP request")
	}

	return http.HandlerFunc(logFn)
}
