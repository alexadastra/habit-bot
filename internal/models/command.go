package models

type Command string

const (
	Checkin   Command = "checkin"
	Gratitude Command = "gratitude"
)

type UserMessage struct {
	Id      int64
	Message string
}

type UserCommand struct {
	Command Command
	UserMessage
}
