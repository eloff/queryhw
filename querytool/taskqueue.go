package querytool

import (
	"sync/atomic"
)

// TaskQueue represents an immutable collection of QueryTask structs
// which can be consumed one at a time by calling Get().
type TaskQueue struct {
	tasks        []QueryTask
	currentIndex uint64
}

// NewTaskQueue constructs a new TaskQueue from the slice of QueryTasks.
// Safety: the caller must not keep a reference to tasks after calling this function.
func NewTaskQueue(tasks []QueryTask) *TaskQueue {
	return &TaskQueue{
		tasks:        tasks,
		currentIndex: 0,
	}
}

// Get returns the next available QueryTask in the queue.
// Returns nil if all tasks have been consumed.
// Safety: Get is safe to call from concurrent goroutines
func (queue *TaskQueue) Get() *QueryTask {
	// Get the next index into the tasks slice atomically.
	// i will be unique to each calling goroutine,
	// so the caller has exclusive access to the task at that index.

	// A Mutex would more commonly be used, or maybe a buffered channel,
	// but this makes the algorithm wait-free. It will perform strictly better.
	// That is unlikely to matter for this program.
	// But it's no more complex to use the more efficient solution.
	// Plus it lets me show off a little.
	i := atomic.AddUint64(&queue.currentIndex, 1) - 1

	if i < uint64(len(queue.tasks)) {
		return &queue.tasks[i]
	}

	return nil
}

// Len returns the total number of tasks in the queue (consumed or not.)
func (queue *TaskQueue) Len() int {
	return len(queue.tasks)
}
