/*
METADATA:
Description: Demonstrates Go goroutines for concurrent programming including creation, synchronization, and communication
Keywords: goroutine, concurrency, go-keyword, sync, waitgroup, mutex, atomic, parallel
Category: concurrency
Concepts: goroutine creation, concurrent execution, synchronization, race conditions, parallel processing
*/

package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// Simple goroutine function
func simpleTask(id int) {
	for i := 0; i < 3; i++ {
		fmt.Printf("Task %d: iteration %d\n", id, i+1)
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Printf("Task %d completed\n", id)
}

// Function with work simulation
func worker(id int, work int) {
	fmt.Printf("Worker %d starting\n", id)
	
	// Simulate work
	time.Sleep(time.Duration(work) * time.Millisecond)
	
	fmt.Printf("Worker %d finished work (%dms)\n", id, work)
}

// Function demonstrating race condition
func unsafeCounter() {
	var counter int
	var wg sync.WaitGroup
	
	// Start 10 goroutines that increment counter
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				counter++ // Race condition!
			}
			fmt.Printf("Goroutine %d finished\n", id)
		}(i)
	}
	
	wg.Wait()
	fmt.Printf("Unsafe counter final value: %d (should be 10000)\n", counter)
}

// Function demonstrating safe counter with mutex
func safeCounterWithMutex() {
	var counter int
	var mu sync.Mutex
	var wg sync.WaitGroup
	
	// Start 10 goroutines that increment counter safely
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				mu.Lock()
				counter++
				mu.Unlock()
			}
			fmt.Printf("Goroutine %d finished\n", id)
		}(i)
	}
	
	wg.Wait()
	fmt.Printf("Safe counter (mutex) final value: %d\n", counter)
}

// Function demonstrating atomic operations
func safeCounterWithAtomic() {
	var counter int64
	var wg sync.WaitGroup
	
	// Start 10 goroutines that increment counter atomically
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				atomic.AddInt64(&counter, 1)
			}
			fmt.Printf("Goroutine %d finished\n", id)
		}(i)
	}
	
	wg.Wait()
	fmt.Printf("Safe counter (atomic) final value: %d\n", counter)
}

// Producer-consumer pattern
func producerConsumer() {
	var wg sync.WaitGroup
	items := make(chan int, 5) // Buffered channel
	
	// Producer goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(items)
		
		for i := 1; i <= 10; i++ {
			fmt.Printf("Producing item %d\n", i)
			items <- i
			time.Sleep(50 * time.Millisecond)
		}
		fmt.Println("Producer finished")
	}()
	
	// Consumer goroutines
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(consumerID int) {
			defer wg.Done()
			
			for item := range items {
				fmt.Printf("Consumer %d consumed item %d\n", consumerID, item)
				time.Sleep(100 * time.Millisecond)
			}
			fmt.Printf("Consumer %d finished\n", consumerID)
		}(i)
	}
	
	wg.Wait()
}

// Worker pool pattern
func workerPool() {
	const numWorkers = 3
	const numJobs = 10
	
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	
	// Start workers
	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range jobs {
				fmt.Printf("Worker %d processing job %d\n", workerID, job)
				
				// Simulate work
				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
				result := job * 2
				
				fmt.Printf("Worker %d finished job %d, result: %d\n", workerID, job, result)
				results <- result
			}
		}(w)
	}
	
	// Send jobs
	go func() {
		for j := 1; j <= numJobs; j++ {
			jobs <- j
		}
		close(jobs)
	}()
	
	// Collect results
	go func() {
		wg.Wait()
		close(results)
	}()
	
	// Print results
	for result := range results {
		fmt.Printf("Received result: %d\n", result)
	}
}

