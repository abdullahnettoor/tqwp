# Task Queue Worker Pool (tqwp)

```
████████╗ ██████╗ ██╗    ██╗██████╗ 
╚══██╔══╝██╔═══██╗██║    ██║██╔══██╗
   ██║   ██║   ██║██║ █╗ ██║██████╔╝
   ██║   ██║▄▄ ██║██║███╗██║██╔═══╝ 
   ██║   ╚██████╔╝╚███╔███╔╝██║     
   ╚═╝    ╚══▀▀═╝  ╚══╝╚══╝ ╚═╝        
```             


[![Go Reference](https://pkg.go.dev/badge/github.com/abdullahnettoor/tqwp.svg)](https://pkg.go.dev/github.com/abdullahnettoor/tqwp)
[![Go Report Card](https://goreportcard.com/badge/github.com/abdullahnettoor/tqwp)](https://goreportcard.com/report/github.com/abdullahnettoor/tqwp)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

</div>

## 📖 Overview

`tqwp` is a Golang package designed to help you manage task processing with a worker pool. It provides an easy-to-use API for enqueuing tasks, processing them concurrently with workers, and retrying failed tasks with configurable retry logic.

## ✨ Features

- 🔄 Concurrent task processing with configurable worker pools
- 🔁 Built-in retry mechanism for failed tasks
- 📊 Task processing metrics and summary
- 📝 Simple logging of task processing, retries, and failures
- 🎯 Custom task implementation through interface
- 🔧 Configurable queue size and worker count

## 🚀 Installation

```bash
go get -u github.com/abdullahnettoor/tqwp
```


## 💡 Quick Start

### 1. Define Your Task

```go
type CustomTask struct {
	tqwp.TaskModel // Embed TaskModel for retry functionality
	Id             uint
	Data           any
}

func (t CustomTask) Process() error {
	// Implement your task logic here
	return nil
}
```

### 2. Create and Configure Worker Pool

```go
wp := tqwp.New(&tqwp.WorkerPoolConfig{
	NumOfWorkers: 10,  // Number of concurrent workers
	MaxRetries:   3,   // Maximum retry attempts
	QueueSize:    100, // Size of task queue buffer
})
```


### 3. Start Processing

```go
wp.Start()
defer wp.Stop()

// Enqueue tasks
wp.EnqueueTask(&CustomTask{
	Id:   1,
	Data: "example",
})

// Get processing summary at the end
defer wp.Summary()
```

## 📚 Examples

Check out our example implementations in the [examples](./examples) directory:

- [Custom Task Processing](./examples/customtask/)
- [Email Sender](./examples/emailsender/)
- [Image Downloader](./examples/imgdownloader/)
- [JSON Processor](./examples/jsonprocessor/)

## ⚙️ Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| NumOfWorkers | Number of concurrent workers | Required |
| MaxRetries | Maximum retry attempts for failed tasks | Required |
| QueueSize | Buffer size for task queue | Required |


## 📋 Requirements

- Go 1.21 or later


## 🧑🏾‍💻 API Overview

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
	QueueSize    int
}
```

- `New(cfg *WorkerPoolConfig)`: Creates a new worker pool with provided configuration.
- `Start()`: Starts the worker pool, distributing tasks to workers.
- `EnqueueTask(task Task)`: Adds a task to the queue.
- `Stop()`: Stops the worker pool and waits for all tasks to be processed.
- `Summary()`: Prints a summary of the processing.

## 📜 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please read our [Contributing Guide](.github/CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## 📊 Project Status

This project is actively maintained. For feature requests and bug reports, please open an issue.

---

<div align="center">
Made with ❤️ by <a href="https://github.com/abdullahnettoor">Abdullah Nettoor</a>
</div>