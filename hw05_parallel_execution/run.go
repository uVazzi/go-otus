package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

var ErrNotGoroutines = errors.New("not goroutines")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var (
		errorCounter atomic.Int64
		taskIndex    int
	)
	lenTasks := len(tasks)

	if lenTasks == 0 {
		return nil
	}
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}
	if n <= 0 {
		return ErrNotGoroutines
	}

	wg := sync.WaitGroup{}
	mutex := new(sync.Mutex)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				mutex.Lock()
				currentTaskIndex := taskIndex
				taskIndex++
				mutex.Unlock()

				if currentTaskIndex >= lenTasks || errorCounter.Load() >= int64(m) {
					return
				}
				if err := tasks[currentTaskIndex](); err != nil {
					errorCounter.Add(1)
				}
			}
		}()
	}
	wg.Wait()

	if errorCounter.Load() >= int64(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
