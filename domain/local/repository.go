package local

import (
	"time"
)

type Repository interface {
	InsertOrUpdate(object Object) error
	RemoveExpired(notSeenPeriod time.Duration) (int, error)
}
