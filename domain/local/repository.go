package local

import (
	"context"
	"time"
)

type Repository interface {
	InsertOrUpdate(ctx context.Context, object Object) error
	RemoveExpired(ctx context.Context, notSeenPeriod time.Duration) error
}
