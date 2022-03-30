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

	BeforeEach(func() {
		pg.ClearAll()
		var err error
		organizationUuid = uuid.NewString()
		_, err = pg.GetUserRepo().InsertOrganization(ctx, organizationUuid, "Gabes Warehouse", "Translyvania")
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
			err := pg.GetUserRepo().InsertUser(ctx, testUser.Uuid, testUser.Email, "", testUser.FirstName, testUser.LastName, "")
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	When("Getting a user", func() {
		It("should not throw an error", func() {
			err := pg.GetUserRepo().InsertUser(ctx, testUser.Uuid, testUser.Email, "", testUser.FirstName, testUser.LastName, "")
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
})
