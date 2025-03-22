//go:build !solution

package waitgroup

type WaitGroup struct {
	counter int64
	locker  chan interface{}
	wait    chan interface{}
}

// New creates WaitGroup.
func New() *WaitGroup {
	return &WaitGroup{
		0,
		make(chan interface{}, 1),
		make(chan interface{}, 1),
	}
}

func (wg *WaitGroup) Add(delta int) {
	wg.locker <- struct{}{}
	defer func() { <-wg.locker }()
	if wg.counter == 0 {
		wg.wait <- struct{}{}
	}
	wg.counter += int64(delta)
	if wg.counter == 0 {
		<-wg.wait
	} else if wg.counter < 0 {
		panic("negative WaitGroup counter")
	}
}

// Done decrements the WaitGroup counter by one.
func (wg *WaitGroup) Done() {
	wg.Add(-1)
}

// Wait blocks until the WaitGroup counter is zero.
func (wg *WaitGroup) Wait() {
	wg.wait <- struct{}{}
	<-wg.wait
}
