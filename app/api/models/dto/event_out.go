package dto

import (
	userRepo "gamma/app/datastore/pg"
	"time"
)

type Event struct {
	ID               string `json:"EventID"`
	EventName        string
	EventDate        time.Time
	EventLocation    string
	EventDescription string
	EventImageUrl    string
}

type EventWithOrganization struct {
	Event
	Organization
}

type EventWithUserStatus struct {
	Event
	Status int
}

func ConvertEvent(event *userRepo.Event) *Event {
	return &Event{
		ID:               event.ID.String(),
		EventName:        event.EventName,
		EventDate:        event.EventDate,
		EventLocation:    event.EventLocation,
		EventDescription: event.EventDescription,
		EventImageUrl:    event.EventImageUrl,
	}
}

func ConvertEvents(events []*userRepo.Event) []*Event {
	ret := make([]*Event, len(events))
	for i, event := range events {
		ret[i] = ConvertEvent(event)
	}
	return ret
}

func ConvertEventWithOrganization(eventWithOrg *userRepo.GetEventsWithOrganizationsRow) *EventWithOrganization {
	return &EventWithOrganization{
		Event: Event{
			ID:               eventWithOrg.ID.String(),
			EventName:        eventWithOrg.EventName,
			EventDate:        eventWithOrg.EventDate,
			EventLocation:    eventWithOrg.EventLocation,
			EventDescription: eventWithOrg.EventDescription,
			EventImageUrl:    eventWithOrg.EventImageUrl,
		},
		Organization: Organization{
			ID:          eventWithOrg.ID_2,
			OrgName:     eventWithOrg.OrgName,
			City:        eventWithOrg.City,
			OrgImageUrl: eventWithOrg.OrgImageUrl,
		},
	}
}

func ConvertEventsWithOrganizations(eventsWithOrgs []*userRepo.GetEventsWithOrganizationsRow) []*EventWithOrganization {
	ret := make([]*EventWithOrganization, len(eventsWithOrgs))
	for i, eo := range eventsWithOrgs {
		ret[i] = ConvertEventWithOrganization(eo)
	}
	return ret
}

func ConvertUserEvent(userEvent *userRepo.GetUserEventsRow) *EventWithUserStatus {
	return &EventWithUserStatus{
		Event: Event{
			ID:               userEvent.EventFk.String(),
			EventName:        userEvent.EventName,
			EventDate:        userEvent.EventDate,
			EventLocation:    userEvent.EventLocation,
			EventDescription: userEvent.EventDescription,
			EventImageUrl:    userEvent.EventImageUrl,
		},
		Status: int(userEvent.ApplicationState),
	}
}

func ConvertUserEvents(userEvents []*userRepo.GetUserEventsRow) []*EventWithUserStatus {
	ret := make([]*EventWithUserStatus, len(userEvents))
	for i, userEvent := range userEvents {
		ret[i] = ConvertUserEvent(userEvent)
	}
	return ret
}
