package pg_test

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gamma/app/datastore/pg"
	"gamma/app/domain/bo"
)

// Ω

// var _ = AfterSuite(func() {
// 	pg.ClearAll()
// })

var _ = Describe("User", func() {

	ctx := context.Background()

	var organizationUuid string

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

	testEventInvite := bo.Invite{
		Uuid: "testEventInvite",
		ExpirationDate: time.Now().Add(time.Hour * 24),
		UseLimit: 1,
		Policy: bo.InvitePolicy {
			InviteType: bo.InviteToEvent,
			InviteTo:  1,
			Constraint: bo.ConstraintEmail,
			Receiver: testUser.Email,
		},
	}

	// testOrgInvite := bo.Invite{
	// 	Uuid: "testOrgInvite",
	// 	ExpirationDate: time.Now().Add(time.Hour * 24),
	// 	UseLimit: 5,
	// 	Policy: bo.InvitePolicy {
	// 		InviteType: bo.InviteToOrg,
	// 		InviteTo:  2,
	// 		Constraint: bo.ConstraintEmail,
	// 		Receiver: testUser.Email,
	// 	},
	// }

	BeforeEach(func() {
		pg.ClearAll()
		var err error
		organizationUuid = uuid.NewString()
		_, err = pg.GetUserRepo().InsertOrganization(ctx, organizationUuid, "Gabes Warehouse", "Translyvania", "http://poggers.com")
		Ω(err).ShouldNot(HaveOccurred())
	})

	When("Inserting an event by organization id", func() {
		It("should not throw an error", func() {
			event, err := pg.GetUserRepo().InsertEventByOrganization(ctx, organizationUuid, &testEvent)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(event.Id).ShouldNot(Equal(0))
		})
	})

	When("Getting a non-existent event", func() {
		It("should throw an error", func() {
			_, err := pg.GetUserRepo().GetEvent(ctx, uuid.NewString())
			Ω(err).Should(HaveOccurred())
		})
	})

	When("Inserting a user", func() {
		It("should not throw an error", func() {
			err := pg.GetUserRepo().InsertUser(ctx, testUser.Uuid, testUser.Email, "123-456-7891", "", testUser.FirstName, testUser.LastName, "", "http://monkey.com")
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	When("Getting a user", func() {
		It("should not throw an error", func() {
			err := pg.GetUserRepo().InsertUser(ctx, testUser.Uuid, testUser.Email, "123-456-7891", "", testUser.FirstName, testUser.LastName, "", "http://monkey.com")
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

	When("Creating an invite", func() {
		It("Should create a new event", func () {
			err := pg.GetUserRepo().InsertInvite(
				ctx,
				testEventInvite.ExpirationDate,
				testEventInvite.UseLimit,
				testEventInvite.Uuid,
				testEventInvite.Policy,
			)

			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	When("Retriving an invite", func() {
		It("Should not fail", func() {

			err := pg.GetUserRepo().InsertInvite(
				ctx,
				testEventInvite.ExpirationDate,
				testEventInvite.UseLimit,
				testEventInvite.Uuid,
				testEventInvite.Policy,
			)

			Ω(err).ShouldNot(HaveOccurred())

			invite, err := pg.GetUserRepo().GetInvite(ctx, testEventInvite.Uuid)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(invite.ExpirationDate.Day()).Should(Equal(testEventInvite.ExpirationDate.Day()))
			Ω(invite.UseLimit).Should(Equal(testEventInvite.UseLimit))
			Ω(invite.Policy.InviteTo).Should(Equal(testEventInvite.Policy.InviteTo))
			Ω(invite.Policy.InviteType).Should(Equal(testEventInvite.Policy.InviteType))
			Ω(invite.Policy.Receiver).Should(Equal(testEventInvite.Policy.Receiver))
			Ω(invite.Policy.Constraint).Should(Equal(testEventInvite.Policy.Constraint))
		})
	})

	When("Decrementing an invite", func() {
		It("should decrease the use limit by 1", func() {
			err := pg.GetUserRepo().InsertInvite(
				ctx,
				testEventInvite.ExpirationDate,
				testEventInvite.UseLimit,
				testEventInvite.Uuid,
				testEventInvite.Policy,
			)
			Ω(err).ShouldNot(HaveOccurred())
			invite, err := pg.GetUserRepo().GetInvite(ctx, testEventInvite.Uuid)
			Ω(err).ShouldNot(HaveOccurred())

			err = pg.GetUserRepo().DecrementInvite(ctx, invite.Id)
			Ω(err).ShouldNot(HaveOccurred())

			invite, err = pg.GetUserRepo().GetInvite(ctx, testEventInvite.Uuid)
			Ω(err).ShouldNot(HaveOccurred())

			fmt.Printf("%v+", invite)
			Ω(invite.UseLimit).Should(Equal(testEventInvite.UseLimit-1))
		})
	})
})
