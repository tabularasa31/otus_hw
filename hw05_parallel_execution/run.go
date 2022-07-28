package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var wg sync.WaitGroup
	var er int64

	queue := make(chan Task, len(tasks))
	for _, task := range tasks {
		queue <- task
	}
	close(queue)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				if task, ok := <-queue; ok && atomic.LoadInt64(&er) < int64(m) {
					if res := task(); res != nil {
						atomic.AddInt64(&er, 1)
					}
				} else {
					return
				}
			}
		}()
	}

	wg.Wait()

	if int(er) >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
