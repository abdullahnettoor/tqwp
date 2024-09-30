package tqwp

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type WorkerPool struct {
	numOfWorkers   int
	queue          *TaskQueue
	wg             *sync.WaitGroup
	taskWg         *sync.WaitGroup
	maxRetries     int
	processedTasks int32
	TaskSuccess    int32
	TaskFailure    int32
	startTime      time.Time
	CompletedIn    time.Duration
}

var logger *customLogger

func New(numOfWorkers int, maxRetries int) *WorkerPool {
	var wg, taskWg sync.WaitGroup

	taskQ := NewTaskQueue(numOfWorkers)

	logger = newCustomLogger()
	return &WorkerPool{
		queue:        taskQ,
		numOfWorkers: numOfWorkers,
		wg:           &wg,
		taskWg:       &taskWg,
		maxRetries:   maxRetries,
	}
}

func NewWorkerPool(taskQ *TaskQueue, numOfTasks, numOfWorkers int, maxRetries int) *WorkerPool {
	var wg, taskWg sync.WaitGroup
	taskWg.Add(len(taskQ.Tasks))
	logger = newCustomLogger()
	return &WorkerPool{
		queue:        taskQ,
		numOfWorkers: numOfWorkers,
		wg:           &wg,
		taskWg:       &taskWg,
		maxRetries:   maxRetries,
	}
}

func (wp *WorkerPool) EnqueueTask(task Task) {
	wp.queue.Enqueue(task)
	wp.taskWg.Add(1)
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

	fmt.Println("-------------------------------------------------------------------------------")
	msg := fmt.Sprintf("\n- Processed %d Tasks \n- Worker Count %d\n- %d Success \n- %d Failed \n- Completed in %v",
		wp.processedTasks,
		wp.numOfWorkers,
		wp.TaskSuccess,
		wp.TaskFailure,
		wp.CompletedIn,
	)
	logger.CustomTag("[SUMMARY] ", msg)
}

func (wp *WorkerPool) worker(id int) {

	defer wp.wg.Done()

	for task := range wp.queue.Tasks {

		go wp.handleTask(id, task)
	}
}

func (wp *WorkerPool) handleTask(id int, task Task) {
	defer wp.taskWg.Done()
	fmt.Println("Q Len:", len(wp.queue.Tasks))

	for {
		err := task.Process()
		if err == nil {
			atomic.AddInt32(&wp.TaskSuccess, 1)
			atomic.AddInt32(&wp.processedTasks, 1)
			// msg := fmt.Sprintf(
			// 	"Worker %d successfully processed task %d",
			// 	id,
			// 	task,
			// )
			// logger.Success(msg)
			return
		}

		if tm, ok := task.(retryableTask); ok {

			if tm.retry(wp.maxRetries) {
				wp.queue.Enqueue(task)
				wp.taskWg.Add(1)
				msg := fmt.Sprintf(
					"Worker %d failed: %s (attempt %d)",
					id,
					err.Error(),
					tm.getRetry(),
				)
				logger.Warn(msg)
				return
			}

			atomic.AddInt32(&wp.TaskFailure, 1)
			atomic.AddInt32(&wp.processedTasks, 1)

			msg := fmt.Sprintf(
				"Worker %d gave up after %d retries: %s",
				id,
				wp.maxRetries,
				err.Error(),
			)
			logger.Error(msg)
			return
		}

		atomic.AddInt32(&wp.processedTasks, 1)
		atomic.AddInt32(&wp.TaskFailure, 1)
		msg := fmt.Sprintf(
			"Worker %d Failed to parse task: %v",
			id,
			task,
		)
		logger.Error(msg)
	}
}
