package httpapi_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
	"resty.dev/v3"

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

	resp, err := clientNoKey.R().Get("/v1/events")
	s.Require().NoError(err)
	s.Require().Equal(http.StatusUnauthorized, resp.StatusCode())
}
