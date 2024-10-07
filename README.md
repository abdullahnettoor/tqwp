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

	// Define the config with number of workers and maximum retries.
	cfg := tqwp.WorkerPoolConfig{
		MaxRetries:   3,
		NumOfWorkers: 10,
	}

	// Set up the worker pool and start processing tasks.
    wp := tqwp.New(&cfg)
	defer wp.Summary()
	defer wp.Stop()

	// Add Tasks to Queue before starting. 
    // Not Recommended.
	wp.EnqueueTask(&CustomTask{
        Id:        uint(111111),
		Data:      rand.Intn(1000),
		TaskModel: tqwp.TaskModel{},
	})

	wp.EnqueueTask(&CustomTask{
        Id:        uint(123124),
		Data:      rand.Intn(1000),
		TaskModel: tqwp.TaskModel{},
	})

	wp.Start()

    // Add Tasks to Queue.
	for i := 1; i <= 1000; i++ {
		t := CustomTask{
			Id:        uint(i),
			Data:      rand.Intn(1000),
			TaskModel: tqwp.TaskModel{},
		}
		wp.EnqueueTask(&t)
	}
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

### Worker Pool

The `WorkerPool` manages task processing across multiple workers:

```go
type WorkerPoolConfig struct {
	NumOfWorkers int
	MaxRetries   int
}
```

- `New(cfg *WorkerPoolConfig)`: Creates a new worker pool with provided configuration.
- `Start()`: Starts the worker pool, distributing tasks to workers.
- `EnqueueTask(task Task)`: Adds a task to the queue.
- `Stop()`: Stops the worker pool and waits for all tasks to be processed.
- `Summary()`: Prints a summary of the processing.


## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

---
<br/>
<br/>

## TODO
- Add interval between retries.
- Add godocs in necessary files.

## Contributors
<a href="https://github.com/abdullahnettoor">
    <img src="https://github.com/abdullahnettoor.png" style="border-radius: 50%; alt="Abdullah Nettoor" width="60" height="60"/>
</a>
