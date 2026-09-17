package httpapi_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"omoidememo.com/telemetry-api/internal/telemetry"
	"omoidememo.com/telemetry-api/internal/testutil"
)

type EventsSuite struct {
	testutil.BaseSuite
}

func TestEventsSuite(t *testing.T) {
	suite.Run(t, new(EventsSuite))
}

func (s *EventsSuite) TestFetchEventsReturnsEmptyByDefault() {
	var out struct {
		Events []telemetry.Event `json:"events"`
	}
	resp, err := s.Client.R().SetResult(&out).Get("/v1/events")
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode())
	s.Require().Empty(out.Events)
}

func (s *EventsSuite) TestFetchEventsReturnsInsertedEvents() {
	now := time.Now().UTC()
	events := []telemetry.Event{
		{
			DeviceID:    testDeviceID,
			Name:        "app_opened",
			TrackedAt:   now.Add(-3 * time.Hour),
			ReceivedAt:  now,
			BatchSentAt: now,
		},
		{
			DeviceID:    testDeviceID,
			Name:        "app_opened",
			TrackedAt:   now.Add(-2 * time.Hour),
			ReceivedAt:  now,
			BatchSentAt: now,
		},
		{
			DeviceID:    testDeviceID,
			Name:        "app_opened",
			TrackedAt:   now.Add(-1 * time.Hour),
			ReceivedAt:  now,
			BatchSentAt: now,
		},
	}
	s.Require().NoError(s.Store.Insert(context.Background(), events))

	var out struct {
		Events []telemetry.Event `json:"events"`
	}
	resp, err := s.Client.R().SetResult(&out).Get("/v1/events")
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode())
	s.Require().Len(out.Events, 3)

	s.Require().True(out.Events[0].TrackedAt.After(out.Events[1].TrackedAt))
	s.Require().True(out.Events[1].TrackedAt.After(out.Events[2].TrackedAt))
	s.Require().Equal(testDeviceID, out.Events[0].DeviceID)
}
