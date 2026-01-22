package go_didactical_apps

/*
* Context:
*	- interface that allows to:
*		> Cancel operations early (e.g., if a client disconnects).
*		> Set timeouts or deadlines for operations.
*		> Pass request-scoped values (like user IDs, trace IDs).
*
*	- defined as:
*		type Context interface {
*    		Deadline() (deadline time.Time, ok bool)
*    		Done() <-chan struct{}
*   		Err() error
*   		Value(key any) any
*		}
*	- methods:
*		1. Deadline()
*			- returns the time when the context will be automatically cancelled.
*			- if no deadline is set, ok is false.
*		2. Done()
*			- returns a receive-only channel (<-chan struct{}), but it's not stored in the interface — it's returned by the underlying implementation
*		that is a struct implementing the interface, named cancelCtx. It's defined inside the context package. 
*			- it is closed when the context is cancelled or times out.
*			- goroutines can select on this to know when to stop.
*		3. Err()
*			- returns the reason for cancellation 
*			- context.Canceled if manually cancelled 
*			- context.DeadlineExceeded if timed out
*			- nil if still active.
*		4. Value(key)
*			- retrieves a value associated with the context.
*			- used for request-scoped data (e.g., user ID, trace ID).
*
*	- context.Background() is considered the default context in Go, the parent context for all contexts: it has no deadline, no cancellation, and no values.
*	Basically, it does nothing, just the parent for other contexts.
*
*	- Go does not implicitly associate any context with main() or goroutines. The context must be explicitly passed as argument if you want cancellation, 
*	deadlines, or scoped values. If goroutines are launched without passing a context, they will run independently and cannot be cancelled externally.
*
*	- Context is an interface, and in Go, interfaces are reference types. Passing it by value means you're passing a reference to the underlying context object, 
*	not copying the entire structure. That said, the context should always be passed by value.
*
*	- context is implemented using a tree-like structure, as each created Context can have a parent. When calling context.WithCancel, WithTimeout, or WithDeadline, 
*	a child context is created. Cancelling a parent context automatically cancels all its children. The cancellation is propagated using channels 
*	(Done() <-chan struct{}), which goroutines can listen to. The implementation is safe for concurrency.
*
* Context types:
* 1. context.Background
*	- use case: used when no cancellation or timeout is needed. Ideal for initializing clients, servers, or long-lived processes.
*	- does nothing, as it has no cancellation, no deadline and carries no values, being parent for other contexts
*	- declaration: ctx := context.Background()
*
* 2. context.TODO
*	- placeholder when you don’t yet know what context to use. Useful during development or when refactoring.
*	- declaration: ctx := context.TODO()
*
* 3. context.WithCancel(parent)
*	- use case: allows manual cancellation of operations. Useful when you want to stop work early (e.g., user aborts a request).
*	- declaration: ctx, cancel := context.WithCancel(context.Background())
*				   defer cancel()
*	- cancel is just a func() — a function with no parameters and no return value, used to manually cancel the context and all its children.
*	- ctx.Err() will return "context canceled" if it was manually cancelled.
*	- Internally, calling cancel():
*		> Closes the Done() channel of the context.
*		> Sets the error returned by ctx.Err() to context.Canceled.
*		> Propagates cancellation to all child contexts.
*		> Signals goroutines waiting on ctx.Done() to stop work.
*
* 4. context.WithTimeout(parent, duration)
*	- use case: automatically cancels context after a timeout, or when cancelFunc is closed, whichever occurs first. 
*	- ideal for network calls, database queries, or time-sensitive operations.
*	- a time interval/duration is passed as second argument
*	- syntax: ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
* 			  defer cancel()
*	- ctx.Err() will return "context.DeadlineExceeded" if the timeout expired. Or "context canceled" if it was manually cancelled.
*
* 5. context.WithDeadline(parent, time)
*	- use case: automatically cancels context at a specific time.
*	- a time point is passed as second argument
*	- syntax: deadline := time.Now().Add(2 * time.Second)
*			  ctx, cancel := context.WithDeadline(context.Background(), deadline)
*			  defer cancel()
*	- when the parent context is another context with deadline, or a context with timeout, the deadline of the newly created context
*	is the earliest of the times in the such-created inheritance chain. Thus, it is ensured the child context doesn't outlive the parent
*
* 6. context.WithValue(parent, value):
*	- context meant for passing request-scoped data (like user IDs, auth tokens, deadlines) downstream
*	- context values are immutable: once set a value in a context using context.WithValue, it cannot be changed
*	- declaration: type keyType string  - use a custom type to avoid key conflicts
*				   key := keyType("context key") - define a key, allocating memory and providing value
*				   ctxValue := context.WithValue(context.Background(), key, 27) - create context specifying the key and the value
*	- a context can carrey multiple key-value pairs, only in cases it is created from a parent context that also carries key-value pair(s).
*	Thus, leading to a sort of inheritance chain in terms of key-values pairs.
*	- if created from context.Background(), a context with value yields only one key-value pair.
*/

