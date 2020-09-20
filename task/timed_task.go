package task

import (
	"errors"
	"time"

	"github.com/LK4D4/trylock"
)

// TimedTask is a timed-task implement for a task type
type TimedTask struct {
	rest int64         // rest operation times, count down if task Run return true
	m    trylock.Mutex // mutex for serialize
	d    time.Duration // duration of every task trigger
	Task               // Task interface
}

// timed task errors
var (
	ErrTryLockFailed error = errors.New("try lock failed")
	ErrNoRest        error = errors.New("no rest run count")
)

// NewTimedTask is constructor of TimedTask
func NewTimedTask(d time.Duration, doCount int64, t Task) *TimedTask {
	if doCount < -1 || d <= 0 {
		return nil
	}

	return &TimedTask{
		rest: doCount,
		m:    trylock.Mutex{},
		d:    d,
		Task: t,
	}
}

// Duration get TimedTask duration
func (t *TimedTask) Duration() time.Duration {
	return t.d
}

// Run is a wrapper of Task.Run
func (t *TimedTask) Run() error {
	if !t.m.TryLock() {
		return ErrTryLockFailed
	}
	defer t.m.Unlock()

	if t.rest == 0 {
		return ErrNoRest
	}

	err := t.Task.Run()

	if err != nil {
		return err
	}

	t.rest--
	return nil
}
