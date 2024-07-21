package login_test

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"

	"github.com/cristure/rancher-go-sandbox/login"
)

var (
	httpClient  *http.Client
	testURL     string
	mappedPorts = map[string]nat.Port{
		"http":  "80/tcp",
		"https": "443/tcp",
	}
)

func TestMain(m *testing.M) {
	scheme, ok := os.LookupEnv("SCHEME")
	if !ok {
		scheme = "http"
	}

	req := testcontainers.ContainerRequest{
		Image:        "rancher/rancher",
		ExposedPorts: []string{"80/tcp", "443/tcp"},
		Privileged:   true,
		Mounts: testcontainers.ContainerMounts{
			{
				Source: testcontainers.GenericVolumeMountSource{
					Name: "my_test_volume",
				},
				Target: "/var/lib/rancher",
			},
		},
	}

	container, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	ip, err := container.Host(context.Background())
	if err != nil {
		panic(err)
	}

	mappedPort, err := container.MappedPort(context.Background(), mappedPorts[scheme])
	if err != nil {
		panic(err)
	}

	// Set TLS config for bad certificates.
	httpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}}

	testURL = fmt.Sprintf("%s://%s:%s", scheme, ip, mappedPort.Port())

	err = waitToBeReady()
	if err != nil {
		panic(err)
	}

	code := m.Run()
	os.Exit(code)
}

func TestLogin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Login Suite")
}

func waitToBeReady() error {
	// usually this should be a /readyz or /healthcheck.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	var err error
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("%w: %w", ctx.Err(), err)

		case <-ticker.C:
			// Marshal set of good credentials into JSON format.
			jsonBytes, jsonErr := json.Marshal(login.Request{
				Description:  "UI Session",
				ResponseType: "Cookie",
				Username:     "admin",
				Password:     "QcDDbr0PDAOm6ee2",
			})
			if jsonErr != nil {
				err = jsonErr
				continue
			}

			// Build the http request with marshalled JSON.
			body := bytes.NewBuffer(jsonBytes)
			req, newReqErr := http.NewRequest(
				http.MethodPost,
				fmt.Sprintf("%s/v3-public/localProviders/local?action=login", testURL),
				body,
			)
			req.Header.Set("Content-Type", "application/json")
			if newReqErr != nil {
				err = newReqErr
				continue
			}

			resp, sendReqErr := httpClient.Do(req)
			if sendReqErr != nil {
				err = sendReqErr
				continue
			}

			if resp.StatusCode == http.StatusCreated {
				// sometimes the first authentication request gives a 500.
				return nil
			}
		}
	}
}
