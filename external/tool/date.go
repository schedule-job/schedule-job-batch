package tool

import "time"

func NewUTCDate(year int, month time.Month, day int, hour int, min int) time.Time {
	return time.Date(year, month, day, hour, min, 0, 0, time.UTC)
}
