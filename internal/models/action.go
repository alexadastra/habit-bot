package models

import (
	"time"

	"github.com/aptible/supercronic/cronexpr"
)

type Action struct {
	ID   string
	Name string
	IsCancelled    bool
	LastExecutedAt time.Time

	IsRepeatable bool
	Priority     int64

	Crontab string
}

func (a *Action) GetNextExecutionTime() time.Time {
	return cronexpr.MustParse(a.Crontab).Next(time.Now().UTC())
}
