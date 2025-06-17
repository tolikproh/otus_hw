package model

import "time"

type Notification struct {
	ID     int       `json:"id"`
	Title  string    `json:"title"`
	Date   time.Time `json:"date"`
	UserID int       `json:"userId"`
}
