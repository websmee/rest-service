package app

import (
	"context"
	"time"

	"github.com/websmee/rest-service/domain/local"
)

const notSeenPeriod = 30 * time.Second

type ExpiredObjectsRemover struct {
	localObjectRepository local.Repository
}

func NewExpiredObjectsRemover(localObjectRepository local.Repository) *ExpiredObjectsRemover {
	return &ExpiredObjectsRemover{localObjectRepository}
}

func (r ExpiredObjectsRemover) RemoveExpired(ctx context.Context) error {
	return r.localObjectRepository.RemoveExpired(ctx, notSeenPeriod)
}
