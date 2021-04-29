package infrastructure

import (
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

func (r localObjectRepository) InsertOrUpdate(object local.Object) error {
	_, err := r.db.Model(&object).
		OnConflict("(id) DO UPDATE").
		Set("last_seen = EXCLUDED.last_seen").
		Insert()

	return errors.Wrap(err, "InsertOrUpdate failed")
}

func (r localObjectRepository) RemoveExpired(notSeenPeriod time.Duration) (int, error) {
	res, err := r.db.Model(&local.Object{}).
		Where("last_seen < ?", time.Now().Add(-notSeenPeriod)).
		Delete()

	if err != nil {
		return 0, errors.Wrap(err, "RemoveExpired failed")
	}

	return res.RowsAffected(), nil
}
