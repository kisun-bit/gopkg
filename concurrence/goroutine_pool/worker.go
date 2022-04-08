package goroutine_pool

import (
	"fmt"
	"github.com/kisunSea/gopkg/logging"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

var logger = logging.GLogger()

type worker struct {
	pool *pool
}

// 全局的任务池
var workerPool = sync.Pool{New: newWorker}

func newWorker() interface{} {
	return &worker{}
}

func (w *worker) run() {
	go func() {
		for {
			var t *task
			w.pool.taskLock.Lock()
			if w.pool.taskHead != nil {
				t = w.pool.taskHead
				w.pool.taskHead = w.pool.taskHead.next
				atomic.AddInt32(&w.pool.taskCount, -1)
			}
			if t == nil {
				w.close()
				w.pool.taskLock.Unlock()
				w.Recycle()
				return
			}
			w.pool.taskLock.Unlock()

			func() {
				defer func() {
					if r := recover(); r != nil {
						msg := fmt.Sprintf("GOPOOL: panic in pool: %s: %v: %s", w.pool.name, r, debug.Stack())
						logger.ErrorF("ctx->%v, err->%v", t.ctx, msg)
						if w.pool.panicHandler != nil {
							w.pool.panicHandler(t.ctx, r)
						}
					}
				}()
				fmt.Println(w.pool.taskCount)
				t.f()
			}()
			t.Recycle()
		}
	}()
}

func (w *worker) close() {
	w.pool.decrWorkerCount()
}

func (w *worker) zero() {
	w.pool = nil
}

func (w *worker) Recycle() {
	w.zero()
	workerPool.Put(w)
}
