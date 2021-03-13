package periodic

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTask(t *testing.T) {
	testCases := []struct {
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
			expected:       ErrNoFunction,
			taskFunc:       "test",
			taskFuncParams: []interface{}{},
		},
		{
			expected:       ErrNotMatchedNumParams,
			taskFunc:       func(val1 int, val2 string) {},
			taskFuncParams: []interface{}{1, "str", 1, []string{"1"}},
		},
		{
			expected:       ErrNotMatchedNumParams,
			taskFunc:       func() {},
			taskFuncParams: []interface{}{1},
		},
	}

	for _, tc := range testCases {
		tc := tc
		_, err := NewTask(tc.taskFunc, tc.taskFuncParams...)
		assert.Equal(t, tc.expected, err)
	}

}

func TestRunTask(t *testing.T) {
	testCases := []struct {
		expected       interface{}
		taskFunc       interface{}
		taskFuncParams []interface{}
	}{
		{
			expected:       123,
			taskFunc:       func() int { return 123 },
			taskFuncParams: []interface{}{},
		},
		{
			expected:       20,
			taskFunc:       func(val1 int, val2 int) int { return val1 * val2 },
			taskFuncParams: []interface{}{10, 2},
		},
		{
			expected:       "string",
			taskFunc:       func(val string) string { return val },
			taskFuncParams: []interface{}{"string"},
		},
		{
			expected:       []float64{1.0, 2.0},
			taskFunc:       func() []float64 { return []float64{1.0, 2.0} },
			taskFuncParams: []interface{}{},
		},
	}

	for _, tc := range testCases {
		tc := tc

		job, _ := NewTask(tc.taskFunc, tc.taskFuncParams...)
		f := job.GetTaskFunc()
		p := job.GetTaskFuncParams()

		actual := f.Call(p)

		assert.Equal(t, reflect.ValueOf(tc.expected).Type(), actual[0].Type())
		assert.Equal(t, reflect.ValueOf(tc.expected).Interface(), actual[0].Interface())
	}
}
