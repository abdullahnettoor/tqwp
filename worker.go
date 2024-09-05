package main

import (
	"fmt"
	"sync"
)

type WorkerPool struct {
	workers     int
	queue       *TaskQueue
	wg          *sync.WaitGroup
	taskWg      *sync.WaitGroup
	maxRetries  int
	TaskSuccess int
	TaskFailure int
}

func NewWorkerPool(taskQ *TaskQueue, workers int, wg, taskWg *sync.WaitGroup, maxRetries int) *WorkerPool {
	return &WorkerPool{
		queue:      taskQ,
		workers:    workers,
		wg:         wg,
		taskWg:     taskWg,
		maxRetries: maxRetries,
	}
}

func (wp *WorkerPool) Start() {
	for i := 1; i <= wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

func (wp *WorkerPool) Stop() {
	close(wp.queue.Tasks)
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for task := range wp.queue.Tasks {
		for attempt := 0; attempt <= wp.maxRetries; attempt++ {
			err := task.Process()
			if err != nil {
				msg := fmt.Sprintf("Worker %d failed task %d: %v (attempt %d)", id, task.Id, err, attempt+1)
				CusLogger.Warn(msg)
				if attempt < wp.maxRetries {
					task.Retries = attempt + 1
					msg = fmt.Sprintf("Worker %d retrying task %d (retry %d)", id, task.Id, task.Retries)
					CusLogger.Warn(msg)
					continue
				} else {
					msg = fmt.Sprintf("Worker %d gave up on task %d after %d retries", id, task.Id, wp.maxRetries)
					CusLogger.Error(msg)
					wp.TaskFailure++
					break
				}
			} else {
				msg := fmt.Sprintf("Worker %d successfully processed task %d", id, task.Id)
				CusLogger.Success(msg)
				wp.TaskSuccess++
				break
			}
		}
		wp.taskWg.Done()
	}
}
