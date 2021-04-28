package app

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/websmee/rest-service/domain/local"
	"github.com/websmee/rest-service/domain/remote"
)

const getObjectTimeout = 5 * time.Second

type ObjectProcessor struct {
	localObjectRepository  local.Repository
	remoteObjectRepository remote.Repository
	limitChan              chan struct{}
}

func NewObjectProcessor(
	localObjectRepository local.Repository,
	remoteObjectRepository remote.Repository,
	maxThreads int,
) *ObjectProcessor {
	return &ObjectProcessor{
		localObjectRepository:  localObjectRepository,
		remoteObjectRepository: remoteObjectRepository,
		limitChan:              make(chan struct{}, maxThreads), // this channel controls how many concurrent goroutines to spawn
	}
}

func (r ObjectProcessor) Process(ctx context.Context, objectIDs []int64) []error {
	errorChan := make(chan error)
	wg := new(sync.WaitGroup)
	for i := range objectIDs {
		wg.Add(1)
		r.limitChan <- struct{}{} // block if max goroutines spawned
		go r.processSingle(ctx, objectIDs[i], wg, errorChan)
	}

	// we wait till all the work is done
	go func() {
		wg.Wait()
		close(errorChan)
	}()

	errs := make([]error, 0, len(objectIDs))
	for err := range errorChan {
		errs = append(errs, err)
	}

	return errs
}

func (r ObjectProcessor) processSingle(ctx context.Context, id int64, wg *sync.WaitGroup, errorChan chan<- error) {
	defer processRecover(errorChan)
	defer func() {
		<-r.limitChan // free buffer for new goroutines
		wg.Done()
	}()

	// cancel fetching remote object by timeout
	ctxTimeout, cancel := context.WithTimeout(ctx, getObjectTimeout)
	defer cancel()

	object, err := r.remoteObjectRepository.GetByID(ctxTimeout, id)
	if err != nil {
		errorChan <- err
		return
	}

	if !object.Online {
		return
	}

	// save every object individually to get proper independent "LastSeen" time
	if err := r.localObjectRepository.InsertOrUpdate(ctx, local.Object{
		ID:       object.ID,
		LastSeen: time.Now(),
	}); err != nil {
		errorChan <- err
		return
	}
}

func processRecover(errorChan chan<- error) {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			errorChan <- err
		} else {
			errorChan <- errors.New(fmt.Sprint(r))
		}
	}
}
