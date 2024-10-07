package tqwp

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// WorkerPool manages the task processing by multiple workers.
// It tracks task success, failure, and processing time.
type WorkerPool struct {

	// ProcessedTasks holds the count of tasks processed by workers.
	ProcessedTasks uint32

	// TaskSuccess holds the count of successfully completed tasks.
	TaskSuccess uint32

	// TaskFailure holds the count of tasks that failed even after retries.
	TaskFailure uint32

	// CompletedIn tracks the time taken for processing tasks.
	// It is available only after the Stop function is called.
	CompletedIn time.Duration

	numOfWorkers uint
	queue        *TaskQueue
	wg           *sync.WaitGroup
	taskWg       *sync.WaitGroup
	maxRetries   uint
	startTime    time.Time
}

// WorkerPoolConfig holds configuration parameters for WorkerPool.
type WorkerPoolConfig struct {
	// NumOfWorkers specifies the number of workers in the pool.
	NumOfWorkers uint

	// MaxRetries specifies the maximum retry attempts for failed tasks.
	MaxRetries uint
}

// DefaultWorkerPoolConfig will give a default configuration of WorkerPool
// with 10 workers and 3 maximum retries
func DefaultWorkerPoolConfig() *WorkerPoolConfig {
	return &WorkerPoolConfig{
		NumOfWorkers: 10,
		MaxRetries: 3,
	}
}

var logger = newCustomLogger()

// New initializes and returns a new WorkerPool instance with the given configuration.
// It sets up the task queue and worker count based on the config.
func New(cfg *WorkerPoolConfig) *WorkerPool {
	var wg, taskWg sync.WaitGroup

	taskQ := NewTaskQueue(cfg.NumOfWorkers)

	return &WorkerPool{
		queue:        taskQ,
		numOfWorkers: cfg.NumOfWorkers,
		wg:           &wg,
		taskWg:       &taskWg,
		maxRetries:   cfg.MaxRetries,
	}
}

// EnqueueTask adds a task to the queue for processing and increments the task wait group counter.
func (wp *WorkerPool) EnqueueTask(task Task) {
	wp.queue.Enqueue(task)
	wp.taskWg.Add(1)
}

// Start begins the task processing by creating worker goroutines.
// It also records the start time for tracking the task completion duration.
func (wp *WorkerPool) Start() {
	logger.Info("Started WorkerPool")
	wp.startTime = time.Now()
	for i := 1; i <= int(wp.numOfWorkers); i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// Stop gracefully stops the WorkerPool by waiting for all tasks to complete.
// It closes the task queue and calculates the total time taken for processing.
func (wp *WorkerPool) Stop() {
	wp.taskWg.Wait()
	close(wp.queue.Tasks)
	wp.wg.Wait()
	wp.CompletedIn = time.Since(wp.startTime)
}

// Summary logs the statistics of the worker pool execution, including
// the number of processed tasks, successes, failures, and total time taken.
func (wp *WorkerPool) Summary() {
	fmt.Println("-------------------------------------------------------------------------------")
	msg := fmt.Sprintf(
		"\n- Processed %d Tasks \n- Worker Count %d\n- %d Success \n- %d Failed \n- Completed in %v",
		wp.ProcessedTasks,
		wp.numOfWorkers,
		wp.TaskSuccess,
		wp.TaskFailure,
		wp.CompletedIn,
	)
	logger.CustomTag("[SUMMARY] ", msg)
}

// worker is the main loop for each worker that pulls tasks from the queue
// and processes them.
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for task := range wp.queue.Tasks {
		go wp.handleTask(id, task)
	}
}

// handleTask processes a single task, handling retries if the task implements
// the retryableTask interface. It logs success, retries, or final failure after
// exhausting retry attempts.
func (wp *WorkerPool) handleTask(id int, task Task) {
	defer wp.taskWg.Done()

	for {
		err := task.Process()
		if err == nil {
			atomic.AddUint32(&wp.TaskSuccess, 1)
			atomic.AddUint32(&wp.ProcessedTasks, 1)
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

			atomic.AddUint32(&wp.TaskFailure, 1)
			atomic.AddUint32(&wp.ProcessedTasks, 1)

			msg := fmt.Sprintf(
				"Worker %d gave up after %d retries: %s",
				id,
				wp.maxRetries,
				err.Error(),
			)
			logger.Error(msg)
			return
		}

		atomic.AddUint32(&wp.ProcessedTasks, 1)
		atomic.AddUint32(&wp.TaskFailure, 1)
		msg := fmt.Sprintf(
			"Worker %d Failed to parse task: %v",
			id,
			task,
		)
		logger.Error(msg)
	}
}
