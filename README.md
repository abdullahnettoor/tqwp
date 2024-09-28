# Task Queue Worker Pool (tqwp)

```
████████╗ ██████╗ ██╗    ██╗██████╗ 
╚══██╔══╝██╔═══██╗██║    ██║██╔══██╗
   ██║   ██║   ██║██║ █╗ ██║██████╔╝
   ██║   ██║▄▄ ██║██║███╗██║██╔═══╝ 
   ██║   ╚██████╔╝╚███╔███╔╝██║     
   ╚═╝    ╚══▀▀═╝  ╚══╝╚══╝ ╚═╝        
```             

`tqwp` is a Golang package designed to help you manage task processing with a worker pool. It provides an easy-to-use API for enqueuing tasks, processing them concurrently with workers, and retrying failed tasks with configurable retry logic.

## Features

- Create custom tasks by implementing the `Task` interface.
- Set up worker pools to process tasks concurrently.
- Configurable retry mechanism for failed tasks.
- Simple logging of task processing, retries, and failures.

## Installation

To install the `tqwp` package, use:

```bash
go get github.com/abdullahnettoor/tqwp
```

## Usage

### 1. Define a Custom Task

To create a custom task, implement the `Task` interface and embed `TaskModel` for retry management.

```go
package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/abdullahnettoor/tqwp"
)

type CustomTask struct {
	tqwp.TaskModel
	Id   uint
	Data int
}

// Implement the Process method to define task behavior.
func (t *CustomTask) Process() error {
	num, isInt := t.Data

	divisor := rand.Intn(2) // Simulating error by making divisor 0
	if divisor == 0 {
		return fmt.Errorf("division by zero") 
	}
	t.Data = num / divisor
	time.Sleep(time.Millisecond * 10) // Simulate processing time
	return nil
}
```

### 2. Use the Worker Pool in the Main Function

In your main function, set up the worker pool, enqueue tasks, and start processing.

```go
package main

import (
	"math/rand"

	"github.com/abdullahnettoor/tqwp"
)

func main() {

	// Define the number of workers, tasks, and maximum retries.
	numOfWorkers := 100
	numTasks := 1000
	maxRetries := 3

	// Create a TaskQueue with capacity for numTasks.
	taskQ := tqwp.NewTaskQueue(numTasks)
	for i := 1; i <= numTasks; i++ {
		t := CustomTask{
			Id:        uint(i),
			Data:      rand.Intn(1000),
			TaskModel: tqwp.TaskModel{},
		}
		taskQ.Enqueue(&t)
	}

	// Set up the worker pool and start processing tasks.
	wp := tqwp.NewWorkerPool(taskQ, numTasks, numOfWorkers, maxRetries)
	defer wp.Summary()
	defer wp.Stop()

	wp.Start()
}
```

### 3. Run the Example

Compile and run your application:

```bash
go run main.go
```

This will create a worker pool with 100 workers to process 1000 tasks concurrently. Failed tasks will be retried up to 3 times. A summary of the processing will be printed at the end.

## API Overview

### Task Interface

The `Task` interface defines the behavior for each task:

```go
type Task interface {
	Process() error
}
```

### Retryable Task

Tasks that embed `TaskModel` automatically gain retry capabilities.

### Worker Pool

The `WorkerPool` manages task processing across multiple workers:

- `NewWorkerPool(taskQ *TaskQueue, numOfTasks int, numOfWorkers int, maxRetries int)`: Creates a new worker pool.
- `Start()`: Starts the worker pool, distributing tasks to workers.
- `Stop()`: Stops the worker pool and waits for all tasks to be processed.
- `Summary()`: Prints a summary of the processing.

### Task Queue

The `TaskQueue` handles enqueuing tasks and distributing them to workers:

- `NewTaskQueue(size int)`: Creates a new task queue with a specific size.
- `Enqueue(task Task)`: Adds a task to the queue.


## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

---
<br/>
<br/>
<br/>

## TODO
- Implement Custom error message for each task
- Coupling the configuration of worker pool in to a struct 
- Implement default configuration for worker pool


## Contributors
<a href="https://github.com/abdullahnettoor">
    <img src="https://github.com/abdullahnettoor.png" style="border-radius: 50%; alt="Abdullah Nettoor" width="60" height="60"/>
</a>