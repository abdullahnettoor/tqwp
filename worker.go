package tqwp

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
	Log := NewCustomLogger()

	defer wp.wg.Done()

	for task := range wp.queue.Tasks {
		err := task.Process()
		if err != nil {
			msg := fmt.Sprintf("Worker %d failed task %d: %v (attempt %d)", id, task.Id, err, task.Retries+1)
			Log.Warn(msg)
			if task.Retries < wp.maxRetries {
				task.Retries++
				msg = fmt.Sprintf("Worker %d retrying task %d (retry %d)", id, task.Id, task.Retries)
				Log.Warn(msg)
				wp.taskWg.Add(1)
				wp.queue.Enqueue(task)
			} else {
				msg = fmt.Sprintf("Worker %d gave up on task %d after %d retries", id, task.Id, wp.maxRetries)
				Log.Error(msg)
				wp.TaskFailure++
			}
		} else {
			msg := fmt.Sprintf("Worker %d successfully processed task %d", id, task.Id)
			Log.Success(msg)
			wp.TaskSuccess++
		}
		wp.taskWg.Done()
	}
}
