package app

import "time"

type Clock interface {
	Now() time.Time
}

type TimeClock struct{}

func (TimeClock) Now() time.Time {
	return time.Now()
}

type SingleTimestampClock struct {
	Timestamp time.Time
}

func (s SingleTimestampClock) Now() time.Time {
	return s.Timestamp
}
