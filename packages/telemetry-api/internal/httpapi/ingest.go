package httpapi

import (
	"context"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/datatypes"

	"omoidememo.com/telemetry-api/internal/telemetry"
)

type IngestInput struct {
	Body *telemetry.IngestBatch
}

type IngestOutput struct {
	Body struct {
		Accepted int `json:"accepted"`
	}
}

func createIngestHandler(s telemetry.Store) func(context.Context, *IngestInput) (*IngestOutput, error) {
	return func(ctx context.Context, input *IngestInput) (*IngestOutput, error) {
		batch := input.Body
		now := time.Now().UTC()
		events := make([]telemetry.Event, len(batch.Events))
		for i, e := range batch.Events {
			props := e.Props
			if len(props) == 0 {
				props = []byte("{}")
			}
			events[i] = telemetry.Event{
				DeviceID:    batch.DeviceID,
				Name:        e.Name,
				Props:       datatypes.JSON(props),
				TrackedAt:   e.TrackedTime(),
				ReceivedAt:  now,
				BatchSentAt: batch.SentAt.UTC(),
			}
		}
		if err := s.Insert(ctx, events); err != nil {
			return nil, huma.Error500InternalServerError("database error", err)
		}

		out := &IngestOutput{}
		out.Body.Accepted = len(events)
		return out, nil
	}
}
