package log

import (
	"context"

	"go.uber.org/zap"
)

// logKey is the key that stores a logger.
type logKey struct{}

// WithLogger stores the logger in the context.
func WithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, logKey{}, logger.Desugar().WithOptions(zap.AddCallerSkip(1)).Sugar())
}

// WithFields stores a logger that has fields included in every log.
func WithFields(ctx context.Context, fields ...interface{}) context.Context {
	logger := getLogger(ctx)
	if logger != nil {
		return context.WithValue(ctx, logKey{}, logger.With(fields...))
	}

	return ctx
}

// Debug will write a debug log if a logger is found within the context.
func Debug(ctx context.Context, msg string, keysAndValues ...interface{}) {
	logger := getLogger(ctx)
	if logger != nil {
		logger.Debugw(msg, keysAndValues...)
	}
}

// Info will write an info log if a logger is found within the context.
func Info(ctx context.Context, msg string, keysAndValues ...interface{}) {
	logger := getLogger(ctx)
	if logger != nil {
		logger.Infow(msg, keysAndValues...)
	}
}

// Warn will write a warning log if a logger is found within the context.
func Warn(ctx context.Context, msg string, keysAndValues ...interface{}) {
	logger := getLogger(ctx)
	if logger != nil {
		logger.Warnw(msg, keysAndValues...)
	}
}

// Error will write an error log if a logger is found within the context.
func Error(ctx context.Context, msg string, keysAndValues ...interface{}) {
	logger := getLogger(ctx)
	if logger != nil {
		logger.Errorw(msg, keysAndValues...)
	}
}

// getLogger will retrieve a logger from the context if found, otherwise nil.
func getLogger(ctx context.Context) *zap.SugaredLogger {
	rawLogger := ctx.Value(logKey{})
	if rawLogger == nil {
		return nil
	}

	return rawLogger.(*zap.SugaredLogger)
}
