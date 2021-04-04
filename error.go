package periodic

import "errors"

var (
	ErrNoFunction          = errors.New("only function can be registered")
	ErrNotMatchedNumParams = errors.New("the number of params is not matched")
	ErrNotRegistered       = errors.New("it is not a registered task")
	ErrDuplicated          = errors.New("task name duplicated")
)
