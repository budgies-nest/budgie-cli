/*
METADATA:
Description: Demonstrates Go testing and benchmarking including unit tests, table tests, benchmarks, and testing best practices
Keywords: testing, benchmark, unit-test, table-test, test-coverage, performance, mock, assert
Category: testing
Concepts: unit testing, table-driven tests, benchmarking, test coverage, testing patterns, performance testing
*/

package main

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"
	"time"
)

// Example functions to test
func Add(a, b int) int {
	return a + b
}

func Multiply(a, b int) int {
	return a * b
}

func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

func IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}
	for i := 3; i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func Factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * Factorial(n-1)
}

func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func FindMax(nums []int) (int, error) {
	if len(nums) == 0 {
		return 0, errors.New("empty slice")
	}
	
	max := nums[0]
	for _, num := range nums[1:] {
		if num > max {
			max = num
		}
	}
	return max, nil
}

func BubbleSort(arr []int) []int {
	n := len(arr)
	result := make([]int, n)
	copy(result, arr)
	
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if result[j] > result[j+1] {
				result[j], result[j+1] = result[j+1], result[j]
			}
		}
	}
	return result
}

// Struct and methods for testing
type Calculator struct {
	memory float64
}

func NewCalculator() *Calculator {
	return &Calculator{}
}

func (c *Calculator) Add(value float64) {
	c.memory += value
}

func (c *Calculator) Subtract(value float64) {
	c.memory -= value
}

func (c *Calculator) Multiply(value float64) {
	c.memory *= value
}

func (c *Calculator) Divide(value float64) error {
	if value == 0 {
		return errors.New("division by zero")
	}
	c.memory /= value
	return nil
}

func (c *Calculator) GetResult() float64 {
	return c.memory
}

func (c *Calculator) Clear() {
	c.memory = 0
}

// Example test functions (these would normally be in _test.go files)

// Basic unit test example
func TestAdd(t *testing.T) {
	result := Add(2, 3)
	expected := 5
	
	if result != expected {
		t.Errorf("Add(2, 3) = %d; expected %d", result, expected)
	}
}

// Table-driven test example
func TestMultiply(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"positive numbers", 3, 4, 12},
		{"zero multiplication", 5, 0, 0},
		{"negative numbers", -3, 4, -12},
		{"negative by negative", -3, -4, 12},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Multiply(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Multiply(%d, %d) = %d; expected %d", 
					tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// Test with error handling
func TestDivide(t *testing.T) {
	// Test successful division
	result, err := Divide(10.0, 2.0)
	if err != nil {
		t.Errorf("Divide(10.0, 2.0) returned unexpected error: %v", err)
	}
	expected := 5.0
	if result != expected {
		t.Errorf("Divide(10.0, 2.0) = %f; expected %f", result, expected)
	}
	
	// Test division by zero
	_, err = Divide(10.0, 0.0)
	if err == nil {
		t.Error("Divide(10.0, 0.0) should return an error")
	}
}

// Test with subtests
func TestIsPrime(t *testing.T) {
	primeTests := []struct {
		input    int
		expected bool
	}{
		{2, true},
		{3, true},
		{4, false},
		{17, true},
		{25, false},
		{29, true},
	}
	
	for _, test := range primeTests {
		t.Run(fmt.Sprintf("IsPrime(%d)", test.input), func(t *testing.T) {
			result := IsPrime(test.input)
			if result != test.expected {
				t.Errorf("IsPrime(%d) = %v; expected %v", 
					test.input, result, test.expected)
			}
		})
	}
}

// Test with setup and teardown
func TestCalculator(t *testing.T) {
	// Setup
	calc := NewCalculator()
	
	t.Run("initial state", func(t *testing.T) {
		if calc.GetResult() != 0 {
			t.Errorf("New calculator should start with 0, got %f", calc.GetResult())
		}
	})
	
	t.Run("addition", func(t *testing.T) {
		calc.Clear()
		calc.Add(5)
		calc.Add(3)
		if calc.GetResult() != 8 {
			t.Errorf("Expected 8, got %f", calc.GetResult())
		}
	})
	
	t.Run("division by zero", func(t *testing.T) {
		calc.Clear()
		calc.Add(10)
		err := calc.Divide(0)
		if err == nil {
			t.Error("Division by zero should return an error")
		}
	})
}

// Benchmark examples
func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add(123, 456)
	}
}

func BenchmarkIsPrime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsPrime(97) // Test with a prime number
	}
}

func BenchmarkFactorial(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Factorial(10)
	}
}

