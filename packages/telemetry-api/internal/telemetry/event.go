package telemetry

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/danielgtaylor/huma/v2"
	"time"

	"gorm.io/datatypes"
)

type Event struct {
	ID          int64          `gorm:"primaryKey"`
	DeviceID    string         `gorm:"type:uuid;not null;index:idx_device_tracked,priority:1"`
	Name        string         `gorm:"type:text;not null;index:idx_name_tracked,priority:1"`
	Props       datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'"`
	TrackedAt   time.Time      `gorm:"not null;index:idx_device_tracked,priority:2;index:idx_name_tracked,priority:2"`
	ReceivedAt  time.Time      `gorm:"not null"`
	BatchSentAt time.Time      `gorm:"not null"`
}

func (Event) TableName() string { return "events" }

var EventValidators = map[string]func(json.RawMessage) error{
	"app_opened":               nil,
	"signup_completed":         nil,
	"login_completed":          nil,
	"review_session_completed": nil,
	"folder_created":           nil,
	"onboarding_complete":      nil,
	// track how the memo was created
	"memo_created": func(props json.RawMessage) error {
		var source struct {
			Source string `json:"source"`
		}
		if err := json.Unmarshal(props, &source); err != nil {
			return fmt.Errorf("invalid JSON: %w", err)
		}
		if source.Source != "manual" && source.Source != "smart_scan" {
			return fmt.Errorf(`source must be "manual" or "smart_scan"`)
		}
		return nil
	},
}

type IngestBatch struct {
	DeviceID string        `json:"deviceId" required:"true" format:"uuid"`
	SentAt   time.Time     `json:"sentAt"   required:"true"`
	Events   []IngestEvent `json:"events"   required:"true" minItems:"1" maxItems:"100"`
}

type IngestEvent struct {
	Name      string          `json:"name"      required:"true"`
	Props     json.RawMessage `json:"props"`
	TrackedAt int64           `json:"trackedAt" required:"true" minimum:"1"`
}

func (batch *IngestBatch) Resolve(ctx huma.Context, prefix *huma.PathBuffer) []error {
	var errs []error
	for i, event := range batch.Events {
		validator, ok := EventValidators[event.Name]
		if !ok {
			errs = append(errs, &huma.ErrorDetail{
				Location: fmt.Sprintf("%s[%d].name", prefix.With("events"), i),
				Message:  fmt.Sprintf("unknown event %q", event.Name),
				Value:    event.Name,
			})
			continue
		}
		if validator != nil {
			if err := validator(event.Props); err != nil {
				errs = append(errs, &huma.ErrorDetail{
					Location: fmt.Sprintf("%s[%d].props", prefix.With("events"), i),
					Message:  err.Error(),
					Value:    string(event.Props),
				})
			}
		}
	}
	return errs
}

func (event IngestEvent) TrackedTime() time.Time {
	return time.UnixMilli(event.TrackedAt).UTC()
}

var ErrNotIngestEvent = errors.New("not an IngestEvent")
