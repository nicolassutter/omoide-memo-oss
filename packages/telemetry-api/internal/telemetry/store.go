package telemetry

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Store interface {
	Insert(ctx context.Context, events []Event) error
	Find(ctx context.Context, startTime time.Time, endTime time.Time) ([]Event, error)
}

type GORMStore struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *GORMStore {
	return &GORMStore{db: db}
}

func (s *GORMStore) Insert(ctx context.Context, events []Event) error {
	if len(events) == 0 {
		return nil
	}
	now := time.Now().UTC()
	for i := range events {
		if events[i].ReceivedAt.IsZero() {
			events[i].ReceivedAt = now
		}
	}
	return s.db.WithContext(ctx).Create(&events).Error
}

func (store *GORMStore) Find(ctx context.Context, startTime time.Time, endTime time.Time) ([]Event, error) {
	var events []Event
	databaseError := store.db.WithContext(ctx).
		Where("tracked_at >= ? AND tracked_at <= ?", startTime, endTime).
		Order("tracked_at DESC").
		Find(&events).Error
	return events, databaseError
}
