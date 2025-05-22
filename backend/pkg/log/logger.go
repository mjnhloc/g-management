package log

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"time"
)

// Logger comment
// en: Logger type struct is a wrapper around slog.Logger
// en: A Logger records structured information about each call to its
// en: Debug, Info, Warn, and Error methods.
// en: For each call, it creates a [Record] and passes it to a [Handler].
type Logger struct {
	logger *slog.Logger
}

var keys []string

var _ slog.Handler = &defaultHandler{}

type defaultHandler struct {
	slog.Handler

	EnableDecoratorNR bool
}

func (t *defaultHandler) Handle(ctx context.Context, r slog.Record) error {
	if v, ok := ctx.Value(logMapCtxKey).(*sync.Map); ok {
		v.Range(func(key, value any) bool {
			if key, ok := key.(string); ok {
				r.AddAttrs(slog.Any(key, value))
			}
			return true
		})
	}

	for _, key := range keys {
		if ctx.Value(key) != nil {
			r.AddAttrs(slog.Any(key, ctx.Value(key)))
		}
	}
	return t.Handler.Handle(ctx, r)
}

func (t *defaultHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &defaultHandler{
		Handler:           t.Handler.WithAttrs(attrs),
		EnableDecoratorNR: t.EnableDecoratorNR,
	}
}

type loggerCtxKey struct{}

var logMapCtxKey = loggerCtxKey{}

// CtxWithValue function
// en: adds a key-value pair to the context in sync.Map for thread safely
// en: this value automatically added to the log record with defaultHandler
func CtxWithValue(ctx context.Context, key string, val any) context.Context {
	m, ok := ctx.Value(logMapCtxKey).(*sync.Map)
	if !ok {
		m = &sync.Map{}
		m.Store(key, val)
		return context.WithValue(ctx, logMapCtxKey, m)
	}
	m.Store(key, val)
	return context.WithValue(ctx, logMapCtxKey, m)
}

const messageKey = "message"

// Initialize function
// en: initializes the logger with defaultHandler
// en: if isDebug is true, it sets the log level to [LevelDebug] otherwise [LevelInfo]
// en: it uses JSONHandler for logging and automatically
// en: adds the key-value pair to the log record using the CtxWithValue function
// en: The fields in the keys are used to retrieve their values from the context and write them to the logger
//
// The ConfigOption arguments allow for configuration of the Logger. They
// are applied in order from first to last, i.e. latter ConfigOptions may
// overwrite the Config fields already set.
func Initialize(ctx context.Context, opts ...ConfigOption) context.Context {
	c := defaultConfig()

	for _, fn := range opts {
		if fn != nil {
			fn(&c)
		}
	}

	resolveAdditionalKeys(&c)

	resolveHandler(&c)

	slog.SetDefault(slog.New(c.Handler))

	return context.WithValue(ctx, logMapCtxKey, &sync.Map{})
}

// Info comment
// en: Info logs at [LevelInfo] with the given context.
// en: msg argument should not be empty because of convention in NewRelic.
// en: args arguments should be passed in pairs.
func (l *Logger) Info(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, args...)
}

// Debug comment
// en: Debug logs at [LevelDebug] with the given context.
// en: msg argument should not be empty because of convention in NewRelic.
// en: args arguments should be passed in pairs.
func (l *Logger) Debug(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, args...)
}

// Warn comment
// en: Warn logs at [LevelWarn] with the given context.
// en: msg argument should not be empty because of convention in NewRelic.
// en: args arguments should be passed in pairs.
func (l *Logger) Warn(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, args...)
}

// Error comment
// en: Error logs at [LevelError] with the given context.
// en: msg argument should not be empty because of convention in NewRelic.
// en: args arguments should be passed in pairs.
func (l *Logger) Error(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, args...)
}

// With comment
// en: With returns a Logger that includes the given attributes
// en: in each output operation.
// en: Attributes should be passed in pairs.
func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		logger: l.logger.With(args...),
	}
}

// WithContext comment
// en: WithContext handles getting fields in context to write to logger
func (l *Logger) WithContext(ctx context.Context) *Logger {
	newLogger := l.With()
	if v, ok := ctx.Value(logMapCtxKey).(*sync.Map); ok {
		v.Range(func(key, value any) bool {
			if key, ok := key.(string); ok {
				newLogger = newLogger.With(key, ctx.Value(key))
			}
			return true
		})
	}
	for _, key := range keys {
		if ctx.Value(key) != nil {
			newLogger = newLogger.With(key, ctx.Value(key))
		}
	}
	return newLogger
}

// Debug comment
// en: Debug calls [Logger.DebugContext] on the default logger.
// en: msg argument should not be empty because of convention in NewRelic.
// en: args arguments should be passed in pairs.
func Debug(ctx context.Context, msg string, args ...any) {
	slog.DebugContext(ctx, msg, args...)
}

// Info comment
// en: Info calls [Logger.InfoContext] on the default logger.
// en: msg argument should not be empty because of convention in NewRelic.
// en: args arguments should be passed in pairs.
func Info(ctx context.Context, msg string, args ...any) {
	slog.InfoContext(ctx, msg, args...)
}

// Warn comment
// en: Warn calls [Logger.DebugContext] on the default logger.
// en: msg argument should not be empty because of convention in NewRelic.
// en: args arguments should be passed in pairs.
func Warn(ctx context.Context, msg string, args ...any) {
	slog.WarnContext(ctx, msg, args...)
}

// Error comment
// en: Error calls [Logger.ErrorContext] on the default logger.
// en: msg argument should not be empty because of convention in NewRelic.
// en: args arguments should be passed in pairs.
func Error(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, args...)
}

// With comment
// en: With calls [Logger.With] on the default logger.
// en: Attributes should be passed in pairs.
func With(args ...any) *Logger {
	return &Logger{
		logger: slog.With(args...),
	}
}

// Group comment
// en: Group returns an Attr for a Group [Value].
// en: The first argument is the key; the remaining arguments
// en: are converted to Attrs as in [Logger.Log].
//
// en: Use Group to collect several key-value pairs under a single
// en: key on a log line, or as the result of LogValue
// en: in order to log a single value as multiple Attrs.
func Group(key string, args ...any) slog.Attr {
	return slog.Group(key, args...)
}

// Fatal comment
// en: Fatal calls [Logger.ErrorContext] on the default logger and exit program.
// en: msg argument should not be empty because of convention in NewRelic.
// en: args arguments should be passed in pairs.
func Fatal(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, args...)
	os.Exit(1)
}

func resolveHandler(c *Config) {
	if c.Handler == nil {
		c.Handler = &defaultHandler{Handler: slog.NewJSONHandler(c.Output, &slog.HandlerOptions{Level: c.Level, ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if v, ok := a.Value.Any().(time.Duration); ok {
				a.Value = slog.StringValue(v.String())
			}
			if a.Key != slog.MessageKey {
				return a
			}
			a.Key = messageKey
			return a
		}}), EnableDecoratorNR: c.EnableDecoratorNR}
	}

	if len(c.HandlerWrappers) == 0 {
		return
	}

	for _, fn := range c.HandlerWrappers {
		c.Handler = fn(c.Handler)
	}
}

func resolveAdditionalKeys(c *Config) {
	keys = append(keys, c.KeysInput...)
}
