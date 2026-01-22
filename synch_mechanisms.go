package go_didactical_apps

/*
* 1. Channels:
*	- not only used to pass data between goroutines, but they can eb used to signal a state, concretely, that an operation is ended
*	This entails either sending an boolean or a zero-sized object (the empty struct: struct{}) via a channel
*
* 2. Semaphores:
*	- Semaphore limits the number of goroutines which can access a critical section, at a given time point, until they release a token, 
*	whilst the other goroutines wait for the token to be released in order to proceed, further.
*	- Go doesn't have built in semaphores, but it uses buffered channels to emulate the behavior of restricting how many goroutines
*	can access a critical section, at a given time. It also entails sending the zero-sized object or a bool status
*	- declaration: sem:= make(chan struct{}, sem_size)
*	- inside goroutines: sem <- struct{}{} => send zero sized object to the buffered channel to acquire (increase the number of buffered elements)
*						 do work
*						 <- sem => read from the buffered channel to emulate semaphore release (decrease buffered elements count)
*
* 3. Wait groups:
*	- acts as a barrier which waits for all added goroutines to finish, before proceeding further with the execution flow, or to program's termination 
*	- it is widely used, as it can happen that the main goroutine exits, before spawned goroutines finish execution => wait for spawned goroutines to finish,
*	before program exits
*	- declaration: var wg sync.WaitGroup
*	- outside goroutine: wg.Add(goroutines_count) => increments a counter that keeps track of all goroutines to wait for
*	- inside goroutine: defer wg.Done() => decrements the counter by one
*	- outside goroutines: wg.Wait() => waits for the counter to reach 0, blocking th execution flow until then. It can be done in a dedicated,
*	separate goroutine, when it is needed to signal other parts of the program (via a channel) without blocking the main goroutine, upon some goroutines have finished.
*
* 4. Mutex:
*	- ensures exclusive access to a critical section, such that only one goroutine can access it
*	- declaration: var mtx sync.Mutex
*	- inside goroutine: mtx.Lock()
*						access critical section
*						mtx.Unlock()
*
* 5. Read-write mutex 
*	- allows multiple reader goroutines to access in non exclusive mode the critical section, for reading 
*	- declaration: var rwmtx sync.RWMutex
*	- inside goroutine: mtx.RLock()
*						access critical section
*						mtx.RUnlock()
*	- since a RWMutex's RLock() and RUnlock() can't be used with cond var, alternatively a signaling channel can be used for broadcast. It is recommended
*	to be buffered such that the producer doesn't block until receivers read the signal. Also, this way multiple consumers can be notified. Moreover, the
*	so called broadcast should be done outside the lock, as the consumers wait on it under RLock(). This way, it is avoided circular lock: 
*		> producer is blocked on channel send, holding Lock().
*		> consumer is blocked on channel receive, holding RLock().
*		type RWSynch struct {
*			rwmtx sync.RWMutex
*			signChan chan struct{}
*			consumersCount uint
*		}
*						
* 6. Condition variable
*	- used to signal, under a mutex lock, when the critical's section state has changed. Also, wait on consumer side on it, till the signal is sent,
*	then access the critical section. It explicitly requires an exclusive lock on waiting end.
*	- sync.Cond.Wait() requires exclusive ownership of the lock. If RWMutex with RLock() is used, multiple goroutines could hold the read lock simultaneously, 
*	which breaks the condition variable's guarantee of atomicity and coordination. This can lead to race conditions, missed signals, or deadlocks. Still, a RWMutex
* can be used, but with Lock(), Unlock() on both producer and consumer sides
*	- created from a pointer to an exclusive mutex lock, invoking NewCond which returns a pointer
*	- declaration: var mtx sync.Mutex
*					condVar := sync.NewCond(&mtx)
*	- a struct can be used to encapsulate synch data, even the shared data
*		type Synch struct {
*			mtx sync.Mutex,
*			condVar *sync.Cond
*			}
* 	- producer: cond.Signal()  - notfies one consumer
*				cond.Broadcast() - notifies multiple consumers
*	- in Go, it is idiomatic and correct to call Signal() or Broadcast() while holding the associated mutex. This is different from C++'s std::condition_variable, 
*	where it's often recommended to unlock before notifying
*	- consumer: for condition_not_met {
*					cond.Wait()
*				}
*	- a spurious wakeup is when a thread (or goroutine, in Go) waiting on a condition variable wakes up without being explicitly signaled, 
*	and the condition it was waiting for hasn't changed. This can happen due to OS scheduling or when a producer broadcasts to multiple goroutines, 
*	waking up more than necessary. The workaround is to always re-check the condition in a for loop that encapsulates the condVar.Wait() call.
*
* 7. sync.Once
*	- ensure a function is executed only once, even if invoked by multiple goroutines
*	- the function is specified alongside it's arguments, but not called in place (using ())
*	- cannot directly retrieve return values from the function passed to Do
*	- declaration: var syncOnce sync.Once
*	- usage: syncOnce.Do(my_func(args)) inside a goroutine
*/

//Sync channel 
func ProducerSynchChannel(sendChan chan<- int, synchChan chan struct{}) {
	//This design only works with buffered channels, as the buffer's capacity is accessed
	for idx := range cap(sendChan) {
		sendChan <- idx * idx
	}

	synchChan <- struct{}{}

	for idx := range cap(sendChan) {
		sendChan <- idx + idx
	}
	close(sendChan)
}

func ConsumerSynchChannel(recvChan <-chan int, synchChan chan struct{}) {
	// wait on synchChan
	<-synchChan
	// thereafter, access the values. 
	for val := range recvChan {
		fmt.Println("[ConsumerSynchChannel] received:", val)
	}
}

