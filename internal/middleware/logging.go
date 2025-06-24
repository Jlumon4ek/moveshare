package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	rw.body.Write(data)
	return rw.ResponseWriter.Write(data)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Читаем тело запроса
		var requestBody []byte
		if r.Body != nil {
			requestBody, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Оборачиваем ResponseWriter
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     200,
			body:           &bytes.Buffer{},
		}

		// Выполняем запрос
		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		// Логируем запрос
		logLevel := slog.LevelInfo
		if rw.statusCode >= 400 {
			logLevel = slog.LevelError
		}

		attrs := []slog.Attr{
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rw.statusCode),
			slog.Duration("duration", duration),
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
		}

		// Добавляем тело запроса если не пустое
		if len(requestBody) > 0 && isJSONContent(r.Header.Get("Content-Type")) {
			var jsonBody interface{}
			if json.Unmarshal(requestBody, &jsonBody) == nil {
				attrs = append(attrs, slog.Any("request_body", jsonBody))
			}
		}

		// Добавляем тело ответа если есть ошибка
		if rw.statusCode >= 400 && rw.body.Len() > 0 {
			attrs = append(attrs, slog.String("response_body", rw.body.String()))
		}

		slog.LogAttrs(r.Context(), logLevel, getStatusMessage(rw.statusCode), attrs...)
	})
}

func isJSONContent(contentType string) bool {
	return contentType == "application/json" ||
		contentType == "application/json; charset=utf-8"
}

func getStatusMessage(status int) string {
	switch {
	case status >= 200 && status < 300:
		return "✅ Request completed successfully"
	case status >= 400 && status < 500:
		return "❌ Client error"
	case status >= 500:
		return "💥 Server error"
	default:
		return "📝 Request processed"
	}
}
