package app

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/websmee/rest-service/domain/local"
	"github.com/websmee/rest-service/domain/remote"
)

type ObjectProcessor struct {
	localObjectRepository  local.Repository
	remoteObjectRepository remote.Repository
}

func NewObjectProcessor(
	localObjectRepository local.Repository,
	remoteObjectRepository remote.Repository,
) *ObjectProcessor {
	return &ObjectProcessor{
		localObjectRepository:  localObjectRepository,
		remoteObjectRepository: remoteObjectRepository,
	}
}

func (r ObjectProcessor) Process(objectIDs []int64) []error {
	errorChan := make(chan error)
	wg := new(sync.WaitGroup)
	for i := range objectIDs {
		wg.Add(1)
		// run all processing concurrently for each ID
		go r.processSingle(objectIDs[i], wg, errorChan)
	}

	// we wait till all the work is done
	go func() {
		wg.Wait()
		close(errorChan)
	}()

	// collect all the errors
	errs := make([]error, 0, len(objectIDs))
	for err := range errorChan {
		errs = append(errs, err)
	}

	return errs
}

func (r ObjectProcessor) processSingle(id int64, wg *sync.WaitGroup, errorChan chan<- error) {
	defer processRecover(errorChan)
	defer wg.Done()

	// get remote object
	object, err := r.remoteObjectRepository.GetByID(id)
	if err != nil {
		errorChan <- err
		return
	}

	// skip offline
	if !object.Online {
		return
	}

	// save every object individually to get proper independent "LastSeen" time
	if err := r.localObjectRepository.InsertOrUpdate(local.Object{
		ID:       object.ID,
		LastSeen: time.Now(),
	}); err != nil {
		errorChan <- err
		return
	}
}

func processRecover(errorChan chan<- error) {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok { // if we got proper error
			errorChan <- err
		} else { // if we got anything else
			errorChan <- errors.New(fmt.Sprint(r))
		}
	}
}
