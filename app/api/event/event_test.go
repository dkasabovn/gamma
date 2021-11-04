package event_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo"
	"go.mongodb.org/mongo-driver/bson/primitive"

	// . "github.com/onsi/gomega"

	"gamma/app/api/event"
	"gamma/app/system/auth/ecJwt"
)

var _ = Describe("Event integration tests", func() {
	id, _ := primitive.ObjectIDFromHex("61784ebf750e3bfb1f849660")
	claims := &ecJwt.GammaClaims{
		Uuid:  id,
		Email: "gkitt@gkitt.gkitt",
	}
	accessToken, _ := ecJwt.ECDSASign(claims)
	req := httptest.NewRequest(http.MethodGet, "/organizations", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.Set("user", jwt.NewWithClaims(jwt.SigningMethodES256, claims))
	event.GetBootstrapData(c)
})
