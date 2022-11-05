package dto

import "time"

type EventUpsert struct {
	EventName        string    `form:"event_name" json:"event_name"`
	EventDate        time.Time `form:"event_date" json:"event_date"`
	EventLocation    string    `form:"event_location" json:"event_location"`
	EventDescription string    `form:"event_description" json:"event_description"`
	OrganizationID   string    `form:"organization_id" json:"organization_id"`
	UserIDs          []string  `form:"user_ids" json:"user_ids"`
	EventImage       []byte    `json:"-" form:"-"`
}

type EventSearch struct {
	OrganizationID *string    `query:"organization_id"`
	DateFloor      *time.Time `query:"date_floor"`
}

type EventCheck struct {
	UserID  string `json:"user_id"`
	EventID string `json:"event_id"`
}
