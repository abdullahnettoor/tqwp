package tqwp

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type workerPool struct {
	ProcessedTasks uint32
	TaskSuccess    uint32
	TaskFailure    uint32
	CompletedIn    time.Duration

	numOfWorkers uint
	queue        *TaskQueue
	wg           *sync.WaitGroup
	taskWg       *sync.WaitGroup
	maxRetries   uint
	startTime    time.Time
}

type WorkerPoolConfig struct {
	NumOfWorkers uint
	MaxRetries   uint
}

var logger = newCustomLogger()

func New(cfg *WorkerPoolConfig) *workerPool {
	var wg, taskWg sync.WaitGroup

	taskQ := NewTaskQueue(cfg.NumOfWorkers)

	return &workerPool{
		queue:        taskQ,
		numOfWorkers: cfg.NumOfWorkers,
		wg:           &wg,
		taskWg:       &taskWg,
		maxRetries:   cfg.MaxRetries,
	}
}

func (wp *workerPool) EnqueueTask(task Task) {
	wp.queue.Enqueue(task)
	wp.taskWg.Add(1)
}

func (wp *workerPool) Start() {

	logger.Info("Started WorkerPool")
	wp.startTime = time.Now()
	for i := 1; i <= int(wp.numOfWorkers); i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

func (wp *workerPool) Stop() {
	wp.taskWg.Wait()
	close(wp.queue.Tasks)
	wp.wg.Wait()
	wp.CompletedIn = time.Since(wp.startTime)
}

func (wp *workerPool) Summary() {

	fmt.Println("-------------------------------------------------------------------------------")
	msg := fmt.Sprintf("\n- Processed %d Tasks \n- Worker Count %d\n- %d Success \n- %d Failed \n- Completed in %v",
		wp.ProcessedTasks,
		wp.numOfWorkers,
		wp.TaskSuccess,
		wp.TaskFailure,
		wp.CompletedIn,
	)
	logger.CustomTag("[SUMMARY] ", msg)
}

func (wp *workerPool) worker(id int) {

	defer wp.wg.Done()

	for task := range wp.queue.Tasks {

		go wp.handleTask(id, task)
	}
}

func (wp *workerPool) handleTask(id int, task Task) {
	defer wp.taskWg.Done()

	for {
		err := task.Process()
		if err == nil {
			atomic.AddUint32(&wp.TaskSuccess, 1)
			atomic.AddUint32(&wp.ProcessedTasks, 1)
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
