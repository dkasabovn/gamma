package user_test

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"syreclabs.com/go/faker"

	"gamma/app/api/models/dto"
	userRepo "gamma/app/datastore/pg"
	"gamma/app/domain/bo"
	"gamma/app/services/user"
)

// Ω
var _ = Describe("Service", func() {
	godotenv.Load("./../../../.env")
	userSvc := user.GetUserService()
	orgUuid := uuid.New()

	BeforeEach(func() {
		userSvc.CreateOrganization(context.Background(), &userRepo.InsertOrganizationParams{
			ID:          orgUuid,
			OrgName:     faker.Company().Name(),
			City:        faker.Address().City(),
			OrgImageUrl: faker.Internet().Url(),
		})
	})

	AfterEach(func() {
		userSvc.DANGER()
	})

	Describe("Signing Up User", func() {
		Context("without errors", func() {
			It("should not fail", func() {
				userParams := &dto.UserSignUp{
					Email:       faker.Internet().Email(),
					PhoneNumber: faker.PhoneNumber().PhoneNumber(),
					RawPassword: faker.Internet().Password(8, 16),
					FirstName:   faker.Name().FirstName(),
					LastName:    faker.Name().LastName(),
					UserName:    faker.Internet().UserName(),
				}
				partialUser, err := userSvc.SignUpUser(context.Background(), userParams)
				Ω(err).ShouldNot(HaveOccurred())

				user, err := userSvc.GetUser(context.Background(), partialUser.UUID)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(user.Email).Should(Equal(userParams.Email))
			})
		})

		// TODO: Test for bad input once validators exist
	})

	Describe("Signing In User", func() {
		userParams := &dto.UserSignUp{
			Email:       faker.Internet().Email(),
			PhoneNumber: faker.PhoneNumber().PhoneNumber(),
			RawPassword: faker.Internet().Password(8, 16),
			FirstName:   faker.Name().FirstName(),
			LastName:    faker.Name().LastName(),
			UserName:    faker.Internet().UserName(),
		}

		var partialUser *bo.PartialUser

		JustBeforeEach(func() {
			var err error
			partialUser, err = userSvc.SignUpUser(context.Background(), userParams)
			Ω(err).ShouldNot(HaveOccurred())
		})

		Context("with matching email and password", func() {
			It("should not fail", func() {
				signedInUser, err := userSvc.SignInUser(context.Background(), userParams.Email, userParams.RawPassword)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(signedInUser.UUID).Should(Equal(partialUser.UUID))
			})
		})

		Context("with matching email but not matching password", func() {
			It("should fail", func() {
				signedInUser, err := userSvc.SignInUser(context.Background(), userParams.Email, faker.Internet().Password(8, 16))
				Ω(err).Should(HaveOccurred())
				Ω(signedInUser).Should(BeNil())
			})
		})

		Context("without matching email or password", func() {
			It("should fail", func() {
				signedInUser, err := userSvc.SignInUser(context.Background(), faker.Internet().Email(), faker.Internet().Password(8, 16))
				Ω(err).Should(HaveOccurred())
				Ω(signedInUser).Should(BeNil())
			})
		})
	})

	Describe("Inserting org user", func() {
		userParams := &dto.UserSignUp{
			Email:       faker.Internet().Email(),
			PhoneNumber: faker.PhoneNumber().PhoneNumber(),
			RawPassword: faker.Internet().Password(8, 16),
			FirstName:   faker.Name().FirstName(),
			LastName:    faker.Name().LastName(),
			UserName:    faker.Internet().UserName(),
		}

		var partialUser *bo.PartialUser

		JustBeforeEach(func() {
			var err error
			partialUser, err = userSvc.SignUpUser(context.Background(), userParams)
			Ω(err).ShouldNot(HaveOccurred())
		})

		Context("With pre-existing organization and user", func() {
			It("should not fail", func() {
				err := userSvc.CreateOrgUser(context.Background(), &userRepo.InsertOrgUserParams{
					PoliciesNum:    int32(bo.Create(bo.CREATE_EVENTS)),
					UserFk:         partialUser.UUID,
					OrganizationFk: orgUuid,
				})
				Ω(err).ShouldNot(HaveOccurred())
			})
		})
	})

	Describe("Inserting event", func() {
		userParams := &dto.UserSignUp{
			Email:       faker.Internet().Email(),
			PhoneNumber: faker.PhoneNumber().PhoneNumber(),
			RawPassword: faker.Internet().Password(8, 16),
			FirstName:   faker.Name().FirstName(),
			LastName:    faker.Name().LastName(),
			UserName:    faker.Internet().UserName(),
		}

		var partialUser *bo.PartialUser

		JustBeforeEach(func() {
			var err error
			partialUser, err = userSvc.SignUpUser(context.Background(), userParams)
			Ω(err).ShouldNot(HaveOccurred())
		})

		Context("With correct rights", func() {
			It("shouldn't fail", func() {
				eventParams := &dto.EventUpsert{
					EventName:        faker.App().Name(),
					EventDate:        faker.Time().Forward(time.Hour * 300),
					EventLocation:    faker.Address().StreetAddress(),
					EventDescription: faker.Lorem().Paragraph(5),
					OrganizationID:   orgUuid.String(),
					EventImage:       []byte{},
				}

				err := userSvc.CreateEvent(context.Background(), &userRepo.GetUserWithOrgRow{
					ID:             0,
					PoliciesNum:    int32(bo.Create(bo.CREATE_EVENTS)),
					UserFk:         partialUser.UUID,
					OrganizationFk: orgUuid,
					ID_2:           orgUuid,
				}, eventParams)

				Ω(err).ShouldNot(HaveOccurred())
			})
		})

		Context("Without correct rights", func() {
			It("should fail", func() {
				eventParams := &dto.EventUpsert{
					EventName:        faker.App().Name(),
					EventDate:        faker.Time().Forward(time.Hour * 300),
					EventLocation:    faker.Address().StreetAddress(),
					EventDescription: faker.Lorem().Paragraph(5),
					OrganizationID:   orgUuid.String(),
					EventImage:       []byte{},
				}

				err := userSvc.CreateEvent(context.Background(), &userRepo.GetUserWithOrgRow{
					ID:             0,
					PoliciesNum:    int32(bo.Create(bo.AFFILIATED)),
					UserFk:         partialUser.UUID,
					OrganizationFk: orgUuid,
					ID_2:           orgUuid,
				}, eventParams)

				Ω(err).Should(HaveOccurred())
			})
		})
	})

	// TODO: test create invite
})
