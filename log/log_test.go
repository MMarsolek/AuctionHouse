package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func TestDebugOutputsDebugLevel(t *testing.T) {
	ctx := WithLogger(context.Background(), getTestLogger(t, zap.DebugLevel))
	Debug(ctx, "my message")
}

func TestInfoOutputsInfoLevel(t *testing.T) {
	ctx := WithLogger(context.Background(), getTestLogger(t, zap.InfoLevel))
	Info(ctx, "my message")
}

func TestWarnOutputsWarnLevel(t *testing.T) {
	ctx := WithLogger(context.Background(), getTestLogger(t, zap.WarnLevel))
	Warn(ctx, "my message")
}

func TestErrorOutputsErrorLevel(t *testing.T) {
	ctx := WithLogger(context.Background(), getTestLogger(t, zap.ErrorLevel))
	Error(ctx, "my message")
}

func TestWithFieldsAddsFieldsToSubsequentCalls(t *testing.T) {
	ctx := WithLogger(context.Background(), getTestLogger(t, zap.InfoLevel))
	ctx = WithFields(ctx, "hello", "world")
	Info(ctx, "my message")
}

func TestGetLoggerReturnsNilOnContextWithNoLogger(t *testing.T) {
	require.Nil(t, getLogger(context.Background()))
}

func getTestLogger(t *testing.T, expectedLevel zapcore.Level) *zap.SugaredLogger {
	return zaptest.NewLogger(t, zaptest.WrapOptions(zap.Hooks(func(e zapcore.Entry) error {
		require.EqualValues(t, expectedLevel, e.Level)
		return nil
	}))).Sugar()
}
