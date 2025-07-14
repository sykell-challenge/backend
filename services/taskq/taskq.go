package taskq

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/antonmashko/taskq"
)

var (
	TaskQueue *taskq.TaskQ
	// Track running jobs for cancellation
	runningJobs = make(map[string]context.CancelFunc)
	jobsMutex   sync.RWMutex
)

// InitTaskQueue initializes the task queue
func InitTaskQueue() {
	// Create task queue with worker limit of 5
	TaskQueue = taskq.New(5)

	// Start the task queue
	if err := TaskQueue.Start(); err != nil {
		log.Printf("Failed to start task queue: %v", err)
		return
	}

	log.Println("Task queue initialized with 5 workers")
}

// ShutdownTaskQueue gracefully shuts down the task queue
func ShutdownTaskQueue() {
	if TaskQueue != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := TaskQueue.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down task queue: %v", err)
		} else {
			log.Println("Task queue shut down gracefully")
		}
	}
}

// EnqueueTask adds a new task to the queue
func EnqueueTask(ctx context.Context, task taskq.Task) (int64, error) {
	if TaskQueue != nil {
		return TaskQueue.Enqueue(ctx, task)
	}
	return 0, nil
}

// RegisterJob registers a job with its cancel function
func RegisterJob(jobID string, cancel context.CancelFunc) {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()
	runningJobs[jobID] = cancel
}

// UnregisterJob removes a job from tracking
func UnregisterJob(jobID string) {
	jobsMutex.Lock()
	defer jobsMutex.Unlock()
	delete(runningJobs, jobID)
}

// CancelJob cancels a running job by its ID
func CancelJob(jobID string) bool {
	fmt.Println("Running Jobs:", runningJobs)
	fmt.Println("Attempting to cancel job:", jobID)
	jobsMutex.RLock()
	cancel, exists := runningJobs[jobID]
	jobsMutex.RUnlock()

	if exists {
		cancel()
		return true
	}
	return false
}

// IsJobRunning checks if a job is currently running
func IsJobRunning(jobID string) bool {
	jobsMutex.RLock()
	defer jobsMutex.RUnlock()
	_, exists := runningJobs[jobID]
	return exists
}
