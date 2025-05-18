package model

import "time"

type Notification struct {
	ID     int
	Title  string
	Date   time.Time
	UserID int
}
