package main

import (
	"context"
	"fmt"
	"time"
)

func withCancelExample() {
	fmt.Println("== withcancel ==")
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("stopped:", ctx.Err())
				return
			default:
				fmt.Println("working")
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()

	time.Sleep(500 * time.Millisecond)
	cancel()
	time.Sleep(100 * time.Millisecond)
}

func withTimeoutExample() {
	fmt.Println("== withtimeout ==")
	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancel()

	select {
	case <-simulateWork(700 * time.Millisecond):
		fmt.Println("done in time")
	case <-ctx.Done():
		fmt.Println("timeout:", ctx.Err())
	}
}

func withDeadlineExample() {
	fmt.Println("== withdeadline ==")
	deadline := time.Now().Add(300 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	select {
	case <-simulateWork(150 * time.Millisecond):
		fmt.Println("done before deadline")
	case <-ctx.Done():
		fmt.Println("deadline hit:", ctx.Err())
	}
}

type ctxKey string

func withValueExample() {
	fmt.Println("== withvalue ==")
	ctx := context.WithValue(context.Background(), ctxKey("request_id"), "req-101")
	ctx = context.WithValue(ctx, ctxKey("user"), "sam")

	processRequest(ctx)
}

func processRequest(ctx context.Context) {
	reqID := ctx.Value(ctxKey("request_id"))
	user := ctx.Value(ctxKey("user"))
	fmt.Printf("request=%v user=%v\n", reqID, user)
}

func doneLoopExample() {
	fmt.Println("== ctx.done in loop ==")
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	steps := []string{"fetch", "validate", "save"}

	for _, step := range steps {
		select {
		case <-ctx.Done():
			fmt.Printf("stopped before %s: %v\n", step, ctx.Err())
			return
		default:
		}

		fmt.Printf("run %s\n", step)
		time.Sleep(220 * time.Millisecond)
	}

	fmt.Println("all steps done")
}

func propagationExample() {
	fmt.Println("== propagation ==")
	parentCtx, parentCancel := context.WithCancel(context.Background())
	childCtx, childCancel := context.WithTimeout(parentCtx, 2*time.Second)
	defer childCancel()

	go func() {
		<-childCtx.Done()
		fmt.Println("child stopped:", childCtx.Err())
	}()

	time.Sleep(200 * time.Millisecond)
	parentCancel()
	time.Sleep(100 * time.Millisecond)
}

func simulateWork(d time.Duration) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		time.Sleep(d)
		close(ch)
	}()
	return ch
}

func main() {
	withCancelExample()
	fmt.Println()
	withTimeoutExample()
	fmt.Println()
	withDeadlineExample()
	fmt.Println()
	withValueExample()
	fmt.Println()
	doneLoopExample()
	fmt.Println()
	propagationExample()
}
