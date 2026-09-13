package httpapi

import (
	"context"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"omoidememo.com/telemetry-api/internal/telemetry"
)

type FetchEventsInput struct {
	StartTime time.Time `query:"startTime" doc:"Start of the tracking time range (defaults to 30 days ago)"`
	EndTime   time.Time `query:"endTime"   doc:"End of the tracking time range (defaults to now)"`
}

type FetchEventsOutput struct {
	Body struct {
		Events []telemetry.Event `json:"events"`
	}
}

func createFetchEventsHandler(telemetryStore telemetry.Store) func(context.Context, *FetchEventsInput) (*FetchEventsOutput, error) {
	return func(ctx context.Context, input *FetchEventsInput) (*FetchEventsOutput, error) {
		endTime := time.Now().UTC()
		if !input.EndTime.IsZero() {
			endTime = input.EndTime
		}

		startTime := endTime.AddDate(0, 0, -30)
		if !input.StartTime.IsZero() {
			startTime = input.StartTime
		}

		events, databaseError := telemetryStore.Find(ctx, startTime, endTime)
		if databaseError != nil {
			return nil, huma.Error500InternalServerError("database error", databaseError)
		}

		output := &FetchEventsOutput{}
		output.Body.Events = events
		return output, nil
	}
}
