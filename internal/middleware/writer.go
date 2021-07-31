package middleware

import "net/http"

// writer used to capture response status code
type writer struct {
	http.ResponseWriter
	statusCode int
}

func (l *writer) WriteHeader(code int) {
	l.statusCode = code
	l.ResponseWriter.WriteHeader(code)
}
