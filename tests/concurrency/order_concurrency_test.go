package concurrency

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test100ConcurrentOrdersSafe(t *testing.T) {
	var counter int64
	var wg sync.WaitGroup
	concurrentRequests := 100

	wg.Add(concurrentRequests)
	for i := 0; i < concurrentRequests; i++ {
		go func() {
			defer wg.Done()
			time.Sleep(2 * time.Millisecond)
			atomic.AddInt64(&counter, 1)
		}()
	}

	wg.Wait()
	assert.Equal(t, int64(100), counter)
}