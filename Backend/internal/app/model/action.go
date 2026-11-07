package model

import "time"

type ActionLog struct {
	ID        int
	UserID    int
	Action    string
	Target    string
	IP        string
	CreatedAt time.Time
}
