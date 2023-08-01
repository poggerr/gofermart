package logger

import (
	"net/http"
	"time"
)

func WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		responseData := &ResponseData{}
		lw := LoggingResponseWriter{
			ResponseWriter: w,
			ResponseData:   responseData,
		}

		//ur := r.WithContext(context.WithValue(r.Context(), "user", "sdcsc"))

		h.ServeHTTP(&lw, r)

		duration := time.Since(time.Now())

		Log.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.Status,
			"duration", duration,
			"size", responseData.Size,
		)
	}

	return http.HandlerFunc(logFn)
}
