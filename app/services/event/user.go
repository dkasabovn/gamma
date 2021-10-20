package event

import (
	"context"
	"database/sql"
	"gamma/app/datastore/event/models"
)

func (e *EventRepository) GetUserByUUID(ctx context.Context, uuid string) (*models.User, error) {
	return models.UserByUserUUID(ctx, e.db, sql.NullString{
		String: uuid,
	})
}
