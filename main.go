package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
)

func main() {

	numOfWorkers := 3
	numTasks := 5
	maxRetries := 3

	tasks := make([]Task, numTasks)
	for i := 0; i < numTasks; i++ {
		tasks[i] = Task{
			Id:         uint(i + 1),
			Data:       rand.Intn(1000),
			Retries:    0,
			MaxRetries: maxRetries,
		}
	}

	var wg sync.WaitGroup
	wp := NewWorkerPool(numOfWorkers, &wg, tasks)
	wp.Start()

	wg.Wait()

	fmt.Println("-----------------------------")
	log.Printf("Processed %d Tasks\n", numTasks)
}
