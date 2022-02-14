package pg_test

import (
	"context"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gamma/app/datastore/pg"
	"gamma/app/domain/bo"
)

// Ω

var _ = AfterSuite(func() {
	pg.ClearAll()
})

var _ = Describe("User", func() {

	ctx := context.Background()

	var organizationId int

	inviteUuid := uuid.NewString()

	testUser := bo.User{
		Id:        0,
		Uuid:      "test",
		Email:     "test@test.com",
		FirstName: "Robert",
		LastName:  "Warwar",
	}

	testEvent := bo.Event{
		EventName:     "The tri-state ipa guy competition",
		EventDate:     time.Now(),
		EventLocation: "Alaska",
		Uuid:          uuid.NewString(),
	}

	BeforeEach(func() {
		pg.ClearAll()
		var err error
		organizationId, err = pg.GetUserRepo().InsertOrganization(ctx, uuid.New().String(), "Gabes Warehouse", "Translyvania")
		Ω(err).ShouldNot(HaveOccurred())
	})

	When("Inserting an event", func() {
		It("should not throw an error", func() {
			err := pg.GetUserRepo().InsertEvent(ctx, testEvent.EventName, testEvent.EventDate, testEvent.EventLocation, testEvent.Uuid, organizationId)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	When("Getting a non-existent event", func() {
		It("should throw an error", func() {
			_, err := pg.GetUserRepo().GetEvent(ctx, uuid.NewString())
			Ω(err).Should(HaveOccurred())
		})
	})

	When("Getting a real event", func() {
		It("should get the event", func() {
			err := pg.GetUserRepo().InsertEvent(ctx, testEvent.EventName, testEvent.EventDate, testEvent.EventLocation, testEvent.Uuid, organizationId)
			Ω(err).ShouldNot(HaveOccurred())
			event, err := pg.GetUserRepo().GetEvent(ctx, testEvent.Uuid)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(event.EventLocation).Should(Equal(testEvent.EventLocation))
		})
	})

	When("Inserting a user", func() {
		It("should not throw an error", func() {
			err := pg.GetUserRepo().InsertUser(ctx, testUser.Uuid, testUser.Email, "", testUser.FirstName, testUser.LastName)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	When("Getting a user", func() {
		It("should not throw an error", func() {
			err := pg.GetUserRepo().InsertUser(ctx, testUser.Uuid, testUser.Email, "", testUser.FirstName, testUser.LastName)
			Ω(err).ShouldNot(HaveOccurred())
			user, err := pg.GetUserRepo().GetUser(ctx, testUser.Uuid)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(user.FirstName).Should(Equal(testUser.FirstName))
			Ω(user.LastName).Should(Equal(testUser.LastName))
		})
	})

	When("Getting a non-existant user", func() {
		It("should throw an error", func() {
			_, err := pg.GetUserRepo().GetUser(ctx, uuid.NewString())
			Ω(err).Should(HaveOccurred())
		})
	})

	When("Inserting an org user", func() {
		It("should not throw an error", func() {
			err := pg.GetUserRepo().InsertUser(ctx, testUser.Uuid, testUser.Email, "", testUser.FirstName, testUser.LastName)
			Ω(err).ShouldNot(HaveOccurred())
			err = pg.GetUserRepo().InsertOrgUser(ctx, testUser.Uuid, organizationId)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	When("Getting org user events", func() {
		It("should get the events within the same organization", func() {
			err := pg.GetUserRepo().InsertUser(ctx, testUser.Uuid, testUser.Email, "", testUser.FirstName, testUser.LastName)
			Ω(err).ShouldNot(HaveOccurred())
			err = pg.GetUserRepo().InsertOrgUser(ctx, testUser.Uuid, organizationId)
			Ω(err).ShouldNot(HaveOccurred())
			err = pg.GetUserRepo().InsertEvent(ctx, testEvent.EventName, testEvent.EventDate, testEvent.EventLocation, testEvent.Uuid, organizationId)
			Ω(err).ShouldNot(HaveOccurred())
			user, err := pg.GetUserRepo().GetUser(ctx, testUser.Uuid)
			Ω(err).ShouldNot(HaveOccurred())
			events, err := pg.GetUserRepo().GetOrgUserEvents(ctx, user.OrgUserFk)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(events).Should(HaveLen(1))
		})
	})

	When("Getting user events", func() {
		It("should get events the user is attending", func() {
			err := pg.GetUserRepo().InsertUser(ctx, testUser.Uuid, testUser.Email, "", testUser.FirstName, testUser.LastName)
			Ω(err).ShouldNot(HaveOccurred())
			err = pg.GetUserRepo().InsertEvent(ctx, testEvent.EventName, testEvent.EventDate, testEvent.EventLocation, testEvent.Uuid, organizationId)
			Ω(err).ShouldNot(HaveOccurred())
			user, err := pg.GetUserRepo().GetUser(ctx, testUser.Uuid)
			Ω(err).ShouldNot(HaveOccurred())
			event, err := pg.GetUserRepo().GetEvent(ctx, testEvent.Uuid)
			Ω(err).ShouldNot(HaveOccurred())
			err = pg.GetUserRepo().InsertUserEvent(ctx, user.Id, event.Id)
			Ω(err).ShouldNot(HaveOccurred())
			events, err := pg.GetUserRepo().GetUserEvents(ctx, user.Id)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(events).Should(HaveLen(1))
			Ω(events[0].EventLocation).Should(Equal(testEvent.EventLocation))
		})
	})

	When("Inserting an invite", func() {
		It("should not fail", func() {
			err := pg.GetUserRepo().InsertEvent(ctx, testEvent.EventName, testEvent.EventDate, testEvent.EventLocation, testEvent.Uuid, organizationId)
			Ω(err).ShouldNot(HaveOccurred())
			err = pg.GetUserRepo().InsertInvite(ctx, inviteUuid, testEvent.Uuid)
			Ω(err).ShouldNot(HaveOccurred())
			invite, err := pg.GetUserRepo().GetInvite(ctx, inviteUuid)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(invite.EventUuid).Should(Equal(testEvent.Uuid))
			Ω(invite.Uuid).Should(Equal(inviteUuid))
		})
	})

	When("Deleting an invite", func() {
		It("should delete the invite", func() {
			err := pg.GetUserRepo().InsertEvent(ctx, testEvent.EventName, testEvent.EventDate, testEvent.EventLocation, testEvent.Uuid, organizationId)
			Ω(err).ShouldNot(HaveOccurred())
			err = pg.GetUserRepo().InsertInvite(ctx, inviteUuid, testEvent.Uuid)
			Ω(err).ShouldNot(HaveOccurred())
			invite, err := pg.GetUserRepo().GetInvite(ctx, inviteUuid)
			Ω(err).ShouldNot(HaveOccurred())
			err = pg.GetUserRepo().DeleteInvite(ctx, invite.Id)
			Ω(err).ShouldNot(HaveOccurred())
			_, err = pg.GetUserRepo().GetInvite(ctx, inviteUuid)
			Ω(err).Should(HaveOccurred())
		})
	})
})
