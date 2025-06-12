package log

import (
	"context"

	customLogger "g-management/pkg/log"

	"github.com/aws/smithy-go/logging"
)

type AwsLogger struct {
	Logger *customLogger.Logger
}

func (l AwsLogger) WithContext(ctx context.Context) logging.Logger {
	return &AwsLogger{
		Logger: l.Logger.WithContext(ctx),
	}
}

func (l AwsLogger) Logf(classification logging.Classification, format string, v ...interface{}) {
	if classification == logging.Warn {
		l.Logger.Warn(context.Background(), format, v...)
	}
	if classification == logging.Debug {
		l.Logger.Debug(context.Background(), format, v...)
	}
}
