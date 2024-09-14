package tqwp

import (
	"fmt"
	"sync"
	"time"
)

type WorkerPool struct {
	numOfTasks   int
	numOfWorkers int
	queue        *TaskQueue
	wg           *sync.WaitGroup
	taskWg       *sync.WaitGroup
	maxRetries   int
	TaskSuccess  int
	TaskFailure  int
	startTime    time.Time
	CompletedIn  time.Duration
}

func NewWorkerPool(taskQ *TaskQueue, numOfTasks, numOfWorkers int, maxRetries int) *WorkerPool {
	var wg, taskWg sync.WaitGroup
	taskWg.Add(len(taskQ.Tasks))
	return &WorkerPool{
		queue:        taskQ,
		numOfTasks:   numOfTasks,
		numOfWorkers: numOfWorkers,
		wg:           &wg,
		taskWg:       &taskWg,
		maxRetries:   maxRetries,
	}
}

func (wp *WorkerPool) Start() {

	wp.startTime = time.Now()
	for i := 1; i <= wp.numOfWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

func (wp *WorkerPool) Stop() {
	wp.taskWg.Wait()
	close(wp.queue.Tasks)
	wp.wg.Wait()
	wp.CompletedIn = time.Since(wp.startTime)
}

func (wp *WorkerPool) Summary() {
	Log := NewCustomLogger()

	fmt.Println("-------------------------------------------------------------------------------")
	msg := fmt.Sprintf("\n- Processed %d Tasks \n- Worker Count %d\n- %d Success \n- %d Failed \n- Completed in %v",
		wp.numOfTasks,
		wp.numOfWorkers,
		wp.TaskSuccess,
		wp.TaskFailure,
		wp.CompletedIn,
	)
	Log.CustomTag("[SUMMARY] ", msg)
}

func (wp *WorkerPool) worker(id int) {
	Log := NewCustomLogger()

	defer wp.wg.Done()

	for task := range wp.queue.Tasks {
		attempt := 0
		for attempt < wp.maxRetries {
			err := task.Process()

			if err == nil {
				msg := fmt.Sprintf("Worker %d successfully processed task %d", id, task)
				Log.Success(msg)
				wp.TaskSuccess++
				break
			}

			msg := fmt.Sprintf("Worker %d failed task %d: %v (attempt %d)", id, task, err, attempt+1)
			Log.Warn(msg)
			if attempt < wp.maxRetries {
				msg = fmt.Sprintf("Worker %d retrying task %d (retry %d)", id, task, attempt)
				attempt++
				Log.Warn(msg)
				if attempt == wp.maxRetries {
					wp.TaskFailure++
					msg = fmt.Sprintf("Worker %d gave up on task %d after %d retries", id, task, wp.maxRetries)
					Log.Error(msg)
				}
			}
		}
		wp.taskWg.Done()
	}
}
