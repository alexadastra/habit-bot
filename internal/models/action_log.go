package models

import "time"

type ActionLog struct {
	ID               string
	ActionID         string
	ExecutedAt       time.Time
	DurationMillisec int64
	Result           interface{}
}
