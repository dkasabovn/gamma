package bo

import (
	"time"
)

type Event struct {
	Id            int
	EventName     string
	EventDate     time.Time
	EventLocation string
	EventImage    string
	Uuid          string
	Organization  int
}
