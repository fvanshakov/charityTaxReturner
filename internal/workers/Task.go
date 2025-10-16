package workers

import (
	"context"
	"errors"
	"log"
	"sync"
)

type Task interface {
	Execute() (interface{}, error)
	Id() string
}

type TaskResult struct {
	TaskId string
	Result interface{}
	Error  error
}

type WorkerPool struct {
	workersCount int
	tasks        chan Task
	results      chan TaskResult
	wg           *sync.WaitGroup
}

func NewWorkerPool(workersCount int) *WorkerPool {
	return &WorkerPool{
		workersCount: workersCount,
		tasks:        make(chan Task, workersCount),
		results:      make(chan TaskResult, workersCount),
	}
}

func (p *WorkerPool) worker(ctx context.Context) {
	defer p.wg.Done()

	for {
		select {
		case task, ok := <-p.tasks:
			if !ok {
				log.Printf("task channel closed %s\n", task.Id())
				return
			}

			result, err := task.Execute()
			p.results <- TaskResult{
				TaskId: task.Id(),
				Result: result,
				Error:  err,
			}
		case <-ctx.Done():
			log.Printf("task context done: %s\n", ctx.Err())
			return
		}
	}
}

func (p *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < p.workersCount; i++ {
		p.wg.Add(1)
		go p.worker(ctx)
	}
}

func (p *WorkerPool) Submit(task Task) error {
	select {
	case p.tasks <- task:
		return nil
	default:
		return errors.New("task channel full")
	}
}

func (p *WorkerPool) Stop() {
	close(p.tasks)
	p.wg.Wait()
	close(p.results)
}

func (p *WorkerPool) Results() <-chan TaskResult {
	return p.results
}
