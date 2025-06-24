package config

import (
	"log/slog"
	"os"
)

func SetupLogger() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				switch level {
				case slog.LevelDebug:
					a.Value = slog.StringValue("🔍 DEBUG")
				case slog.LevelInfo:
					a.Value = slog.StringValue("ℹ️  INFO")
				case slog.LevelWarn:
					a.Value = slog.StringValue("⚠️  WARN")
				case slog.LevelError:
					a.Value = slog.StringValue("🚨 ERROR")
				}
			}
			return a
		},
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
