package event_test

import (
	"context"
	"gamma/app/services/event"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	USERID = "61784ebf750e3bfb1f849660"
)

var _ = Describe("Event", func() {
	Context("initially", func() {
		user, err := event.EventRepo().GetUserByUUID(context.Background(), USERID)
		Expect(err).Should(BeNil())
		It("exists in database", func() {
			Expect(user.Userid).Should(Equal(1))
		})
		It("has a connection to SigmaChiV2", func() {
			join, err := event.EventRepo().GetOrganizationsByUserID(context.Background(), user.Userid)
			Expect(err).Should(BeNil())
			Expect(len(join)).Should(Equal(1))
			firstJoin := join[0]
			Expect(firstJoin.Permissionscode.Value()).Should(Equal(int64(128)))
			Expect(firstJoin.Name.Value()).Should(Equal("SigmaChiV2"))
		})
		It("has no events", func() {
			events, err := event.EventRepo().GetUserAttendingEvents(context.Background(), user.Userid)
			Expect(err).Should(BeNil())
			Expect(len(events)).Should(Equal(0))
		})
		It("has no applications", func() {
			apps, err := event.EventRepo().GetUserEventApplications(context.Background(), user.Userid)
			Expect(err).Should(BeNil())
			Expect(len(apps)).Should(Equal(0))
		})
	})
})
