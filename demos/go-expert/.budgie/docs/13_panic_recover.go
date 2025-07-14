/*
METADATA:
Description: Demonstrates Go panic and recover mechanisms for handling exceptional situations
Keywords: panic, recover, defer, stack-trace, exceptional-handling, graceful-degradation
Category: error-handling
Concepts: panic function, recover function, defer in panic scenarios, graceful error handling
*/

package main

import (
	"fmt"
	"runtime"
)

// Function that panics
func riskyFunction(shouldPanic bool) {
	if shouldPanic {
		panic("something terrible happened!")
	}
	fmt.Println("Function completed successfully")
}

// Function with defer and panic
func functionWithDefer() {
	defer fmt.Println("Defer 1: This will execute even if panic occurs")
	defer fmt.Println("Defer 2: This will also execute")
	
	fmt.Println("About to panic...")
	panic("intentional panic")
	
	fmt.Println("This line will never execute")
}

// Function that recovers from panic
func safeFunction() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
		}
	}()
	
	fmt.Println("Starting risky operation...")
	panic("oops!")
	fmt.Println("This won't be reached")
}

// Function that recovers and returns error
func safeFunctionWithError() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("function panicked: %v", r)
		}
	}()
	
	panic("something went wrong")
	return nil
}

// Function that demonstrates selective recovery
func selectiveRecovery(input string) (result string, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch v := r.(type) {
			case string:
				if v == "expected panic" {
					err = fmt.Errorf("handled expected panic: %s", v)
				} else {
					panic(r) // re-panic if unexpected
				}
			case error:
				err = fmt.Errorf("recovered from error panic: %v", v)
			default:
				panic(r) // re-panic for unknown types
			}
		}
	}()
	
	switch input {
	case "expected":
		panic("expected panic")
	case "error":
		panic(fmt.Errorf("error type panic"))
	case "unexpected":
		panic(42) // This will re-panic
	default:
		return "success", nil
	}
}

// Function that demonstrates stack unwinding
func stackUnwindingDemo() {
	defer fmt.Println("Level 0 defer")
	level1()
	fmt.Println("This won't execute")
}

func level1() {
	defer fmt.Println("Level 1 defer")
	level2()
	fmt.Println("This won't execute")
}

func level2() {
	defer fmt.Println("Level 2 defer")
	level3()
	fmt.Println("This won't execute")
}

func level3() {
	defer fmt.Println("Level 3 defer")
	panic("panic from level 3")
}

// Function that recovers and provides stack trace
func functionWithStackTrace() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic: %v\n", r)
			
			// Print stack trace
			buf := make([]byte, 1024)
			stack := runtime.Stack(buf, false)
			fmt.Printf("Stack trace:\n%s", stack)
		}
	}()
	
	deepFunction()
}

func deepFunction() {
	evenDeeperFunction()
}

func evenDeeperFunction() {
	panic("deep panic")
}

// Function demonstrating panic in goroutine
func panicInGoroutine() {
	// Note: This is just for demonstration
	// In real code, panics in goroutines will crash the program
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered in main goroutine: %v\n", r)
		}
	}()
	
	// This recover won't catch panic from goroutine
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovered in goroutine: %v\n", r)
			}
		}()
		panic("panic in goroutine")
	}()
	
	// Give goroutine time to execute
	fmt.Println("Main goroutine continuing...")
}

// Resource management with panic
func resourceManagement() {
	resource := "database connection"
	fmt.Printf("Acquiring resource: %s\n", resource)
	
	defer func() {
		fmt.Printf("Releasing resource: %s\n", resource)
		if r := recover(); r != nil {
			fmt.Printf("Panic during resource usage: %v\n", r)
		}
	}()
	
	// Simulate work that might panic
	doRiskyWork(true)
}

func doRiskyWork(shouldPanic bool) {
	fmt.Println("Doing risky work...")
	if shouldPanic {
		panic("work failed")
	}
	fmt.Println("Work completed successfully")
}

// Function that demonstrates when NOT to use recover
func badRecoveryExample() {
	defer func() {
		// Bad: Recovering from everything without proper handling
		if r := recover(); r != nil {
			fmt.Printf("Something went wrong, but continuing anyway: %v\n", r)
		}
	}()
	
	// This would hide a real programming error
	var slice []int
	fmt.Println(slice[10]) // This should panic (index out of range)
}

// Function that demonstrates proper error handling vs panic
func properErrorHandling(data []int, index int) (int, error) {
	if index < 0 || index >= len(data) {
		return 0, fmt.Errorf("index %d out of range for slice of length %d", index, len(data))
	}
	return data[index], nil
}

func panicVersion(data []int, index int) int {
	if index < 0 || index >= len(data) {
		panic(fmt.Sprintf("index %d out of range", index))
	}
	return data[index]
}

