/*
METADATA:
Description: Demonstrates Go function declaration, parameters, return values, and basic function patterns
Keywords: function, func, parameters, return, multiple-return, named-return, recursion, scope
Category: functions
Concepts: function definition, parameter passing, return values, function scope, recursion
*/

package main

import "fmt"

// Simple function with no parameters and no return value
func sayHello() {
	fmt.Println("Hello, World!")
}

// Function with parameters
func greet(name string) {
	fmt.Printf("Hello, %s!\n", name)
}

// Function with multiple parameters
func add(a int, b int) int {
	return a + b
}

// Function with parameters of same type (shortened syntax)
func multiply(a, b int) int {
	return a * b
}

// Function with multiple return values
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

// Function with named return values
func rectangle(length, width float64) (area, perimeter float64) {
	area = length * width
	perimeter = 2 * (length + width)
	return // naked return
}

// Function with variadic parameters
func sum(numbers ...int) int {
	total := 0
	for _, num := range numbers {
		total += num
	}
	return total
}

// Function with mixed parameters and variadic
func processNumbers(multiplier int, numbers ...int) []int {
	result := make([]int, len(numbers))
	for i, num := range numbers {
		result[i] = num * multiplier
	}
	return result
}

// Recursive function
func factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}

// Function that returns a function (closure)
func makeMultiplier(factor int) func(int) int {
	return func(x int) int {
		return x * factor
	}
}

// Function with slice parameter
func printSlice(slice []string) {
	fmt.Printf("Slice: %v\n", slice)
	for i, value := range slice {
		fmt.Printf("  [%d]: %s\n", i, value)
	}
}

// Function with map parameter
func printMap(m map[string]int) {
	fmt.Printf("Map: %v\n", m)
	for key, value := range m {
		fmt.Printf("  %s: %d\n", key, value)
	}
}

// Function with pointer parameter
func increment(x *int) {
	*x++
}

func main() {
	// Calling simple functions
	sayHello()
	greet("Alice")

	// Functions with return values
	result := add(5, 3)
	fmt.Printf("5 + 3 = %d\n", result)

	product := multiply(4, 7)
	fmt.Printf("4 × 7 = %d\n", product)

	// Function with multiple return values
	quotient, err := divide(10.0, 3.0)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("10.0 ÷ 3.0 = %.2f\n", quotient)
	}

	// Error case
	_, err = divide(10.0, 0.0)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Named return values
	area, perimeter := rectangle(5.0, 3.0)
	fmt.Printf("Rectangle 5×3: Area=%.2f, Perimeter=%.2f\n", area, perimeter)

	// Variadic functions
	total := sum(1, 2, 3, 4, 5)
	fmt.Printf("Sum of 1,2,3,4,5 = %d\n", total)

	// Passing slice to variadic function
	numbers := []int{10, 20, 30}
	total2 := sum(numbers...) // spread operator
	fmt.Printf("Sum of slice %v = %d\n", numbers, total2)

	// Mixed parameters with variadic
	doubled := processNumbers(2, 1, 2, 3, 4, 5)
	fmt.Printf("Numbers doubled: %v\n", doubled)

	// Recursive function
	fact := factorial(5)
	fmt.Printf("5! = %d\n", fact)

	// Function closure
	double := makeMultiplier(2)
	triple := makeMultiplier(3)
	fmt.Printf("double(5) = %d\n", double(5))
	fmt.Printf("triple(5) = %d\n", triple(5))

	// Functions with complex parameters
	fruits := []string{"apple", "banana", "cherry"}
	printSlice(fruits)

	ages := map[string]int{"Alice": 30, "Bob": 25}
	printMap(ages)

	// Function with pointer
	value := 10
	fmt.Printf("Before increment: %d\n", value)
	increment(&value)
	fmt.Printf("After increment: %d\n", value)

	// Anonymous function
	square := func(x int) int {
		return x * x
	}
	fmt.Printf("Square of 6 = %d\n", square(6))

	// Immediately invoked function expression (IIFE)
	result3 := func(a, b int) int {
		return a*a + b*b
	}(3, 4)
	fmt.Printf("3² + 4² = %d\n", result3)
}