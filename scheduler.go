package periodic

import (
	"sync"
)

//Scheduler struct keep task informations
type Scheduler struct {
	taskList map[string]*Task
	rwMutex  sync.RWMutex
}
