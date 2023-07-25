package middlewares

import (
	"context"
	"github.com/poggerr/gophermart/internal/logger"
	"net/http"
	"time"
)

func WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		responseData := &logger.ResponseData{}
		lw := logger.LoggingResponseWriter{
			ResponseWriter: w,
			ResponseData:   responseData,
		}

		ur := r.WithContext(context.WithValue(r.Context(), "user", "sdcsc"))

		h.ServeHTTP(&lw, ur)

		duration := time.Since(time.Now())

		logger.Log.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.Status,
			"duration", duration,
			"size", responseData.Size,
		)
	}

	return http.HandlerFunc(logFn)
}
