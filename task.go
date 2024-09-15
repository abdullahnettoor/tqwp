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

type CustomTask struct {
	Id      uint
	Data    any
	Retries int
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
