package main

import (
	"fmt"
	"time"

	"github.com/chyezh/easysche/task"

	"github.com/chyezh/easysche/scheduler"
)

type CountTask struct {
	count int
	name  string
}

func (t *CountTask) Run() error {
	fmt.Printf("task: %s, count: %d\n", t.name, t.count)
	t.count++
	return nil
}

func (t *CountTask) Retry() (bool, error) {
	return false, nil
}

func (t *CountTask) Kill() (bool, error) {
	return true, nil
}

func (t *CountTask) Percent() (int64, error) {
	return 0, nil
}

func main() {
	s := scheduler.NewScheduler()

	go s.Run()

	s.RegisterTimedTask("count1", task.NewTimedTask(time.Duration(1*time.Second), -1, &CountTask{name: "count1"}))
	s.RegisterTimedTask("count2", task.NewTimedTask(time.Duration(2*time.Second), -1, &CountTask{name: "count2"}))
	s.RegisterTimedTask("count3", task.NewTimedTask(time.Duration(10*time.Second), -1, &CountTask{name: "count3"}))

	time.Sleep(time.Duration(100 * time.Second))
}
