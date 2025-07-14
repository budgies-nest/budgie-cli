/*
METADATA:
Description: Demonstrates advanced Go function concepts including closures, higher-order functions, defer, and function types
Keywords: closure, higher-order-function, defer, function-type, callback, first-class-function, anonymous-function
Category: functions
Concepts: closures, function variables, defer statement, function composition, callbacks
*/

package main

import (
	"fmt"
	"strings"
)

// Function type declaration
type MathOperation func(int, int) int
type StringProcessor func(string) string

// Higher-order function that takes a function as parameter
func calculate(a, b int, operation MathOperation) int {
	return operation(a, b)
}

// Function that returns a function
func createValidator(minLength int) func(string) bool {
	return func(s string) bool {
		return len(s) >= minLength
	}
}

// Closure that captures and modifies external variables
func createCounter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

// Function that demonstrates defer
func deferExample() {
	fmt.Println("Function start")
	defer fmt.Println("Deferred: This prints last")
	defer fmt.Println("Deferred: This prints second to last")
	
	fmt.Println("Function middle")
	fmt.Println("Function end")
}

// Function with defer and loop
func deferInLoop() {
	fmt.Println("Defer in loop example:")
	for i := 1; i <= 3; i++ {
		defer fmt.Printf("Deferred: %d\n", i)
	}
	fmt.Println("Loop finished")
}

// Function demonstrating defer with function parameters
func deferWithParams() {
	x := 10
	defer func(val int) {
		fmt.Printf("Deferred with captured value: %d\n", val)
	}(x)
	
	x = 20
	defer func() {
		fmt.Printf("Deferred with closure value: %d\n", x)
	}()
	
	fmt.Println("Function body executed")
}

// Higher-order function for slice processing
func processSlice(slice []int, processor func(int) int) []int {
	result := make([]int, len(slice))
	for i, v := range slice {
		result[i] = processor(v)
	}
	return result
}

// Function that takes multiple function parameters
func chainProcessors(input string, processors ...StringProcessor) string {
	result := input
	for _, processor := range processors {
		result = processor(result)
	}
	return result
}

// Method that returns multiple functions
func createMathOperations() (MathOperation, MathOperation, MathOperation, MathOperation) {
	add := func(a, b int) int { return a + b }
	subtract := func(a, b int) int { return a - b }
	multiply := func(a, b int) int { return a * b }
	divide := func(a, b int) int {
		if b != 0 {
			return a / b
		}
		return 0
	}
	return add, subtract, multiply, divide
}

// Decorator pattern with functions
func withLogging(fn func(string) string) func(string) string {
	return func(input string) string {
		fmt.Printf("Processing: %s\n", input)
		result := fn(input)
		fmt.Printf("Result: %s\n", result)
		return result
	}
}

// Memoization with closures
func memoize(fn func(int) int) func(int) int {
	cache := make(map[int]int)
	return func(x int) int {
		if result, exists := cache[x]; exists {
			fmt.Printf("Cache hit for %d\n", x)
			return result
		}
		result := fn(x)
		cache[x] = result
		fmt.Printf("Cache miss for %d, computed %d\n", x, result)
		return result
	}
}

// Expensive function for memoization demo
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func main() {
	// Function variables
	var addOp MathOperation = func(a, b int) int { return a + b }
	var multiplyOp MathOperation = func(a, b int) int { return a * b }

	// Using higher-order function
	result1 := calculate(5, 3, addOp)
	result2 := calculate(5, 3, multiplyOp)
	fmt.Printf("Addition: %d, Multiplication: %d\n", result1, result2)

	// Creating and using validators
	emailValidator := createValidator(5)
	passwordValidator := createValidator(8)
	
	fmt.Printf("'test' is valid email: %t\n", emailValidator("test"))
	fmt.Printf("'test@example.com' is valid email: %t\n", emailValidator("test@example.com"))
	fmt.Printf("'123' is valid password: %t\n", passwordValidator("123"))
	fmt.Printf("'password123' is valid password: %t\n", passwordValidator("password123"))

	// Counter closure
	counter1 := createCounter()
	counter2 := createCounter()
	
	fmt.Printf("Counter1: %d, %d, %d\n", counter1(), counter1(), counter1())
	fmt.Printf("Counter2: %d, %d\n", counter2(), counter2())

	// Defer examples
	fmt.Println("\nDefer example:")
	deferExample()
	
	fmt.Println("\nDefer in loop:")
	deferInLoop()
	
	fmt.Println("\nDefer with parameters:")
	deferWithParams()

	// Slice processing with functions
	numbers := []int{1, 2, 3, 4, 5}
	squared := processSlice(numbers, func(x int) int { return x * x })
	doubled := processSlice(numbers, func(x int) int { return x * 2 })
	
	fmt.Printf("Original: %v\n", numbers)
	fmt.Printf("Squared: %v\n", squared)
	fmt.Printf("Doubled: %v\n", doubled)

	// Chaining string processors
	toUpper := func(s string) string { return strings.ToUpper(s) }
	addPrefix := func(s string) string { return ">>> " + s }
	addSuffix := func(s string) string { return s + " <<<" }
	
	result := chainProcessors("hello world", toUpper, addPrefix, addSuffix)
	fmt.Printf("Processed string: %s\n", result)

	// Multiple function returns
	add, sub, mul, div := createMathOperations()
	fmt.Printf("10 + 5 = %d\n", add(10, 5))
	fmt.Printf("10 - 5 = %d\n", sub(10, 5))
	fmt.Printf("10 * 5 = %d\n", mul(10, 5))
	fmt.Printf("10 / 5 = %d\n", div(10, 5))

	// Decorator pattern
	simpleProcessor := func(s string) string { return strings.ToUpper(s) }
	loggedProcessor := withLogging(simpleProcessor)
	
	fmt.Println("\nDecorator example:")
	loggedProcessor("hello")

	// Memoization example
	fmt.Println("\nMemoization example:")
	memoizedFib := memoize(fibonacci)
	
	fmt.Printf("fib(10) = %d\n", memoizedFib(10))
	fmt.Printf("fib(10) = %d\n", memoizedFib(10)) // Should hit cache
	fmt.Printf("fib(11) = %d\n", memoizedFib(11))

	// Function slice
	operations := []MathOperation{add, sub, mul}
	fmt.Println("\nFunction slice:")
	for i, op := range operations {
		fmt.Printf("Operation %d: 8 op 4 = %d\n", i, op(8, 4))
	}
}