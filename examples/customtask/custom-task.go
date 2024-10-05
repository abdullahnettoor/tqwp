package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/abdullahnettoor/tqwp"
)

// CustomTask represents a custom task to test the tqwp package.
// It embeds the TaskModel from tqwp and holds an Id and Data.
type CustomTask struct {
	tqwp.TaskModel
	Id   uint
	Data any
}

// Process implements the Process method from the Task interface.
// It simulates a task that attempts to divide a number by a random divisor.
// If the divisor is zero or if the data is of an invalid type, it returns an error.
func (t *CustomTask) Process() error {
	num, isInt := t.Data.(int)
	if !isInt {
		return fmt.Errorf("invalid type")
	}
	divisor := rand.Intn(2) // Simulating error cases by making divisor 0
	if divisor == 0 {
		return fmt.Errorf("division by zero")
	}
	t.Data = num / divisor
	time.Sleep(time.Millisecond * 10) // Simulating some processing time
	return nil
}

func main() {

	// Creating a new instance of workerPool with the provided configuration.
	// This worker pool is set to use 10 workers and retry tasks up to 3 times.
	wp := tqwp.New(&tqwp.WorkerPoolConfig{MaxRetries: 3, NumOfWorkers: 10})
	defer wp.Summary()
	defer wp.Stop()

	// Enqueueing tasks before starting the worker pool.
	wp.EnqueueTask(&CustomTask{
		Id:        uint(111111),
		Data:      rand.Intn(1000),
	})
	wp.EnqueueTask(&CustomTask{
		Id:        uint(123124),
		Data:      rand.Intn(1000),
	})

	// Start the worker pool to process tasks.
	wp.Start()

	// Populate the task queue with multiple tasks for processing.
	for i := 1; i <= 1000; i++ {
		t := CustomTask{
			Id:        uint(i),
			Data:      rand.Intn(1000),
		}
		wp.EnqueueTask(&t)
	}

}
