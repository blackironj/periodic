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
	task     *Task
	interval time.Duration
	status   TaskStatus
	stopSig  chan struct{}
}

//Scheduler struct has task informations
type Scheduler struct {
	rwMutex sync.RWMutex

	taskMap map[string]*scheduledTask
}

//RegisterTask register a task that runs periodically
func (s *Scheduler) RegisterTask(taskName string, interval time.Duration, taskFunc interface{}, params ...interface{}) error {
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
		task:     newTask,
		interval: interval,
		status:   Stopped,
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

//GetNumOfTasks returns the number of registered tasks in scheduler
func (s *Scheduler) GetNumOfTasks() int {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	return len(s.taskMap)
}

//Run excutes tasks that status are "Stopped"
//If there are no parameters the scheduler runs all tasks.
//On the other hand, If there are parameters the scehduler runs specific tasks
func (s *Scheduler) Run(taskNames ...string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if len(taskNames) == 0 {
		for _, t := range s.taskMap {
			schedule(t)
		}
		return
	}

	for _, taskName := range taskNames {
		if t, ok := s.taskMap[taskName]; ok {
			schedule(t)
		}
	}
}

func schedule(st *scheduledTask) {
	if st.status == Running {
		return
	}
	taskFunc := st.task.GetTaskFunc()
	params := st.task.GetTaskFuncParams()
	f := func() {
		taskFunc.Call(params)
	}

	stopSigChan := make(chan struct{})

	go func(do func(), interval time.Duration) {
		for {
			start := time.Now()

			do()

			end := time.Now()
			elapsed := end.Sub(start)
			calculatedInterval := interval - elapsed

			if calculatedInterval < 0 {
				calculatedInterval = 0
			}

			timer := time.NewTimer(calculatedInterval)
			select {
			case <-timer.C:
			case <-stopSigChan:
				releaseTimer(timer)
				return
			}
			releaseTimer(timer)
		}
	}(f, st.interval)

	st.stopSig = stopSigChan
	st.status = Running
}

func releaseTimer(timer *time.Timer) {
	if !timer.Stop() {
		<-timer.C
	}
}

//Stop stops tasks that status are "Running"
//If there are no parameters the scheduler stops all tasks.
//On the other hand, If there are parameters the scehduler stops specific tasks
func (s *Scheduler) Stop(taskNames ...string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if len(taskNames) == 0 {
		for _, t := range s.taskMap {
			stop(t)
		}
		return
	}

	for _, taskName := range taskNames {
		if t, ok := s.taskMap[taskName]; ok {
			stop(t)
		}
	}
}

func stop(st *scheduledTask) {
	if st.status != Running {
		return
	}
	st.stopSig <- struct{}{}
	st.status = Stopped
	close(st.stopSig)
}
