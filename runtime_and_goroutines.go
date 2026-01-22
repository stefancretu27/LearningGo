package go_didactical_apps

import (
	"fmt"
	"runtime"
	"sync"
)

	/**
	* 1. The Go runtime environment is a set of components that support the execution of Go programs. It includes:
	*	- Garbage Collector (GC): Automatically frees unused memory.
	* 	- Scheduler: Manages goroutines (lightweight threads).
	*	- Stack management: Dynamically grows/shrinks goroutine stacks.
	*	- Concurrency primitives: Channels, mutexes, etc.
	*	- Standard library support: For networking, I/O, etc.
	*
	* -> The runtime is an OS-specific binary, which is linked into the binary of the built program, during compilation, making Go programs self-contained.
	* -> Unlike Java, Go compiles directly to native machine code, so there's no virtual machine (VM) involved.
	 */

	/**
	* 2. Goroutines' management by the Go Scheduler
	* -> Go uses an M:N scheduler:
	*	- M OS threads
	* 	- N goroutines
	* 	- P logical processors (set by GOMAXPROCS). In Go's runtime scheduler, P stands for "Processor", but it's not a physical CPU core —
	* it's a logical processor used by the Go runtime to manage execution of goroutines. P is associated with a OS thread (M), at a given time, that
	* executes the goroutines.
	* Each P maintains a local run FIFO queue of goroutines that are ready to run. When an OS thread (M) is assigned to a P, it pulls goroutines from this queue and 
	* executes them one at a time.
	* The number of Ps is controlled by runtime.GOMAXPROCS(n). At any moment, only GOMAXPROCS goroutines can run in parallel, but many more can be concurrent 
	* (waiting, blocked, etc.).
	* The number of Ms can be greater than P — Go may spawn more OS threads (M) to handle blocking operations like I/O or syscalls.
	*
	* -> The Go runtime:
	*	- Maps many goroutines (N) onto a smaller number of OS threads (M)
	*	- Uses 3 concepts to ensure fairness and flexibility amongst goroutines' execution, avoiding starvation:
	*		i) work-stealing (if a Ps run out of goroutines, it can take others from another P)
	*		ii) preemption (long running goroutines, with blocking I/O, for example, are rescheduled, as the compiler inserts safe points in the code, 
	* where the goroutine can be paused)
	*		iii) and cooperative scheduling (goroutines are executed till a blocking operation is encountered, then replaced by another goroutine) 
	*	- Handles blocking syscalls by detaching the thread and spinning up another to keep the system responsive. Goroutines use non-blocking I/O and runtime-managed blocking
	* (e.g., netpoller), so blocking one goroutine doesn’t block the underlying OS thread.
	*	- blocked goroutines (mutex, channel recv, I/O operation etc) are replaced by ready goroutines.
	*
	* -> Stack Management:
	*	- Goroutines start with a small stack (~2 KB)
	*	- The stack grows and shrinks dynamically
	*	- This is in contrast to C++ threads, which typically allocate a fixed-size stack (e.g., 1 MB)
	*/
	
	/**
	* 3. Goroutines memory management.
	* They are more lightweight than OS threads in terms of:
	*	- Memory usage
	*	- Startup cost
	*	- Scalability
	*	- Scheduling flexibility
	* Thus, there can be spawned up to hundreds of thousand of goroutines, as there is no limit. Ther limits are given by the available memory, OS thread limits and scheduler overhead.
	* The OS thread limit per process can be around 32768/65536, but it can be increased/decreased. This number is dependent on stack size and amount of virtual memory
	*	- number of threads = total virtual memory / (stack size * 1024 * 1024)
	*
	* Local Variables & Arguments & return values of a goroutine:
	*	- Stored on the goroutine's stack, not the thread's.
 	*	- Each goroutine has its own stack, so local variables are isolated. The stack is managed by Go runtime, not by OS (as it happens with OS threads)
	* 	- Arguments passed to a goroutine function are copied onto its stack.
	*	- Even if multiple goroutines run on the same OS thread (via scheduling), their local variables are not shared.
	*
	* Thread specific resources and goroutines:
	*	- goroutine inherits the thread’s signal mask and thread-local state while executing
	*	- when a goroutine is executing on a thread, it uses the thread’s CPU registers — because registers are part of the CPU state during execution.
	*	- when a goroutine is executing on a thread, the Go runtime switches the stack pointer to the goroutine’s own stack. So, it executes using its own stack — not the thread’s stack.
	*	The thread's stack segment memory, allocated by OS, is used for: thread-local data, system-level operations, Go runtime internals (e.g., scheduling, syscall wrappers)
	*
	* Execution Flow:
	*	1. Goroutine is created:
	*		- It gets its own stack and execution context.
	*		- It’s placed in a queue managed by a P.
	*		- Scheduler picks a goroutine (G) from the queue of a P.
	*	2. An OS thread (M) is assigned to the P and begins executing the goroutine.
	*	3. The goroutine runs inside the thread, using the thread’s CPU and system resources.
	*	4. If the goroutine blocks (e.g., on I/O or syscall):
	*		- The thread may be parked.
	*		- Another thread may be spun up to keep the system responsive.
	*	5. When the goroutine finishes:
	*		- The thread can pick another goroutine from the queue.
	*		- The goroutine’s stack and resources are cleaned up.
	 */

	 /**
	 * 4. Goroutines' execution:
	 *	- when a Go program starts, it begins execution in the main function, which runs in the main goroutine. This is analogous to the "main thread" in C/C++.
	 *	- the main goroutine is special in that when it exits, the entire program terminates — regardless of whether other goroutines are still running.
	 *	- are not automatically joined like threads in some other languages. If the main goroutine exits, all other goroutines are abruptly terminated, 
	 *	even if they haven't finished.
	 *	- to ensure all gorotuines finish their execution before the main goroutine exits, synch mechanisms are used
	 */

