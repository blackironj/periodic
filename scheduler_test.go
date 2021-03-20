package periodic

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScheduler(t *testing.T) {
	var firstVal, secondVal, thirdVal int32
	testFunc := func(val *int32) {
		time.Sleep(time.Millisecond * 50)
		atomic.AddInt32(val, 1)
	}

	jobs := []struct {
		taskName    string
		interval    time.Duration
		taskFunc    interface{}
		taskParams  *int32
		expectedErr error
	}{
		{
			taskName:    "first",
			interval:    time.Millisecond * 500,
			taskFunc:    testFunc,
			taskParams:  &firstVal,
			expectedErr: nil,
		},
		{
			taskName:    "first",
			interval:    time.Millisecond * 500,
			taskFunc:    testFunc,
			taskParams:  &firstVal,
			expectedErr: ErrDuplicated,
		},
		{
			taskName:    "second",
			interval:    time.Millisecond * 1000,
			taskFunc:    testFunc,
			taskParams:  &secondVal,
			expectedErr: nil,
		},
		{
			taskName:    "third",
			interval:    time.Millisecond * 2000,
			taskFunc:    testFunc,
			taskParams:  &thirdVal,
			expectedErr: nil,
		},
	}

	scheduler := NewScheduler()

	t.Run("Register tasks to scheduler", func(t *testing.T) {
		for _, j := range jobs {
			j := j
			newTask, err := NewTask(j.taskFunc, j.taskParams)
			assert.NoError(t, err)

			err = scheduler.RegisterTask(j.taskName, j.interval, newTask)
			assert.Equal(t, j.expectedErr, err)
		}
	})

	scheduler.Run()
	time.Sleep(time.Millisecond * 2900)
	scheduler.Stop()

	time.Sleep(time.Millisecond * 100)
	t.Run("Check scheduler runs accurately", func(t *testing.T) {
		cases := []struct {
			expected int32
			actual   int32
		}{
			{
				expected: 6,
				actual:   atomic.LoadInt32(&firstVal),
			},
			{
				expected: 3,
				actual:   atomic.LoadInt32(&secondVal),
			},
			{
				expected: 2,
				actual:   atomic.LoadInt32(&thirdVal),
			},
		}

		for _, tc := range cases {
			tc := tc
			assert.Equal(t, tc.expected, tc.actual)
		}

	})
}
