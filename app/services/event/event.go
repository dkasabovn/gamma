package event

import (
	"database/sql"
	"gamma/app/datastore/event"
	"sync"
)

var (
	eventSingle sync.Once
	eventRepo   *EventRepository
)

type EventRepository struct {
	db *sql.DB
}

func Events() *EventRepository {
	eventSingle.Do(func() {
		eventRepo = &EventRepository{
			db: event.EventDB(),
		}
	})
	return eventRepo
}
