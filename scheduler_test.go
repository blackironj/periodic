package periodic

import (
	"strconv"
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

func TestScheduler_GetTaskStatus(t *testing.T) {
	scheduler := NewScheduler()

	testTask, err := NewTask(func() { /*Do nothing*/ })
	assert.NoError(t, err)

	testCases := []struct {
		taskKey        string
		run            bool
		register       bool
		expectedStatus TaskStatus
		expectedErr    error
	}{
		{
			taskKey:        "running",
			run:            true,
			register:       true,
			expectedStatus: Running,
			expectedErr:    nil,
		},
		{
			taskKey:        "stopped",
			run:            false,
			register:       true,
			expectedStatus: Stopped,
			expectedErr:    nil,
		},
		{
			taskKey:        "unregistered",
			register:       false,
			expectedStatus: Unknown,
			expectedErr:    ErrNotRegistered,
		},
	}

	for _, tc := range testCases {
		tc := tc
		if !tc.register {
			continue
		}
		err = scheduler.RegisterTask(tc.taskKey, time.Second*1, testTask)
		assert.NoError(t, err)
		if tc.run {
			scheduler.Run(tc.taskKey)
		}
	}

	for _, tc := range testCases {
		tc := tc
		status, err := scheduler.GetTaskStatus(tc.taskKey)
		assert.Equal(t, tc.expectedErr, err)
		assert.Equal(t, tc.expectedStatus, status)
	}
	scheduler.Stop()
}

func TestScheduler_GetNumOfTasks(t *testing.T) {
	NumOfTasks := 100

	scheduler := NewScheduler()

	testTask, err := NewTask(func() { /*Do nothing*/ })
	assert.NoError(t, err)

	for i := 0; i < NumOfTasks; i++ {
		err = scheduler.RegisterTask(strconv.Itoa(i), time.Second*1, testTask)
		assert.NoError(t, err)
	}

	t.Run("Get num of running tasks after running", func(t *testing.T) {
		//Run specific tasks
		for i := 0; i < NumOfTasks; {
			scheduler.Run(strconv.Itoa(i))
			i = i + 2
		}
		num := scheduler.GetNumOfTasks(Running)
		assert.Equal(t, NumOfTasks/2, num)
		num = scheduler.GetNumOfTasks(Stopped)
		assert.Equal(t, NumOfTasks/2, num)

		time.Sleep(time.Second * 1)

		//Run tasks remained
		scheduler.Run()

		num = scheduler.GetNumOfTasks(Running)
		assert.Equal(t, NumOfTasks, num)
		num = scheduler.GetNumOfTasks(Stopped)
		assert.Equal(t, 0, num)
	})

	time.Sleep(time.Millisecond * 500)

	t.Run("Get num of stopped tasks after stopping", func(t *testing.T) {
		//Stop specific tasks
		for i := 0; i < NumOfTasks; {
			scheduler.Stop(strconv.Itoa(i))
			i = i + 2
		}
		num := scheduler.GetNumOfTasks(Stopped)
		assert.Equal(t, NumOfTasks/2, num)
		num = scheduler.GetNumOfTasks(Running)
		assert.Equal(t, NumOfTasks/2, num)

		time.Sleep(time.Second * 1)
		//Stop tasks remained
		scheduler.Stop()

		num = scheduler.GetNumOfTasks(Stopped)
		assert.Equal(t, NumOfTasks, num)
		num = scheduler.GetNumOfTasks(Running)
		assert.Equal(t, 0, num)
	})
}

func wait(wg *sync.WaitGroup) chan struct{} {
	ch := make(chan struct{})
	go func() {
		wg.Wait()
		ch <- struct{}{}
	}()
	return ch
}
