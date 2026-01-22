package go_didactical_apps

import (
	"fmt"
	"time"
	"sync"
)

/*
* Time package
*	- provides functionality for measuring and displaying time.
*	- the calendrical calculations always assume a Gregorian calendar, with no leap seconds.
*	- the functions calls of the API provided by the Time package rely on go runtime implementations, which, in their
* turn have OS specific implementations. Thus, the latter perform specific system calls to retrieve time
*	- that said, OSs provide 2 distinct types of clocks: monotonic (steady_clock) and wall clocks (system_clock).
*
* 	Wall clocks:
*		- are system clocks, whose value can be adjusted for clock synchronization
*		- used for telling time
*
* 	Monotonic clocks:
*		- are steady clocks, whose value cannot be adjusted
*		- used for measuring time
*
*	- the Time returned by time.Now contains both a wall clock reading and a monotonic clock reading. Therefater, the functions
* entailing time indication, called on the time instance use the wall clock, whereas the functions involving time measuring (such
* as additions, substractions, comparisons) use the monotonic clock of the t instance. Thus the API is not split.  
*	- eg: operations t.After(u), t.Before(u), t.Equal(u), and t.Sub(u) are carried out using the monotonic clock readings alone, 
* ignoring the wall clock readings, if both t and u embed monotonic clocks.
*	- t.AddDate(y, m, d), t.Round(d), and t.Truncate(d) are wall time computations, they always strip any monotonic clock reading from their results
*
* Types
* 1. Time
*	- represents an instant in time with nanosecond precision.
* 	- type Time struct
*	- even if it is a struct, programs should typically store and pass them as values, not pointers
*	- a Time value can be used by multiple goroutines simultaneously
*	- the zero value of type Time is January 1, year 1, 00:00:00.000000000 UTC
*	- each Time has associated with it a Location, consulted when computing the presentation form of the time, such as in the Format, Hour, and Year methods.
* Changing the location leads to changin the representation, but not the value.
*	- the Go == operator compares not just the time instant but also the Location and the monotonic clock reading => difficult to use as map/DB key
*	- func Now() Time: returns current local time
*	- func (t Time) After/Before(u Time) bool: indicates if time instant t is after/before u.
*	- func (t Time) Clock() (hour, min, sec int): returns the year, month, and day in which t occurs.
*	- func (t Time) Day/Hour/Minute/Second() int: returns the day/hour/minute/second of the month/day/hour/minute specified by t.
*	- func (t Time) String() string: returns the time formatted using the format string "2006-01-02 15:04:05.999999999 -0700 MST"
*	- func (t Time) Add/Sub(d Duration) Time: returns the time t+d/t-d
*
* 2. Duration
* 	- represents the elapsed time between two instants (Time instances) as an int64 nanosecond count.
*	- type Duration int64
*	- valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h". Can be used for suffixes
*	- func ParseDuration(s string) (Duration, error): parses a duration string, that is a possibly signed sequence of decimal numbers, 
* each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m"
*	- func Since(t Time) Duration: returns the time elapsed since t. It is shorthand for time.Now().Sub(t).
*	- func Until(t Time) Duration: returns the duration until t. It is shorthand for t.Sub(time.Now()).
*
* 3. Ticker
*	- holds a channel that sends `ticks` of a clock at intervals, speicifed as a duration upon creation, on a internal channel
*	- for receiving on that channel, it is listened repeatedly to in empty or range for loops
*	- type Ticker struct {
*    	C <-chan Time // The channel on which the ticks are delivered.
*		//other fields
*   } 
*	- func NewTicker(d Duration) *Ticker: returns pointer to a new Ticker containing a channel that will send the time with a period specified by the duration argument.
*	- eg: ticker := time.NewTicker(time.Second)
*		  defer ticker.Stop()
*		  for t := range ticker.C
*
* 4. Timer
*	- represents a single event. When the Timer expires, the current time will be sent on C, unless the Timer was created by AfterFunc. 
*	- a Timer must be created with NewTimer or AfterFunc.
*	- type Timer struct {
*    		C <-chan Time
*    		// contains filtered or unexported fields
*	 }
*	- func AfterFunc(d Duration, f func()) *Timer: waits for the duration to elapse and then calls f in its own goroutine. 
*												   It returns a Timer that can be used to cancel the call using its Stop method.
*	- func NewTimer(d Duration) *Timer: returns a pointer a new Timer that will send the current time on its channel after at least duration d.
*	- func (t *Timer) Stop() bool: prevents the Timer from firing. It returns true if the call stops the timer, false if the timer has already expired or been stopped. 
*								   Stop does not close the channel, to prevent a read from the channel succeeding incorrectly.
*
* Funcs:
* 1. func After(d Duration) <-chan Time: waits for the duration to elapse and then sends the current time on the returned (receiving) channel, of type Time.
* 										 It is equivalent to NewTimer(d).C. The underlying Timer is not recovered by the garbage collector until the timer fires.
*										 Lexx flexible than using the explicit timer, as it cannot be stopped or reset, but convenient when such operations aren't necessary
*										 eg: case <-time.After(10 * time.Second):
* 2. func Sleep(d Duration): Sleep pauses the current goroutine for at least the duration d
*							 eg: time.Sleep(100 * time.Millisecond)
* 3. func Tick(d Duration) <-chan Time: tick is a convenience wrapper for NewTicker providing access to the ticking channel only.
*										 Without a way to shut down the underlying Ticker, it cannot be recovered by the garbage collector; it "leaks"
*/

