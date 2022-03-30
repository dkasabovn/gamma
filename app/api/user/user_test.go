package user_api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gamma/app/api/core"
	"gamma/app/api/models/auth"
	user_api "gamma/app/api/user"
	"gamma/app/datastore/pg"
	"gamma/app/domain/bo"
	"gamma/app/system/util/tests"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	URL = "http://localhost:6969"
	BAD_TOKEN = "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1dWlkIjoiYmJlZWI5NzUtNWE5My00NzU4LWExN2QtMWMwMGNiYzI4ZmYzIiwiYXVkIjoidXNlci5nYW1tYSIsImV4cCI6MTY0OTIxMjE2MCwiaXNzIjoiYXV0aC5nYW1tYSJ9.ikSuerG_t-WgXFMKi9ReaW2PVYDC6tHrmfixYFxdV4KJ2HBfrB_vvdAirPGbWEGhqDj_RGHM7BEuZdwwvowPW3Q"
)

var _ = BeforeSuite(func() {
	tests.LoadTestKeys()
	pg.ClearAll()

	go user_api.StartAPI(":6969")
})

var _ = Describe("API", func() {

	var access_token string
	var refresh_token string

	var CLIENT *http.Client

	var _ = BeforeEach(func() {
		CLIENT = &http.Client{}
	})
	
	_signUp := auth.UserSignup{
		Email : "new_email@email.com",
		RawPassword: "securePassword",
		FirstName: "gabriel",
		LastName: "diaz",
		UserName: "XxBOBxX",
	}

	_user := bo.User{
		Email: _signUp.Email,
		FirstName: _signUp.FirstName,
		LastName: _signUp.LastName,
		UserName: _signUp.UserName,
	}

	_signIn := auth.UserSignIn{
		Email: _signUp.Email,
		RawPassword: _signUp.RawPassword,
	}
	


	When("Signing Up with a new email", func() {
		It("Should create a new user", func() {

			data, _ := json.Marshal(_signUp)
			reader := bytes.NewReader(data)
			req, err := http.NewRequest(
				echo.POST,
				fmt.Sprintf("%s/auth/signup", URL),
				reader,
			)

			Ω(err).ShouldNot(HaveOccurred())
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			response, err := CLIENT.Do(req)

			Ω(err).ShouldNot(HaveOccurred())
			Ω(response.StatusCode).Should(Equal(http.StatusOK))

			var api_response core.ApiResponse
			body, _ := ioutil.ReadAll(response.Body)
			json.Unmarshal(body, &api_response)

			access_token = fmt.Sprint(api_response.Data["bearer_token"])
			Ω(len(access_token)).ShouldNot(Equal(0))

			refresh_token = core.GetCookie(response.Cookies(), "refresh_token").Value
			Ω(len(refresh_token)).ShouldNot(Equal(0))
		})
	})

	When("Signing up with a used email", func() {
		It("should fail", func() {
			data, _ := json.Marshal(_signUp)
			reader := bytes.NewReader(data)
			req, err := http.NewRequest(
				echo.POST,
				fmt.Sprintf("%s/auth/signup", URL),
				reader,
			)
			
			Ω(err).ShouldNot(HaveOccurred())
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			response, err := CLIENT.Do(req)

			
			Ω(err).ShouldNot(HaveOccurred())
			Ω(response.StatusCode).Should(Equal(http.StatusInternalServerError))
		})
	})

	When("Logging in with a valid account and correct password", func() {
		It("Should return new cookies", func() {
			
			data, _ := json.Marshal(_signIn)
			reader := bytes.NewReader(data)
			req, err := http.NewRequest(
				echo.POST,
				 fmt.Sprintf("%s/auth/signin", URL),
				reader,
			)

			Ω(err).ShouldNot(HaveOccurred())
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			response, err := CLIENT.Do(req)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(response.StatusCode).Should(Equal(http.StatusOK))

			
			var api_response core.ApiResponse
			body, _ := ioutil.ReadAll(response.Body)
			json.Unmarshal(body, &api_response)

			access_token = fmt.Sprint(api_response.Data["bearer_token"])			
			Ω(len(access_token)).ShouldNot(Equal(0))

			refresh_token = core.GetCookie(response.Cookies(), "refresh_token").Value
			Ω(len(refresh_token)).ShouldNot(Equal(0))
		})
	})

	When("Loggin in with an account that does not exist",func() {
		It("Should Faild", func() {
			data, _ := json.Marshal(auth.UserSignIn{
				Email: "not-an-account@mail.com",
				RawPassword: "password",
			})

			reader := bytes.NewReader(data)
			req,err := http.NewRequest(
				echo.POST,
				fmt.Sprintf("%s/auth/signin", URL),
				reader,
			)

			Ω(err).ShouldNot(HaveOccurred())
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			response, err := CLIENT.Do(req)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(response.StatusCode).Should(Equal(http.StatusInternalServerError))
		})
	})

	When("Logging-in with an account that does exist with an incorrect password", func() {
		It("Should fail", func() {
			data, _ := json.Marshal(auth.UserSignIn{
				Email: _signIn.Email,
				RawPassword: "incorrect",
			})

			reader := bytes.NewReader(data)
			req, err := http.NewRequest(
				echo.POST,
				fmt.Sprintf("%s/auth/signin",URL),
				reader,
			)

			Ω(err).ShouldNot(HaveOccurred())
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			response, err := CLIENT.Do(req)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(response.StatusCode).Should(Equal(http.StatusUnauthorized))
		})
	})

	When("Refreshing with a valid refresh token", func() {
		It("Should Create a new token", func() {

			req, err := http.NewRequest(
				echo.GET,
				fmt.Sprintf("%s/auth/refresh",URL),
				nil,
			)
			Ω(err).ShouldNot(HaveOccurred())

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set("refresh_token", refresh_token)
			req.AddCookie(&http.Cookie{
				Name:     "refresh_token",
				Value:    refresh_token,
				HttpOnly: true,
			})
			
			response, err := CLIENT.Do(req)
			
			Ω(err).ShouldNot(HaveOccurred())
			Ω(response.StatusCode).Should(Equal(http.StatusOK))

			var api_response core.ApiResponse
			body, _ := ioutil.ReadAll(response.Body)
			json.Unmarshal(body, &api_response)

			new_access_token := fmt.Sprint(api_response.Data["bearer_token"])
			Ω(len(new_access_token)).ShouldNot(Equal(0))

			refresh_cookie := core.GetCookie(response.Cookies(), "refresh_token")
			Ω(refresh_cookie).ShouldNot(BeNil())
			Ω(refresh_cookie.Value).ShouldNot(Equal(refresh_token))

			refresh_token = refresh_cookie.Value
			access_token = new_access_token
		})
	})

	When("Attempting to access endpoints that require auth without tokens", func() {
		It("Should fail", func() {

			req, err := http.NewRequest(
				echo.GET,
				fmt.Sprintf("%s/api/user",URL),
				nil,
			)

			Ω(err).ShouldNot(HaveOccurred())
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Add(echo.HeaderAuthorization, "Bearer " + BAD_TOKEN)

			response, err := CLIENT.Do(req)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(response.StatusCode).Should(Equal(http.StatusUnauthorized))
		})
	})

	When("Getting user with valid jwt", func() {
		It("Should return the users info", func() {
			req,err := http.NewRequest(echo.GET,
				fmt.Sprintf("%s/api/user", URL),
				nil,
			)
			
			Ω(err).ShouldNot(HaveOccurred())
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Add(echo.HeaderAuthorization, "Bearer " + access_token)

			response, err := CLIENT.Do(req)
			Ω(err).ShouldNot(HaveOccurred())

			var api_response core.ApiResponse
			body, _ := ioutil.ReadAll(response.Body)
			json.Unmarshal(body, &api_response)


			resp_user, ok := api_response.Data["user"]
			Ω(ok).Should(BeTrue())
			u_map, _ := resp_user.(map[string]interface{})
			jsonString, _ := json.Marshal(u_map)
			u := bo.User{}
			json.Unmarshal(jsonString, &u)
			
			Ω(u.Email).Should(Equal(_user.Email))
			Ω(u.FirstName).Should(Equal(_user.FirstName))
			Ω(u.LastName).Should(Equal(_user.LastName))
			Ω(u.UserName).Should(Equal(_user.UserName))
			

		})
	})

})