// Fan-in pattern
func fanIn() {
	// Create multiple input channels
	input1 := make(chan string)
	input2 := make(chan string)
	input3 := make(chan string)
	
	// Output channel
	output := make(chan string)
	
	// Fan-in function
	fanInFunc := func(inputs ...chan string) chan string {
		var wg sync.WaitGroup
		out := make(chan string)
		
		// Start a goroutine for each input channel
		for _, input := range inputs {
			wg.Add(1)
			go func(ch chan string) {
				defer wg.Done()
				for value := range ch {
					out <- value
				}
			}(input)
		}
		
		// Close output when all inputs are done
		go func() {
			wg.Wait()
			close(out)
		}()
		
		return out
	}
	
	// Start the fan-in
	output = fanInFunc(input1, input2, input3)
	
	// Send data to input channels
	go func() {
		for i := 1; i <= 3; i++ {
			input1 <- fmt.Sprintf("Input1-%d", i)
			time.Sleep(100 * time.Millisecond)
		}
		close(input1)
	}()
	
	go func() {
		for i := 1; i <= 3; i++ {
			input2 <- fmt.Sprintf("Input2-%d", i)
			time.Sleep(150 * time.Millisecond)
		}
		close(input2)
	}()
	
	go func() {
		for i := 1; i <= 3; i++ {
			input3 <- fmt.Sprintf("Input3-%d", i)
			time.Sleep(200 * time.Millisecond)
		}
		close(input3)
	}()
	
	// Receive from output
	for value := range output {
		fmt.Printf("Fan-in received: %s\n", value)
	}
}

// Fan-out pattern
func fanOut() {
	input := make(chan int)
	
	// Create multiple output channels
	output1 := make(chan int)
	output2 := make(chan int)
	output3 := make(chan int)
	
	// Fan-out function
	go func() {
		defer close(output1)
		defer close(output2)
		defer close(output3)
		
		for value := range input {
			// Send to all outputs
			output1 <- value
			output2 <- value
			output3 <- value
		}
	}()
	
	// Start consumers
	var wg sync.WaitGroup
	
	wg.Add(1)
	go func() {
		defer wg.Done()
		for value := range output1 {
			fmt.Printf("Consumer 1 received: %d\n", value)
			time.Sleep(50 * time.Millisecond)
		}
	}()
	
	wg.Add(1)
	go func() {
		defer wg.Done()
		for value := range output2 {
			fmt.Printf("Consumer 2 received: %d\n", value)
			time.Sleep(100 * time.Millisecond)
		}
	}()
	
	wg.Add(1)
	go func() {
		defer wg.Done()
		for value := range output3 {
			fmt.Printf("Consumer 3 received: %d\n", value)
			time.Sleep(150 * time.Millisecond)
		}
	}()
	
	// Send data to input
	go func() {
		defer close(input)
		for i := 1; i <= 5; i++ {
			fmt.Printf("Sending: %d\n", i)
			input <- i
			time.Sleep(200 * time.Millisecond)
		}
	}()
	
	wg.Wait()
}

// Timeout pattern
func timeoutPattern() {
	work := make(chan string)
	
	// Worker that might take too long
	go func() {
		time.Sleep(2 * time.Second) // Simulate slow work
		work <- "work completed"
	}()
	
	// Timeout after 1 second
	select {
	case result := <-work:
		fmt.Printf("Work result: %s\n", result)
	case <-time.After(1 * time.Second):
		fmt.Println("Work timed out after 1 second")
	}
}

