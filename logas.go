package logas

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/logging"
)

// ContextKey logas context key
type ContextKey string

// LoggerKey context key for logger
var LoggerKey ContextKey = "logas logger ctx"

// SpanKey context key for span ID
var SpanKey ContextKey = "logas span ctx"

// TraceKey context key for trace ID
var TraceKey ContextKey = "logas trace ctx"

// Debugf logs a message at DEBUG level
func Debugf(ctx context.Context, format string, args ...interface{}) {
	if ctx.Value(LoggerKey) == nil {
		log.Printf("Debug. "+format, args...)
		return
	}

	logger := ctx.Value(LoggerKey).(*logging.Logger)
	entry := createCommonEntry(ctx, logging.Debug, format, args...)
	logger.Log(*entry)
}

// Infof logs a message at INFO level
func Infof(ctx context.Context, format string, args ...interface{}) {
	if ctx.Value(LoggerKey) == nil {
		log.Printf("Info. "+format, args...)
		return
	}

	logger := ctx.Value(LoggerKey).(*logging.Logger)
	entry := createCommonEntry(ctx, logging.Info, format, args...)
	logger.Log(*entry)
}

// Warningf logs a message at WARNING level
func Warningf(ctx context.Context, format string, args ...interface{}) {
	if ctx.Value(LoggerKey) == nil {
		log.Printf("Warning. "+format, args...)
		return
	}

	logger := ctx.Value(LoggerKey).(*logging.Logger)
	entry := createCommonEntry(ctx, logging.Warning, format, args...)
	logger.Log(*entry)
}

// Errorf logs a message at ERROR level
func Errorf(ctx context.Context, format string, args ...interface{}) {
	if ctx.Value(LoggerKey) == nil {
		log.Printf("Error. "+format, args...)
		return
	}

	logger := ctx.Value(LoggerKey).(*logging.Logger)
	entry := createCommonEntry(ctx, logging.Error, format, args...)
	logger.Log(*entry)
}

// Criticalf logs a message at CRITICAL level
func Criticalf(ctx context.Context, format string, args ...interface{}) {
	if ctx.Value(LoggerKey) == nil {
		log.Printf("Critical. "+format, args...)
		return
	}

	logger := ctx.Value(LoggerKey).(*logging.Logger)
	entry := createCommonEntry(ctx, logging.Critical, format, args...)
	logger.Log(*entry)
}

func createCommonEntry(ctx context.Context, severity logging.Severity, format string, args ...interface{}) *logging.Entry {
	entry := &logging.Entry{
		Payload:   fmt.Sprintf(format, args...),
		Severity:  severity,
		Timestamp: time.Now(),
	}
	addContext(ctx, entry)
	return entry
}

func addContext(ctx context.Context, entry *logging.Entry) {
	if ctx.Value(TraceKey) != nil {
		entry.Trace = fmt.Sprintf("projects/%s/traces/%s", getProjectID(), ctx.Value(TraceKey).(string))
	}
	if ctx.Value(SpanKey) != nil {
		entry.SpanID = ctx.Value(SpanKey).(string)
	}
}

var projectID string

func getProjectID() string {
	if projectID != "" {
		return projectID
	}
	if metadata.OnGCE() {
		pid, err := metadata.ProjectID()
		if err == nil {
			projectID = pid
		}
	}
	return projectID
}
