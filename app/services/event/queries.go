package event

import (
	"context"
	"gamma/app/datastore/events/models"
)

func (e *EventRepository) GetUserByUUID(ctx context.Context, uuid string) (*models.User, error) {
	return models.UserByUseruuid(ctx, e.db, uuid)
}

func (e *EventRepository) GetOrganizationsByUserID(ctx context.Context, id int) ([]*models.UserOrgJoin, error) {
	return models.UserOrgJoinsByUserID(ctx, e.db, id)
}

func (e *EventRepository) GetUserAttendingEvents(ctx context.Context, id int) ([]*models.Userevent, error) {
	return models.UsereventsByUserfk(ctx, e.db, id)
}

func (e *EventRepository) GetUserEventApplications(ctx context.Context, id int) ([]*models.Eventapplication, error) {
	return models.EventapplicationsByUserfk(ctx, e.db, id)
}
