package app

import (
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

func (r ExpiredObjectsRemover) RemoveExpired() (int, error) {
	return r.localObjectRepository.RemoveExpired(notSeenPeriod)
}
