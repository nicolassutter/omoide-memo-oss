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
	var output struct {
		Events []telemetry.Event `json:"events"`
	}
	response, err := s.AdminClient.R().SetResult(&output).Get("/v1/events")
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, response.StatusCode())
	s.Require().Empty(output.Events)
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

	var output struct {
		Events []telemetry.Event `json:"events"`
	}
	response, err := s.AdminClient.R().SetResult(&output).Get("/v1/events")
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, response.StatusCode())
	s.Require().Len(output.Events, 3)

	s.Require().True(output.Events[0].TrackedAt.After(output.Events[1].TrackedAt))
	s.Require().True(output.Events[1].TrackedAt.After(output.Events[2].TrackedAt))
	s.Require().Equal(testDeviceID, output.Events[0].DeviceID)
}
