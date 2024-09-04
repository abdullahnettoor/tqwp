package main

import (
	"fmt"
	"sync"
)

type WorkerPool struct {
	workers int
	tasks   []Task
	wg      *sync.WaitGroup
}

func NewWorkerPool(workers int, wg *sync.WaitGroup, tasks []Task) *WorkerPool {
	return &WorkerPool{
		workers: workers,
		tasks:   tasks,
		wg:      wg,
	}
}

func (wp *WorkerPool) Start() {
	for i := 1; i <= wp.workers; i++ {
		wp.wg.Add(1)
		go worker(i, wp.tasks, wp.wg)
	}

}

func worker(id int, tasks []Task, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, task := range tasks {
		for attempt := 0; attempt <= task.MaxRetries; attempt++ {
			err := task.Process()
			if err != nil {
				fmt.Printf("Worker %d failed task %d: %v (attempt %d)\n", id, task.Id, err, attempt+1)
				if attempt < task.MaxRetries {
					task.Retries = attempt + 1
					fmt.Printf("Worker %d retrying task %d (retry %d)\n", id, task.Id, task.Retries)
					continue
				} else {
					fmt.Printf("Worker %d gave up on task %d after %d retries\n", id, task.Id, task.MaxRetries)
					break
				}
			} else {
				fmt.Printf("Worker %d successfully processed task %d\n", id, task.Id)
				break
			}
		}
	}
}
