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
	task        *Task
	interval    time.Duration
	immediately bool
	status      TaskStatus
}

//Scheduler struct has task informations
type Scheduler struct {
	taskMap map[string]*scheduledTask
	rwMutex sync.RWMutex
}

//RegisterTask register a task that runs periodically
func (s *Scheduler) RegisterTask(taskName string, interval time.Duration, imediately bool, taskFunc interface{}, params ...interface{}) error {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if _, ok := s.taskMap[taskName]; ok {
		return ErrDuplicated
	}

	newTask, err := NewTask(taskFunc, params)
	if err != nil {
		return err
	}

	s.taskMap[taskName] = &scheduledTask{
		task:        newTask,
		interval:    interval,
		immediately: imediately,
		status:      Stopped,
	}
	return nil
}

//GetTaskStatus returns status of specific task
func (s *Scheduler) GetTaskStatus(taskName string) (TaskStatus, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	if task, ok := s.taskMap[taskName]; ok {
		return task.status, nil
	}
	return Unknown, ErrNotRegistered
}
