package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type counterErr struct {
	mu  sync.RWMutex
	err int
}

func (c *counterErr) Read() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.err
}

func (c *counterErr) Write() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.err++
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	countErr := counterErr{}
	wg := sync.WaitGroup{}
	wg.Add(n)
	ch := make(chan Task)

	for i := 0; i < n; i++ {
		go func() {
			for t := range ch {
				if t() != nil {
					countErr.Write()
				}
			}
			wg.Done()
		}()
	}

	for _, t := range tasks {
		if countErr.Read() >= m {
			break
		}
		ch <- t
	}

	close(ch)
	wg.Wait()

	if countErr.Read() >= m || m <= 0 {
		fmt.Println("Errors count: ", countErr.err)
		return ErrErrorsLimitExceeded
	}

	return nil
}
