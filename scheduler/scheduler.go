package scheduler

import (
	"reflect"
	"sync"
	"time"

	"github.com/chyezh/easysche/task"
)

// Scheduler is data model for a Scheduler instance
type Scheduler struct {
	TimedTaskMap sync.Map
	notifier     chan int
}

type TimedTask struct {
	Name   string
	Ticker *time.Ticker
	*task.TimedTask
}

const (
	IndexNotifier = iota
	IndexOffSet
	SigUpdate
)

func NewScheduler() *Scheduler {
	return &Scheduler{
		notifier: make(chan int),
	}
}

// Run is the main process loop of a Scheduler
func (s *Scheduler) Run() {
	for {
		cases, timedTaskList := s.newListenCases()

		chosen, recv, _ := reflect.Select(cases)
		if chosen == IndexNotifier {
			if recv.Int() == SigUpdate {
				continue
			}
		}

		task := timedTaskList[chosen-IndexOffSet]
		go func() {
			for {
				err := task.Run()
				if err == nil {
					break
				}
				retry, err := task.Retry()
				if err != nil || !retry {
					break
				}
			}
		}()
	}
}

// Stop notify the stop signal for Scheduler
func (s *Scheduler) Stop() {
}

// SubmitTask add new Task for Scheduler to Run
func (s *Scheduler) SubmitTask(t *task.Task) {
}

// RegisterTimedTask add new timed task and notify update
func (s *Scheduler) RegisterTimedTask(name string, t *task.TimedTask) {
	newTask := TimedTask{
		Name:      name,
		Ticker:    time.NewTicker(t.Duration()),
		TimedTask: t,
	}

	s.TimedTaskMap.LoadOrStore(name, &newTask)
	s.NotifyUpdate()
}

// NotifyUpdate notify scheduler to update listen config
func (s *Scheduler) NotifyUpdate() {
	s.notifier <- SigUpdate
}

// NotifyKill notify scheduler to kill
func (s *Scheduler) NotifyKill() {
	close(s.notifier)
}

// newListenList create the new listen cases witch scheduler config
func (s *Scheduler) newListenCases() ([]reflect.SelectCase, []*TimedTask) {
	selectCases := make([]reflect.SelectCase, 0)
	timedTaskList := make([]*TimedTask, 0)

	// add update listen cases notifier
	selectCases = append(selectCases, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(s.notifier),
	})

	// add timed task to listen list
	s.TimedTaskMap.Range(func(key, value interface{}) bool {
		selectCases = append(selectCases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(value.(*TimedTask).Ticker.C),
		})
		timedTaskList = append(timedTaskList, value.(*TimedTask))
		return true
	})

	return selectCases, timedTaskList
}
