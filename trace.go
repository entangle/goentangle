package goentangle

import (
	"time"
	"sync"
)

const (
	nanosecondDivisor = 1000000000
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

	// Serialize.
	Serialize() interface{}
}

// Trace implementation.
type traceImpl struct {
	sync.Mutex

	// Description.
	description string

	// Start time.
	startTime time.Time

	// End time.
	endTime time.Time

	// Parent trace.
	parent *traceImpl

	// Sub-traces.
	subTraces []Trace
}

func (t *traceImpl) Description() string {
	return t.description
}

func (t *traceImpl) Begin(description string) Trace {
	return &traceImpl {
		description: description,
		startTime: time.Now().UTC(),
		parent: t,
		subTraces: make([]Trace, 0),
	}
}

func (t *traceImpl) End() {
	t.endTime = time.Now().UTC()

	if t.parent != nil {
		t.parent.Lock()
		t.parent.subTraces = append(t.parent.subTraces, t)
		t.parent.Unlock()
	}
}

func (t *traceImpl) SubTraces() []Trace {
	t.parent.Lock()
	defer t.parent.Unlock()

	traces := make([]Trace, len(t.subTraces))
	copy(traces, t.subTraces)
	return traces
}

func (t *traceImpl) Serialize() (ser interface{}) {
	serSubTraces := make([]interface{}, len(t.subTraces))

	for i, subTrace := range t.subTraces {
		serSubTraces[i] = subTrace.Serialize()
	}

	return []interface{} {
		t.description,
		t.startTime.UnixNano(),
		t.endTime.UnixNano(),
		serSubTraces,
	}
}

// New trace.
func NewTrace(description string) Trace {
	return &traceImpl {
		description: description,
		startTime: time.Now().UTC(),
		subTraces: make([]Trace, 0),
	}
}

// Deserialize trace.
func DeserializeTrace(ser interface{}) (t Trace, err error) {
	serArr, serOk := ser.([]interface{})
	if !serOk || len(serArr) < 4 {
		return nil, ErrDeserializationError
	}

	var description string
	if description, err = DeserializeString(serArr[0]); err != nil {
		return
	}

	var startTimeNano, endTimeNano int64
	if startTimeNano, err = DeserializeInt64(serArr[1]); err != nil {
		return
	}

	if endTimeNano, err = DeserializeInt64(serArr[2]); err != nil {
		return
	}

	serSubTraces, serOk := serArr[3].([]interface{})
	if !serOk {
		return nil, ErrDeserializationError
	}

	subTraces := make([]Trace, len(serSubTraces))
	for i, serTrace := range serSubTraces {
		if subTraces[i], err = DeserializeTrace(serTrace); err != nil {
			return
		}
	}

	return &traceImpl {
		description: description,
		startTime: time.Unix(startTimeNano / nanosecondDivisor, startTimeNano % nanosecondDivisor),
		endTime: time.Unix(endTimeNano / nanosecondDivisor, endTimeNano % nanosecondDivisor),
		subTraces: subTraces,
	}, nil
}
