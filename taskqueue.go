package tqwp

import (
	"sync"
)

type TaskQueue struct {
	Tasks chan Task
	mu    sync.Mutex
}

func NewTaskQueue(size int) *TaskQueue {
	return &TaskQueue{
		Tasks: make(chan Task, size),
	}
}

func (tq *TaskQueue) Enqueue(task Task) {
	tq.mu.Lock()
	defer tq.mu.Unlock()
	tq.Tasks <- task
}
