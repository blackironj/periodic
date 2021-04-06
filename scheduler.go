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
	task     Task
	interval time.Duration
	status   TaskStatus
	stopSig  chan struct{}
}

//Scheduler struct has task informations
type Scheduler struct {
	rwMutex sync.RWMutex

	taskMap map[string]*scheduledTask
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		taskMap: make(map[string]*scheduledTask),
	}
}

//RegisterTask register a task that runs periodically
func (s *Scheduler) RegisterTask(taskName string, interval time.Duration, task Task) error {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	if _, ok := s.taskMap[taskName]; ok {
		return ErrDuplicated
	}

	s.taskMap[taskName] = &scheduledTask{
		task:     task,
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
func (s *Scheduler) GetNumOfTasks(status ...TaskStatus) int {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	if len(status) == 0 {
		return len(s.taskMap)
	}

	count := 0
	for _, st := range s.taskMap {
		if st.status == status[0] {
			count++
		}
	}
	return count
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
		timer := time.NewTimer(interval)
		calculatedInterval := interval
		for {
			//TODO: In the case of immediate execution, it should also be handled.
			timer.Reset(calculatedInterval)
			select {
			case <-timer.C:
			case <-stopSigChan:
				releaseTimer(timer)
				return
			}
			start := time.Now()

			do()

			end := time.Now()
			elapsed := end.Sub(start)
			calculatedInterval = interval - elapsed

			if calculatedInterval < 0 {
				calculatedInterval = 0
			}
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
	st.status = Stopped
	close(st.stopSig)
}

// Reset resets a interval time of scheduled task.
// if the task is running, it will stop the task and then cahnge a interval
func (s *Scheduler) Reset(taskName string, interval time.Duration) error {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	task, ok := s.taskMap[taskName]
	if !ok {
		return ErrNotRegistered
	}

	stop(task)
	task.interval = interval

	return nil
}
