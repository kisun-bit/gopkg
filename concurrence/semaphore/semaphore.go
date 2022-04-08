package semaphore

import (
	"container/list"
	"context"
	"sync"
)

// Weighted provides a way to bound concurrent access to a resource.
// The callers can request access with a given weight.
type Weighted struct {
	mu      sync.Mutex
	size    int64
	cur     int64
	waiters list.List // FIFO
}

type waiter struct {
	n     int64
	ready chan<- struct{} // Closed when semaphore acquired.
}

// NewWeighted creates a new weighted semaphore with a given maximum weight for concurrent access.
func NewWeighted(n int64) *Weighted {
	w := &Weighted{size: n}
	return w
}

// Acquire acquires the semaphore with a weight of n, blocking until resources
// are available or ctx is done. On success, returns nil. On failure, returns
// ctx.Err() and leaves the semaphore unchanged.
//
// If ctx is already done, Acquire may still succeed without blocking.
func (s *Weighted) Acquire(ctx context.Context, n int64) error {
	s.mu.Lock()
	if s.size-s.cur >= n && s.waiters.Len() == 0 {
		s.cur += n
		s.mu.Unlock()
		return nil
	}

	if n > s.size {
		// exceeding resources
		s.mu.Unlock()
		<-ctx.Done()
		return ctx.Err()
	}

	// the resource number is less than `n` or `waiters` is not None.
	// blocking current task and adding task to `waiters`
	ready := make(chan struct{})
	w := waiter{n: n, ready: ready}
	e := s.waiters.PushBack(w)
	s.mu.Unlock()

	// polling to wake up executable task in `waiters`
	select {
	case <-ctx.Done():
		// context has done
		err := ctx.Err()
		s.mu.Lock()
		select {
		case <-ready:
			// (highly unlikely) wake up a split second after context ends, just acquired the semaphore
			err = nil
		default:
			// (high probability) `e` is the first element
			isFront := s.waiters.Front() == e
			s.waiters.Remove(e)
			// trying to wake up the task after `e` TODO
			if isFront && s.size > s.cur {
				s.notifyWaiters()
			}
		}
		s.mu.Unlock()
		return err
	// ready channel has closed
	case <-ready:
		return nil
	}
}

// TryAcquire acquires the semaphore with a weight of n without blocking.
// On success, returns true. On failure, returns false and leaves the semaphore unchanged.
func (s *Weighted) TryAcquire(n int64) bool {
	s.mu.Lock()
	success := s.size-s.cur >= n && s.waiters.Len() == 0
	if success {
		s.cur += n
	}
	s.mu.Unlock()
	return success
}

// Release releases the semaphore with a weight of n.
func (s *Weighted) Release(n int64) {
	s.mu.Lock()
	s.cur -= n
	if s.cur < 0 {
		s.mu.Unlock()
		panic("semaphore: released more than held")
	}
	s.notifyWaiters()
	s.mu.Unlock()
}

func (s *Weighted) notifyWaiters() {
	for {
		next := s.waiters.Front()
		if next == nil {
			break // No more waiters blocked.
		}

		w := next.Value.(waiter)
		if s.size-s.cur < w.n {
			// Not enough tokens for the next waiter.  We could keep going (to try to
			// find a waiter with a smaller request), but under load that could cause
			// starvation for large requests; instead, we leave all remaining waiters
			// blocked.
			//
			// Consider a semaphore used as a read-write lock, with N tokens, N
			// readers, and one writer.  Each reader can Acquire(1) to obtain a read
			// lock.  The writer can Acquire(N) to obtain a write lock, excluding all
			// of the readers.  If we allow the readers to jump ahead in the queue,
			// the writer will starve — there is always one token available for every
			// reader.
			break
		}

		s.cur += w.n
		s.waiters.Remove(next)
		close(w.ready)
	}
}
