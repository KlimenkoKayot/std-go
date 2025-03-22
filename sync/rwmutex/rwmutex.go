//go:build !solution

package rwmutex

// A RWMutex is a reader/writer mutual exclusion lock.
// The lock can be held by an arbitrary number of readers or a single writer.
// The zero value for a RWMutex is an unlocked mutex.
//
// If a goroutine holds a RWMutex for reading and another goroutine might
// call Lock, no goroutine should expect to be able to acquire a read lock
// until the initial read lock is released. In particular, this prohibits
// recursive read locking. This is to ensure that the lock eventually becomes
// available; a blocked Lock call excludes new readers from acquiring the
// lock.
type RWMutex struct {
	readers int64
	locker  chan interface{}
	read    chan interface{}
	write   chan interface{}
}

// New creates *RWMutex.
func New() *RWMutex {
	locker := make(chan interface{}, 1)
	read := make(chan interface{}, 1)
	write := make(chan interface{}, 1)
	return &RWMutex{
		0,
		locker,
		read,
		write,
	}
}

// RUnlock undoes a single RLock call;
// it does not affect other simultaneous readers.
// It is a run-time error if rw is not locked for reading
// on entry to RUnlock.
func (rw *RWMutex) RUnlock() {
	rw.locker <- struct{}{}
	rw.readers--
	if rw.readers == 0 {
		<-rw.write
	}
	<-rw.locker
}

// RLock locks rw for reading.
//
// It should not be used for recursive read locking; a blocked Lock
// call excludes new readers from acquiring the lock. See the
// documentation on the RWMutex type.
func (rw *RWMutex) RLock() {
	rw.read <- struct{}{}
	rw.locker <- struct{}{}
	rw.readers += 1
	if rw.readers == 1 {
		rw.write <- struct{}{}
	}
	<-rw.locker
	<-rw.read
}

// Lock locks rw for writing.
// If the lock is already locked for reading or writing,
// Lock blocks until the lock is available.
func (rw *RWMutex) Lock() {
	rw.read <- struct{}{}
	rw.write <- struct{}{}
}

// Unlock unlocks rw for writing. It is a run-time error if rw is
// not locked for writing on entry to Unlock.
//
// As with Mutexes, a locked RWMutex is not associated with a particular
// goroutine. One goroutine may RLock (Lock) a RWMutex and then
// arrange for another goroutine to RUnlock (Unlock) it.
func (rw *RWMutex) Unlock() {
	<-rw.write
	<-rw.read
}
