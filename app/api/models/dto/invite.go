package dto

import "time"

type InviteCreate struct {
	User           *string   `json:"user_uuid"`
	Organization   *string   `json:"org_uuid"`
	Capacity       int       `json:"capacity"`
	ExpirationDate time.Time `json:"exp_date"`
	OrgUserUuid    string    `json:"org_user_uuid"`
}
