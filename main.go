package main

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"

	"charm.land/fang/v2"
	"github.com/notes/zesh/zeshcli"
)

var version = "dev"

func main() {
	setupLogger()

	cmd := zeshcli.NewRootCommand(version)
	if err := fang.Execute(context.Background(), cmd, fang.WithColorSchemeFunc(fang.AnsiColorScheme), fang.WithoutVersion()); err != nil {
		slog.Error("command failed", "error", err)
		os.Exit(1)
	}
}

func setupLogger() {
	level, ok := parseLogLevel(os.Getenv("ZESH_LOG"))
	if !ok {
		slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
		return
	}

	writer := io.Writer(os.Stderr)
	if path := os.Getenv("ZESH_LOG_FILE"); path != "" {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
		if err != nil {
			slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})))
			slog.Warn("couldn't open log file", "path", path, "error", err)
			return
		}
		writer = file
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(writer, &slog.HandlerOptions{Level: level})))
}

func parseLogLevel(value string) (slog.Level, bool) {
	switch strings.ToLower(value) {
	case "debug":
		return slog.LevelDebug, true
	case "info":
		return slog.LevelInfo, true
	case "warn", "warning":
		return slog.LevelWarn, true
	case "error":
		return slog.LevelError, true
	default:
		return 0, false
	}
}
