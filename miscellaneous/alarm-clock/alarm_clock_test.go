package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// newTestClock creates an isolated AlarmClock for testing.
func newTestClock(bufSize int) *AlarmClock {
	if bufSize == 0 {
		bufSize = 10_000
	}
	return NewAlarmClock(bufSize)
}

// TestSetAlarmBasic verifies a basic alarm is stored with the correct fire time.
func TestSetAlarmBasic(t *testing.T) {
	ac := newTestClock(0)
	ac.Start()
	defer ac.Stop()

	before := time.Now()
	ac.SetAlarm("test-alarm", 5*time.Second)
	time.Sleep(30 * time.Millisecond) // let processor handle it

	cb := ac.GetAlarmCallback("test-alarm")
	fireAt, exists := cb()

	if !exists {
		t.Fatal("alarm 'test-alarm' should exist after SetAlarm")
	}

	expected := before.Add(5 * time.Second)
	tolerance := 200 * time.Millisecond
	if fireAt.Before(expected.Add(-tolerance)) || fireAt.After(expected.Add(tolerance)) {
		t.Errorf("fire time %v not within %v of expected %v", fireAt, tolerance, expected)
	}
}

// TestSetAlarmUpdateToSoonerTime verifies the alarm is updated when
// the new time (now + d) is earlier than the existing fire time.
func TestSetAlarmUpdateToSoonerTime(t *testing.T) {
	ac := newTestClock(0)
	ac.Start()
	defer ac.Stop()

	// Set alarm far in the future.
	ac.SetAlarm("my-alarm", 10*time.Second)
	time.Sleep(30 * time.Millisecond)

	cb := ac.GetAlarmCallback("my-alarm")
	originalFireAt, _ := cb()

	// Update to fire much sooner. Since 1s < 10s, this should update.
	ac.SetAlarm("my-alarm", 1*time.Second)
	time.Sleep(30 * time.Millisecond)

	newFireAt, exists := cb()
	if !exists {
		t.Fatal("alarm should still exist")
	}
	if !newFireAt.Before(originalFireAt) {
		t.Errorf("alarm should have moved earlier: original=%v, new=%v", originalFireAt, newFireAt)
	}
}

// TestSetAlarmKeepOlderTime verifies the alarm is NOT updated when the
// new time would be later than the existing fire time.
func TestSetAlarmKeepOlderTime(t *testing.T) {
	ac := newTestClock(0)
	ac.Start()
	defer ac.Stop()

	// Set alarm to fire in 1s.
	ac.SetAlarm("keep-old", 1*time.Second)
	time.Sleep(30 * time.Millisecond)

	cb := ac.GetAlarmCallback("keep-old")
	originalFireAt, _ := cb()

	// Attempt to update to a LATER time (10s). Should be rejected.
	ac.SetAlarm("keep-old", 10*time.Second)
	time.Sleep(30 * time.Millisecond)

	currentFireAt, exists := cb()
	if !exists {
		t.Fatal("alarm 'keep-old' should still exist")
	}
	if !currentFireAt.Equal(originalFireAt) {
		t.Errorf("alarm time changed unexpectedly: original=%v, current=%v",
			originalFireAt, currentFireAt)
	}
}

// TestGetAlarmCallbackNotFound verifies the callback returns false for
// non-existent alarms.
func TestGetAlarmCallbackNotFound(t *testing.T) {
	ac := newTestClock(0)
	ac.Start()
	defer ac.Stop()

	cb := ac.GetAlarmCallback("nonexistent")
	_, exists := cb()
	if exists {
		t.Error("callback for non-existent alarm should return exists=false")
	}
}

// TestCallbackReturnsCorrectTime verifies the callback returns the exact
// scheduled fire time.
func TestCallbackReturnsCorrectTime(t *testing.T) {
	ac := newTestClock(0)
	ac.Start()
	defer ac.Stop()

	before := time.Now()
	ac.SetAlarm("cb-test", 3*time.Second)
	time.Sleep(30 * time.Millisecond)

	cb := ac.GetAlarmCallback("cb-test")
	fireAt, exists := cb()

	if !exists {
		t.Fatal("alarm 'cb-test' should exist")
	}

	expected := before.Add(3 * time.Second)
	tolerance := 200 * time.Millisecond
	if fireAt.Before(expected.Add(-tolerance)) || fireAt.After(expected.Add(tolerance)) {
		t.Errorf("fire time %v not within %v of expected %v", fireAt, tolerance, expected)
	}
}

// TestMultipleAlarms verifies that multiple distinct alarms are stored
// independently.
func TestMultipleAlarms(t *testing.T) {
	ac := newTestClock(0)
	ac.Start()
	defer ac.Stop()

	ac.SetAlarm("alarm-1", 1*time.Second)
	ac.SetAlarm("alarm-2", 2*time.Second)
	ac.SetAlarm("alarm-3", 3*time.Second)
	time.Sleep(50 * time.Millisecond)

	for _, name := range []string{"alarm-1", "alarm-2", "alarm-3"} {
		cb := ac.GetAlarmCallback(name)
		if _, exists := cb(); !exists {
			t.Errorf("alarm '%s' should exist", name)
		}
	}
}

