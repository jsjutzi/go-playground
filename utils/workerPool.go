package utils

import (
	"context"
)

// Task defninition
type Task struct {
	Ctx    context.Context
	Func   func(context.Context) Event
	Result chan interface{}
}

// Worker pool definition
type WorkerPool struct {
	workerCount int
	TaskQueue   chan Task
}

// Create a new worker pool
func NewWorkerPool(workerCount int) *WorkerPool {
	pool := &WorkerPool{
		TaskQueue:   make(chan Task),
		workerCount: workerCount,
	}

	pool.StartWorkers()

	return pool
}

// Start the worker pool
func (wp *WorkerPool) StartWorkers() {
	// Create workerCount number of goroutines
	for i := 0; i < wp.workerCount; i++ {
		go wp.worker()
	}
}

func (wp *WorkerPool) worker() {
	for task := range wp.TaskQueue {
		// Check if context is already done before starting task
		select {
		case <-task.Ctx.Done():
			task.Result <- Event{Name: "Error", Data: "N/A", Error: task.Ctx.Err()}
		default:
			result := task.Func(task.Ctx)
			task.Result <- result
		}
	}
}

func (wp *WorkerPool) Submit(t Task) {
	wp.TaskQueue <- t
}
