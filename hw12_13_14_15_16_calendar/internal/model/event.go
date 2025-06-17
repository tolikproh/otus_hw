package model

import "time"

type Event struct {
	ID                   int            `json:"id"`
	Title                string         `json:"title"`
	DateTimeStart        time.Time      `json:"dateTimeStart"`
	DateTimeEnd          time.Time      `json:"dateTimeEnd"`
	Description          string         `json:"description"`
	UserID               int            `json:"userId"`
	NotificationDuration *time.Duration `json:"notificationDuration"`
}
