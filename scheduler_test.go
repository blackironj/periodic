package periodic

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSchedulerRun(t *testing.T) {
	interval := time.Millisecond * 300
	intervalPlusTiny := interval + time.Millisecond*30

	wg := &sync.WaitGroup{}
	wg.Add(2)

	testTask, err := NewTask(func() { wg.Done() })
	assert.NoError(t, err)

	scheduler := NewScheduler()
	err = scheduler.RegisterTask("test", interval, testTask)
	assert.NoError(t, err)

	select {
	case <-time.After(intervalPlusTiny):
		//No error
	case <-wait(wg):
		t.Fatal()
	}
}

func wait(wg *sync.WaitGroup) chan struct{} {
	ch := make(chan struct{})
	go func() {
		wg.Wait()
		ch <- struct{}{}
	}()
	return ch
}
