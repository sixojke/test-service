package domain

import "time"

type Project struct {
	Id        int
	Name      string
	CreatedAt time.Time
}
