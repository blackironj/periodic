package periodic

import "errors"

var (
	NoFunction          = errors.New("only function can be registered")
	NotMatchedNumParams = errors.New("the number of params is not matched")
)
