package pg_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gamma/app/datastore/pg"
	"gamma/app/domain/bo"
)

// Ω

var _ = Describe("User", func() {

	fakeUser := bo.User{
		Id:        0,
		Uuid:      "test",
		Email:     "test@test.com",
		FirstName: "Robert",
		LastName:  "Warwar",
	}

	BeforeEach(func() {
		pg.ClearAll()
		err := pg.GetUserRepo().InsertOrganization(context.Background(), "Gabe's E-girl Warehouse", "Translyvania", "asdfasdf")
		Ω(err).ShouldNot(HaveOccurred())
		err = pg.GetUserRepo().InsertUser(context.Background(), fakeUser.Uuid, fakeUser.Email, fakeUser.FirstName, fakeUser.LastName)
		Ω(err).ShouldNot(HaveOccurred())
	})

	When("inserting an org user", func() {
		It("should insert without a problem", func() {
			pg.GetUserRepo().InsertOrgUser(context.Background(), fakeUser.Uuid)
		})
	})
})
