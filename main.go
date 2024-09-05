package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var CusLogger *customLogger

func main() {
	start := time.Now()

	numOfWorkers := 3
	numTasks := 10000
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

	fmt.Println("-------------------------------------------------------------------------------")
	msg := fmt.Sprintf("\n- Processed %d Tasks \n- %d Success \n- %d Failed \n- Process took %v",
		numTasks,
		wp.TaskSuccess,
		wp.TaskFailure,
		time.Since(start),
	)
	CusLogger.CustomeTag("[SUMMARY]", msg)
}
