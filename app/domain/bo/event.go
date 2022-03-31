package bo

import (
	"time"
)

type Event struct {
	Id            int
	Uuid          string
	EventName     string
	EventDate     time.Time
	EventLocation string
	EventImage    string
	Organization  int
}
