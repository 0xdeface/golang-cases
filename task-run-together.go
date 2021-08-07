package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/*
	Задача, написать функцию запускающую N задач паралельно,
	Выполнение задач должно быть прервано при достижении Q количества ошибок
*/
func main() {
	tasks := makeSomeTask(15)
	newRunner(tasks, 10, 15)
}

func makeSomeTask(n int) (tasks []func() error) {
	tasks = make([]func() error, n)
	for i := 0; i < n; i++ {
		chance := rand.Intn(10)
		tasks[i] = func(x int) func() error {
			return func() error {
				fmt.Printf("task %v is work\n", x)
				time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
				fmt.Printf("task %v is finish\n", x)
				if chance > 5 {
					return errors.New("error raised")
				}
				return nil
			}
		}(i)
	}
	return
}

func worker(ctx context.Context, wg *sync.WaitGroup, jobChan chan func() error, resultChan chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		case resultChan <- (<-jobChan)():
			wg.Done()
		}
	}
}

func newRunner(tasks []func() error, parallelTask int, criticalErrorCount int) {
	wg := &sync.WaitGroup{}
	wg.Add(len(tasks))
	ctx, cancel := context.WithCancel(context.TODO())

	jobChannel := make(chan func() error)
	resultChannel := make(chan error)
	done := make(chan struct{})

	defer close(jobChannel)
	defer close(resultChannel)
	defer close(done)

	// run workers
	for i := 0; i < parallelTask; i++ {
		go worker(ctx, wg, jobChannel, resultChannel)
	}
	// send jobs to workers
	go func() {
		for _, task := range tasks {
			jobChannel <- task
		}
		wg.Wait()
		done <- struct{}{}
	}()
	// wait for result
	for {
		select {
		case workError := <-resultChannel:
			fmt.Printf("work finished, %v\n", workError)
			if workError != nil {
				criticalErrorCount--
				fmt.Println("error count is ", criticalErrorCount)
			}
			if criticalErrorCount == 0 {
				fmt.Println("error count limit!")
				cancel()
				return
			}
		case <-done:
			fmt.Println("all task a finished")
			cancel()
			return
		}
	}

}
