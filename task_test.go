package periodic

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTask(t *testing.T) {
	cases := []struct {
		expected       error
		taskFunc       interface{}
		taskFuncParams []interface{}
	}{
		{
			expected:       nil,
			taskFunc:       func() {},
			taskFuncParams: []interface{}{},
		},
		{
			expected:       NoFunction,
			taskFunc:       "test",
			taskFuncParams: []interface{}{},
		},
		{
			expected:       NotMatchedNumParams,
			taskFunc:       func(val1 int, val2 string) {},
			taskFuncParams: []interface{}{1, "str", 1, []string{"1"}},
		},
		{
			expected:       NotMatchedNumParams,
			taskFunc:       func() {},
			taskFuncParams: []interface{}{1},
		},
		// {
		// 	expected:       nil,
		// 	taskFunc:       func(val1 string, val2 int64) {},
		// 	taskFuncParams: []interface{}{"str", 123},
		// },
	}
	for _, tc := range cases {
		tc := tc
		_, err := NewTask(time.Second, true, "NewTask", tc.taskFunc, tc.taskFuncParams...)
		assert.Equal(t, tc.expected, err)
	}
}

func TestRunTask(t *testing.T) {
	cases := []struct {
		interval    time.Duration
		sleep       time.Duration
		expected    int32
		immediately bool
	}{
		{
			interval:    time.Millisecond * 500,
			sleep:       time.Millisecond * 750,
			expected:    2,
			immediately: true,
		},
		{
			interval:    time.Millisecond * 500,
			sleep:       time.Millisecond * 750,
			expected:    1,
			immediately: false,
		},
		{
			interval:    time.Millisecond * 100,
			sleep:       time.Millisecond * 1050,
			expected:    10,
			immediately: false,
		},
	}

	for _, tc := range cases {
		tc := tc
		var counter int32
		f := func() {
			atomic.AddInt32(&counter, 1)
		}

		task, err := NewTask(tc.interval, tc.immediately, "run-task", f)
		assert.NoError(t, err)

		task.Run()
		time.Sleep(tc.sleep)
		task.Stop()

		assert.Equal(t, tc.expected, counter)
	}
}
