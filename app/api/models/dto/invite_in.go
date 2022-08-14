package dto

import (
	"time"
)

type InviteCreate struct {
	ExpirationDate time.Time `json:"expiration_date"`
	Capacity       int       `json:"invite_capacity"`
	UserUuid       string    `json:"user_uuid"`
	EntityUuid     string    `json:"entity_uuid"`
}