func BenchmarkReverseString(b *testing.B) {
	s := "Hello, World! This is a test string for benchmarking."
	b.ResetTimer() // Reset timer after setup
	
	for i := 0; i < b.N; i++ {
		ReverseString(s)
	}
}

// Benchmark with different input sizes
func BenchmarkBubbleSort(b *testing.B) {
	sizes := []int{10, 100, 1000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			// Generate test data
			data := make([]int, size)
			for i := 0; i < size; i++ {
				data[i] = size - i // Reverse order (worst case)
			}
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				testData := make([]int, len(data))
				copy(testData, data)
				b.StartTimer()
				
				BubbleSort(testData)
			}
		})
	}
}

// Memory allocation benchmark
func BenchmarkStringConcatenation(b *testing.B) {
	words := []string{"hello", "world", "from", "go", "benchmarking"}
	
	b.Run("using-plus-operator", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result := ""
			for _, word := range words {
				result += word + " "
			}
		}
	})
	
	b.Run("using-strings-builder", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var builder strings.Builder
			for _, word := range words {
				builder.WriteString(word)
				builder.WriteString(" ")
			}
			_ = builder.String()
		}
	})
}

// Helper functions for testing
func assertEqual(t *testing.T, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertError(t *testing.T, err error, wantError bool) {
	t.Helper()
	if (err != nil) != wantError {
		t.Errorf("error = %v, wantError = %v", err, wantError)
	}
}

func assertFloatEqual(t *testing.T, got, want, tolerance float64) {
	t.Helper()
	if math.Abs(got-want) > tolerance {
		t.Errorf("got %f, want %f (tolerance %f)", got, want, tolerance)
	}
}

// Example test using helper functions
func TestFindMax(t *testing.T) {
	t.Run("normal case", func(t *testing.T) {
		result, err := FindMax([]int{1, 5, 3, 9, 2})
		assertError(t, err, false)
		assertEqual(t, result, 9)
	})
	
	t.Run("empty slice", func(t *testing.T) {
		_, err := FindMax([]int{})
		assertError(t, err, true)
	})
	
	t.Run("single element", func(t *testing.T) {
		result, err := FindMax([]int{42})
		assertError(t, err, false)
		assertEqual(t, result, 42)
	})
}

// Performance comparison benchmark
func BenchmarkStringOperations(b *testing.B) {
	data := strings.Repeat("a", 1000)
	
	b.Run("Contains", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			strings.Contains(data, "z")
		}
	})
	
	b.Run("Index", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			strings.Index(data, "z")
		}
	})
}

// Parallel benchmark
func BenchmarkParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Simulate some work
			IsPrime(97)
		}
	})
}

