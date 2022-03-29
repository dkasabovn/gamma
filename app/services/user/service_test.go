package user_test

import (
	"context"
	"fmt"
	"gamma/app/datastore/pg"
	"gamma/app/domain/bo"
	"gamma/app/services/user"
	"gamma/app/system/auth/ecJwt"
	"gamma/app/system/util/tests"

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
	_user := bo.User{
		Email : "new_email@email.com",
		FirstName: "joe",
		LastName: "bob",
		UserName: "XxJoeBOBxX",
	}

	When("Creating a user", func() {
		It("should not throw an error", func() {

			jwt, err := user.GetUserService().CreateUser(
				ctx,
				_user.Email,
				_password,
				_user.FirstName,
				_user.LastName,
				_user.UserName,
			)

			Ω(err).ShouldNot(HaveOccurred())
			token, ok := ecJwt.ECDSAVerify(jwt.BearerToken)
			Ω(ok).Should(BeTrue())

			_, ok = ecJwt.ECDSAVerify(jwt.RefreshToken)
			Ω(ok).Should(BeTrue())

			claims := token.Claims.(*ecJwt.GammaClaims)
			_user.Uuid = claims.Uuid
			fmt.Println(_user.Uuid)
		})
	})

	When("Signing in with a real user", func() {
		It("should not fail", func() {
			jwt, err := user.GetUserService().SignInUser(
				ctx,
				_user.Email,
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
				_user.Email,
				"incorrectPassword",
			)

			Ω(err).Should(HaveOccurred())
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
			u, err := user.GetUserService().GetUser(ctx, _user.Uuid)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(u.Email).Should(Equal(_user.Email))
			Ω(u.Uuid).Should(Equal(_user.Uuid))
			Ω(u.FirstName).Should(Equal(_user.FirstName))
			Ω(u.LastName).Should(Equal(_user.LastName))
			Ω(u.UserName).Should(Equal(_user.UserName))
		})
	})

	When("Getting a user that does not exist", func() {
		It("Should return an error", func() {
			u, err := user.GetUserService().GetUser(ctx, uuid.NewString())
			Ω(err).Should(HaveOccurred())
			Ω(u).Should(BeNil())
		})
	})

})
