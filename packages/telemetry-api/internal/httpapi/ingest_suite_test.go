package httpapi_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"omoidememo.com/telemetry-api/internal/telemetry"
	"omoidememo.com/telemetry-api/internal/testutil"
)

const testDeviceID = "11111111-1111-1111-1111-111111111111"

type IngestSuite struct {
	testutil.BaseSuite
}

func TestIngestSuite(t *testing.T) {
	suite.Run(t, new(IngestSuite))
}

func (s *IngestSuite) TestIngestAcceptsValidBatch() {
	batch := telemetry.IngestBatch{
		DeviceID: testDeviceID,
		SentAt:   time.Now().UTC(),
		Events: []telemetry.IngestEvent{
			{Name: "app_opened", TrackedAt: time.Now().UnixMilli()},
			{
				Name:      "memo_created",
				Props:     json.RawMessage(`{"source":"manual"}`),
				TrackedAt: time.Now().UnixMilli(),
			},
		},
	}

	var output struct {
		Accepted int `json:"accepted"`
	}
	response, err := s.PublicClient.R().
		SetBody(batch).
		SetResult(&output).
		Post("/v1/events")
	s.Require().NoError(err)
	s.Require().Equal(http.StatusAccepted, response.StatusCode())
	s.Require().Equal(2, output.Accepted)

	var stored []telemetry.Event
	s.Require().NoError(s.DB.Order("id").Find(&stored).Error)
	s.Require().Len(stored, 2)
	s.Require().Equal(testDeviceID, stored[0].DeviceID)
	s.Require().Equal("app_opened", stored[0].Name)
	s.Require().Equal(testDeviceID, stored[1].DeviceID)
	s.Require().Equal("memo_created", stored[1].Name)
}

func (s *IngestSuite) TestIngestRejectsUnknownEvent() {
	batch := telemetry.IngestBatch{
		DeviceID: testDeviceID,
		SentAt:   time.Now().UTC(),
		Events: []telemetry.IngestEvent{
			{Name: "bogus_event", TrackedAt: time.Now().UnixMilli()},
		},
	}

	response, err := s.PublicClient.R().SetBody(batch).Post("/v1/events")
	s.Require().NoError(err)
	s.Require().Equal(http.StatusUnprocessableEntity, response.StatusCode())
}
