package tqwp

// Task is an interface which will be passed
// through the channel in Taskqueue
type Task interface {
	// Process function will be called in workers to process the task.
	// This functions must be implemented when creating custom tasks.
	Process() error
}

type RetryableTask interface {
	Task
	Retry(maxRetries int) bool
	GetRetry() int
}

type TaskModel struct {
	retries int
}

func (tm *TaskModel) Retry(maxRetries int) bool {
	if tm.retries < maxRetries {
		tm.retries++
		return true
	}
	return false
}

func (tm *TaskModel) GetRetry() int {
	return tm.retries
}
