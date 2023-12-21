package formatter

import (
	"time"
)

func SecondsInEngineeringNotation(d time.Duration) string {
	s := d.Seconds()
	if s >= 100 {
		return d.Round(time.Second).String()
	} else if s >= 10 {
		return d.Round(100 * time.Millisecond).String()
	} else if s >= 1 {
		return d.Round(10 * time.Millisecond).String()
	}

	ms := d.Milliseconds()
	if ms >= 100 {
		return d.Round(time.Millisecond).String()
	} else if ms >= 10 {
		return d.Round(100 * time.Microsecond).String()
	} else if ms >= 1 {
		return d.Round(10 * time.Microsecond).String()
	}

	us := d.Microseconds()
	if us >= 100 {
		return d.Round(time.Microsecond).String()
	} else if us >= 10 {
		return d.Round(100 * time.Nanosecond).String()
	} else if us >= 1 {
		return d.Round(10 * time.Nanosecond).String()
	}

	return d.String()
}
