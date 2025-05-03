package chron

import (
	"log/slog"
	"time"
)

func CreateChron(logger *slog.Logger, chronDuration time.Duration, fn func()) {
	ticker := time.NewTicker(chronDuration)
	defer ticker.Stop()

	for range ticker.C {
		fn()
	}
}