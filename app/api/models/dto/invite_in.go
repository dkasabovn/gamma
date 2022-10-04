package dto

import (
	"time"
)

type InviteCreate struct {
	ExpirationDate time.Time `json:"expiration_date"`
	Capacity       int       `json:"invite_capacity"`
	OrganizationID string    `json:"organization_id"`
	UserUuid       string    `json:"user_uuid"`
	EntityUuid     string    `json:"entity_uuid"`
	EntityType     int       `json:"entity_type"`
}

type InviteGet struct {
	InviteID string `json:"invite_id" param:"invite_id"`
}