func TimePackage() {
		var wg sync.WaitGroup
	var duration time.Duration
	duration = 1 * time.Second
	syncChan := make(chan struct{}, 1)

	//Timeouts: time.Sleep and time.After functions publicly defined in time package
	wg.Add(2)
	go func() {
		defer wg.Done()
		time.Sleep(duration)
		fmt.Println("[timeout][producer] Sleep done")
		syncChan <- struct{}{}
	}()

	go func() {
		defer wg.Done()
		select {
		case <-time.After(duration):
			fmt.Println("[timeout][consumer] Timeout after: ", duration.String())
		case <-syncChan:
			fmt.Println("[timeout][consumer] Notified on syncChan")
		case <-time.After(duration * 2):
			fmt.Println("[timeout][consumer] Timeout after: ", (duration * 2).String())
		}
	}()
	wg.Wait()
	close(syncChan)

	//Ticker: repeatedly sends time over channel C, using the duration interval specified upon creation
	var tickDuration time.Duration
	tickDuration = 500 * time.Millisecond
	ticker := time.NewTicker(tickDuration)
	syncChan = make(chan struct{}, 1)
	wg.Add(2)
	go func(duration time.Duration) {
		defer wg.Done()
		time.Sleep(duration)
		fmt.Println("[ticker][producer] Sleep done", duration.String())
		syncChan <- struct{}{}
	}(duration * 2)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-syncChan:
				fmt.Println("[ticker][consumer] syncChan")
				return
			case t := <-ticker.C:
				fmt.Println("[ticker][consumer] ticker sent time on channel C:", t)
			case <-time.After(duration):
				fmt.Println("[ticker][consumer] timout after:", duration.String())
			}
		}
	}()
	wg.Wait()
	ticker.Stop()

	//Timer: send time over channel C, once, after the duration specified upon creation
	var timerDuration time.Duration
	timerDuration = 2 * time.Second
	timer := time.NewTimer(timerDuration)

	wg.Add(2)
	go func() {
		defer wg.Done()
		t := <-timer.C
		fmt.Println("{[timer] fired after:", t.String())
	}()
	go func() {
		defer wg.Done()
		t := <-time.After(timerDuration)
		fmt.Println("{[timer] timeout after:", t.String())
	}()
	wg.Wait()
	timer.Stop()
}