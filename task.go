package tqwp

// Task represents the interface that defines a unit of work.
// The Process function must be implemented by users to define custom task behavior.
type Task interface {
	// Process is the function executed by workers to process the task.
	// It should return an error if processing fails.
	Process() error
}

// RetryableTask is an interface that extends Task and manages retries for failed tasks.
// It allows tasks to be retried up to a specified maxRetries value.
type retryableTask interface {
	// Task is embedded to ensure users only need to implement the Process method for custom tasks.
	Task

	// Retry attempts to retry the task, returning true if retries are still allowed.
	retry(maxRetries uint) bool

	// GetRetry returns the current retry count for the task.
	getRetry() uint
}

// TaskModel is a base struct that users can embed in their custom tasks
// to manage retry logic by keeping track of retry attempts.
type TaskModel struct {
	retries uint
}

// Retry increments the retry count and returns true if the task
// can still be retried (i.e., the retry count is below maxRetries).
func (tm *TaskModel) retry(maxRetries uint) bool {
	if tm.retries < maxRetries {
		tm.retries++
		return true
	}
	return false
}

// GetRetry returns the number of retries that have been attempted for the task.
func (tm *TaskModel) getRetry() uint {
	return tm.retries
}
