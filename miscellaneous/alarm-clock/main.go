package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== AlarmClock Demo ===")

	clock := GetInstance(100_000)
	clock.Start()
	defer clock.Stop()

	// Set three alarms.
	clock.SetAlarm("morning-coffee", 2*time.Second)
	clock.SetAlarm("standup-meeting", 3*time.Second)
	clock.SetAlarm("lunch-break", 4*time.Second)
	time.Sleep(30 * time.Millisecond)

	// Callback: get exact fire time.
	cb := clock.GetAlarmCallback("morning-coffee")
	if fireAt, exists := cb(); exists {
		fmt.Printf("morning-coffee fires at: %s\n", fireAt.Format("15:04:05.000"))
	}

	// Try to push morning-coffee 10s later → should be ignored.
	clock.SetAlarm("morning-coffee", 10*time.Second)
	time.Sleep(30 * time.Millisecond)
	if fireAt, exists := cb(); exists {
		fmt.Printf("morning-coffee still at: %s (unchanged)\n", fireAt.Format("15:04:05.000"))
	}

	// Move lunch-break sooner (1s < 4s) → should update.
	clock.SetAlarm("lunch-break", 1*time.Second)
	time.Sleep(30 * time.Millisecond)
	lunchCb := clock.GetAlarmCallback("lunch-break")
	if fireAt, exists := lunchCb(); exists {
		fmt.Printf("lunch-break moved to: %s\n", fireAt.Format("15:04:05.000"))
	}

	fmt.Println("=== Done ===")
}
