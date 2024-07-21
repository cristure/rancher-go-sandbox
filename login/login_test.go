package login_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cristure/rancher-go-sandbox/login"
)

var _ = Describe("Login", func() {
	var (
		goodCredentials login.Request
		badCredentials  login.Request
	)

	BeforeEach(func() {
		goodCredentials = login.Request{
			Description:  "UI Session",
			ResponseType: "Cookie",
			Username:     "admin",
			Password:     "QcDDbr0PDAOm6ee2",
		}
		badCredentials = login.Request{
			Description:  "UI Session",
			ResponseType: "Cookie",
			Username:     "not_admin",
			Password:     "bad_password",
		}
	})

	Describe("Logging in with set credentails", func() {
		Context("With good credentials", func() {
			It("Should log in successfully", func() {
				// Marshal set of good credentials into JSON format.
				jsonBytes, err := json.Marshal(goodCredentials)
				Expect(err).To(BeNil())

				// Build the http request with marshalled JSON.
				body := bytes.NewBuffer(jsonBytes)
				req, err := http.NewRequest(
					http.MethodPost,
					fmt.Sprintf("%s/v3-public/localProviders/local?action=login", testURL),
					body,
				)
				req.Header.Set("Content-Type", "application/json")
				Expect(err).To(BeNil())

				// Send the request.
				resp, err := httpClient.Do(req)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))

				// Read the response body.
				respBody, err := io.ReadAll(resp.Body)
				Expect(err).To(BeNil())

				// Unmarshal the response into login.Response
				var loginResponse login.Response
				err = json.Unmarshal(respBody, &loginResponse)
				Expect(err).To(BeNil())

				// Assert that the login response token is not empty.
				Expect(loginResponse.Token).ToNot(BeEmpty())
				Expect(loginResponse.Type).To(Equal("token"))

				// Build and send a request to an authorized endpoint.
				// First check that the request can be authenticated with the Authorization Bearer header.
				req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v3/users?me=true", testURL), nil)
				req.Header.Set("Authorization", "Bearer "+loginResponse.Token)
				Expect(err).To(BeNil())

				resp, err = httpClient.Do(req)
				Expect(err).To(BeNil())

				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				// Second check that the authentication works by setting the correct cookie with the provided token.
				req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v3/users?me=true", testURL), nil)
				req.AddCookie(&http.Cookie{Name: login.SessionCookieName, Value: loginResponse.Token})
				Expect(err).To(BeNil())

				resp, err = httpClient.Do(req)
				Expect(err).To(BeNil())

				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})

			It("Shouldn't log in successfully", func() {
				// Marshal set of good credentials into JSON format.
				jsonBytes, err := json.Marshal(badCredentials)
				Expect(err).To(BeNil())

				// Build the http request with marshalled JSON.
				body := bytes.NewBuffer(jsonBytes)
				req, err := http.NewRequest(
					http.MethodPost,
					fmt.Sprintf("%s/v3-public/localProviders/local?action=login", testURL),
					body,
				)
				req.Header.Set("Content-Type", "application/json")
				Expect(err).To(BeNil())

				// Send the request.
				resp, err := httpClient.Do(req)
				Expect(err).To(BeNil())

				// Assert that the request's status is 401.
				Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))

				// Read the response body.
				respBody, err := io.ReadAll(resp.Body)
				Expect(err).To(BeNil())

				// Unmarshal the response into login.Response
				var loginResponse login.Response
				err = json.Unmarshal(respBody, &loginResponse)
				Expect(err).To(BeNil())

				// Assert that the response does not contain a token, and it is marked accordingly.
				Expect(loginResponse.Token).To(BeEmpty())
				Expect(loginResponse.Type).To(Equal("error"))

				// Second check that the authentication doesn't work by setting a wrong token.
				req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v3/users?me=true", testURL), nil)
				req.AddCookie(&http.Cookie{Name: login.SessionCookieName, Value: loginResponse.Token})
				Expect(err).To(BeNil())

				resp, err = httpClient.Do(req)
				Expect(err).To(BeNil())

				Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
			})
		})
	})
})
