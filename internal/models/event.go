package models

import "time"

type CheckinEvent struct {
	UserID    int64
	CreatedAt time.Time
}

type GratitudeEvent struct {
	UserID    int64
	Message   string
	CreatedAt time.Time
}