func CtxCancel() {
	ctx, cancelFunc := context.WithCancel(context.Background())

	fmt.Println("[CtxCancel] Created context with cancel, whose cancel func is used to close the goroutines")

	bufChan := make(chan int, 5)
	signalChan := make(chan struct{}, 1)

	go func() {
		if ctxD, hasDeadline := ctx.Deadline(); hasDeadline {
			fmt.Println("[CtxCancel] Context with cancel, whose deadline is", ctxD.String())
		} else {
			fmt.Println("[CtxCancel] Context with cancel has no deadline")
		}

		for idx := 0; idx < cap(bufChan); idx++ {
			bufChan <- (idx + 1) * (idx + 1)
			time.Sleep(1 * time.Millisecond)
		}
		close(bufChan)
		cancelFunc()
		signalChan <- struct{}{}
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("[CtxCancel] ctx.Done:", ctx.Err())
			return
		case val, isOpen := <-bufChan:
			if isOpen {
				fmt.Println("[CtxCancel] recv val:", val)
			} else {
				fmt.Println("[CtxCancel] channel is closed")
				return
			}
		case <-signalChan:
			fmt.Println("[CtxCancel] notified on signalChan")
			return
		}
	}
}

func CtxTimeout(timeoutVal time.Duration) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeoutVal)
	fmt.Println("[CtxTimeout] Created context with timeout, which is of type time.Duration, representing a time interval for the context's lifespan")

	var wg sync.WaitGroup
	wg.Add(2)
	unbufChan := make(chan int)

	go func(sendChan chan<- int, size int) {
		defer wg.Done()

		if ctxD, hasDeadline := ctx.Deadline(); hasDeadline {
			fmt.Println("[CtxTimeout] Context has timeout: ", ctxD.String())
		}

		for i := 0; i < size; i++ {
			sendChan <- i + 1
			time.Sleep(1 * time.Millisecond)
		}
		close(unbufChan)
	}(unbufChan, 7)

	go func(recvChan <-chan int) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("[CtxTimeout] ctx.Done(): ", ctx.Err())
				return

			case <-time.After(timeoutVal):
				fmt.Println("[CtxTimeout] timeout issued by time.After: ", ctx.Err())
				return

			case val, isOpen := <-recvChan:
				if isOpen {
					fmt.Println("[CtxTimeout] recv val: ", val)
				} else {
					fmt.Println("[CtxTimeout] recv channel is closed")
					return
				}
			}
		}
	}(unbufChan)

	wg.Wait()
	cancelFunc()
}

func CtxDeadline(timeoutVal time.Duration) {
	var timePoint time.Time = time.Now().Add(timeoutVal)
	var timeoutVal2 time.Duration = time.Until(timePoint)
	ctx, cancelFunc := context.WithDeadline(context.Background(), timePoint)
	defer cancelFunc()
	fmt.Println("[CtxDeadline] Created context with deadline, which is of type time.Time, representing a time point in the future when the context timeouts:", timePoint.String(), timeoutVal.String(), timeoutVal2.String())

	bufChan := make(chan int, 8)
	var wg sync.WaitGroup
	wg.Add(2)

	go func(sendChan chan<- int) {
		defer wg.Done()

		if ctxD, hasDeadline := ctx.Deadline(); hasDeadline {
			fmt.Println("[CtxDeadline] Context has deadline: ", ctxD.String())
		}

		for idx := 0; idx < cap(sendChan); idx++ {
			sendChan <- idx + idx
			time.Sleep(1 * time.Millisecond)
		}
		close(bufChan)
	}(bufChan)

	go func(recvChan <-chan int) {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				fmt.Println("[CtxDeadline] ctx.Done(): ", ctx.Err())
				return
			case <-time.After(timeoutVal):
				fmt.Println("[CtxTimeout] timeout issued by time.After: ", ctx.Err())
				return
			case val, isOpen := <-recvChan:
				if isOpen {
					fmt.Println("[CtxDeadline] recv val: ", val)

				} else {
					fmt.Println("[CtxDeadline] recv channel is closed")
					return
				}
				//default:
				//time.Sleep(1 * time.Millisecond)
			}
		}

	}(bufChan)

	wg.Wait()
}

func CtxValue() {
	//define custom type for key
	type keyT string

	key := keyT("context key")
	ctx := context.WithValue(context.Background(), key, 27)

	signalChan := make(chan struct{}, 1)
	defer close(signalChan)

	go func() {
		value := ctx.Value(key)
		fmt.Println("[CtxValue] value:", value)
		signalChan <- struct{}{}
	}()

	<-signalChan
}

// context. Value()
// define custom type to avoid key collisions
type keyType string

func accessCtxValue(ctx context.Context, key keyType, grtnId int) {
	value := ctx.Value(key)
	fmt.Println("[accessCtxValue]:", value, grtnId)
}

func AboutContext() {
	const timeoutVal = 6 * time.Millisecond
	CtxCancel()
	fmt.Println("--------------------------")
	CtxDeadline(timeoutVal)
	fmt.Println("--------------------------")
	CtxTimeout(timeoutVal)
	fmt.Println("--------------------------")
	CtxValue()

	//context.WithValue
	key := keyType("context key")
	ctxValue := context.WithValue(context.Background(), key, "context value")
	waitGroup.Add(6)
	for idx := 0; idx < 5; idx++ {
		go func() {
			defer waitGroup.Done()
			accessCtxValue(ctxValue, key, idx)
		}()
	}
	key2 := keyType("child context key")
	ctxValue2 := context.WithValue(ctxValue, key2, "child context value")
	go func() {
		defer waitGroup.Done()
		fmt.Println("child get parent context value:", ctxValue2.Value(key), " | child get context value:", ctxValue2.Value(key2))
	}()
	waitGroup.Wait()
}