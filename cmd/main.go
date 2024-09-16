package main

import (
	"math/rand"

	"github.com/abdullahnettoor/tqwp"
)

func main() {

	numOfWorkers := 100
	numTasks := 1000
	maxRetries := 3

	taskQ := tqwp.NewTaskQueue(numTasks)
	for i := 1; i <= numTasks; i++ {
		t := tqwp.CustomTask{
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
