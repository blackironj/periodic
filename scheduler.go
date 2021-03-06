package periodic

import (
	"sync"
	"time"
)

type TaskStatus int

const (
	Running TaskStatus = iota
	Stopped
)

type scheduledTask struct {
	task     task
	interval time.Duration
	ticker   *time.Ticker

	immediately bool
	status      TaskStatus
}

//Scheduler struct keep task informations
type Scheduler struct {
	taskList map[string]*scheduledTask
	rwMutex  sync.RWMutex
}
