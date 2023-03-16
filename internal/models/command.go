package models

import "time"

type Command string

const (
	Start     Command = "start"
	Checkin   Command = "checkin"
	Gratitude Command = "gratitude"
	Stats	  Command = "stats"
)

type UserMessage struct {
	UserID  int64
	Message string
	SentAt  time.Time
}

type UserCommand struct {
	Command Command
	UserMessage
}
