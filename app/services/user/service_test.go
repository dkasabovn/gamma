package user_test

import (
	"context"
	"gamma/app/datastore/pg"
	"gamma/app/domain/bo"
	"gamma/app/services/user"
	"gamma/app/system/auth/ecJwt"
	"gamma/app/system/util/tests"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Ω
var _ = BeforeSuite(func() {
	tests.LoadTestKeys()
	pg.ClearAll()
})

var _ = Describe("User", func() {

	ctx := context.Background()

	_password := "securePassword"
	testUser := bo.User{
		Email:     "new_email@email.com",
		FirstName: "joe",
		LastName:  "bob",
		UserName:  "XxJoeBOBxX",
	}

	testOrgUser := bo.OrgUser {
		PoliciesNum: 1,
	}

	testOrg := bo.Organization{
		Uuid: uuid.NewString(),
		OrganizationName: "myOrg",
		City: "Austin",
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

	testOrgInvite := bo.Invite{
		Uuid: "testOrgInvite",
		ExpirationDate: time.Now().Add(time.Hour * 24),
		UseLimit: 5,
		Policy: bo.InvitePolicy {
			InviteType: bo.InviteToOrg,
			InviteTo:  2,
			Constraint: bo.ConstraintEmail,
			Receiver: testUser.Email,
		},
	}





	When("Creating a user", func() {
		It("should not throw an error", func() {

			jwt, err := user.GetUserService().CreateUser(
				ctx,
				_password,
				testUser.Email,
				testUser.PhoneNumber,
				testUser.FirstName,
				testUser.LastName,
				testUser.UserName,
				testUser.ImageUrl,
			)

			Ω(err).ShouldNot(HaveOccurred())
			token, ok := ecJwt.ECDSAVerify(jwt.BearerToken)
			Ω(ok).Should(BeTrue())

			_, ok = ecJwt.ECDSAVerify(jwt.RefreshToken)
			Ω(ok).Should(BeTrue())

			claims := token.Claims.(*ecJwt.GammaClaims)
			testUser.Uuid = claims.Uuid
		})
	})

	When("Signing in with a real user", func() {
		It("should not fail", func() {
			jwt, err := user.GetUserService().SignInUser(
				ctx,
				testUser.Email,
				_password,
			)

			Ω(err).ShouldNot(HaveOccurred())
			_, ok := ecJwt.ECDSAVerify(jwt.BearerToken)
			Ω(ok).Should(BeTrue())

			_, ok = ecJwt.ECDSAVerify(jwt.RefreshToken)
			Ω(ok).Should(BeTrue())
		})
	})

	When("Siging in with an incorrect password", func() {
		It("Should return unathorized", func() {
			jwt, err := user.GetUserService().SignInUser(
				ctx,
				testUser.Email,
				"incorrectPassword",
			)

			Ω(err).ShouldNot(HaveOccurred())
			Ω(jwt).Should(BeNil())

		})
	})

	When("Signing in with a fake user", func() {
		It("should fail", func() {
			jwt, err := user.GetUserService().SignInUser(
				ctx,
				"fake@email.com",
				"fake-password",
			)

			Ω(jwt).Should(BeNil())
			Ω(err).Should(HaveOccurred())
		})
	})

	When("Getting a user with a valid uuid", func() {
		It("Should return the user", func() {
			u, err := user.GetUserService().GetUser(ctx, testUser.Uuid)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(u.Email).Should(Equal(testUser.Email))
			Ω(u.Uuid).Should(Equal(testUser.Uuid))
			Ω(u.FirstName).Should(Equal(testUser.FirstName))
			Ω(u.LastName).Should(Equal(testUser.LastName))
			Ω(u.UserName).Should(Equal(testUser.UserName))
		})
	})

	When("Getting a user that does not exist", func() {
		It("Should return an error", func() {
			u, err := user.GetUserService().GetUser(ctx, uuid.NewString())
			Ω(err).Should(HaveOccurred())
			Ω(u).Should(BeNil())
		})
	})

	// TODO: EVENT INVITES
	When("Inserting an event", func () {
		
	})

	When("Creating an invite to an event with valid credentials", func () {
		It("Should return a new uuid", func() {


			
		})
	})

	When("Creating an invite to an event with invalid credentials", func () {
		It("Should fail", func() {
			Ω(0).Should(Equal(1))
		})
	})

	When("Accepting an invite to an event with a use limit > 0, and the users email is in the policy", func () {
		It("Should successfully insert into the user event table", func() {
			Ω(0).Should(Equal(1))
		})
	})

	When("Accepting an invite to an event with a use limit > 0, and the users email is NOT in the policy", func () {
		It("Should fail to insert into the user event table", func() {
			Ω(0).Should(Equal(1))
		})
	})

	When("Accepting an invite to an event with a use limit > 0, and the the users Org is in the policy", func () {
		It("Should successfully insert into the user event table", func() {
			Ω(0).Should(Equal(1))
		})
	})

	When("Accepting an invite to an event with a use limit > 0, and the the users Org is NOT in the policy", func () {
		It("Should successfully insert into the user event table", func() {
			Ω(0).Should(Equal(1))
		})
	})

	When("Accepting an invite to an event with a use limit <= 0, and the the users Org is in the policy", func () {
		It("Should fail to insert into the user event table", func() {
			Ω(0).Should(Equal(1))
		})
	})

	//TODO: ORG INVITES

	When("Accepting an invite to an org with a use limit < 0, and the the users email is in the policy", func () {
		It("Should successfully insert into the user org table", func() {
			Ω(0).Should(Equal(1))
		})
	})

	When("Accepting an invite to an org with a use limit < 0, and the the users email is NOT in the policy", func () {
		It("Should Fail to insert into the user org table", func() {
			Ω(0).Should(Equal(1))
		})
	})

	When("Accepting an invite to an org with a use limit <= 0, and the the users email is in the policy", func () {
		It("Should fail insert into the user event table", func() {
			Ω(0).Should(Equal(1))
		})
	})



})
