package infrastructure

import (
	"context"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/pkg/errors"

	"github.com/websmee/rest-service/domain/local"
)

type localObjectRepository struct {
	db *pg.DB
}

func NewLocalObjectRepository(db *pg.DB) local.Repository {
	return &localObjectRepository{db}
}

func (r localObjectRepository) InsertOrUpdate(_ context.Context, object local.Object) error {
	_, err := r.db.Model(&object).
		OnConflict("(ID) DO UPDATE").
		Set("last_seen = EXCLUDED.last_seen").
		Insert()

	return errors.Wrap(err, "InsertOrUpdate failed")
}

func (r localObjectRepository) RemoveExpired(_ context.Context, notSeenPeriod time.Duration) error {
	_, err := r.db.Model(&local.Object{}).
		Where("last_seen < ?", time.Now().Add(-notSeenPeriod)).
		Delete()

	return errors.Wrap(err, "RemoveExpired failed")
}
