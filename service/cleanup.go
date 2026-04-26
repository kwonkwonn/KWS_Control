package service

// cleanupChain will execute registered cleanup functions in reverse order when run() is called.
// cleanup chain will register cleanup functions for each step of the process,
// so that if any step fails, all previous steps can be undone to maintain system consistency.
type cleanupChain struct {
	steps []func()
}

func (c *cleanupChain) push(fn func()) {
	c.steps = append(c.steps, fn)
}

func (c *cleanupChain) run() {
	for i := len(c.steps) - 1; i >= 0; i-- {
		c.steps[i]()
	}
	//for Idempotency
	c.steps = nil
}
