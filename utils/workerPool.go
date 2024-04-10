package utils

import (
	"fmt"
	"sync"
	"time"
)

// Task defninition
// TODO: Add more fields to the task i.e. a function to execute
type Task struct {
	ID int
}

// Way to process tasks
func (t *Task) Process() {
	fmt.Printf("Processing task with ID: %d\n", t.ID)
	time.Sleep(2 * time.Second)
}

// Worker pool definition
type WorkerPool struct {
	Tasks []Task
	Concurrency int
	tasksChan chan Task
	wg sync.WaitGroup
}

// Functions to execute worker pool
func (wp *WorkerPool) worker() {
	for task := range wp.tasksChan {
		task.Process()
		wp.wg.Done()
	}
}

func (wp *WorkerPool) Run() {
	// Initialize the tasks channel
	wp.tasksChan = make(chan Task, len(wp.Tasks))

	// Start workers
	for i := 0; i < wp.Concurrency; i++ {
		go wp.worker()
	}

	// Send tasks to the channel
	wp.wg.Add(len(wp.Tasks))
	for _, task := range wp.Tasks {
		wp.tasksChan <- task
	}

	close(wp.tasksChan)

	// Wait for all tasks to be processed
	wp.wg.Wait()
}
