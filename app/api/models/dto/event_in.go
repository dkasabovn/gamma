package dto

import (
	"time"
)

type EventUpsert struct {
	EventName        string    `form:"event_name"`
	EventDate        time.Time `form:"event_date"`
	EventLocation    string    `form:"event_location"`
	EventDescription string    `form:"event_description"`
	OrganizationID   string    `form:"organization_id"`
	EventImage       []byte    `json:"-" form:"-"`
}

type EventSearch struct {
	OrganizationID *string    `query:"org_id"`
	DateFloor      *time.Time `query:"date_floor"`
}

type EventGet struct {
	EventID string `json:"event_id" param:"event_id"`
}

type EventDelete struct {
	EventID        string `json:"event_id" param:"event_id"`
	OrganizationID string `form:"organization_id"`
}

type EventUpdate struct {
	EventID          string    `json:"event_id" param:"event_id"`
	EventName        string    `form:"event_name"`
	EventDate        time.Time `form:"event_date"`
	EventLocation    string    `form:"event_location"`
	EventDescription string    `form:"event_description"`
	OrganizationID   string    `form:"organization_id"`
	EventImage       []byte    `json:"-" form:"-"`
}
