// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package userRepo

import (
	"encoding/json"
	"time"
)

type Event struct {
	ID               int32
	EventName        string
	EventDate        time.Time
	EventLocation    string
	EventDescription string
	Uuid             string
	EventImageUrl    string
	OrganizationFk   int32
}

type Invite struct {
	ID             int32
	ExpirationDate time.Time
	Capacity       int32
	PolicyJson     json.RawMessage
	Uuid           string
	OrgUserUuid    string
	OrgFk          int32
}

type OrgUser struct {
	ID             int32
	PoliciesNum    int32
	UserFk         int32
	OrganizationFk int32
}

type Organization struct {
	ID          int32
	OrgName     string
	City        string
	Uuid        string
	OrgImageUrl string
}

type User struct {
	ID           int32
	Uuid         string
	Email        string
	PasswordHash string
	PhoneNumber  string
	FirstName    string
	LastName     string
	ImageUrl     string
	Validated    bool
	RefreshToken string
}

type UserEvent struct {
	ID               int32
	UserFk           int32
	EventFk          int32
	ApplicationState string
}
