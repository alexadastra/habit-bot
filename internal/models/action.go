package models

import (
	"time"
)

type Action struct {
	ID          string
	IsCancelled bool

	ActionExecutionerProperties

	ActionSchedulerProperties
}

type ActionExecutionerProperties struct {
	Name     string
	Priority int64
}

type ActionSchedulerProperties struct {
	Crontab      string
	IsRepeatable bool
	ScheduledAt  time.Time
}
