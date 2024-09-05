package main

import (
	"fmt"
	"math/rand"
	"sync"
)

var CusLogger *customLogger

func main() {

	numOfWorkers := 3
	numTasks := 15
	maxRetries := 3

	var wg, taskWg sync.WaitGroup

	CusLogger = NewCustomLogger()

	taskQ := NewTaskQueue(numTasks)
	for i := 1; i <= numTasks; i++ {
		t := Task{
			Id:      uint(i),
			Data:    rand.Intn(1000),
			Retries: 0,
		}
		taskQ.Enqueue(t)
		taskWg.Add(1)
	}
	wp := NewWorkerPool(taskQ, numOfWorkers, &wg, &taskWg, maxRetries)

	wp.Start()

	taskWg.Wait()
	wp.Stop()
	wg.Wait()

	fmt.Println("-----------------------------------------------------------------")
	msg := fmt.Sprintf("Processed %d Tasks | %d Success | %d Failed |", numTasks, wp.TaskSuccess, wp.TaskFailure)
	CusLogger.CustomeTag("SUMMARY ", msg)
}
