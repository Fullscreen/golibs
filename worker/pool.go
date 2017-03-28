package worker

import (
	"sync"
)

type Task interface {
	Execute()
}

type Pool struct {
	mu    sync.Mutex
	Size  int
	tasks chan Task
	kill  chan struct{}
	wg    sync.WaitGroup
}

func NewPool(size int) *Pool {
	pool := &Pool{
		kill:  make(chan struct{}),
		tasks: make(chan Task, 128),
	}
	pool.Resize(size)
	return pool
}

func (p *Pool) worker() {
	defer p.wg.Done()
	for {
		select {
		case task, ok := <-p.tasks:
			if !ok {
				return
			}
			task.Execute()
		case <-p.kill:
			return
		}
	}
}

func (p *Pool) Resize(size int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for p.Size < size {
		p.Size++
		p.wg.Add(1)
		go p.worker()
	}
	for p.Size > size {
		p.Size--
		p.kill <- struct{}{}
	}
}

func (p *Pool) Execute(task Task) {
	p.tasks <- task
}

func (p *Pool) Close() {
	close(p.tasks)
}

func (p *Pool) Wait() {
	p.wg.Wait()
}
