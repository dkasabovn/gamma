package definition

import (
	"context"
	"gamma/app/domain/bo"
	"time"
)

type UserRepository interface {
	GetUser(ctx context.Context, uuid string) (*bo.User, error)
	GetUserByEmail(ctx context.Context, email string) (*bo.User, error)
	InsertUser(ctx context.Context, uuid string, email string, hash string, firstName string, lastName string, userName string) error
	GetUserEvents(ctx context.Context, userId int) ([]bo.Event, error)
	InsertUserEvent(ctx context.Context, userId, eventId int) error
	GetEvent(ctx context.Context, eventUuid string) (*bo.Event, error)
	InsertEvent(ctx context.Context, event_name string, event_date time.Time, event_location string, uuid string, orgFk int) error
	InsertOrganization(ctx context.Context, uuid, name, city string) (int, error)
	InsertOrgUser(ctx context.Context, userId int, organizationId int, policiesNum int) error
	GetUserOrganizations(ctx context.Context, userId int) ([]bo.OrganizationUser, error)
	InsertEventByOrganization(ctx context.Context, orgUuid string, event *bo.Event) (*bo.Event, error)
}
