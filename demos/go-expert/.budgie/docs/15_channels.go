/*
METADATA:
Description: Demonstrates Go channels for communication between goroutines including buffered/unbuffered channels and patterns
Keywords: channel, make, send, receive, close, buffered-channel, unbuffered-channel, select, range
Category: concurrency
Concepts: channel creation, channel operations, channel direction, buffered vs unbuffered, channel patterns
*/

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Function demonstrating basic channel operations
func basicChannels() {
	fmt.Println("=== BASIC CHANNEL OPERATIONS ===")
	
	// Create unbuffered channel
	ch := make(chan string)
	
	// Send and receive in separate goroutines
	go func() {
		ch <- "Hello, Channel!"
	}()
	
	message := <-ch
	fmt.Printf("Received: %s\n", message)
	
	// Create buffered channel
	bufferedCh := make(chan int, 3)
	
	// Send multiple values without blocking
	bufferedCh <- 1
	bufferedCh <- 2
	bufferedCh <- 3
	
	// Receive values
	fmt.Printf("Buffered channel values: %d, %d, %d\n", <-bufferedCh, <-bufferedCh, <-bufferedCh)
}

// Function demonstrating channel directions
func channelDirections() {
	fmt.Println("\n=== CHANNEL DIRECTIONS ===")
	
	// Function that only sends
	sender := func(ch chan<- string, messages []string) {
		for _, msg := range messages {
			ch <- msg
		}
		close(ch)
	}
	
	// Function that only receives
	receiver := func(ch <-chan string) {
		for msg := range ch {
			fmt.Printf("Received: %s\n", msg)
		}
	}
	
	ch := make(chan string)
	messages := []string{"Hello", "World", "From", "Channels"}
	
	go sender(ch, messages)
	receiver(ch)
}

// Function demonstrating channel closing
func channelClosing() {
	fmt.Println("\n=== CHANNEL CLOSING ===")
	
	ch := make(chan int, 5)
	
	// Send some values
	go func() {
		for i := 1; i <= 5; i++ {
			ch <- i
			time.Sleep(100 * time.Millisecond)
		}
		close(ch) // Signal that no more values will be sent
	}()
	
	// Receive until channel is closed
	for {
		value, ok := <-ch
		if !ok {
			fmt.Println("Channel is closed")
			break
		}
		fmt.Printf("Received: %d\n", value)
	}
	
	// Alternative: range over channel
	fmt.Println("\nUsing range over channel:")
	ch2 := make(chan string, 3)
	
	go func() {
		defer close(ch2)
		for _, fruit := range []string{"apple", "banana", "cherry"} {
			ch2 <- fruit
		}
	}()
	
	for fruit := range ch2 {
		fmt.Printf("Fruit: %s\n", fruit)
	}
}

// Function demonstrating select statement
func selectStatement() {
	fmt.Println("\n=== SELECT STATEMENT ===")
	
	ch1 := make(chan string)
	ch2 := make(chan string)
	
	// Send to channels at different times
	go func() {
		time.Sleep(200 * time.Millisecond)
		ch1 <- "Channel 1"
	}()
	
	go func() {
		time.Sleep(100 * time.Millisecond)
		ch2 <- "Channel 2"
	}()
	
	// Select receives from whichever channel is ready first
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Printf("Received from ch1: %s\n", msg1)
		case msg2 := <-ch2:
			fmt.Printf("Received from ch2: %s\n", msg2)
		}
	}
}

// Function demonstrating select with default case
func selectWithDefault() {
	fmt.Println("\n=== SELECT WITH DEFAULT ===")
	
	ch := make(chan string)
	
	// Non-blocking receive
	select {
	case msg := <-ch:
		fmt.Printf("Received: %s\n", msg)
	default:
		fmt.Println("No message received, continuing...")
	}
	
	// Non-blocking send
	select {
	case ch <- "Hello":
		fmt.Println("Sent message")
	default:
		fmt.Println("Could not send message, channel not ready")
	}
}

