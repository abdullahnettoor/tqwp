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

var Logger *customLogger

func NewWorkerPool(taskQ *TaskQueue, numOfTasks, numOfWorkers int, maxRetries int) *WorkerPool {
	var wg, taskWg sync.WaitGroup
	taskWg.Add(len(taskQ.Tasks))
	Logger = NewCustomLogger()
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

	fmt.Println("-------------------------------------------------------------------------------")
	msg := fmt.Sprintf("\n- Processed %d Tasks \n- Worker Count %d\n- %d Success \n- %d Failed \n- Completed in %v",
		wp.numOfTasks,
		wp.numOfWorkers,
		wp.numOfTasks-wp.TaskFailure,
		wp.TaskFailure,
		wp.CompletedIn,
	)
	Logger.CustomTag("[SUMMARY] ", msg)
}

func (wp *WorkerPool) worker(id int) {

	defer wp.wg.Done()

	for task := range wp.queue.Tasks {

		go wp.handleTask(id, task)
	}
}

func (wp *WorkerPool) handleTask(id int, task Task) {
	defer wp.taskWg.Done()

	for {
		err := task.Process()
		if err == nil {
			msg := fmt.Sprintf(
				"Worker %d successfully processed task %d",
				id,
				task)
			Logger.Success(msg)
			return
		}

		if tm, ok := task.(RetryableTask); ok && tm.Retry(wp.maxRetries) {
			wp.queue.Enqueue(task)
			wp.taskWg.Add(1)
			msg := fmt.Sprintf(
				"Worker %d failed on task %d: %s (attempt %d)",
				id,
				task,
				err.Error(),
				tm.GetRetry()+1)
			Logger.Warn(msg)
			return
		}

		wp.TaskFailure++
		msg := fmt.Sprintf(
			"Worker %d gave up on task %d after %d retries",
			id,
			task,
			wp.maxRetries)
		Logger.Error(msg)

	}

}
