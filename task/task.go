package task

// Task is interface difination of task for scheduler to manage
type Task interface {
	Percent() (int64, error) // query the task processing percent
	Retry() (bool, error)    // query if task need retry
	Run() error              // run a task
	Kill() (bool, error)     // kill a task
}
