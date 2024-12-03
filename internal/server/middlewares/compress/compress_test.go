package compress

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewCompressWriter(t *testing.T) {

	rr := httptest.NewRecorder()

	cw := newCompressWriter(rr)

	if cw == nil {
		t.Fatal("newCompressWriter() возвращает nil")
	}

	if cw.w != rr {
		t.Errorf("newCompressWriter() = %v, want %v", cw.w, rr)
	}

	if cw.zw == nil {
		t.Fatal("newCompressWriter().zw возвращает nil")
	}

	if err := cw.zw.Close(); err != nil {
		t.Errorf("не удалось закрыть gzip.Writer: %v", err)
	}
}

func TestCompressWriter_Header(t *testing.T) {

	rr := httptest.NewRecorder()
	cw := &compressWriter{w: rr}

	expectedKey := "Content-Encoding"
	expectedValue := "gzip"
	rr.Header().Set(expectedKey, expectedValue)

	header := cw.Header()

	if header.Get(expectedKey) != expectedValue {
		t.Errorf("Header() = %v, want %v", header.Get(expectedKey), expectedValue)
	}
}

func TestCompressWriter_Write(t *testing.T) {

	rr := httptest.NewRecorder()
	cw := newCompressWriter(rr)

	inputData := []byte("Hello, Gzip!")
	expectedLength := len(inputData)

	n, err := cw.Write(inputData)

	if err != nil {
		t.Fatalf("Write() вернул ошибку: %v", err)
	}

	if n != expectedLength {
		t.Errorf("Write() вернул %d байтов, ожидалось %d", n, expectedLength)
	}

	if err := cw.zw.Close(); err != nil {
		t.Errorf("не удалось закрыть gzip.Writer: %v", err)
	}

	result := rr.Body.Bytes()

	if len(result) == 0 {
		t.Error("Тело ответа пустое, ожидались сжатые данные")
	}
}

func TestCompressWriter_WriteHeader(t *testing.T) {

	rr := httptest.NewRecorder()
	cw := newCompressWriter(rr)

	statusCode := http.StatusOK

	cw.WriteHeader(statusCode)

	if got := rr.Code; got != statusCode {
		t.Errorf("WriteHeader() = %d, ожидалось %d", got, statusCode)
	}

	cw.WriteHeader(statusCode)

	if got := rr.Code; got != statusCode {
		t.Errorf("WriteHeader() = %d, ожидалось %d", got, statusCode)
	}
}

func TestCompressWriter_Close(t *testing.T) {

	rr := httptest.NewRecorder()
	cw := newCompressWriter(rr)

	cw.WriteHeader(http.StatusOK)

	if err := cw.Close(); err != nil {
		t.Fatalf("Close() вернул ошибку: %v", err)
	}

	result := rr.Body.Bytes()

	if len(result) == 0 {
		t.Error("Тело ответа пустое, ожидались сжатые данные")
	}

	if got := rr.Header().Get("Content-Encoding"); got != "gzip" {
		t.Errorf("Заголовок Content-Encoding = %v, ожидался gzip", got)
	}
}

func TestNewCompressReader(t *testing.T) {

	var b bytes.Buffer
	writer := gzip.NewWriter(&b)
	if _, err := writer.Write([]byte("Hello, Gzip!")); err != nil {
		t.Fatalf("не удалось записать данные в gzip: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("не удалось закрыть gzip.Writer: %v", err)
	}

	reader := io.NopCloser(bytes.NewReader(b.Bytes()))

	cr, err := newCompressReader(reader)
	if err != nil {
		t.Fatalf("newCompressReader вернул ошибку: %v", err)
	}

	if cr == nil {
		t.Fatal("newCompressReader вернул nil compressReader")
	}

	if cr.zr == nil {
		t.Fatal("gzip.Reader не был инициализирован")
	}

	decompressedData, err := io.ReadAll(cr.zr)
	if err != nil {
		t.Fatalf("не удалось прочитать decompressed данные: %v", err)
	}

	expected := "Hello, Gzip!"
	if string(decompressedData) != expected {
		t.Errorf("ожидалось: %s, получено: %s", expected, string(decompressedData))
	}
}

func TestCompressRead(t *testing.T) {
	gzipData := []byte{0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0x06, 0x00, 0x42, 0x43, 0x02, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	gzipReader, err := gzip.NewReader(bytes.NewReader(gzipData))
	if err != nil {
		t.Fatal(err)
	}
	cr := compressReader{r: gzipReader, zr: gzipReader}
	gzipReader.Close()
	buf := make([]byte, 10)
	_, err = cr.Read(buf)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestCompressReader_Close(t *testing.T) {

	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	_, err := writer.Write([]byte("Hello, Gzip!"))
	if err != nil {
		t.Fatalf("failed to write to gzip: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close gzip.Writer: %v", err)
	}

	reader := io.NopCloser(bytes.NewReader(buf.Bytes()))

	cr, err := newCompressReader(reader)
	if err != nil {
		t.Fatalf("newCompressReader returned an error: %v", err)
	}

	if err := cr.Close(); err != nil {
		t.Fatalf("Close() should not return an error: %v", err)
	}
}

func TestGzipMiddleware(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	middleware := GzipMiddleware(next)
	middleware.ServeHTTP(w, req)
	if w.Header().Get("Content-Encoding") == "gzip" {
		t.Errorf("Expected no Content-Encoding, got %s", w.Header().Get("Content-Encoding"))
	}

	req, err = http.NewRequest("POST", "/", strings.NewReader("Hello, World!"))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Encoding", "gzip")
	w = httptest.NewRecorder()
	next = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
	middleware = GzipMiddleware(next)
	middleware.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code 500, got %d", w.Code)
	}
}
