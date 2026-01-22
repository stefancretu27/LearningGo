package go_didactical_apps

/*
* 1. Channels:
*	- act like a thread-safe ( entail synchronized read-write), typed (all elements have same type) FIFO queue, used for goroutines communication
*	- allow to send and receive data between goroutines. The FIFO behavior is tangled with when multiple consumers/producers are involved, due to scheduling
*	- another way to perceived them: producer-consumer pipe. 
*	- a channel is reusable until it is closed. Only the sender must close it. After it is closed, another one can be allocated into that channel variable
*	- there is no need to close a channel unless it is desired to signal that no more values will be sent. Channels can remain open indefinitely
*	- the communication synchronization is implemented within channels structs using mutex and cond variable
*
* 2. Unbuffered channels
*	- declaration: var ubufChan chan
*	- declaration + allocation: unbufChan := make(chan int)
*	- only one value is in transit at a time
*	- the sender blocks until the receiver reads the value => tight coordination between goroutines, synchronous communication
*	- no explicit buffer is allocated, as they require a concurrent receiver
*	- indefinite for loop with if+ret status can be used to determine for a value's receiving
*		for {
			if data, ok := <-recvChan; ok {
*
* 3. Buffered channels
*	- no simple declaration is possible 
*	- declaration + allocation: unbufChan := make(chan int, bufSize)
*	- low coupling of producer and consumer => asynch communication, as the synchronization takes place when the buffer is full/empty
*	- the sender blocks when the buffer is full, waiting for at least a value to be read
*	- whereas the receiver blocks when the buffer is empty, thus waiting for a value to be sent
*	- range for can be used over receiving buffered channel, to read data until the sender explicitly closes it
*		for val := range recvChan -> value is received
*	- If the sender stops sending but doesn't close the channel, receivers will block when trying to read from an empty channel.
*	This can lead to goroutine leaks if receivers are waiting forever.
*
* 4. "select"
*	- is like a switch case for channels, allowing a goroutine to wait on multiple channels operations. Can be same channel, or multiple channels
*	- the select chooses only one case: the first one that is ready. If multiple cases are ready at the same time, one is chosen randomly.
*	- it blocks if no case is ready.
*	- use select with case for timeout ( <-time.After(timeout)) to avoid waiting forever (blocking). If value is available within that time, this case won't be executed
*	- use select with default case to avoid waiting forever (blocking). Executed if no other case is ready.
*/

// Show how a channel can be reused for bidirectional communication, in the same select, of a goroutine. 
func ReuseChan(reusedChan chan int, idx int) {
	duration := time.Second * time.Duration(idx)
	time.Sleep(duration)
	select {
	case val := <-reusedChan:
		fmt.Println("[ReuseChan]",idx, "Recv from reused chan:", val)
	case reusedChan <- idx:
		fmt.Println("[ReuseChan]", idx, "Sent to reused chan:", idx)
	}
}

// Close send chan to prevent receiver from blocking + range over bufchan upon receiving
func ProducerBufChan(sendChan chan<- int) {
	for idx := range cap(sendChan) {
		sendChan <- idx + 1
	}
	close(sendChan)
}

func ConsumerBufChan(recvChan <-chan int) {
	for val := range recvChan {
		fmt.Println("[ConsumerBufChan]:", val)
	}
}

// Use indefinite for loop, if and ret status upon receiving to determine whena  value is received on unbuf channel
func ProducerUnbufChan(sendChan chan<- int) {
	len := 7
	for idx := range len {
		sendChan <- idx
	}

	close(sendChan)
}

func ConsumerUnbufChan(recvChan <-chan int) {
	for {
		if data, ok := <-recvChan; ok {
			fmt.Println("[ConsumerUnbufChan]", data)
		} else {
			fmt.Println("[ConsumerUnbufChan] closed channel")
			return
		}
	}
}

func ProducersConsumerUnbufChan() {
	type Data struct {
		idx int
		val int
	}

	const producersNo int = 3
	const elementsNo int = 5
	unbufChan := make(chan Data)

	var wg sync.WaitGroup
	wg.Add(producersNo + 1)

	var goroutinesCounter atomic.Int32
	for idx := 0; idx < producersNo; idx++ {
		go func(index int) {
			defer wg.Done()
			for idy := 0; idy < elementsNo; idy++ {
				unbufChan <- Data{
					idx: idx,
					val: idy,
				}
			}

			fmt.Println("[ConsumerUnbufChan] Sent data in goroutine: ", index)
			goroutinesCounter.Add(1)
			if goroutinesCounter.Load() == int32(producersNo) {
				fmt.Println("[ConsumerUnbufChan] close channel")
				close(unbufChan)
			}
		}(idx)

	}

	go func() {
		defer wg.Done()
		for {
			select {
			case val, isOpen := <-unbufChan:
				if isOpen {
					fmt.Println("[ConsumerUnbufChan] val: ", val)
				} else {
					fmt.Println("[ConsumerUnbufChan] channel is closed")
					return
				}
			}
		}
	}()

	wg.Wait()
}

func Channels() {
	//channel declaration. At this point it can be buffered or unbuffered, as this is determined upon allocation time
	var unbufChan chan int
	unbufChan = make(chan int)

	// Unbuffered channel basic usage
	go func() {
		unbufChan <- 42
	}()

	val := <-unbufChan
	fmt.Println("Recv from unbufChan: ", val)

	// Buffered channels cannot be declared without allocating them
	bufChan := make(chan int, 5)

	go func() {
		for i := range 5 {
			bufChan <- i
		}
	}()

	for i := range bufChan {
		fmt.Println("Recv from bufChan", i)
	}

	// Reuse same channel for bidirectional data sharing. The timing of operations show one sends, one receives, but a gorotuine cannot both send and receive, this way.
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		ReuseChan(unbufChan, 0)
	}()
	go func() {
		defer wg.Done()
		ReuseChan(unbufChan, 1)
	}()

	wg.Wait()

	// Allocate new buffered channel in the existing variable. Use range for loops with bufChan capacity and receiver channel
	bufChan = make(chan int, 7)
	wg.Add(2)
	go func() {
		defer wg.Done()
		Producer(bufChan)
	}()
	go func() {
		defer wg.Done()
		Consumer(bufChan)
	}()
	wg.Wait()

	ProducersConsumerUnbufChan()

}