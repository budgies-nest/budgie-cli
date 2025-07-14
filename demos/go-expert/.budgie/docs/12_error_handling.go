/*
METADATA:
Description: Demonstrates Go error handling including error creation, custom errors, error wrapping, and best practices
Keywords: error, error-handling, custom-error, error-wrapping, fmt.Errorf, errors.New, panic, recover
Category: error-handling
Concepts: error interface, error creation, error propagation, error wrapping, custom error types
*/

package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Custom error types
type ValidationError struct {
	Field   string
	Message string
}

func (ve ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", ve.Field, ve.Message)
}

type NotFoundError struct {
	Resource string
	ID       interface{}
}

func (nfe NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %v not found", nfe.Resource, nfe.ID)
}

// Error with additional methods
type DatabaseError struct {
	Operation string
	Err       error
	Code      int
}

func (de DatabaseError) Error() string {
	return fmt.Sprintf("database error during %s (code %d): %v", de.Operation, de.Code, de.Err)
}

func (de DatabaseError) Unwrap() error {
	return de.Err
}

func (de DatabaseError) Temporary() bool {
	return de.Code >= 500 && de.Code < 600
}

// Functions demonstrating different error creation methods
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

func parseInt(s string) (int, error) {
	if s == "" {
		return 0, fmt.Errorf("empty string cannot be converted to integer")
	}
	
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("failed to parse '%s' as integer: %w", s, err)
	}
	
	return value, nil
}

func validateUser(name, email string, age int) error {
	if name == "" {
		return ValidationError{Field: "name", Message: "cannot be empty"}
	}
	
	if !strings.Contains(email, "@") {
		return ValidationError{Field: "email", Message: "must contain @ symbol"}
	}
	
	if age < 0 || age > 150 {
		return ValidationError{Field: "age", Message: "must be between 0 and 150"}
	}
	
	return nil
}

func findUser(id int) (string, error) {
	users := map[int]string{
		1: "Alice",
		2: "Bob",
		3: "Carol",
	}
	
	if name, exists := users[id]; exists {
		return name, nil
	}
	
	return "", NotFoundError{Resource: "User", ID: id}
}

func simulateDatabaseOperation(operation string, shouldFail bool) error {
	if shouldFail {
		baseErr := errors.New("connection timeout")
		return DatabaseError{
			Operation: operation,
			Err:       baseErr,
			Code:      503,
		}
	}
	return nil
}

// Function that demonstrates error wrapping
func processFile(filename string) error {
	if filename == "" {
		return errors.New("filename cannot be empty")
	}
	
	// Simulate file reading error
	if !strings.HasSuffix(filename, ".txt") {
		return fmt.Errorf("invalid file format for %s: %w", filename, errors.New("only .txt files are supported"))
	}
	
	// Simulate permission error
	if strings.Contains(filename, "restricted") {
		return fmt.Errorf("access denied to %s: %w", filename, errors.New("insufficient permissions"))
	}
	
	return nil
}

// Function chain demonstrating error propagation
func level3Function() error {
	return errors.New("error from level 3")
}

func level2Function() error {
	err := level3Function()
	if err != nil {
		return fmt.Errorf("level 2 failed: %w", err)
	}
	return nil
}

func level1Function() error {
	err := level2Function()
	if err != nil {
		return fmt.Errorf("level 1 failed: %w", err)
	}
	return nil
}

// Multiple return values with errors
func calculateStats(numbers []float64) (float64, float64, error) {
	if len(numbers) == 0 {
		return 0, 0, errors.New("cannot calculate stats for empty slice")
	}
	
	sum := 0.0
	for _, num := range numbers {
		sum += num
	}
	
	mean := sum / float64(len(numbers))
	
	// Calculate variance
	variance := 0.0
	for _, num := range numbers {
		diff := num - mean
		variance += diff * diff
	}
	variance /= float64(len(numbers))
	
	return mean, variance, nil
}

// Error handling with cleanup
func processWithCleanup(data string) (result string, err error) {
	// Simulate resource allocation
	fmt.Println("Allocating resources...")
	
	// Use defer for cleanup
	defer func() {
		fmt.Println("Cleaning up resources...")
		if r := recover(); r != nil {
			err = fmt.Errorf("panic during processing: %v", r)
		}
	}()
	
	if data == "panic" {
		panic("simulated panic")
	}
	
	if data == "error" {
		return "", errors.New("simulated error")
	}
	
	return strings.ToUpper(data), nil
}

// Panic and recover demonstration
func safeOperation(operation func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()
	
	operation()
	return nil
}

