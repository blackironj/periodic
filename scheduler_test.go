package periodic

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScheduler_SameTaskKeyErr(t *testing.T) {
	task, _ := NewTask(func() {})

	scheduler := NewScheduler()
	err := scheduler.RegisterTask("test", time.Second, task)
	assert.NoError(t, err)

	err = scheduler.RegisterTask("test", time.Second, task)
	assert.Equal(t, ErrDuplicated, err)
}

func TestScheduler_Run(t *testing.T) {
	interval := time.Millisecond * 300
	waitingInterval := interval + time.Millisecond*30

	wg := sync.WaitGroup{}
	wg.Add(1)

	testTask, err := NewTask(func() { wg.Done() })
	assert.NoError(t, err)

	scheduler := NewScheduler()
	err = scheduler.RegisterTask("test", interval, testTask)
	assert.NoError(t, err)

	scheduler.Run()

	select {
	case <-time.After(waitingInterval):
		t.Fatal()
	case <-wait(&wg):
		//No error
		scheduler.Stop()
	}
}

func TestScheduler_RunAndStop(t *testing.T) {
	interval := time.Millisecond * 500
	waitingInterval := time.Millisecond * 550

	wg := &sync.WaitGroup{}
	wg.Add(1)
	testTask, err := NewTask(func() { wg.Done() })
	assert.NoError(t, err)

	scheduler := NewScheduler()
	err = scheduler.RegisterTask("test", interval, testTask)
	assert.NoError(t, err)

	scheduler.Run()
	time.Sleep(time.Millisecond * 200)
	scheduler.Stop()

	select {
	case <-time.After(waitingInterval):
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
