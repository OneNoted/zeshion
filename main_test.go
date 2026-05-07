package main

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		name  string
		input string
		level slog.Level
		ok    bool
	}{
		{name: "empty", input: "", ok: false},
		{name: "debug", input: "debug", level: slog.LevelDebug, ok: true},
		{name: "info", input: "info", level: slog.LevelInfo, ok: true},
		{name: "warn", input: "warn", level: slog.LevelWarn, ok: true},
		{name: "warning", input: "warning", level: slog.LevelWarn, ok: true},
		{name: "error", input: "error", level: slog.LevelError, ok: true},
		{name: "case insensitive", input: "DEBUG", level: slog.LevelDebug, ok: true},
		{name: "unknown", input: "trace", ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, ok := parseLogLevel(tt.input)

			assert.Equal(t, tt.ok, ok)
			assert.Equal(t, tt.level, level)
		})
	}
}
