package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type ctxKey int

const requestKey ctxKey = ctxKey(42)

// Logger handler to log the start and end of a request
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		writer := &writer{ResponseWriter: w}

		start := time.Now()
		route := mux.CurrentRoute(r)

		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = uuid.New().String()
			r.Header.Set("X-Request-ID", id)
		}
		ctx = context.WithValue(ctx, requestKey, id)

		remoteAddr, _ := getIP(r)
		referer := r.Referer()
		userAgent := r.UserAgent()
		method := r.Method
		pathTemplate, _ := route.GetPathTemplate()
		vars := mux.Vars(r)
		query := r.URL.Query()

		e := log.Info() // Log before handling the call
		e.Str("X-Request-ID", id)
		e.Str("user", GetUser(r).SAMAccountName)
		e.Str("RemoteAddress", remoteAddr).Str("Referer", referer).Str("User-Agent", userAgent)
		e.Str("method", method).Str("path", pathTemplate).Interface("vars", vars).Interface("query", query)
		e.Msg("START")

		defer func(start time.Time) {
			e := log.Info() // Log after handling the call
			e.Str("X-Request-ID", id)
			e.Str("user", GetUser(r).SAMAccountName)
			e.Str("RemoteAddress", remoteAddr).Str("Referer", referer).Str("User-Agent", userAgent)
			e.Str("method", method).Str("path", pathTemplate).Interface("vars", vars).Interface("query", query)
			e.Int("code", writer.statusCode).AnErr("context", ctx.Err())
			e.Dur("duration", time.Since(start)).Int64("resp_time", time.Since(start).Milliseconds())
			e.Msg("END")
		}(start)
		next.ServeHTTP(writer, r.WithContext(ctx))
	})
}

// see: https://golangbyexample.com/golang-ip-address-http-request/
func getIP(r *http.Request) (string, error) {
	// Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	// Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	// Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", fmt.Errorf("No valid ip found")
}

func GetRequestID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestKey).(string)
	return id, ok
}