func main() {
	fmt.Println("=== BASIC ERROR HANDLING ===")
	
	// Simple error handling
	result, err := divide(10, 2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("10 / 2 = %.2f\n", result)
	}
	
	// Error case
	_, err = divide(10, 0)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println("\n=== ERROR WRAPPING ===")
	
	// Error wrapping with fmt.Errorf
	_, err = parseInt("abc")
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		
		// Check if it wraps another error
		if unwrapped := errors.Unwrap(err); unwrapped != nil {
			fmt.Printf("Unwrapped error: %v\n", unwrapped)
		}
	}
	
	// Successful parsing
	value, err := parseInt("123")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Parsed value: %d\n", value)
	}

	fmt.Println("\n=== CUSTOM ERROR TYPES ===")
	
	// Validation errors
	err = validateUser("", "invalid-email", 200)
	if err != nil {
		fmt.Printf("Validation error: %v\n", err)
		
		// Type assertion to access custom fields
		if ve, ok := err.(ValidationError); ok {
			fmt.Printf("  Field: %s\n", ve.Field)
			fmt.Printf("  Message: %s\n", ve.Message)
		}
	}
	
	// Valid user
	err = validateUser("Alice", "alice@example.com", 30)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("User validation passed")
	}

	fmt.Println("\n=== NOT FOUND ERRORS ===")
	
	// User found
	name, err := findUser(1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Found user: %s\n", name)
	}
	
	// User not found
	_, err = findUser(999)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		
		// Type assertion for custom error
		if nfe, ok := err.(NotFoundError); ok {
			fmt.Printf("  Resource: %s\n", nfe.Resource)
			fmt.Printf("  ID: %v\n", nfe.ID)
		}
	}

	fmt.Println("\n=== DATABASE ERRORS ===")
	
	// Successful operation
	err = simulateDatabaseOperation("SELECT", false)
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
	} else {
		fmt.Println("Database operation successful")
	}
	
	// Failed operation
	err = simulateDatabaseOperation("INSERT", true)
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		
		if de, ok := err.(DatabaseError); ok {
			fmt.Printf("  Operation: %s\n", de.Operation)
			fmt.Printf("  Code: %d\n", de.Code)
			fmt.Printf("  Temporary: %t\n", de.Temporary())
			fmt.Printf("  Underlying error: %v\n", de.Unwrap())
		}
	}

	fmt.Println("\n=== ERROR CHAINS ===")
	
	// Demonstrate error propagation through call stack
	err = level1Function()
	if err != nil {
		fmt.Printf("Chain error: %v\n", err)
		
		// Walk through error chain
		current := err
		level := 0
		for current != nil {
			fmt.Printf("  Level %d: %v\n", level, current)
			current = errors.Unwrap(current)
			level++
		}
	}

	fmt.Println("\n=== FILE PROCESSING ERRORS ===")
	
	testFiles := []string{"data.txt", "config.json", "restricted.txt", ""}
	
	for _, filename := range testFiles {
		err := processFile(filename)
		if err != nil {
			fmt.Printf("File '%s' error: %v\n", filename, err)
		} else {
			fmt.Printf("File '%s' processed successfully\n", filename)
		}
	}

	fmt.Println("\n=== MULTIPLE RETURN VALUES ===")
	
	numbers := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	mean, variance, err := calculateStats(numbers)
	if err != nil {
		fmt.Printf("Stats error: %v\n", err)
	} else {
		fmt.Printf("Mean: %.2f, Variance: %.2f\n", mean, variance)
	}
	
	// Empty slice case
	_, _, err = calculateStats([]float64{})
	if err != nil {
		fmt.Printf("Empty slice error: %v\n", err)
	}

	fmt.Println("\n=== ERROR HANDLING WITH CLEANUP ===")
	
	testData := []string{"hello", "error", "panic"}
	
	for _, data := range testData {
		fmt.Printf("Processing '%s':\n", data)
		result, err := processWithCleanup(data)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
		} else {
			fmt.Printf("  Result: %s\n", result)
		}
		fmt.Println()
	}

	fmt.Println("\n=== PANIC AND RECOVER ===")
	
	// Safe operation that might panic
	err = safeOperation(func() {
		fmt.Println("Executing safe operation...")
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Operation completed successfully")
	}
	
	// Operation that panics
	err = safeOperation(func() {
		panic("something went wrong!")
	})
	if err != nil {
		fmt.Printf("Recovered error: %v\n", err)
	}

	fmt.Println("\n=== ERROR HANDLING PATTERNS ===")
	
	// Pattern 1: Early return
	processValue := func(value string) error {
		if value == "" {
			return errors.New("value cannot be empty")
		}
		
		if len(value) < 3 {
			return fmt.Errorf("value '%s' is too short", value)
		}
		
		fmt.Printf("Processing value: %s\n", value)
		return nil
	}
	
	testValues := []string{"", "hi", "hello"}
	for _, val := range testValues {
		if err := processValue(val); err != nil {
			fmt.Printf("Error processing '%s': %v\n", val, err)
		}
	}
	
	// Pattern 2: Error accumulation
	var errs []error
	
	if err := processValue(""); err != nil {
		errs = append(errs, err)
	}
	
	if err := processValue("hi"); err != nil {
		errs = append(errs, err)
	}
	
	if len(errs) > 0 {
		fmt.Printf("Multiple errors occurred:\n")
		for i, err := range errs {
			fmt.Printf("  %d: %v\n", i+1, err)
		}
	}
	
	// Pattern 3: Error sentinel values
	var (
		ErrInvalidInput = errors.New("invalid input")
		ErrNotFound     = errors.New("not found")
		ErrTimeout      = errors.New("timeout")
	)
	
	checkError := func(err error) {
		switch {
		case errors.Is(err, ErrInvalidInput):
			fmt.Println("Handling invalid input error")
		case errors.Is(err, ErrNotFound):
			fmt.Println("Handling not found error")
		case errors.Is(err, ErrTimeout):
			fmt.Println("Handling timeout error")
		default:
			fmt.Printf("Unknown error: %v\n", err)
		}
	}
	
	testErrors := []error{ErrInvalidInput, ErrNotFound, errors.New("other error")}
	for _, err := range testErrors {
		checkError(err)
	}
}