func producerFor(sentChan chan<- int) {
	fmt.Println("[producerFor] len:", len(sentChan), " cap:", cap(sentChan))
	for  idx := 0; idx < cap(sentChan); idx++ {
		sentChan <- idx
	}
	close(sentChan)
}

func consumerFor(recvChan <-chan int) {
	for data := range recvChan {
		fmt.Println("[consumerFor] Received:", data)
	}
}

func consumerInfiniteLoop(recvChan <-chan int) {
	ok := true
	for ok {
		val, ok := <-recvChan
		if !ok {
			fmt.Println("[consumerInfiniteLoop] Channel closed")
			return
		}
		fmt.Println("[consumerInfiniteLoop] Received:", val)
	}
}

func consumerSelectInfiniteLoop(recvChan <-chan int) {
	for {
		select {
		case val, ok := <-recvChan:
			if !ok {
				fmt.Println("[consumerSelectInfiniteLoop] Channel closed")
				return
			} else {
				fmt.Println("[consumerSelectInfiniteLoop] Received:", val)
			}
		default:
			//will print many times, till values are received. Normally, do nothing, or can be ommitted -> blocking select
			fmt.Println("[consumerSelectInfiniteLoop] Non-blocking select:")
		}
	}
}

func consumerSelectTimeoutNonBlocking(recvChan <-chan int) {
	select {
		case val, ok := <-recvChan:
			if !ok {
				fmt.Println("[consumerSelectTimeout] Channel closed")
				return
			} else {
				fmt.Println("[consumerSelectTimeout] Received:", val)
			}
		case <-time.After(1*time.Second):
			fmt.Println("[consumerSelectTimeout] Timeout:")
		}
}

func producer(sharedData *int) {
	*sharedData = 27
	fmt.Println("[producer]:", *sharedData)
}

func consumer(sharedData int) {
	fmt.Println("[consumer]:", sharedData)
}

func producerUnbufChan(sendChan chan<- int) {
	len := 7
	for idx := range len {
		sendChan <- idx
	}

	close(sendChan)
}

func consumerUnbufChan(recvChan <-chan int) {
	for {
		if data, ok := <-recvChan; ok {
			fmt.Println("[consumerUnbufChan]", data)
		} else {
			fmt.Println("[consumerUnbufChan] closed channel")
			return
		}
	}
}

