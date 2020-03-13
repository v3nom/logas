package logas

import (
	"context"
	"net/http"
	"strings"

	"github.com/v3nom/pipes"
)

// Middleware middleware to capture request context for Logas
func Middleware(ctx context.Context, w http.ResponseWriter, r *http.Request, next pipes.Next) {
	span, trace := getTraceContext(r)
	ctx = context.WithValue(ctx, SpanKey, span)
	ctx = context.WithValue(ctx, TraceKey, trace)
	next(ctx)
}

func getTraceContext(r *http.Request) (span, trace string) {
	header := r.Header.Get("X-Cloud-Trace-Context")
	if header == "" {
		return "", ""
	}
	parts := strings.Split(header, ";")[0]
	return strings.Split(parts, "/")[1], strings.Split(parts, "/")[0]
}