// TestDuplicateAlarmNames verifies that setting the same alarm name multiple
// times with decreasing durations results in the earliest fire time.
func TestDuplicateAlarmNames(t *testing.T) {
	ac := newTestClock(0)
	ac.Start()
	defer ac.Stop()

	// Set the same alarm 5 times with decreasing durations.
	for i := 5; i >= 1; i-- {
		ac.SetAlarm("dup", time.Duration(i)*time.Second)
	}
	time.Sleep(50 * time.Millisecond)

	cb := ac.GetAlarmCallback("dup")
	fireAt, exists := cb()
	if !exists {
		t.Fatal("alarm 'dup' should exist")
	}

	// The fire time should be ~1s from now (the shortest duration).
	expected := time.Now().Add(1 * time.Second)
	tolerance := 500 * time.Millisecond
	if fireAt.Before(expected.Add(-tolerance)) || fireAt.After(expected.Add(tolerance)) {
		t.Errorf("fire time %v not near expected ~1s from now (%v)", fireAt, expected)
	}
}

// TestSetAlarmAfterStop verifies that SetAlarm returns ErrStopped after Stop.
func TestSetAlarmAfterStop(t *testing.T) {
	ac := newTestClock(0)
	ac.Start()
	ac.Stop()

	if err := ac.SetAlarm("after-stop", time.Second); err != ErrStopped {
		t.Errorf("expected ErrStopped after Stop(), got %v", err)
	}
}

// TestGracefulStop verifies that Stop() returns promptly.
func TestGracefulStop(t *testing.T) {
	ac := newTestClock(0)
	ac.Start()

	for i := 0; i < 100; i++ {
		ac.SetAlarm(fmt.Sprintf("stop-%d", i), time.Hour)
	}

	done := make(chan struct{})
	go func() {
		ac.Stop()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Stop() did not return within 5 seconds")
	}
}

// TestSingleton verifies that GetInstance returns the same instance.
func TestSingleton(t *testing.T) {
	ResetInstance()
	defer ResetInstance()

	a := GetInstance(100)
	b := GetInstance(999) // different size — should be ignored

	if a != b {
		t.Error("GetInstance must return the same instance on every call")
	}
}

// TestHeapOrdering verifies the min-heap keeps the earliest alarm at the top.
func TestHeapOrdering(t *testing.T) {
	ac := newTestClock(0)
	ac.Start()
	defer ac.Stop()

	ac.SetAlarm("far", 10*time.Second)
	ac.SetAlarm("near", 1*time.Second)
	ac.SetAlarm("mid", 5*time.Second)
	time.Sleep(50 * time.Millisecond)

	// The heap root should be the "near" alarm.
	if ac.hp.Len() == 0 {
		t.Fatal("heap should not be empty")
	}
	if ac.hp[0].Name != "near" {
		t.Errorf("heap root should be 'near', got '%s'", ac.hp[0].Name)
	}
}

// TestHighConcurrency verifies that the system handles 1,000,000 concurrent
// SetAlarm calls from 100 goroutines without errors.
func TestHighConcurrency(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping high-concurrency test in short mode")
	}

	const numGoroutines = 100
	const alarmsPerGoroutine = 10_000
	const totalAlarms = numGoroutines * alarmsPerGoroutine // 1,000,000

	ac := NewAlarmClock(totalAlarms)
	ac.Start()
	defer ac.Stop()

	start := time.Now()
	var wg sync.WaitGroup
	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(gid int) {
			defer wg.Done()
			for i := 0; i < alarmsPerGoroutine; i++ {
				name := fmt.Sprintf("g%d-a%d", gid, i)
				if err := ac.SetAlarm(name, time.Hour); err != nil {
					t.Errorf("SetAlarm failed: %v", err)
					return
				}
			}
		}(g)
	}
	wg.Wait()
	elapsed := time.Since(start)

	t.Logf("Ingested %d alarms in %v (%.0f alarms/sec)",
		totalAlarms, elapsed, float64(totalAlarms)/elapsed.Seconds())

	// Let processor drain.
	time.Sleep(500 * time.Millisecond)

	if ac.hp.Len() != totalAlarms {
		t.Errorf("expected %d alarms in heap, got %d", totalAlarms, ac.hp.Len())
	}
}

// BenchmarkSetAlarm measures single-goroutine SetAlarm throughput.
func BenchmarkSetAlarm(b *testing.B) {
	ac := NewAlarmClock(b.N + 1)
	ac.Start()
	defer ac.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ac.SetAlarm(fmt.Sprintf("bench-%d", i), time.Hour)
	}
}

// BenchmarkSetAlarmParallel measures multi-goroutine SetAlarm throughput.
func BenchmarkSetAlarmParallel(b *testing.B) {
	ac := NewAlarmClock(2_000_000)
	ac.Start()
	defer ac.Stop()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			ac.SetAlarm(fmt.Sprintf("pbench-%d", i), time.Hour)
			i++
		}
	})
}

// BenchmarkGetAlarmCallback measures callback creation + invocation throughput.
func BenchmarkGetAlarmCallback(b *testing.B) {
	ac := NewAlarmClock(10_000)
	ac.Start()
	defer ac.Stop()

	ac.SetAlarm("bench-cb", time.Hour)
	time.Sleep(20 * time.Millisecond)

	cb := ac.GetAlarmCallback("bench-cb")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cb()
	}
}
