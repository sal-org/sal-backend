package util

import "time"

// GetCurrentTime - current time in utc
func GetCurrentTime() time.Time {
	return time.Now().UTC()
}
