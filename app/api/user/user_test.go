package user_api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gamma/app/api/core"
	"gamma/app/api/models/auth"
	user "gamma/app/api/user"
	"gamma/app/datastore/pg"
	"gamma/app/system/util/tests"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = BeforeSuite(func() {
	fmt.Println("BEFORE")
	tests.LoadTestKeys()
	pg.ClearAll()
})

var _ = Describe("API", func() {
	fmt.Print("BEGIN")
	API := user.API()
	REC := httptest.NewRecorder()

	var access_token string
	var refresh_token string

	_signUp := auth.UserSignup{
		Email:       "new_email@email.com",
		RawPassword: "securePassword",
		FirstName:   "gabriel",
		LastName:    "diaz",
		UserName:    "XxBOBxX",
	}

	// _user := bo.User{
	// 	Email: _signUp.Email,
	// 	FirstName: _signUp.FirstName,
	// 	LastName: _signUp.LastName,
	// 	UserName: _signUp.UserName,
	// }

	_signIn := auth.UserSignIn{
		Email:       _signUp.Email,
		RawPassword: _signUp.RawPassword,
	}

	When("Signing Up with a new email that already exists", func() {
		It("Should create a new user", func() {

			data, _ := json.Marshal(_signUp)
			reader := bytes.NewReader(data)
			req := httptest.NewRequest(http.MethodPost, "/auth/signup", reader)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := API.E().NewContext(req, REC)

			err := user.SignUpController(c)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(REC.Code).Should(Equal(http.StatusOK))

			access_token = REC.Result().Header.Get("bearer_token")
			Ω(len(access_token)).ShouldNot(Equal(0))

			refresh_token = core.GetCookie(REC.Result().Cookies(), "refresh_token").Value
			Ω(len(refresh_token)).ShouldNot(Equal(0))
		})
	})

	When("Signing up with a used email", func() {
		It("should fail", func() {
			data, _ := json.Marshal(_signUp)
			reader := bytes.NewReader(data)
			req := httptest.NewRequest(http.MethodPost, "/auth/signup", reader)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := API.E().NewContext(req, REC)

			err := user.SignUpController(c)
			Ω(err).Should(HaveOccurred())

			Ω(REC.Code).Should(Equal(http.StatusInternalServerError))
		})
	})

	When("Logging in with a valid account and correct password", func() {
		It("Should return new cookies", func() {

			data, _ := json.Marshal(_signIn)
			reader := bytes.NewReader(data)
			req := httptest.NewRequest(http.MethodPost, "/auth/signin", reader)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := API.E().NewContext(req, REC)

			err := user.SignInController(c)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(REC.Code).Should(Equal(http.StatusOK))
			access_token = REC.Result().Header.Get("bearer_token")
			Ω(len(access_token)).ShouldNot(Equal(0))

			refresh_token = core.GetCookie(REC.Result().Cookies(), "refresh_token").Value
			Ω(len(refresh_token)).ShouldNot(Equal(0))
		})
	})

	When("Loggin in with an account that does not exist", func() {
		It("Should Faild", func() {
			data, _ := json.Marshal(auth.UserSignIn{
				Email:       "not-an-account@mail.com",
				RawPassword: "password",
			})

			reader := bytes.NewReader(data)
			req := httptest.NewRequest(http.MethodPost, "/auth/signin", reader)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := API.E().NewContext(req, REC)

			err := user.SignInController(c)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(REC.Code).Should(Equal(http.StatusInternalServerError))
		})
	})

	When("Logging-in with an account that does exist with an incorrect password", func() {
		It("Should fail", func() {
			data, _ := json.Marshal(auth.UserSignIn{
				Email:       _signIn.Email,
				RawPassword: "incorrect",
			})

			reader := bytes.NewReader(data)
			req := httptest.NewRequest(http.MethodPost, "/auth/signin", reader)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := API.E().NewContext(req, REC)

			err := user.SignInController(c)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(REC.Code).Should(Equal(http.StatusInternalServerError))
		})
	})

	When("Refreshing with a valid refresh token", func() {
		It("Should Create a new token", func() {

			req := httptest.NewRequest(http.MethodGet, "/auth/refresh", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := API.E().NewContext(req, REC)
			c.SetCookie(&http.Cookie{
				Name:     "refresh_token",
				Value:    refresh_token,
				HttpOnly: true,
			})

			err := user.RefreshTokenController(c)

			Ω(err).ShouldNot(HaveOccurred())
			Ω(REC.Code).Should(Equal(http.StatusOK))

			refresh_cookie := core.GetCookie(REC.Result().Cookies(), "refresh_token")
			Ω(refresh_cookie).ShouldNot(BeNil())
			Ω(refresh_cookie.Value).ShouldNot(Equal(refresh_token))

			new_access_token := REC.Result().Header.Get("bearer_token")
			Ω(new_access_token).ShouldNot(Equal(access_token))

			refresh_token = refresh_cookie.Value
			access_token = new_access_token
		})
	})

	When("Attempting to access endpoints that require auth without tokens", func() {
		It("Should fail", func() {

			req := httptest.NewRequest(http.MethodGet, "/getUser", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := API.E().NewContext(req, REC)

			err := user.RefreshTokenController(c)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(REC.Code).Should(Equal(http.StatusUnauthorized))

		})
	})

	When("Getting user with valid jwt", func() {
		It("Should return the users info", func() {
			req := httptest.NewRequest(http.MethodGet, "/auth/refresh", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(echo.HeaderAuthorization, access_token)
			c := API.E().NewContext(req, REC)

			err := user.GetUserController(c)
			Ω(err).ShouldNot(HaveOccurred())

			var response core.ApiResponse
			body, _ := ioutil.ReadAll(REC.Result().Request.Body)
			json.Unmarshal(body, &response)

			_, ok := response.Data["uuid"]
			Ω(ok).ShouldNot(BeFalse())

		})
	})

})
