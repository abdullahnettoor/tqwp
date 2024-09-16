package tqwp

import (
	"fmt"
	"time"

	"math/rand"
)

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

type CustomTask struct {
	TaskModel
	Id   uint
	Data any
}

func (t *CustomTask) Process() error {
	num, isInt := t.Data.(int)
	if !isInt {
		return fmt.Errorf("invalid type")
	}
	divisor := rand.Intn(2)
	if divisor == 0 {
		return fmt.Errorf("division by zero")
	}
	t.Data = num / divisor
	time.Sleep(time.Millisecond * 10)
	return nil
}
