package requesto

import (
	"context"
	"sync"
)

// Task represents a single, independent request to be executed concurrently.
type Task struct {
	ID      string
	Request *Request
}

// Result encapsulates the outcome of a concurrent request, containing either a
// Response or an Error, along with the ID of the original Task.
type Result struct {
	TaskID   string
	Response *Response
	Error    error
}

// managerConfig is a private struct used to aggregate configuration for a
// ConcurrencyManager via the functional options pattern.
type managerConfig struct {
	ctx      context.Context
	poolSize int
}

// ManagerOption is a function type for configuring a ConcurrencyManager.
type ManagerOption func(*managerConfig)

// WithContext sets a master context for the entire set of concurrent tasks,
// useful for global timeouts or cancellations.
func WithContext(ctx context.Context) ManagerOption {
	return func(c *managerConfig) {
		if ctx != nil {
			c.ctx = ctx
		}
	}
}

// WithPoolSize sets the size of the concurrency pool, which determines the
// number of worker goroutines running in parallel.
func WithPoolSize(size int) ManagerOption {
	return func(c *managerConfig) {
		if size > 0 {
			c.poolSize = size
		}
	}
}

// ConcurrencyManager manages and executes a pool of concurrent HTTP requests.
type ConcurrencyManager struct {
	client *Client
	config *managerConfig
	tasks  []Task
}

// NewManager creates and returns a new ConcurrencyManager.
// It takes a configured requesto.Client and a set of options to initialize its settings.
func NewManager(client *Client, opts ...ManagerOption) *ConcurrencyManager {
	config := &managerConfig{
		ctx:      context.Background(),
		poolSize: 10,
	}

	// Apply all user-provided options.
	for _, opt := range opts {
		opt(config)
	}

	return &ConcurrencyManager{
		client: client,
		config: config,
		tasks:  make([]Task, 0),
	}
}

// AddTasks adds one or more pre-configured Tasks to the manager's queue.
func (cm *ConcurrencyManager) AddTasks(tasks ...Task) *ConcurrencyManager {
	cm.tasks = append(cm.tasks, tasks...)
	return cm
}

// AddURLs is a helper function to quickly create and add tasks from a list of URL strings.
func (cm *ConcurrencyManager) AddURLs(urls ...string) *ConcurrencyManager {
	for _, u := range urls {
		// Use NewRequestWithContext to ensure each request is associated with the manager's context.
		req := cm.client.NewRequestWithContext(cm.config.ctx).SetURL(u)

		cm.tasks = append(cm.tasks, Task{
			ID:      u,
			Request: req,
		})
	}
	return cm
}

// Run starts the worker pool, executes all queued tasks, and blocks until they
// are all completed. It returns a slice of Results. The order of the results is
// not guaranteed to match the order in which tasks were added.
func (cm *ConcurrencyManager) Run() []Result {
	var wg sync.WaitGroup
	// Use buffered channels to prevent blocking when queuing tasks.
	tasksCh := make(chan Task, len(cm.tasks))
	resultsCh := make(chan Result, len(cm.tasks))

	// Start the worker goroutines.
	for range cm.config.poolSize {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasksCh {
				// Check if the master context has been canceled; if so, stop sending new requests.
				if cm.config.ctx.Err() != nil {
					resultsCh <- Result{TaskID: task.ID, Error: cm.config.ctx.Err()}
					continue
				}

				resp, err := task.Request.Get()
				resultsCh <- Result{TaskID: task.ID, Response: resp, Error: err}
			}
		}()
	}

	// Push all tasks into the tasks channel.
	for _, task := range cm.tasks {
		tasksCh <- task
	}
	// All tasks have been sent; close the channel so workers will exit when idle.
	close(tasksCh)

	// Wait for all workers to finish.
	wg.Wait()
	close(resultsCh)

	// Collect all results from the results channel into a slice.
	results := make([]Result, 0, len(cm.tasks))
	for result := range resultsCh {
		results = append(results, result)
	}
	return results
}
