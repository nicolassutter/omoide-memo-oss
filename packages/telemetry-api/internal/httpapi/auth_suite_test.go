package httpapi_test

import (
	"net/http"
	"net/http/httptest"
	"resty.dev/v3"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"omoidememo.com/telemetry-api/internal/config"
	"omoidememo.com/telemetry-api/internal/httpapi"
	"omoidememo.com/telemetry-api/internal/telemetry"
	"omoidememo.com/telemetry-api/internal/testutil"
)

type AuthSuite struct {
	testutil.BaseSuite
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}

func (s *AuthSuite) TestRequestsWithoutApiKeyAreRejected() {
	clientNoKey := resty.New().SetBaseURL(s.ServerURL)

	response, err := clientNoKey.R().Get("/v1/events")
	s.Require().NoError(err)
	s.Require().Equal(http.StatusUnauthorized, response.StatusCode())
}

func (s *AuthSuite) TestPublicKeyForbiddenOnFetch() {
	response, err := s.PublicClient.R().Get("/v1/events")
	s.Require().NoError(err)
	s.Require().Equal(http.StatusUnauthorized, response.StatusCode())
}

func (s *AuthSuite) TestAdminKeyRequiredOnFetch() {
	response, err := s.AdminClient.R().Get("/v1/events")
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, response.StatusCode())
}

func (s *AuthSuite) TestEitherKeyAcceptedOnIngest() {
	batch := telemetry.IngestBatch{
		DeviceID: "11111111-1111-1111-1111-111111111111",
		SentAt:   time.Now().UTC(),
		Events: []telemetry.IngestEvent{
			{Name: "app_opened", TrackedAt: time.Now().UnixMilli()},
		},
	}

	s.Run("public key accepted", func() {
		response, err := s.PublicClient.R().SetBody(batch).Post("/v1/events")
		s.Require().NoError(err)
		s.Require().Equal(http.StatusAccepted, response.StatusCode())
	})

	s.Run("admin key accepted", func() {
		response, err := s.AdminClient.R().SetBody(batch).Post("/v1/events")
		s.Require().NoError(err)
		s.Require().Equal(http.StatusAccepted, response.StatusCode())
	})

	s.Run("missing key rejected", func() {
		noKeyClient := resty.New().SetBaseURL(s.ServerURL)
		response, err := noKeyClient.R().SetBody(batch).Post("/v1/events")
		s.Require().NoError(err)
		s.Require().Equal(http.StatusUnauthorized, response.StatusCode())
	})
}

func (s *AuthSuite) TestWrongKeyIsRejected() {
	wrongClient := resty.New().
		SetBaseURL(s.ServerURL).
		SetHeader("X-Api-Public-Key", "definitely-not-the-key").
		SetHeader("X-Admin-Api-Key", "also-not-the-key")

	batch := telemetry.IngestBatch{
		DeviceID: "11111111-1111-1111-1111-111111111111",
		SentAt:   time.Now().UTC(),
		Events: []telemetry.IngestEvent{
			{Name: "app_opened", TrackedAt: time.Now().UnixMilli()},
		},
	}

	s.Run("POST /v1/events with wrong keys", func() {
		response, err := wrongClient.R().SetBody(batch).Post("/v1/events")
		s.Require().NoError(err)
		s.Require().Equal(http.StatusUnauthorized, response.StatusCode())
	})

	s.Run("GET /v1/events with wrong admin key", func() {
		response, err := wrongClient.R().Get("/v1/events")
		s.Require().NoError(err)
		s.Require().Equal(http.StatusUnauthorized, response.StatusCode())
	})

	s.Run("admin key in public header still accepted", func() {
		adminValue := s.AdminClient.Header().Get("X-Admin-Api-Key")
		adminInPublicHeader := resty.New().
			SetBaseURL(s.ServerURL).
			SetHeader("X-Api-Public-Key", adminValue)
		response, err := adminInPublicHeader.R().SetBody(batch).Post("/v1/events")
		s.Require().NoError(err)
		s.Require().Equal(http.StatusAccepted, response.StatusCode())
	})
}

func (s *AuthSuite) TestDevModeBypassesAuth() {
	handler, _ := httpapi.NewRouter(config.Config{DevMode: true}, s.Store)
	devServer := httptest.NewServer(handler)
	defer devServer.Close()

	noKeyClient := resty.New().SetBaseURL(devServer.URL)

	s.Run("GET /v1/events without key", func() {
		response, err := noKeyClient.R().Get("/v1/events")
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, response.StatusCode())
	})

	batch := telemetry.IngestBatch{
		DeviceID: "11111111-1111-1111-1111-111111111111",
		SentAt:   time.Now().UTC(),
		Events: []telemetry.IngestEvent{
			{Name: "app_opened", TrackedAt: time.Now().UnixMilli()},
		},
	}
	s.Run("POST /v1/events without key", func() {
		response, err := noKeyClient.R().SetBody(batch).Post("/v1/events")
		s.Require().NoError(err)
		s.Require().Equal(http.StatusAccepted, response.StatusCode())
	})
}
