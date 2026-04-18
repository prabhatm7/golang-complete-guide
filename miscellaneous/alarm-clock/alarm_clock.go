package main

import (
	"container/heap"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// ErrStopped is returned when an operation is attempted on a stopped clock.
var ErrStopped = errors.New("alarmclock: clock is stopped")

// Alarm holds the name, scheduled fire time, and internal heap index.
type Alarm struct {
	Name   string
	FireAt time.Time
	index  int // Position in the min-heap; managed by heap.
}

// alarmHeap implements heap.Interface as a min-heap ordered by FireAt.
// Earliest alarm sits at index 0 → O(1) peek, O(log n) push/pop/fix.
type alarmHeap []*Alarm

func (h alarmHeap) Len() int           { return len(h) }
func (h alarmHeap) Less(i, j int) bool { return h[i].FireAt.Before(h[j].FireAt) }
func (h alarmHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *alarmHeap) Push(x any) {
	a := x.(*Alarm)
	a.index = len(*h)
	*h = append(*h, a)
}

func (h *alarmHeap) Pop() any {
	old := *h
	n := len(old)
	a := old[n-1]
	old[n-1] = nil // avoid memory leak
	a.index = -1
	*h = old[:n-1]
	return a
}

// setAlarmReq is the message sent through the buffered channel on SetAlarm.
type setAlarmReq struct {
	name     string
	duration time.Duration
}

// AlarmClock is a high-throughput, lock-free (single-writer) alarm scheduler.
//
// Design:
//   - Buffered channel (default 1M) absorbs burst traffic without blocking callers.
//   - Single processor goroutine (actor model) owns the heap + map exclusively,
//     so no mutex is needed for writes.
//   - GetAlarmCallback returns a closure that captures the alarm name. The actual
//     read from the map happens whenever the caller invokes the closure.
type AlarmClock struct {
	alarms map[string]*Alarm // O(1) lookup by name. Owned by processor goroutine.
	hp     alarmHeap         // Min-heap ordered by FireAt.

	reqChan  chan setAlarmReq // Buffered channel for high-throughput ingestion.
	stopChan chan struct{}    // Closed to signal graceful shutdown.
	done     chan struct{}    // Closed when the processor goroutine exits.
	stopped  atomic.Bool      // Set after Stop(); prevents writes to closed channel.

	nowFunc func() time.Time // Clock source; overridable for testing.
}

// ---------------------------------------------------------------------------
// Singleton
// ---------------------------------------------------------------------------

var (
	singletonInstance *AlarmClock
	singletonOnce     sync.Once
)

// GetInstance returns the process-wide singleton AlarmClock.
// First call creates it with the given buffer size; subsequent calls return
// the same instance.
func GetInstance(channelBufSize int) *AlarmClock {
	singletonOnce.Do(func() {
		singletonInstance = NewAlarmClock(channelBufSize)
	})
	return singletonInstance
}

// ResetInstance destroys the singleton (testing only).
func ResetInstance() {
	singletonInstance = nil
	singletonOnce = sync.Once{}
}

// NewAlarmClock creates an independent AlarmClock. Use for tests or when
// multiple isolated schedulers are needed.
func NewAlarmClock(channelBufSize int) *AlarmClock {
	if channelBufSize <= 0 {
		channelBufSize = 1_000_000
	}
	ac := &AlarmClock{
		alarms:   make(map[string]*Alarm, 1024),
		hp:       make(alarmHeap, 0, 1024),
		reqChan:  make(chan setAlarmReq, channelBufSize),
		stopChan: make(chan struct{}),
		done:     make(chan struct{}),
		nowFunc:  time.Now,
	}
	heap.Init(&ac.hp)
	return ac
}

// Start launches the background processor goroutine.
func (ac *AlarmClock) Start() {
	go ac.processLoop()
}

// Stop signals the processor to shut down and blocks until it exits.
func (ac *AlarmClock) Stop() {
	ac.stopped.Store(true)
	close(ac.stopChan)
	<-ac.done
}

// SetAlarm schedules or conditionally updates an alarm.
//
// Rules:
//   - New alarm → scheduled at now + d.
//   - Existing alarm → updated ONLY if (now + d) < current fire time.
//     We never push an alarm further into the future.
func (ac *AlarmClock) SetAlarm(name string, d time.Duration) error {
	if ac.stopped.Load() {
		return ErrStopped
	}
	select {
	case ac.reqChan <- setAlarmReq{name: name, duration: d}:
		return nil
	case <-ac.stopChan:
		return ErrStopped
	}
}

// GetAlarmCallback returns a closure that, when invoked, returns the fire
// time of the named alarm and whether it exists.
//
// Usage:
//
//	cb := clock.GetAlarmCallback("morning-alarm")
//	fireAt, exists := cb()
func (ac *AlarmClock) GetAlarmCallback(name string) func() (time.Time, bool) {
	return func() (time.Time, bool) {
		if a, ok := ac.alarms[name]; ok {
			return a.FireAt, true
		}
		return time.Time{}, false
	}
}

// ---------------------------------------------------------------------------
// Processor loop (actor model) — single goroutine.
//
// Batch-drains up to maxBatchSize requests before looping to keep throughput
// high under burst traffic.
// ---------------------------------------------------------------------------

const maxBatchSize = 4096

func (ac *AlarmClock) processLoop() {
	defer close(ac.done)

	for {
		select {
		case req := <-ac.reqChan:
			ac.handleSet(req)

			// Batch-drain: absorb all readily available requests.
		drain:
			for i := 0; i < maxBatchSize; i++ {
				select {
				case r := <-ac.reqChan:
					ac.handleSet(r)
				default:
					break drain
				}
			}

		case <-ac.stopChan:
			return
		}
	}
}

// handleSet processes a single SetAlarm request inside the processor goroutine.
func (ac *AlarmClock) handleSet(req setAlarmReq) {
	newFireAt := ac.nowFunc().Add(req.duration)

	if existing, ok := ac.alarms[req.name]; ok {
		// Only move the alarm closer, never further away.
		if newFireAt.Before(existing.FireAt) {
			existing.FireAt = newFireAt
			heap.Fix(&ac.hp, existing.index)
		}
		return
	}

	// New alarm.
	alarm := &Alarm{Name: req.name, FireAt: newFireAt}
	heap.Push(&ac.hp, alarm)
	ac.alarms[req.name] = alarm
}
