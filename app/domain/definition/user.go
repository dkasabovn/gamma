package definition

import (
	"context"
	"database/sql"
	"gamma/app/domain/bo"
	"time"
)

type UserRepository interface {
	GetUser(ctx context.Context, uuid string) (*bo.User, error)
	InsertUser(ctx context.Context, uuid string, email string, firstName string, lastName string) error
	GetOrgUserEvents(ctx context.Context, orgUserFk sql.NullInt64) ([]bo.Event, error)
	InsertInvite(ctx context.Context, inviteUuid string, eventUuid string) error
	GetInvite(ctx context.Context, inviteUuid string) (*bo.UserEventInvite, error)
	DeleteInvite(ctx context.Context, inviteId int) error
	GetUserEvents(ctx context.Context, userId int) ([]bo.Event, error)
	InsertUserEvent(ctx context.Context, userId, eventId int) error
	GetEvent(ctx context.Context, eventUuid string) (*bo.Event, error)
	InsertEvent(ctx context.Context, event_name string, event_date time.Time, event_location string, uuid string, orgFk int) error
	InsertOrganization(ctx context.Context, uuid, name, city string) (int, error)
	InsertOrgUser(ctx context.Context, userUuid string, organizationId int) error
}
