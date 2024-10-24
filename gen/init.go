package gen

import (
	"log/slog"
	"os"
	"time"

	"github.com/MatusOllah/slogcolor"
	"github.com/fatih/color"
)

func init() {
	level := os.Getenv("slogLevel")
	slog.Info("slogLevel", slog.String("level", level))

	defaultLevel := slog.LevelError
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
	// handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: defaultLevel})
	log := slog.New(slogcolor.NewHandler(os.Stdout, &slogcolor.Options{
		Level:         defaultLevel,
		TimeFormat:    time.DateTime,
		SrcFileMode:   slogcolor.ShortFile,
		SrcFileLength: 0,
		MsgPrefix:     color.HiWhiteString("| "),
		MsgLength:     0,
		MsgColor:      color.New(),
		NoColor:       false,
	}))
	slog.SetDefault(log)
}
