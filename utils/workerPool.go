package utils

// Task defninition
type Task func(chan string)

// Worker pool definition
type WorkerPool struct {
	workerCount int
	taskQueue chan Task
}

// Create a new worker pool
func NewWorkerPool(workerCount int, queueSize int) *WorkerPool {
	return &WorkerPool{
		workerCount: workerCount,
		taskQueue: make(chan Task, queueSize),
	}
}

// Start the worker pool
func (wp *WorkerPool) Start() {
	// Create workerCount number of goroutines
	for i := 0; i < wp.workerCount; i++ {
		go wp.worker()
	}
}

func (wp *WorkerPool) worker() {
    for task := range wp.taskQueue {
        responseChan := make(chan string)
        task(responseChan)
        close(responseChan)
    }
}

func (wp *WorkerPool) Submit(t Task) {
	wp.taskQueue <- t
}
