package goentangle

import (
	"time"
	"sync"
)

// Trace.
type Trace interface {
	// Description.
	Description() string

	// Begin tracing.
	//
	// Returns a sub-trace.
	Begin(description string) Trace

	// End the trace.
	End()

	// Sub-traces.
	SubTraces() []Trace
}

// Trace implementation.
type trace struct {
	sync.Mutex

	// Description.
	description string

	// Start time.
	startTime time.Time

	// End time.
	endTime time.Time

	// Parent trace.
	parent *trace

	// Sub-traces.
	subTraces []Trace
}

func (t *trace) Description() string {
	return t.description
}

func (t *trace) Begin(description string) Trace {
	return &trace {
		description: description,
		startTime: time.Now().UTC(),
		parent: t,
		subTraces: make([]Trace, 0),
	}
}

func (t *trace) End() {
	t.endTime = time.Now().UTC()

	if t.parent != nil {
		t.parent.Lock()
		t.parent.subTraces = append(t.parent.subTraces, t)
		t.parent.Unlock()
	}
}

func (t *trace) SubTraces() []Trace {
	t.parent.Lock()
	defer t.parent.Unlock()

	traces := make([]Trace, len(t.subTraces))
	copy(traces, t.subTraces)
	return traces
}

// New trace.
func NewTrace(description string) Trace {
	return &trace {
		description: description,
		startTime: time.Now().UTC(),
		subTraces: make([]Trace, 0),
	}
}