// Function demonstrating timeout with select
func selectTimeout() {
	fmt.Println("\n=== SELECT WITH TIMEOUT ===")
	
	ch := make(chan string)
	
	// Simulate slow operation
	go func() {
		time.Sleep(2 * time.Second)
		ch <- "Slow operation result"
	}()
	
	// Timeout after 1 second
	select {
	case result := <-ch:
		fmt.Printf("Received: %s\n", result)
	case <-time.After(1 * time.Second):
		fmt.Println("Operation timed out")
	}
}

// Function demonstrating worker pools with channels
func workerPoolWithChannels() {
	fmt.Println("\n=== WORKER POOL WITH CHANNELS ===")
	
	const numWorkers = 3
	const numJobs = 10
	
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	
	// Start workers
	for w := 1; w <= numWorkers; w++ {
		go func(workerID int) {
			for job := range jobs {
				fmt.Printf("Worker %d processing job %d\n", workerID, job)
				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
				results <- job * job
			}
		}(w)
	}
	
	// Send jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)
	
	// Collect results
	for r := 1; r <= numJobs; r++ {
		result := <-results
		fmt.Printf("Result: %d\n", result)
	}
}

// Function demonstrating ping-pong pattern
func pingPong() {
	fmt.Println("\n=== PING-PONG PATTERN ===")
	
	ping := make(chan string)
	pong := make(chan string)
	
	// Ping goroutine
	go func() {
		for i := 0; i < 3; i++ {
			ping <- "ping"
			response := <-pong
			fmt.Printf("Ping received: %s\n", response)
		}
		close(ping)
	}()
	
	// Pong goroutine
	go func() {
		for message := range ping {
			fmt.Printf("Pong received: %s\n", message)
			pong <- "pong"
		}
		close(pong)
	}()
	
	// Wait for completion
	time.Sleep(500 * time.Millisecond)
}

// Function demonstrating channel pipeline
func channelPipeline() {
	fmt.Println("\n=== CHANNEL PIPELINE ===")
	
	// Stage 1: Generate numbers
	numbers := make(chan int)
	go func() {
		defer close(numbers)
		for i := 1; i <= 5; i++ {
			numbers <- i
		}
	}()
	
	// Stage 2: Square the numbers
	squares := make(chan int)
	go func() {
		defer close(squares)
		for num := range numbers {
			squares <- num * num
		}
	}()
	
	// Stage 3: Format as strings
	formatted := make(chan string)
	go func() {
		defer close(formatted)
		for square := range squares {
			formatted <- fmt.Sprintf("Square: %d", square)
		}
	}()
	
	// Final stage: Print results
	for result := range formatted {
		fmt.Println(result)
	}
}

// Function demonstrating fan-in with channels
func fanInChannels() {
	fmt.Println("\n=== FAN-IN WITH CHANNELS ===")
	
	// Create input channels
	input1 := make(chan string)
	input2 := make(chan string)
	
	// Fan-in function
	fanIn := func(ch1, ch2 <-chan string) <-chan string {
		out := make(chan string)
		go func() {
			defer close(out)
			for {
				select {
				case msg, ok := <-ch1:
					if !ok {
						ch1 = nil
					} else {
						out <- msg
					}
				case msg, ok := <-ch2:
					if !ok {
						ch2 = nil
					} else {
						out <- msg
					}
				}
				if ch1 == nil && ch2 == nil {
					break
				}
			}
		}()
		return out
	}
	
	// Start input goroutines
	go func() {
		defer close(input1)
		for i := 1; i <= 3; i++ {
			input1 <- fmt.Sprintf("Input1-%d", i)
			time.Sleep(100 * time.Millisecond)
		}
	}()
	
	go func() {
		defer close(input2)
		for i := 1; i <= 3; i++ {
			input2 <- fmt.Sprintf("Input2-%d", i)
			time.Sleep(150 * time.Millisecond)
		}
	}()
	
	// Receive from fan-in
	output := fanIn(input1, input2)
	for msg := range output {
		fmt.Printf("Fan-in output: %s\n", msg)
	}
}