func GoRuntimeAndGoroutines() {
	//The GOMAXPROCS(n) variable limits the number of operating system threads that can execute user-level Go code simultaneously. If n<1, no change is made 
	runtime.GOMAXPROCS(32)
	fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0))

	// Goroutines idioms
	// 1. launch and forget (fire and forget): start a groutine and don't wait for it, as no result or value is wated from it
	go func() {
		fmt.Println("Launch and forget: do not wait for the goroutine, from the goroutine it was started from!")
	}

	/* 2. Encapsulate concurrent work in a separate function, whose call is managed by an anonymous function started as a goroutine
	* Adavantages:
	*  - closure: can access surrounding variables without passing them as parameters to the function, but as arguments to anonymous function
	*  - synch primitives management: synch mechanisms deal with the anonymous function, whereas the functions itself can be agnostic of concurrency management 
	* 	go func() {
	* 		doWork()
	* 	}
	*/

	sharedData := 9
	syncChan := make(chan struct{})
	var waitGroup sync.WaitGroup
	waitGroup.Add(2)
	go func() {
		defer waitGroup.Done()
		producer(&sharedData)
		syncChan <- struct{}{}
	}()
	// Channels can be used to synch: struct{} is the empty struct type — a special type that has zero size and no fields, perfect for signaling without carrying data
	go func() {
		defer waitGroup.Done()
		<-syncChan
		consumer(sharedData)
	}()
	waitGroup.Wait()

	/*
	* 3. Defer cleanup and closing. Advantages:
	*	- ensures cleanup is scheduled and executed even if subsequent computations panic/fail, thus exiting prematurely, as opposed to cleanup up at the end of function
	*	- starting the logic with cleanup makes the developer to focus firstly on resources' management, thus less liely to forget some
	*/

	intBufferedChan := make(chan int, 5)
	waitGroup.Add(2)
	/*
	* 4.1 Close channel idiom: only the sender must close the channel, to signal the sending is done
	*/
	go func(bufferedChan chan int) {
		defer waitGroup.Done()
		producerFor(intBufferedChan)
	}(intBufferedChan)
	/*
	* 4.2 Use for range over the receiving channel, to receive until the channel is closed by the sender
	*/
	go func(bufferedChan chan int) {
		defer waitGroup.Done()
		consumerFor(intBufferedChan)
	}(intBufferedChan)
	waitGroup.Wait()

	/* 5. Wait for goroutines to finish work.
	*	- as the main goroutine might exit before all goroutines finish work, their results can be discarded
	*	- an approach is to use sync.WaitGroup, which acts as a barrier, as shown above and beneath
	*/
	// Allocate a new channel, as the previous one is closed
	intBufferedChan := make(chan int, 5)
	waitGroup.Add(2)
	go func(bufferedChan chan int) {
		defer waitGroup.Done()
		producerFor(intBufferedChan)
	}(intBufferedChan)
	/*
	* 6. Comma, ok <- recvChannel, then check ok bool status for received values (true), or fi the channel was closed (false)
	*/
	go func(bufferedChan chan int) {
		defer waitGroup.Done()
		consumerInfiniteLoop(intBufferedChan)
	}(intBufferedChan)
	waitGroup.Wait()

	// Create a new channel
	intBufferedChan = make(chan int, 5)
	waitGroup.Add(2)
	go func(bufferedChan chan int) {
		defer waitGroup.Done()
		producer(intBufferedChan)
	}(intBufferedChan)
	/*
	* 7. Use select with default case to avoid waiting forever (blocking)
	*/
	go func(bufferedChan chan int) {
		defer waitGroup.Done()
		consumerSelectInfiniteLoop(intBufferedChan)
	}(intBufferedChan)
	waitGroup.Wait()

	// Create a new channel
	intBufferedChan = make(chan int, 5)
	waitGroup.Add(2)
	go func(bufferedChan chan int) {
		defer waitGroup.Done()
		producer(intBufferedChan)
	}(intBufferedChan)
	/*
	* 8. Use select with case for timeout ( <-time.After(timeout)) to avoid waiting forever (blocking)
	*/
	go func(bufferedChan chan int) {
		defer waitGroup.Done()
		consumerSelectTimeoutNonBlocking(intBufferedChan)
	}(intBufferedChan)
	waitGroup.Wait()

	

}