package pg_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gamma/app/datastore/pg"
	"gamma/app/domain/bo"
)

// Ω

var _ = Describe("User", func() {

	testUser := bo.User{
		Id:        0,
		Uuid:      "test",
		Email:     "test@test.com",
		FirstName: "Robert",
		LastName:  "Warwar",
	}

	testOrg := bo.Organization {
		Id: 0,
		Uuid: "asdfasdf",
		OrganizationName: "Gabe's warehouse",
		City: "Translyvania",
	}

	testInvite := bo.UserEventInvite{
		Id: 0,
		Uuid: "0",
		EventUuid: "0",
	}

	testUserEvent := bo.UserEvent{
		Id: 0,
		UserFk: 0, 
		EventFk: 0,
	}

	testEvent := bo.Event{
		Id: 0,
		EventName: "big party",
		EventDate: time.Now(),
		EventLocation: "my house",
		Uuid: "biguuid",
		Organization: 0,
	}
	

	BeforeEach(func() {
		pg.ClearAll()
		err := pg.GetUserRepo().InsertOrganization(context.Background(), testOrg.Uuid, testOrg.OrganizationName, testOrg.City)
		Ω(err).ShouldNot(HaveOccurred())
		err = pg.GetUserRepo().InsertUser(context.Background(), testUser.Uuid, testUser.Email, testUser.FirstName, testUser.LastName)
		Ω(err).ShouldNot(HaveOccurred())
	})

	// ORG USERS

	When("inserting an org user", func() {
		Context("into an org that doesn't exist", func() {
			fakeOrgId := 1
		
			It("should return an error", func() {
				err := pg.GetUserRepo().InsertOrgUser(context.Background(), testUser.Uuid, fakeOrgId)
				Ω(err).Should(HaveOccurred())
			})
		})

		Context("the user does not exist", func() {
			fakeUserUuid := "1"

			It("should return an error", func() {
				err := pg.GetUserRepo().InsertOrgUser(context.Background(), fakeUserUuid, testOrg.Id)
				Ω(err).Should(HaveOccurred())
			})
		})

		Context("org and user is valid", func() {
			It("should insert without error", func() {
				err := pg.GetUserRepo().InsertOrgUser(context.Background(), testUser.Uuid, testOrg.Id)
				Ω(err).ShouldNot(HaveOccurred())
			})
		})

	})

	// INVITES

	When("Inviting user", func() {

		eventId := "0"

		Context("user does not exist", func() {
			eventId := "0"
			fakeUser := bo.User{
				Id: 1,
			}

			It("Should return error", func() {
				err := pg.GetUserRepo().InsertInvite(context.Background(),fakeUser.Uuid, eventId)
				Ω(err).Should(HaveOccurred())
			})
		})

		Context("User does exist", func() {
			It("Should Insert without issue", func() {
				err:=pg.GetUserRepo().InsertInvite(context.Background(), testUser.Uuid, eventId)
				Ω(err).ShouldNot(HaveOccurred())
			})
		})
	})

	When("Geting Invite", func() {

		BeforeEach(func() {
			Expect(pg.GetUserRepo().InsertInvite(context.Background(),testUser.Uuid, testInvite.Uuid)).To(Succeed())		
		})

		Context("Invite does not exist", func() {
			fakeInviteId := "1"
			It("Should return an error", func() {
				_, err := pg.GetUserRepo().GetInvite(context.Background(), fakeInviteId)
				Ω(err).Should(HaveOccurred())
			})
		})

		Context("Invite does exist", func() {
			It("Should return the proper invite with the status of true", func() {
				invite, err := pg.GetUserRepo().GetInvite(context.Background(), testInvite.Uuid)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(invite.Valid).Should(Equal(true))
				Ω(invite.Uuid).Should(Equal(testInvite.Uuid))

			})
		})

	})

	When("Deleting Invites", func() {
		BeforeEach( func() {
			Expect(pg.GetUserRepo().InsertInvite(context.Background(),testUser.Uuid, testInvite.Uuid)).To(Succeed())
		})

		Context("Invite does not exist", func() {
			fakeInviteId := 1
			It("Should return an error", func() {
				err := pg.GetUserRepo().DeleteInvite(context.Background(), fakeInviteId)
				Ω(err).Should(HaveOccurred())
			})
		})

		Context("Invite exists", func() {
			It("Should be deleted", func() {
				err := pg.GetUserRepo().DeleteInvite(context.Background(), testInvite.Id)
				Ω(err).ShouldNot(HaveOccurred())

				_, err = pg.GetUserRepo().GetInvite(context.Background(), testInvite.Uuid)
				Ω(err).Should(HaveOccurred())
			})
		})
	})

	// EVENTS

	When("Inserting Events", func() {

		Context("Organization does not exist", func() {
			fakeOrgId := 1

			It("Should fail", func() {
				err := pg.GetUserRepo().InsertEvent(context.Background(), testEvent.Id, testEvent.EventName, testEvent.EventDate, testEvent.EventLocation, testEvent.Uuid , fakeOrgId)
				Ω(err).Should(HaveOccurred())
			})
		})

		Context("Organization does exist", func() {
			It("should insert event", func() {
				err := pg.GetUserRepo().InsertEvent(context.Background(), testEvent.Id, testEvent.EventName, testEvent.EventDate, testEvent.EventLocation, testEvent.Uuid , testEvent.Organization)
				Ω(err).ShouldNot(HaveOccurred())
			})
		})

	})

	When("Getting Events", func() {
		BeforeEach(func() {
			err := pg.GetUserRepo().InsertEvent(context.Background(), testEvent.Id, testEvent.EventName, testEvent.EventDate, testEvent.EventLocation, testEvent.Uuid , fakeOrgId)
			Ω(err).ShouldNot(HaveOccurred())
		})

		Context("Event does not exist", func() {
			fakeEventUuid := "fake"
			It("Should fail", func() { 
				_, err := pg.GetUserRepo().GetEvent(context.Background(), fakeEventUuid)
				Ω(err).ShouldNot(HaveOccurred())
			})
		})

		Context("Event does exist", func() {
			It("Should pass and match", func() { 
				event, err := pg.GetUserRepo().GetEvent(context.Background(), testEvent.Uuid)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(event.Id)           .Should(Equal(testEvent.Id))
				Ω(event.EventName)    .Should(Equal(testEvent.EventName))
				Ω(event.EventDate)    .Should(Equal(testEvent.EventDate))
				Ω(event.EventLocation).Should(Equal(testEvent.EventLocation))
				Ω(event.Uuid)         .Should(Equal(testEvent.Uuid))
				Ω(event.Organization) .Should(Equal(testEvent.Organization))
			})
		})
	})

	// User Events

	When("Inserting User Events", func() {

		BeforeEach(func() {
			err := pg.GetUserRepo().InsertEvent(context.Background(), testEvent.Id, testEvent.EventName, testEvent.EventDate, testEvent.EventLocation, testEvent.Uuid , fakeOrgId)
			Ω(err).ShouldNot(HaveOccurred())
		})

		Context("user does not exist", func() {
			fakeUserId := 1
			It("Should fail", func() {
				err := pg.GetUserRepo().InsertUserEvent(context.Background(), fakeUserId, testEvent.Id)
				Ω(err).Should(HaveOccurred())
			})
		})

		Context("event does not exist", func() {
			fakeEventId := 1
			It("Should Fail", func() {
				err := pg.GetUserRepo().InsertUserEvent(context.Background(), testUser.Id, fakeEventId)
				Ω(err).Should(HaveOccurred())
			})
		})

		Context("user and event do exist", func() {
			It("Should pass", func() {
				err := pg.GetUserRepo().InsertUserEvent(context.Background(), testUser.Id, testEvent.Id)
				Ω(err).Should(HaveOccurred())
			})
		})
	})

	When("Getting User Events", func() {

	testUsers := []bo.User{
		{Id: 0, Uuid: "uuid1", Email: "test@test.com",  FirstName:"Robert",      LastName: "WarWar"},
		{Id: 1, Uuid: "uuid2", Email: "EMAIL@MAIL.COM", FirstName:"Shyanne",     LastName: "WarWar"},
		{Id: 2, Uuid: "uuid3", Email: "test@gmail.com", FirstName:"Hiram",       LastName: "WarWar"},
		{Id: 3, Uuid: "uuid4", Email: "test@yahoo.com", FirstName:"Caitlynbert", LastName: "WarWar"},
	}

	testEvents := []bo.Event{
		{Id: 1, EventName: "big party", EventDate: time.Now(), EventLocation: "Your moms house", Uuid: "uuid1", Organization: 1},
		{Id: 2, EventName: "halloween", EventDate: time.Now(), EventLocation: "commons",         Uuid: "uuid2", Organization: 2},
		{Id: 3, EventName: "christman", EventDate: time.Now(), EventLocation: "sbisa",           Uuid: "uuid3", Organization: 1},
		{Id: 4, EventName: "4th july",  EventDate: time.Now(), EventLocation: "kyle field",      Uuid: "uuid4", Organization: 3},
		{Id: 5, EventName: "birth",     EventDate: time.Now(), EventLocation: "your dads",       Uuid: "uuid5", Organization: 1},
	}

	testUserEvents := []bo.UserEvent{
		{Id: 0, UserFk: 0, EventFk: 1},
		{Id: 1, UserFk: 1, EventFk: 1},
		{Id: 2, UserFk: 0, EventFk: 2},
		{Id: 3, UserFk: 1, EventFk: 3},
		{Id: 4, UserFk: 2, EventFk: 4},
		{Id: 5, UserFk: 2, EventFk: 5},
	}

	testUsersWithEvents := map[int]map[int]bool{}

		BeforeEach(func() {

			pg.ClearAll()

			for _, user := range testUsers{
				testUsersWithEvents[user.Id] = make(map[int]bool)
				err := pg.GetUserRepo().InsertUser(context.Background(), user.Uuid, user.Email, user.FirstName, user.LastName)
				Ω(err).ShouldNot(HaveOccurred())
			}

			for _, event := range testEvents{
				err := pg.GetUserRepo().InsertEvent(context.Background(), event.Id, event.EventName, event.EventDate, event.EventLocation, event.Uuid , event.Organization)
				Ω(err).ShouldNot(HaveOccurred())
			}

			for _, userEvent := range testUserEvents{
				testUsersWithEvents[userEvent.UserFk][userEvent.EventFk] = true
				err := pg.GetUserRepo().InsertUserEvent(context.Background(), userEvent.UserFk, userEvent.EventFk)
				Ω(err).Should(HaveOccurred())
			}

		})

		Context("user does not exist", func() {
			fakeUserId := -1
			It("Should fail", func() {
				_, err := pg.GetUserRepo().GetUserEvents(context.Background(), fakeUserId)
				Ω(err).Should(HaveOccurred()) 
			})
		})

		Context("user does and is attending to at least 1 event", func() {
			It("Should pass", func() {
				for _, user := range(testUsers) {
					usersEvents, err := pg.GetUserRepo().GetUserEvents(context.Background(), user.Id)
					Ω(err).ShouldNot(HaveOccurred())

					// convert UsersEvents to hashmap
					events := map[int]bool{}
					for _, event := range usersEvents{
						_, contain := events[event.Id]
						Ω(contain).Should(Equal(false), "Returned duplicate event %d", event.Id)
						events[event.Id] = true
					}

					// compares actual against list
					for eventId, _ := range testUsersWithEvents[user.Id]{
						_, contain := events[eventId]
						Ω(contain).Should(Equal(true), "Missing event %d", eventId)
					}
				}
			})
		})

		Context("user exists but is not any events", func(){
			usersEvents, err := pg.GetUserRepo().GetUserEvents(context.Background(), 5)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(len(usersEvents)).Should(Equal(0))
		})
	})

	// TODO: GetOrgUserEvents
	

	
})
