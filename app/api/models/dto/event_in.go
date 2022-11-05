package dto

import "time"

type EventUpsert struct {
	EventName        string    `json:"event_name"`
	EventDate        time.Time `json:"event_date"`
	EventLocation    string    `json:"event_location"`
	EventDescription string    `json:"event_description"`
	OrganizationID   string    `json:"organization_id"`
	UserIDs          []string  `json:"user_ids"`
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
