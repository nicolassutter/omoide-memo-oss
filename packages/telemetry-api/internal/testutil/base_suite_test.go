package testutil

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gorm.io/datatypes"

	"omoidememo.com/telemetry-api/internal/telemetry"
)

type HarnessSmokeSuite struct {
	BaseSuite
}

func TestHarnessSmokeSuite(t *testing.T) {
	suite.Run(t, new(HarnessSmokeSuite))
}

func (s *HarnessSmokeSuite) TestHarnessIsWired() {
	resp, err := s.Client.R().Get("/health")
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode())

	var decoded map[string]string
	s.Require().NoError(json.Unmarshal(resp.Bytes(), &decoded))
	s.Require().Equal("ok", decoded["status"])
}

func (s *HarnessSmokeSuite) TestResetClearsEvents() {
	event := telemetry.Event{
		DeviceID:    "00000000-0000-0000-0000-000000000001",
		Name:        "app_opened",
		Props:       datatypes.JSON([]byte(`{}`)),
		TrackedAt:   time.Now().UTC(),
		ReceivedAt:  time.Now().UTC(),
		BatchSentAt: time.Now().UTC(),
	}
	s.Require().NoError(s.Store.Insert(context.Background(), []telemetry.Event{event}))

	var before []telemetry.Event
	s.Require().NoError(s.DB.Find(&before).Error)
	s.Require().Len(before, 1)

	s.Reset()

	var after []telemetry.Event
	s.Require().NoError(s.DB.Find(&after).Error)
	s.Require().Empty(after)
}
