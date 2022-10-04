package dto

import "time"

type EventUpsert struct {
	EventName           string    `form:"event_name"`
	EventDate           time.Time `form:"event_date"`
	EventLocation       string    `form:"event_location"`
	EventDescription    string    `form:"event_description"`
	OrganizationID      string    `form:"organization_id"`
	OrganizationUserIDs []string  `form:"user_ids"`
	EventImage          []byte    `json:"-" form:"-"`
}

type EventSearch struct {
	OrganizationID *string    `query:"organization_id"`
	DateFloor      *time.Time `query:"date_floor"`
}

type EventValidate struct {
	OtherJWT       string `json:"user_jwt"`
	EventID        string `json:"event_id"`
	OrganizationID string `json:"organization_id"`
}
