package goroutine_pool

import (
	"context"
	"sync"
)

type task struct {
	ctx  context.Context
	f    func()
	next *task
}

// 全局的任务对象池
var taskPool = sync.Pool{New: newTask}

func (t *task) zero() {
	t.ctx = nil
	t.f = nil
	t.next = nil
}

func (t *task) Recycle() {
	t.zero()
	taskPool.Put(t)
}

func newTask() interface{} {
	return &task{}
}
