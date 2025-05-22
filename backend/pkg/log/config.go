package log

import (
	"io"
	"log/slog"
	"os"
)

type Level int

// Handler comment.
// en: create new type to make log independent as much as possible.
type Handler slog.Handler

// Config comment.
// en: contains log application behavior settings.
type Config struct {
	Level slog.Level

	Output io.Writer

	KeysInput []string

	EnableDecoratorNR bool

	Handler Handler

	HandlerWrappers []func(Handler) Handler
}

const (
	DebugLevel Level = Level(slog.LevelDebug)
	InfoLevel  Level = Level(slog.LevelInfo)
)

// defaultConfig comment
// en: creates a Config populated with default settings.
func defaultConfig() Config {
	c := Config{}

	c.Level = slog.LevelInfo
	c.Output = os.Stdout
	c.EnableDecoratorNR = false

	return c
}
