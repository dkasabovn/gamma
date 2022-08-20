package dto

import (
	userRepo "gamma/app/datastore/pg"
	"time"

	"github.com/google/uuid"
)

type Invite struct {
	ID             uuid.UUID
	ExpirationDate time.Time
	Capacity       int32
	EntityUuid     uuid.UUID
	EntityType     int32
}

func ConvertInvite(invite *userRepo.Invite) *Invite {
	return &Invite{
		ID:             invite.ID,
		ExpirationDate: invite.ExpirationDate,
		Capacity:       invite.Capacity,
		EntityUuid:     invite.EntityUuid,
		EntityType:     invite.EntityType,
	}
}

func ConvertInvites(invites []*userRepo.Invite) []*Invite {
	ret := make([]*Invite, len(invites))
	for i, invite := range invites {
		ret[i] = ConvertInvite(invite)
	}
	return ret
}
