# Task Queue Worker Pool 

A Golang-based concurrent task processing system with a task queue and worker pool. The system allows for processing tasks with automatic retries on failure, and custom logging is used to track the progress, successes, and failures of tasks.

## Features

- **Task Queue**: A buffered channel that holds tasks to be processed.
- **Worker Pool**: A pool of workers that processes tasks concurrently.
- **Task Retry Mechanism**: Automatically retries failed tasks based on a configurable limit.
- **Custom Logging**: Logs task progress and system status with different log levels (Info, Warning, Error, Success).
- **Error Handling**: Tasks can fail (e.g., division by zero), and the system retries until the maximum retry limit is reached.

## Project Structure

```plaintext
├── LICENSE
├── README.md
├── customlog.go
├── go.mod
├── main.go
├── task.go
├── task.md
├── taskqueue.go
└── worker.go
```

- **main.go**: The entry point of the application that initializes the task queue, worker pool, and processes tasks.
- **task.go**: Defines the `Task` struct and its `Process` method for simulating task execution.
- **taskqueue.go**: Defines the `TaskQueue` struct and its `Enqueue` method for adding tasks to the queue.
- **worker.go**: Defines the `WorkerPool` struct and manages task processing by workers.
- **customlog.go**: Implements a custom logging system with different log levels.

## How It Works

The application simulates processing tasks in a worker pool. Each task is an integer division operation that may fail (e.g., division by zero). Failed tasks are retried a limited number of times, and all processing results (successes, failures, retries) are logged.

1. **Task Creation**: The main application generates a fixed number of tasks and adds them to the task queue.
2. **Worker Pool**: The worker pool starts with a predefined number of workers that dequeue tasks and process them concurrently.
3. **Task Processing**: Each task is processed using the `Process` method, which simulates a division operation. If an error occurs, the worker retries the task until it either succeeds or reaches the maximum number of retries.
4. **Logging**: Task status (success, failure, retries) is logged using the custom logging system.

## Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/abdullahnettoor/tqwp.git
    cd taskQworkerpool
    ```

2. Install dependencies:

    The project uses Go modules, so dependencies are handled automatically. Ensure you have Go installed:

    ```bash
    go mod tidy
    ```

3. Run the project:

    ```bash
    go run main.go
    ```

## Configuration

- **Number of Workers**: Defined in the `main.go` file (default is 3 workers).
- **Number of Tasks**: Also defined in `main.go` (default is 10,000 tasks).
- **Maximum Retries**: The maximum number of retries for each task is set to 3 in the `main.go` file. This can be adjusted based on your requirements.

## Example Output

The output logs will contain messages for each task's progress:

```plaintext
[INFO]: 2024/09/12 12:00:00 Worker 1 successfully processed task 1
[WARN]: 2024/09/12 12:00:01 Worker 2 failed task 2: division by zero (attempt 1)
[WARN]: 2024/09/12 12:00:01 Worker 2 retrying task 2 (retry 2)
[ERROR]: 2024/09/12 12:00:02 Worker 2 gave up on task 2 after 3 retries
```

At the end, a summary is logged:

```plaintext
[SUMMARY]: 2024/09/12 13:45:51 
- Processed 10000 Tasks 
- 9355 Success 
- 645 Failed 
- Process took 348.657875ms
```

## Usage

To customize the worker pool or task behavior, modify the following parameters in `main.go`:

- `numOfWorkers`: Number of workers in the pool.
- `numTasks`: Total number of tasks to be processed.
- `maxRetries`: Maximum number of retries allowed for each task.

You can also adjust the task processing logic in the `Process` method within `task.go`.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

<!-- 
## Contributors

- [Abdullah Nettoor](https://github.com/abdullahNettoor)
 -->
