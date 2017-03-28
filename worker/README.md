worker.Pool
===========

Elastic worker pools for concurrent applications.

Usage
=====
```go
type MyTask struct {
	Message string
}

func (t *MyTask) Execute() {
	fmt.Println(t.Message)
}

func main() {
	messages := []string{"ohai"}

	// create a new worker pool
	pool = worker.NewPool(5)
	pool.Resize(10)

	for _, msg := range messages {
		pool.Execute(&MyTask{Message: msg})
	}

	pool.Close()
	pool.Wait()
}
```
