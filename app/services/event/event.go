package event

import (
	"database/sql"
	"gamma/app/datastore/events"
	"gamma/app/datastore/mongodb"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	eventSingle sync.Once
	eventRepo   *EventRepository
)

type EventRepository struct {
	db *sql.DB
	invites *mongo.Collection
}

func EventRepo() *EventRepository {
	eventSingle.Do(func() {
		eventRepo = &EventRepository{
			db: events.EventDB(),
			invites: mongodb.MongoDB().Collection(mongodb.InvitesCollectionName),
		}
	})
	return eventRepo
}