// Custom panic recovery middleware
func withPanicRecovery(fn func()) func() error {
	return func() error {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Middleware recovered from panic: %v\n", r)
			}
		}()
		
		fn()
		return nil
	}
}

func main() {
	fmt.Println("=== BASIC PANIC DEMONSTRATION ===")
	
	// This will panic and terminate if not recovered
	fmt.Println("Calling risky function with safe=true:")
	riskyFunction(false)
	
	fmt.Println("\nCalling risky function with safe=false (would panic):")
	// riskyFunction(true) // Uncomment to see panic

	fmt.Println("\n=== DEFER WITH PANIC ===")
	
	// Demonstrate that defer runs even with panic
	fmt.Println("Calling function with defer and panic:")
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Main recovered: %v\n", r)
		}
	}()
	
	// Uncomment to see defers execute before panic propagation
	// functionWithDefer()

	fmt.Println("\n=== BASIC RECOVERY ===")
	
	fmt.Println("Calling safe function:")
	safeFunction()
	fmt.Println("Program continues after recovery")

	fmt.Println("\n=== RECOVERY WITH ERROR RETURN ===")
	
	err := safeFunctionWithError()
	if err != nil {
		fmt.Printf("Function returned error: %v\n", err)
	}

	fmt.Println("\n=== SELECTIVE RECOVERY ===")
	
	testInputs := []string{"success", "expected", "error"}
	
	for _, input := range testInputs {
		fmt.Printf("Testing input '%s':\n", input)
		result, err := selectiveRecovery(input)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
		} else {
			fmt.Printf("  Result: %s\n", result)
		}
	}
	
	// This would re-panic and crash the program
	fmt.Println("Testing unexpected panic (would re-panic):")
	// result, err := selectiveRecovery("unexpected")

	fmt.Println("\n=== STACK UNWINDING ===")
	
	fmt.Println("Demonstrating stack unwinding:")
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered at main level: %v\n", r)
		}
	}()
	stackUnwindingDemo()

	fmt.Println("\n=== STACK TRACE ===")
	
	fmt.Println("Function with stack trace:")
	functionWithStackTrace()

	fmt.Println("\n=== PANIC IN GOROUTINE ===")
	
	fmt.Println("Panic in goroutine:")
	panicInGoroutine()
	
	// Sleep to let goroutine complete
	import "time"
	time.Sleep(100 * time.Millisecond)

	fmt.Println("\n=== RESOURCE MANAGEMENT ===")
	
	fmt.Println("Resource management with panic:")
	resourceManagement()

	fmt.Println("\n=== ERROR HANDLING VS PANIC ===")
	
	data := []int{1, 2, 3, 4, 5}
	
	// Proper error handling
	fmt.Println("Using proper error handling:")
	value, err := properErrorHandling(data, 10)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Value: %d\n", value)
	}
	
	// Using panic (with recovery)
	fmt.Println("Using panic version:")
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovered from panic: %v\n", r)
			}
		}()
		value := panicVersion(data, 10)
		fmt.Printf("Value: %d\n", value)
	}()

	fmt.Println("\n=== PANIC RECOVERY MIDDLEWARE ===")
	
	riskyOperation := func() {
		fmt.Println("Doing risky operation...")
		panic("operation failed")
	}
	
	safeOperation := withPanicRecovery(riskyOperation)
	err = safeOperation()
	if err != nil {
		fmt.Printf("Operation error: %v\n", err)
	}

	fmt.Println("\n=== BEST PRACTICES ===")
	
	fmt.Println("Best practices for panic/recover:")
	fmt.Println("1. Use panic for truly exceptional situations")
	fmt.Println("2. Prefer returning errors for expected failure cases")
	fmt.Println("3. Use recover to prevent program termination")
	fmt.Println("4. Don't use recover to ignore programming errors")
	fmt.Println("5. Always recover in the same goroutine as the panic")
	fmt.Println("6. Use defer for cleanup that must happen")
	
	// Example of when to use panic vs error
	fmt.Println("\nWhen to use panic:")
	fmt.Println("- Programming errors (array bounds, nil pointer)")
	fmt.Println("- Unrecoverable situations (out of memory)")
	fmt.Println("- Initialization failures in init() functions")
	
	fmt.Println("\nWhen to use errors:")
	fmt.Println("- Expected failure conditions")
	fmt.Println("- Input validation failures")
	fmt.Println("- Network timeouts")
	fmt.Println("- File not found")

	fmt.Println("\n=== PANIC WITH DIFFERENT TYPES ===")
	
	// Panic can be called with any type
	panicTypes := []interface{}{
		"string panic",
		42,
		fmt.Errorf("error panic"),
		[]string{"slice", "panic"},
	}
	
	for i, panicValue := range panicTypes {
		fmt.Printf("Testing panic type %d:\n", i+1)
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("  Recovered %T: %v\n", r, r)
				}
			}()
			panic(panicValue)
		}()
	}
}