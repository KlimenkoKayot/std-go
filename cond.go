//go:build !solution

package cond

// A Locker represents an object that can be locked and unlocked.
type Locker interface {
	Lock()
	Unlock()
}

// Cond implements a condition variable, a rendezvous point
// for goroutines waiting for or announcing the occurrence
// of an event.
//
// Each Cond has an associated Locker L (often a *sync.Mutex or *sync.RWMutex),
// which must be held when changing the condition and
// when calling the Wait method.
type Cond struct {
	L      Locker
	signal chan interface{}
	wait   chan interface{}
}

// New returns a new Cond with Locker l.
func New(l Locker) *Cond {
	return &Cond{
		L:      l,
		signal: make(chan interface{}, 1),
		wait:   make(chan interface{}),
	}
}

func (c *Cond) Wait() {
	c.L.Unlock()
	c.wait <- struct{}{}
	c.L.Lock()
}

func (c *Cond) Signal() {
	select {
	case <-c.wait:
		return
	default:
		return
	}
}

func (c *Cond) Broadcast() {
	for {
		select {
		case <-c.wait:
			continue
		default:
			return
		}
	}
}
