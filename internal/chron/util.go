package chron

import (
	"time"
)

func CreateChron(chronDuration time.Duration, fn func()) {
	ticker := time.NewTicker(chronDuration)
	defer ticker.Stop()

	for range ticker.C {
		fn()
	}
}