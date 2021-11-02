package event

import (
	"database/sql"
	"gamma/app/datastore/events"
	"sync"
)

var (
	eventSingle sync.Once
	eventRepo   *EventRepository
)

type EventRepository struct {
	db *sql.DB
}

func EventRepo() *EventRepository {
	eventSingle.Do(func() {
		eventRepo = &EventRepository{
			db: events.EventDB(),
		}
	})
	return eventRepo
}
