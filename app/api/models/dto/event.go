package dto

import (
	userRepo "gamma/app/datastore/pg"
	"time"
)

type ResEvent struct {
	EventName        string    `json:"event_name"`
	EventDate        time.Time `json:"event_date"`
	EventLocation    string    `json:"event_location"`
	EventDescription string    `json:"event_description"`
	Uuid             string    `json:"uuid"`
	EventImageUrl    string    `json:"event_image"`
}

func ConvertEvent(event *userRepo.Event) *ResEvent {
	return &ResEvent{
		EventName:        event.EventName,
		EventDate:        event.EventDate,
		EventLocation:    event.EventLocation,
		EventDescription: event.EventDescription,
		Uuid:             event.Uuid,
		EventImageUrl:    event.EventImageUrl,
	}
}

func ConvertEvents(events []*userRepo.Event) []*ResEvent {
	event_list := make([]*ResEvent, len(events))
	for i, event := range events {
		event_list[i] = ConvertEvent(event)
	}
	return event_list
}