// Function demonstrating channel synchronization
func channelSynchronization() {
	fmt.Println("\n=== CHANNEL SYNCHRONIZATION ===")
	
	done := make(chan bool)
	
	// Worker goroutine
	go func() {
		fmt.Println("Worker: Starting work...")
		time.Sleep(500 * time.Millisecond)
		fmt.Println("Worker: Work completed")
		done <- true
	}()
	
	// Wait for worker to complete
	fmt.Println("Main: Waiting for worker...")
	<-done
	fmt.Println("Main: Worker completed, continuing...")
}

// Function demonstrating rate limiting with channels
func rateLimiting() {
	fmt.Println("\n=== RATE LIMITING ===")
	
	// Create rate limiter
	rate := time.Tick(200 * time.Millisecond)
	
	// Process requests with rate limiting
	requests := []string{"req1", "req2", "req3", "req4", "req5"}
	
	for _, req := range requests {
		<-rate // Wait for rate limiter
		fmt.Printf("Processing request: %s at %s\n", req, time.Now().Format("15:04:05.000"))
	}
}

// Function demonstrating channel buffering effects
func bufferingEffects() {
	fmt.Println("\n=== BUFFERING EFFECTS ===")
	
	// Unbuffered channel
	fmt.Println("Unbuffered channel:")
	unbuffered := make(chan int)
	
	start := time.Now()
	go func() {
		for i := 1; i <= 3; i++ {
			fmt.Printf("Sending %d at %v\n", i, time.Since(start))
			unbuffered <- i
		}
		close(unbuffered)
	}()
	
	time.Sleep(100 * time.Millisecond) // Delay receiving
	for value := range unbuffered {
		fmt.Printf("Received %d at %v\n", value, time.Since(start))
		time.Sleep(50 * time.Millisecond)
	}
	
	// Buffered channel
	fmt.Println("\nBuffered channel:")
	buffered := make(chan int, 3)
	
	start = time.Now()
	go func() {
		for i := 1; i <= 3; i++ {
			fmt.Printf("Sending %d at %v\n", i, time.Since(start))
			buffered <- i
		}
		close(buffered)
	}()
	
	time.Sleep(200 * time.Millisecond) // Delay receiving
	for value := range buffered {
		fmt.Printf("Received %d at %v\n", value, time.Since(start))
		time.Sleep(50 * time.Millisecond)
	}
}

// Function demonstrating channel as semaphore
func channelSemaphore() {
	fmt.Println("\n=== CHANNEL AS SEMAPHORE ===")
	
	// Semaphore with capacity of 2
	semaphore := make(chan struct{}, 2)
	var wg sync.WaitGroup
	
	// Start 5 goroutines, but only 2 can run concurrently
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }() // Release semaphore
			
			fmt.Printf("Goroutine %d: acquired semaphore\n", id)
			time.Sleep(500 * time.Millisecond)
			fmt.Printf("Goroutine %d: releasing semaphore\n", id)
		}(i)
	}
	
	wg.Wait()
	fmt.Println("All goroutines completed")
}

func main() {
	basicChannels()
	channelDirections()
	channelClosing()
	selectStatement()
	selectWithDefault()
	selectTimeout()
	workerPoolWithChannels()
	pingPong()
	channelPipeline()
	fanInChannels()
	channelSynchronization()
	rateLimiting()
	bufferingEffects()
	channelSemaphore()
	
	fmt.Println("\n=== CHANNEL BEST PRACTICES ===")
	fmt.Println("1. Close channels when no more values will be sent")
	fmt.Println("2. Only the sender should close the channel")
	fmt.Println("3. Use buffered channels to prevent goroutine blocking")
	fmt.Println("4. Use select for non-blocking operations")
	fmt.Println("5. Use range to receive all values from a channel")
	fmt.Println("6. Check channel status with comma ok idiom")
	fmt.Println("7. Use channels for synchronization, not just data passing")
	
	fmt.Println("\nChannel examples completed!")
}