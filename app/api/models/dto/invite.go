package dto

import (
	"time"

	userRepo "gamma/app/datastore/pg"
)

type InviteCreate struct {
	ExpirationDate time.Time `json:"expiration_date"`
	Capacity       int       `json:"invite_capacity"`
	UserUuid       string    `json:"user_uuid"`
	EntityUuid     string    `json:"entity_uuid"`
}

type InviteGetEntity struct {
	EntityUuid string `json:"entity_uuid" query:"entity_uuid"`
}

type InviteGet struct {
	InviteUuid string `json:"invite_uuid" param:"invite_uuid"`
}

type ResInvite struct {
	ExpirationDate time.Time `json:"expiration_date"`
	Capacity       int       `json:"invite_capacity"`
	InviteUuid     string    `json:"uuid"`
	EntityType     int       `json:"entity_type"`
	EntityUuid     string    `json:"entity_uuid"`
}

func ConvertInvite(invite *userRepo.Invite) *ResInvite {
	return &ResInvite{
		ExpirationDate: invite.ExpirationDate,
		Capacity:       int(invite.Capacity),
		InviteUuid:     invite.Uuid,
		EntityUuid:     invite.EntityUuid,
		EntityType:     int(invite.EntityType),
	}
}

func ConvertInvites(invites []*userRepo.Invite) []*ResInvite {
	res := make([]*ResInvite, len(invites))
	for i, invite := range invites {
		res[i] = ConvertInvite(invite)
	}
	return res
}
