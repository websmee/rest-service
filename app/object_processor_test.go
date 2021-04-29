package app

import (
	"sync"
	"testing"
	"time"

	"github.com/websmee/rest-service/domain/local"
	"github.com/websmee/rest-service/domain/remote"
)

// play with these to test your specific conditions
const (
	dbDelayMs            = 10
	remoteServiceDelayMs = 50
	requestsAtOnce       = 10000
	idsPerRequest        = 100
)

func BenchmarkObjectProcessorProcess_10000r_100ids(b *testing.B) {
	processor := NewObjectProcessor(
		localObjectRepositoryMock{dbDelayMs},
		remoteObjectRepositoryMock{remoteServiceDelayMs},
	)

	for n := 0; n < b.N; n++ {
		wg := new(sync.WaitGroup)
		for g := 0; g < requestsAtOnce; g++ {
			wg.Add(1)
			go func() {
				processor.Process(generateObjectIDs())
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func generateObjectIDs() []int64 {
	ids := make([]int64, idsPerRequest)
	for i := range ids {
		ids[i] = int64(i)
	}

	return ids
}

type localObjectRepositoryMock struct {
	delay int64
}

func (r localObjectRepositoryMock) InsertOrUpdate(_ local.Object) error {
	time.Sleep(time.Duration(r.delay) * time.Millisecond)
	return nil
}

func (r localObjectRepositoryMock) RemoveExpired(_ time.Duration) (int, error) {
	time.Sleep(time.Duration(r.delay) * time.Millisecond)
	return 0, nil
}

type remoteObjectRepositoryMock struct {
	delay int64
}

func (r remoteObjectRepositoryMock) GetByID(id int64) (*remote.Object, error) {
	time.Sleep(time.Duration(r.delay) * time.Millisecond)
	return &remote.Object{
		ID:     id,
		Online: true,
	}, nil
}
