package model

import "time"

type Event struct {
	ID                   int
	Title                string
	DateTimeStart        time.Time
	DateTimeEnd          time.Time
	Description          string
	UserID               int
	NotificationDuration *time.Duration
}
