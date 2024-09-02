package gen

import (
	"log/slog"
	"os"
)

func init() {
	level := os.Getenv("slogLevel")
	slog.Info("slogLevel", slog.String("level", level))

	defaultLevel := slog.LevelInfo
	if level != "" {
		switch level {
		case "debug":
			defaultLevel = slog.LevelDebug
		case "info":
			defaultLevel = slog.LevelInfo
		case "warn":
			defaultLevel = slog.LevelWarn
		case "error":
			defaultLevel = slog.LevelError
		}
	}
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: defaultLevel})
	log := slog.New(handler)
	slog.SetDefault(log)
}
