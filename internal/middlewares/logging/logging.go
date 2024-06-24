package logging

import (
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Logging struct {
	RequestID     string
	RequestMethod string
	RequestURL    *url.URL
	RequestHeader http.Header
}

func NewLogging(r *http.Request) *Logging {
	return &Logging{
		RequestID:     uuid.NewString(),
		RequestURL:    r.URL,
		RequestMethod: r.Method,
		RequestHeader: r.Header,
	}
}

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trace := NewLogging(r)

		logrus.WithFields(logrus.Fields{
			"request_id":     trace.RequestID,
			"request_url":    trace.RequestURL.String(),
			"request_method": trace.RequestMethod,
			"request_header": trace.RequestHeader,
		}).Info("Request received")

		next.ServeHTTP(w, r)
	}
}
