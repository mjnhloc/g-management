package log

import (
	"io"
	"log/slog"
)

// ConfigOption comment.
// en: configures the Config when provided to NewApplication.
type ConfigOption func(*Config)

// ConfigLevel comment.
// en: set the level of log.
func ConfigLevel(level Level) ConfigOption {
	return func(cfg *Config) { cfg.Level = slog.Level(level) }
}

// ConfigOutput comment.
// en: set how logger outputs log.
func ConfigOutput(w io.Writer) ConfigOption {
	return func(cfg *Config) { cfg.Output = w }
}

// ConfigAdditionalKeys comment.
// en: set additional keys will be log from context.
func ConfigAdditionalKeys(k []string) ConfigOption {
	return func(cfg *Config) { cfg.KeysInput = k }
}

// ConfigEnableDecoratorNR comment.
// en: decorator for new relic attribute.
// en: new relic groups all metadata to attribute "newrelic".
// en: split the attribute to get metadata and create attr our own.
func ConfigEnableDecoratorNR(e bool) ConfigOption {
	return func(cfg *Config) { cfg.EnableDecoratorNR = e }
}

// ConfigHandleWrapper comment.
// en: sets wrapper wraps the slog handler.
func ConfigHandleWrapper(wrapper func(Handler) Handler) ConfigOption {
	return func(cfg *Config) { cfg.HandlerWrappers = append(cfg.HandlerWrappers, wrapper) }
}

// ConfigHandleWrapper comment.
// en: override slog handler.
// en: otherwise use default handler.
func ConfigHandler(h Handler) ConfigOption {
	return func(cfg *Config) { cfg.Handler = h }
}