//Semaphore
func ProducerSem(idx int, arr *[10]int, sem chan struct{}) {
	// acquire sem token, by sending a zero-sized struct over the channel
	sem <- struct{}{}

	//do work on shared data
	arr[idx] = idx + idx

	// release sem token, by reading from the same channel
	<-sem
}


// Mutex + cond var
type Synch struct {
	mtx sync.Mutex
	condVar *sync.Cond
}

func NewSync() *Synch {
	s := &Synch{}
	s.condVar = sync.NewCond(&s.mtx)
	return s
}

func ProducerMtxCondVar(s *Synch, slice *[]int) {
	s.mtx.Lock()
	for idx := range 4 {
		*slice = append(*slice, 100+idx)
	}
	fmt.Println("[ProducerMtxCondVar]", len(*slice))
	s.condVar.Signal()
	s.mtx.Unlock()
}

func ConsumerMtxCondVar(s *Synch, slice *[]int) {
	s.mtx.Lock()
	// wait on condition variable until the condition related to the critical section is changed
	for len(*slice) == 0 {
		s.condVar.Wait()
	}
	fmt.Println("[ConsumerMtxCondVar]", len(*slice))
	for idx := range len(*slice) {
		fmt.Println("[ConsumerMtxCondVar]:", (*slice)[idx])
	}
	//alternatively
	// for idx, val := range *slice {
	// 	fmt.Println("[ConsumerMtxCondVar]:", idx, val)
	// }
	s.mtx.Unlock()
}

// RWMutex
type RWSynch struct {
	rwmtx          sync.RWMutex
	signalChan     chan struct{}
	consumersCount uint
}

func NewRWSynch(consumersCount uint) *RWSynch {
	s := &RWSynch{consumersCount: consumersCount}
	s.signalChan = make(chan struct{}, consumersCount)
	return s
}

func ProducerRWMtx(rws *RWSynch, slice *[]int) {
	rws.rwmtx.Lock()
	for idx := range 4 {
		*slice = append(*slice, 20+idx)
	}
	fmt.Println("[ProducerRWMtx] slice length:", len(*slice))
	rws.rwmtx.Unlock()

	//notify outside mutex <=> channel not protected by mutex
	for idx := uint(0); idx < rws.consumersCount; idx++ {
		rws.signalChan <- struct{}{}
	}
}

func ConsumerRWMtx(rws *RWSynch, slice *[]int, title string) {
	// wait on notification channel before mutex lock. Otherwise, the producer is blocked until the consumer receives on the channel, if the channel is unbuffered => deadlock
	fmt.Println("[ConsumerRWMtx] Title:", title)
	rws.rwmtx.RLock()
	// buffered channel and signaling on producer's side outside Lock() scope makes safe to wait on channel under RLock() scope
	<-rws.signalChan
	fmt.Println("[ConsumerRWMtx] slice length", len(*slice))
	for _, val := range *slice {
		fmt.Println("[ConsumerRWMtx]:", val)
	}
	rws.rwmtx.RUnlock()
}

func SynchMechanisms() {
	var wg sync.WaitGroup
	synchChan := make(chan struct{})
	buffChan := make(chan uint, 4)

	// Synch Channel + wait group
	wg.Add(2)
	go func() {
		defer wg.Done()
		ProducerSynchChannel(buffChan, synchChan)
	}
	go func() {
		defer wg.Done()
		ConsumerSynchChannel(buffChan, synchChan)
	}
	wg.Wait()

	// Semaphore + wait group
	semChan := make(chan struct{}, 5)
	var arr [10]int
	wg.Add(10)
	for idx := range len(arr) {
		go func() {
			defer wg.Done()
			ProducerSem(idx, &arr, semChan)
		}()
	}
	go func() {
		wg.Wait()
		fmt.Println("The 10 producers finished work!")
		//signal to consumer goroutines that all producers finished work, without blocking the main goroutine => use separate goroutine for waiting and signaling
		synchChan <- struct{}{}
	}()
	go func() {
		// wait on signaling channel
		<-synchChan
		for idx, val := range arr {
			fmt.Println("arr: ", idx, val)
		}
	}()
	time.Sleep(1 * time.Second)

	//Mutex + condition variable
	var waitGroup sync.WaitGroup
	synch := NewSynch()
	var slice []int
	waitGroup.Add(2)
	go func() {
		defer waitGroup.Done()
		ProducerMtxCondVar(synch, &slice)
	}()
	go func() {
		defer waitGroup.Done()
		ConsumerMtxCondVar(synch, &slice)
	}()
	waitGroup.Wait()

	//RWMutex with buffered channel used for broadcast signalling to avoid blocking upon sending receiving
	rws := NewRWSynch(2)
	slice = []int{}
	waitGroup.Add(3)
	go func() {
		defer waitGroup.Done()
		ProducerRWMtx(rws, &slice)
	}()
	go func() {
		defer waitGroup.Done()
		ConsumerRWMtx(rws, &slice, "consumer 0")
	}()
	go func() {
		defer waitGroup.Done()
		ConsumerRWMtx(rws, &slice, "Consumer 1")
	}()
	waitGroup.Wait()

	//sync.once
	var syncOnce sync.Once
	counter := 0
	waitGroup.Add(10)
	for idx := 0 ; idx <10; idx++ {
		go func() {
			defer waitGroup.Done()
			//target function not invoked in place. rather a function literal is passed as arg to sync.Do
			syncOnce.Do(func(){counter++})
		}
	}
	waitGroup.Wait()

}