func main() {
	fmt.Println("=== TESTING AND BENCHMARKING EXAMPLES ===")
	
	fmt.Println("\nThis file demonstrates Go testing patterns:")
	fmt.Println("1. Basic unit tests")
	fmt.Println("2. Table-driven tests")
	fmt.Println("3. Error handling tests")
	fmt.Println("4. Subtests")
	fmt.Println("5. Setup and teardown")
	fmt.Println("6. Benchmarking")
	fmt.Println("7. Memory benchmarks")
	fmt.Println("8. Parallel benchmarks")
	fmt.Println("9. Helper functions")
	fmt.Println("10. Test assertions")
	
	fmt.Println("\n=== RUNNING TESTS ===")
	fmt.Println("To run these tests, you would typically:")
	fmt.Println("1. Create files ending with '_test.go'")
	fmt.Println("2. Use 'go test' command")
	fmt.Println("3. Use 'go test -v' for verbose output")
	fmt.Println("4. Use 'go test -cover' for coverage")
	fmt.Println("5. Use 'go test -bench=.' for benchmarks")
	
	// Demonstrate manual test execution
	fmt.Println("\n=== MANUAL TEST DEMONSTRATION ===")
	
	// Mock testing framework for demonstration
	type mockT struct {
		failed bool
		logs   []string
	}
	
	func (m *mockT) Errorf(format string, args ...interface{}) {
		m.failed = true
		m.logs = append(m.logs, fmt.Sprintf("ERROR: "+format, args...))
	}
	
	func (m *mockT) Error(args ...interface{}) {
		m.failed = true
		m.logs = append(m.logs, fmt.Sprintf("ERROR: %v", args...))
	}
	
	func (m *mockT) Run(name string, fn func(*mockT)) {
		subTest := &mockT{}
		fmt.Printf("Running test: %s\n", name)
		fn(subTest)
		if subTest.failed {
			fmt.Printf("  FAIL: %s\n", name)
			for _, log := range subTest.logs {
				fmt.Printf("    %s\n", log)
			}
		} else {
			fmt.Printf("  PASS: %s\n", name)
		}
	}
	
	func (m *mockT) Helper() {} // Mock helper
	
	// Run some tests manually for demonstration
	t := &mockT{}
	
	// Test Add function
	fmt.Println("\nTesting Add function:")
	if Add(2, 3) == 5 {
		fmt.Println("  PASS: Add(2, 3) = 5")
	} else {
		fmt.Println("  FAIL: Add(2, 3) != 5")
	}
	
	// Test Divide function
	fmt.Println("\nTesting Divide function:")
	result, err := Divide(10, 2)
	if err == nil && result == 5.0 {
		fmt.Println("  PASS: Divide(10, 2) = 5.0")
	} else {
		fmt.Printf("  FAIL: Divide(10, 2) = %f, error: %v\n", result, err)
	}
	
	_, err = Divide(10, 0)
	if err != nil {
		fmt.Println("  PASS: Divide(10, 0) returns error")
	} else {
		fmt.Println("  FAIL: Divide(10, 0) should return error")
	}
	
	// Test IsPrime function with table
	fmt.Println("\nTesting IsPrime function:")
	primeTests := []struct {
		input    int
		expected bool
	}{
		{2, true}, {3, true}, {4, false}, {17, true}, {25, false},
	}
	
	allPassed := true
	for _, test := range primeTests {
		result := IsPrime(test.input)
		if result == test.expected {
			fmt.Printf("  PASS: IsPrime(%d) = %v\n", test.input, result)
		} else {
			fmt.Printf("  FAIL: IsPrime(%d) = %v, expected %v\n", 
				test.input, result, test.expected)
			allPassed = false
		}
	}
	
	if allPassed {
		fmt.Println("  All IsPrime tests passed!")
	}
	
	// Demonstrate benchmark timing
	fmt.Println("\n=== MANUAL BENCHMARK DEMONSTRATION ===")
	
	// Benchmark IsPrime function
	fmt.Println("Benchmarking IsPrime(97):")
	start := time.Now()
	iterations := 1000000
	for i := 0; i < iterations; i++ {
		IsPrime(97)
	}
	duration := time.Since(start)
	
	fmt.Printf("  %d iterations in %v\n", iterations, duration)
	fmt.Printf("  %v per operation\n", duration/time.Duration(iterations))
	
	// Benchmark string operations
	fmt.Println("\nBenchmarking string concatenation:")
	testString := "Hello, World!"
	
	// Using + operator
	start = time.Now()
	iterations = 10000
	for i := 0; i < iterations; i++ {
		result := ""
		for j := 0; j < 10; j++ {
			result += testString
		}
	}
	plusDuration := time.Since(start)
	
	// Using strings.Builder
	start = time.Now()
	for i := 0; i < iterations; i++ {
		var builder strings.Builder
		for j := 0; j < 10; j++ {
			builder.WriteString(testString)
		}
		_ = builder.String()
	}
	builderDuration := time.Since(start)
	
	fmt.Printf("  Plus operator: %v (%v per op)\n", 
		plusDuration, plusDuration/time.Duration(iterations))
	fmt.Printf("  strings.Builder: %v (%v per op)\n", 
		builderDuration, builderDuration/time.Duration(iterations))
	fmt.Printf("  Builder is %.2fx faster\n", 
		float64(plusDuration)/float64(builderDuration))
	
	fmt.Println("\n=== TESTING BEST PRACTICES ===")
	fmt.Println("1. Write tests for all public functions")
	fmt.Println("2. Use table-driven tests for multiple test cases")
	fmt.Println("3. Test both success and error cases")
	fmt.Println("4. Use subtests for better organization")
	fmt.Println("5. Write helper functions to reduce duplication")
	fmt.Println("6. Test edge cases and boundary conditions")
	fmt.Println("7. Use benchmarks to measure performance")
	fmt.Println("8. Aim for high test coverage")
	fmt.Println("9. Keep tests simple and focused")
	fmt.Println("10. Use mocks for external dependencies")
	
	fmt.Println("\n=== BENCHMARKING BEST PRACTICES ===")
	fmt.Println("1. Use b.ResetTimer() after setup")
	fmt.Println("2. Use b.StopTimer() and b.StartTimer() for complex setup")
	fmt.Println("3. Run benchmarks multiple times")
	fmt.Println("4. Compare different implementations")
	fmt.Println("5. Measure both time and memory allocations")
	fmt.Println("6. Use parallel benchmarks for concurrent code")
	fmt.Println("7. Test with realistic data sizes")
	fmt.Println("8. Profile when needed (go tool pprof)")
	
	fmt.Println("\nTesting and benchmarking examples completed!")
}