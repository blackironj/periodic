package periodic

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScheduler(t *testing.T) {
	var firstVal, secondVal, thirdVal int
	jobs := []struct {
		taskName    string
		interval    time.Duration
		taskFunc    interface{}
		expectedErr error
	}{
		{
			taskName:    "first",
			interval:    time.Millisecond * 500,
			taskFunc:    func() { firstVal++ },
			expectedErr: nil,
		},
		{
			taskName:    "first",
			interval:    time.Millisecond * 500,
			taskFunc:    func() { firstVal++ },
			expectedErr: ErrDuplicated,
		},
		{
			taskName:    "second",
			interval:    time.Millisecond * 200,
			taskFunc:    func() { secondVal++ },
			expectedErr: nil,
		},
		{
			taskName:    "third",
			interval:    time.Millisecond * 1000,
			taskFunc:    func() { thirdVal++ },
			expectedErr: nil,
		},
	}

	scheduler := NewScheduler()

	t.Run("Register tasks to scheduler", func(t *testing.T) {
		for _, j := range jobs {
			j := j
			newTask, err := NewTask(j.taskFunc)
			assert.NoError(t, err)

			err = scheduler.RegisterTask(j.taskName, j.interval, newTask)
			assert.Equal(t, j.expectedErr, err)
		}
	})

	scheduler.Run()
	time.Sleep(time.Millisecond * 1850)
	scheduler.Stop()

	time.Sleep(time.Millisecond * 500)
	t.Run("Check scheduler runs accurately", func(t *testing.T) {
		cases := []struct {
			expected int
			actual   int
		}{
			{
				expected: 4,
				actual:   firstVal,
			},
			{
				expected: 10,
				actual:   secondVal,
			},
			{
				expected: 2,
				actual:   thirdVal,
			},
		}
		for _, tc := range cases {
			tc := tc
			assert.Equal(t, tc.expected, tc.actual)
		}
	})
}
