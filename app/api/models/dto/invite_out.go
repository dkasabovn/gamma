package dto

import (
	"time"

	"github.com/google/uuid"
)

type Invite struct {
	ID             uuid.UUID
	ExpirationDate time.Time
	Capacity       int32
	OrgUserFk      int32
	EntityUuid     uuid.UUID
	EntityType     int32
}
