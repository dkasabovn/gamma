package userRepo_test

import (
	"context"
	"gamma/app/datastore"
	userRepo "gamma/app/datastore/pg"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"syreclabs.com/go/faker"
)

// Ω

var _ = Describe("./App/Datastore/Pg/User", func() {
	repo := userRepo.New(datastore.RwInstance())
	userData := &userRepo.InsertUserParams{
		Uuid:         uuid.NewString(),
		Email:        faker.Internet().Email(),
		PasswordHash: uuid.NewString(),
		PhoneNumber:  faker.PhoneNumber().PhoneNumber(),
		FirstName:    faker.Name().FirstName(),
		LastName:     faker.Name().LastName(),
		ImageUrl:     faker.Internet().Url(),
		Validated:    false,
		RefreshToken: uuid.NewString(),
	}

	var userId int

	orgData := &userRepo.InsertOrganizationParams{
		OrgName:     faker.Company().Name(),
		City:        faker.Address().StreetAddress(),
		Uuid:        uuid.NewString(),
		OrgImageUrl: faker.Internet().Url(),
	}

	var orgId int

	eventData := &userRepo.InsertEventParams{
		EventName:        faker.Company().Name(),
		EventDate:        faker.Time().Between(time.Now(), time.Now().Add(time.Hour*4000)),
		EventLocation:    faker.Address().StreetAddress(),
		EventDescription: faker.RandomString(200),
		Uuid:             uuid.NewString(),
		EventImageUrl:    faker.Internet().Url(),
		OrganizationFk:   0,
	}

	BeforeEach(func() {
		err := repo.InsertUser(context.Background(), userData)
		Ω(err).ShouldNot(HaveOccurred())
		err = repo.InsertOrganization(context.Background(), orgData)
		Ω(err).ShouldNot(HaveOccurred())

		user, err := repo.GetUserByUuid(context.Background(), userData.Uuid)
		Ω(err).ShouldNot(HaveOccurred())
		userId = int(user.ID)
		organization, err := repo.GetOrganizationByUuid(context.Background(), orgData.Uuid)
		Ω(err).ShouldNot(HaveOccurred())
		orgId = int(organization.ID)

		err = repo.InsertOrgUser(context.Background(), &userRepo.InsertOrgUserParams{
			PoliciesNum:    69,
			UserFk:         user.ID,
			OrganizationFk: organization.ID,
		})
		Ω(err).ShouldNot(HaveOccurred())

		eventData.OrganizationFk = int32(orgId)
		err = repo.InsertEvent(context.Background(), eventData)
		Ω(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		repo.TruncateAll(context.Background())
	})

	When("getting a user by uuid", func() {
		It("should not fail", func() {
			user, err := repo.GetUserByUuid(context.Background(), userData.Uuid)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(user.FirstName).Should(Equal(userData.FirstName))
		})
	})

	When("getting a user by email", func() {
		It("should not fail", func() {
			user, err := repo.GetUserByEmail(context.Background(), userData.Email)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(user.FirstName).Should(Equal(userData.FirstName))
		})
	})

	When("getting a user and org user by uuid", func() {
		It("should properly fetch both the org user and user", func() {
			user_org_user, err := repo.GetUserOrgUserJoin(context.Background(), userData.Uuid)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(user_org_user.PoliciesNum).Should(Equal(int32(69)))
			Ω(user_org_user.FirstName).Should(Equal(userData.FirstName))
		})
	})

	When("getting user organizations", func() {
		It("should return the users' organizations", func() {
			orgs, err := repo.GetUserOrganizations(context.Background(), int32(userId))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(orgs[0].City).Should(Equal(orgData.City))
			Ω(len(orgs)).Should(Equal(1))
		})
	})

	When("getting organization events", func() {
		It("should get all events that are part of an organization", func() {
			events, err := repo.GetOrganizationEvents(context.Background(), orgData.Uuid)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(len(events)).Should(Equal(1))
			Ω(events[0].EventDescription).Should(Equal(eventData.EventDescription))
		})
	})

	// TODO: test GetUserEvents

	When("getting events by uuid and id", func() {
		It("should not fail", func() {
			event, err := repo.GetEventByUuid(context.Background(), eventData.Uuid)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(event.EventDescription).Should(Equal(eventData.EventDescription))
			event, err = repo.GetEventById(context.Background(), event.ID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(event.EventDescription).Should(Equal(eventData.EventDescription))
		})
	})
})
