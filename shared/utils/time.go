package utils

import "time"

func Now() *time.Time {
	now := time.Now()
	return &now
}
