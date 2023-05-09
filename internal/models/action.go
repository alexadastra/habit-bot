package models

import "time"

type Action struct {
	ID             string
	Name           string
	Priority       int64
	ScheduledAt    time.Time
	IsCancelled    bool
	LastExecutedAt time.Time
}
