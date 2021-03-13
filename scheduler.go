package periodic

import (
	"sync"
	"time"
)

type TaskStatus int

const (
	Unknown TaskStatus = iota
	Running
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
	taskMap map[string]*scheduledTask
	rwMutex sync.RWMutex
}

//GetTaskStatus returns status of specific task
func (s *Scheduler) GetTaskStatus(taskName string) (TaskStatus, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.Unlock()

	if task, ok := s.taskMap[taskName]; ok {
		return task.status, nil
	}
	return Unknown, ErrNotRegistered
}
