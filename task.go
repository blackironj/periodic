package periodic

import (
	"reflect"
	"time"
)

type TaskStatus int

const (
	Running TaskStatus = iota
	Stopped
)

type task struct {
	taskFunction interface{}
	taskParams   []interface{}
	interval     time.Duration
	ticker       *time.Ticker

	immediately bool
	status      TaskStatus
}

//NewTask creates a new task to register with scheduler
func NewTask(interval time.Duration, immediately bool, taskName string, taskFunc interface{}, taskFuncParams ...interface{}) (*task, error) {
	typ := reflect.TypeOf(taskFunc)
	if typ.Kind() != reflect.Func {
		return nil, NoFunction
	}

	f := reflect.ValueOf(taskFunc)
	if len(taskFuncParams) != f.Type().NumIn() {
		return nil, NotMatchedNumParams
	}

	return &task{
		taskFunction: taskFunc,
		taskParams:   taskFuncParams,
		interval:     interval,
		immediately:  immediately,
		status:       Stopped,
		ticker:       time.NewTicker(interval),
	}, nil
}

//Run runs task periodically
func (t *task) Run() {
	f := reflect.ValueOf(t.taskFunction)
	in := make([]reflect.Value, len(t.taskParams))

	for i, param := range t.taskParams {
		in[i] = reflect.ValueOf(param)
	}

	t.ticker = time.NewTicker(t.interval)
	t.status = Running
	go func() {
		if t.immediately {
			for ; true; <-t.ticker.C {
				f.Call(in)
			}
		} else {
			for range t.ticker.C {
				f.Call(in)
			}
		}
	}()
}

//Stop stops doing task
func (t *task) Stop() {
	t.ticker.Stop()
	t.status = Stopped
}
