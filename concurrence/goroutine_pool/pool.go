package goroutine_pool

import (
	"context"
	"sync"
	"sync/atomic"
)

type Pool interface {
	// Name returns the corresponding pool name.
	Name() string
	// SetCap sets the goroutine capacity of the pool.
	SetCap(cap int32)
	// Go executes f.
	Go(f func())
	// CtxGo executes f and accepts the context.
	CtxGo(ctx context.Context, f func())
	// SetPanicHandler sets the panic handler.
	SetPanicHandler(f func(context.Context, interface{}))
}

type pool struct {
	// The name of the pool
	name string

	// capacity of the pool, the maximum number of goroutines that are actually working
	cap int32
	// Configuration information
	config *Config
	// linked list of tasks
	taskHead  *task
	taskTail  *task
	taskLock  sync.Mutex
	taskCount int32

	// Record the number of running workers
	workerCount int32

	// This method will be called when the worker panic
	panicHandler func(context.Context, interface{})
}

func NewPool(name string, cap int32, config *Config) Pool {
	p := &pool{
		name:   name,
		cap:    cap,
		config: config,
	}
	return p
}

func (p *pool) Name() string {
	return p.name
}

func (p *pool) SetCap(cap int32) {
	atomic.StoreInt32(&p.cap, cap)
}

func (p *pool) Go(f func()) {
	p.CtxGo(context.Background(), f)
}

func (p *pool) CtxGo(ctx context.Context, f func()) {
	t := taskPool.Get().(*task)
	t.ctx = ctx
	t.f = f

	p.PutTask(t)

	if p.IsTrigger() {
		p.incrWorkerCount()
		w := workerPool.Get().(*worker)
		w.pool = p
		w.run()
	}
}

func (p *pool) SetPanicHandler(f func(context.Context, interface{})) {
	p.panicHandler = f
}

func (p *pool) taskCount_() int32 {
	return atomic.LoadInt32(&p.taskCount)
}

func (p *pool) Capacity() int32 {
	return atomic.LoadInt32(&p.cap)
}

func (p *pool) WorkerCount() int32 {
	return atomic.LoadInt32(&p.workerCount)
}

func (p *pool) incrWorkerCount() {
	atomic.AddInt32(&p.workerCount, 1)
}

func (p *pool) incrTaskCount() {
	atomic.AddInt32(&p.taskCount, 1)
}

func (p *pool) decrWorkerCount() {
	atomic.AddInt32(&p.workerCount, -1)
}

func (p *pool) PutTask(task_ *task) {
	p.taskLock.Lock()
	if p.taskHead == nil {
		p.taskHead = task_
		p.taskTail = task_
	} else {
		p.taskTail.next = task_
		p.taskTail = task_
	}
	p.taskLock.Unlock()
	p.incrTaskCount()
}

func (p *pool) IsTrigger() bool {
	// 满足下述两个条件才触发任务执行:
	// 1. 任务数`taskCount`达到阈值且当前活跃的任务数小于p.cap
	// 2. 没有活跃的任务
	return (p.taskCount_() >= p.config.ScaleThreshold && p.WorkerCount() < p.Capacity()) || p.WorkerCount() == 0
}
