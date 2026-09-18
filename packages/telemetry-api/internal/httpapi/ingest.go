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
		for i, event := range batch.Events {
			props := event.Props
			if len(props) == 0 {
				props = []byte("{}")
			}
			events[i] = telemetry.Event{
				DeviceID:    batch.DeviceID,
				Name:        event.Name,
				Props:       datatypes.JSON(props),
				TrackedAt:   event.TrackedTime(),
				ReceivedAt:  now,
				BatchSentAt: batch.SentAt.UTC(),
			}
		}
		if err := s.Insert(ctx, events); err != nil {
			return nil, huma.Error500InternalServerError("database error", err)
		}

		output := &IngestOutput{}
		output.Body.Accepted = len(events)
		return output, nil
	}
}
