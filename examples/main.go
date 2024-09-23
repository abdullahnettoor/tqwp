package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/abdullahnettoor/tqwp"
)

type CustomTask struct {
	tqwp.TaskModel
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

func main() {

	numOfWorkers := 100
	numTasks := 1000
	maxRetries := 3

	taskQ := tqwp.NewTaskQueue(numTasks)
	for i := 1; i <= numTasks; i++ {
		t := CustomTask{
			Id:        uint(i),
			Data:      rand.Intn(1000),
			TaskModel: tqwp.TaskModel{},
		}
		taskQ.Enqueue(&t)
	}

	wp := tqwp.NewWorkerPool(taskQ, numTasks, numOfWorkers, maxRetries)
	defer wp.Summary()
	defer wp.Stop()

	wp.Start()
}
