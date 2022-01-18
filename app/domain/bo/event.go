package bo

import (
	"time"
)

type Event struct {
	Id            int
	EventName     string
	EventDate     time.Time
	EventLocation string
	Uuid          string
	Organization  int
}
