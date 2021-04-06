package periodic

import "reflect"

type Task struct {
	taskFunction interface{}
	taskParams   []interface{}
}

//NewTask creates a new task to register with scheduler
func NewTask(taskFunc interface{}, taskFuncParams ...interface{}) (Task, error) {
	typ := reflect.TypeOf(taskFunc)
	if typ.Kind() != reflect.Func {
		return Task{}, ErrNoFunction
	}

	f := reflect.ValueOf(taskFunc)
	if len(taskFuncParams) != f.Type().NumIn() {
		return Task{}, ErrNotMatchedNumParams
	}

	return Task{
		taskFunction: taskFunc,
		taskParams:   taskFuncParams,
	}, nil
}

//GetTaskFuncParams return task function parameters
//if no parameters, it returns zero length slice
func (t *Task) GetTaskFuncParams() []reflect.Value {
	params := make([]reflect.Value, len(t.taskParams))
	for i, param := range t.taskParams {
		params[i] = reflect.ValueOf(param)
	}
	return params
}

//GetTaskFunc returns task function
/*
	f := GetTaskFunc()
	params := GetTaskFuncParams()
	f.Call(params)
*/
func (t *Task) GetTaskFunc() reflect.Value {
	return reflect.ValueOf(t.taskFunction)
}