// Graceful shutdown pattern
func gracefulShutdown() {
	done := make(chan bool)
	quit := make(chan bool)
	
	// Worker goroutine
	go func() {
		for {
			select {
			case <-quit:
				fmt.Println("Worker received quit signal")
				done <- true
				return
			default:
				fmt.Println("Worker is working...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()
	
	// Simulate running for 2 seconds, then shutdown
	time.Sleep(2 * time.Second)
	fmt.Println("Sending quit signal...")
	quit <- true
	
	// Wait for worker to finish
	<-done
	fmt.Println("Worker finished gracefully")
}

func main() {
	fmt.Printf("Number of CPUs: %d\n", runtime.NumCPU())
	fmt.Printf("Number of goroutines at start: %d\n", runtime.NumGoroutine())
	
	fmt.Println("\n=== BASIC GOROUTINES ===")
	
	// Sequential execution
	fmt.Println("Sequential execution:")
	start := time.Now()
	for i := 1; i <= 3; i++ {
		simpleTask(i)
	}
	fmt.Printf("Sequential time: %v\n", time.Since(start))
	
	// Concurrent execution
	fmt.Println("\nConcurrent execution:")
	start = time.Now()
	var wg sync.WaitGroup
	
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			simpleTask(id)
		}(i)
	}
	
	wg.Wait()
	fmt.Printf("Concurrent time: %v\n", time.Since(start))
	fmt.Printf("Number of goroutines after basic demo: %d\n", runtime.NumGoroutine())

	fmt.Println("\n=== WORKER GOROUTINES ===")
	
	var wg2 sync.WaitGroup
	workloads := []int{200, 100, 300, 150, 250}
	
	for i, work := range workloads {
		wg2.Add(1)
		go func(id, workTime int) {
			defer wg2.Done()
			worker(id, workTime)
		}(i+1, work)
	}
	
	wg2.Wait()
	fmt.Println("All workers completed")

	fmt.Println("\n=== RACE CONDITION DEMONSTRATION ===")
	
	fmt.Println("Unsafe counter (race condition):")
	unsafeCounter()
	
	fmt.Println("\nSafe counter with mutex:")
	safeCounterWithMutex()
	
	fmt.Println("\nSafe counter with atomic operations:")
	safeCounterWithAtomic()

	fmt.Println("\n=== PRODUCER-CONSUMER PATTERN ===")
	
	producerConsumer()

	fmt.Println("\n=== WORKER POOL PATTERN ===")
	
	workerPool()

	fmt.Println("\n=== FAN-IN PATTERN ===")
	
	fanIn()

	fmt.Println("\n=== FAN-OUT PATTERN ===")
	
	fanOut()

	fmt.Println("\n=== TIMEOUT PATTERN ===")
	
	timeoutPattern()

	fmt.Println("\n=== GRACEFUL SHUTDOWN PATTERN ===")
	
	gracefulShutdown()

	fmt.Println("\n=== GOROUTINE LEAKS PREVENTION ===")
	
	// Example of potential goroutine leak and how to prevent it
	preventLeak := func() {
		ch := make(chan int)
		done := make(chan bool)
		
		// Goroutine that could leak
		go func() {
			defer func() {
				done <- true
			}()
			
			select {
			case value := <-ch:
				fmt.Printf("Received value: %d\n", value)
			case <-time.After(1 * time.Second):
				fmt.Println("Goroutine timed out, preventing leak")
			}
		}()
		
		// Don't send anything to ch, causing potential leak
		// ch <- 42  // Uncomment to prevent timeout
		
		// Wait for goroutine to finish (with timeout)
		select {
		case <-done:
			fmt.Println("Goroutine finished properly")
		case <-time.After(2 * time.Second):
			fmt.Println("Goroutine cleanup timed out")
		}
	}
	
	preventLeak()

	fmt.Println("\n=== RUNTIME GOROUTINE INFORMATION ===")
	
	fmt.Printf("Number of goroutines at end: %d\n", runtime.NumGoroutine())
	fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))
	
	// Force garbage collection to clean up
	runtime.GC()
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Number of goroutines after GC: %d\n", runtime.NumGoroutine())

	fmt.Println("\n=== ANONYMOUS GOROUTINES ===")
	
	// Anonymous goroutine
	go func() {
		fmt.Println("Anonymous goroutine executing")
		time.Sleep(100 * time.Millisecond)
		fmt.Println("Anonymous goroutine finished")
	}()
	
	// Wait for anonymous goroutine
	time.Sleep(200 * time.Millisecond)

	fmt.Println("\n=== GOROUTINE WITH CLOSURE ===")
	
	message := "Hello from closure"
	var wg3 sync.WaitGroup
	
	wg3.Add(1)
	go func() {
		defer wg3.Done()
		fmt.Printf("Goroutine with closure: %s\n", message)
		
		// Modify the captured variable
		message = "Modified by goroutine"
	}()
	
	wg3.Wait()
	fmt.Printf("Message after goroutine: %s\n", message)

	fmt.Println("\nGoroutine examples completed!")
}