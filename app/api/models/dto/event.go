package dto

import (
	userRepo "gamma/app/datastore/pg"
	"gamma/app/domain/bo"
	"time"
)

type ResEvent struct {
	EventName        string    `json:"event_name"`
	EventDate        time.Time `json:"event_date"`
	EventLocation    string    `json:"event_location"`
	EventDescription string    `json:"event_description"`
	Uuid             string    `json:"uuid"`
	EventImageUrl    string    `json:"event_image"`
	OrgUuid          *string   `json:"org_uuid,omitempty"`
	OrgName          *string   `json:"org_name,omitempty"`
	OrgImage         *string   `json:"org_image,omitempty"`
	ApplicationState string    `json:"state,omitempty"`
}

type ReqEvent struct {
	EventName        string    `form:"event_name"`
	EventDate        time.Time `form:"event_date"`
	EventLocation    string    `form:"event_location"`
	EventDescription string    `form:"event_description"`
	// EventImage handled separately
}

func ConvertOrgEvent(event *userRepo.Event) *ResEvent {
	return &ResEvent{
		EventName:        event.EventName,
		EventDate:        event.EventDate,
		EventLocation:    event.EventLocation,
		EventDescription: event.EventDescription,
		Uuid:             event.Uuid,
		EventImageUrl:    event.EventImageUrl,
	}
}

func ConvertEvent(event *userRepo.GetEventsRow) *ResEvent {
	return &ResEvent{
		EventName:        event.EventName,
		EventDate:        event.EventDate,
		EventLocation:    event.EventLocation,
		EventDescription: event.EventDescription,
		Uuid:             event.Uuid,
		EventImageUrl:    event.EventImageUrl,
		OrgUuid:          &event.Uuid_2,
		OrgName:          &event.OrgName,
		OrgImage:         &event.OrgImageUrl,
		ApplicationState: string(bo.ParseApplicationState(event.ApplicationState)),
	}
}

func ConvertOrgEvents(events []*userRepo.Event) []*ResEvent {
	event_list := make([]*ResEvent, len(events))
	for i, event := range events {
		event_list[i] = ConvertOrgEvent(event)
	}
	return event_list
}

func ConvertEvents(events []*userRepo.GetEventsRow) []*ResEvent {
	event_list := make([]*ResEvent, len(events))
	for i, event := range events {
		event_list[i] = ConvertEvent(event)
	}
	return event_list
}
