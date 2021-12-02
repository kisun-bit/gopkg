package traceback

import (
	"github.com/kisunSea/gopkg/stream/buffer"
	"runtime"
	"sync"
)

var (
	_stacktracePool = sync.Pool{
		New: func() interface{} {
			return newProgramCounters(64)
		},
	}
)

func TakeStacktrace(skip int) string {
	buf_ := buffer.Get()
	defer buf_.Free()
	programCounters := _stacktracePool.Get().(*programCounters)
	defer _stacktracePool.Put(programCounters)

	var numFrames int
	for {
		// Skip the call to runtime.Callers and TakeStacktrace so that the
		// program counters start at the caller of TakeStacktrace.
		numFrames = runtime.Callers(skip+2, programCounters.pcs)
		if numFrames < len(programCounters.pcs) {
			break
		}
		// Don't put the too-short counter slice back into the pool; this lets
		// the pool adjust if we consistently take deep stacktraces.
		programCounters = newProgramCounters(len(programCounters.pcs) * 2)
	}

	i := 0
	frames := runtime.CallersFrames(programCounters.pcs[:numFrames])

	// Note: On the last iteration, frames.Next() returns false, with a valid
	// frame, but we ignore this frame. The last frame is a a runtime frame which
	// adds noise, since it's only either runtime.main or runtime.goexit.
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		if i != 0 {
			buf_.AppendByte('\n')
		}
		i++
		buf_.AppendString(frame.Function)
		buf_.AppendByte('\n')
		buf_.AppendByte('\t')
		buf_.AppendString(frame.File)
		buf_.AppendByte(':')
		buf_.AppendInt(int64(frame.Line))
	}

	return buf_.String()
}

type programCounters struct {
	pcs []uintptr
}

func newProgramCounters(size int) *programCounters {
	return &programCounters{make([]uintptr, size)}
}
