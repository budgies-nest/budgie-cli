/*
METADATA:
Description: Demonstrates Go for loops, range loops, and loop control statements (break, continue)
Keywords: for, loop, range, break, continue, iteration, while, infinite-loop, nested-loops
Category: control-flow
Concepts: iteration, loop control, range iteration, infinite loops, loop patterns
*/

package main

import "fmt"

func main() {
	// Basic for loop
	fmt.Println("Basic for loop:")
	for i := 0; i < 5; i++ {
		fmt.Printf("Iteration %d\n", i)
	}

	// For loop with different increment
	fmt.Println("\nCounting by 2s:")
	for i := 0; i <= 10; i += 2 {
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// For loop counting backwards
	fmt.Println("\nCountdown:")
	for i := 5; i > 0; i-- {
		fmt.Printf("%d ", i)
	}
	fmt.Println("Go!")

	// While-style loop (condition only)
	fmt.Println("\nWhile-style loop:")
	count := 0
	for count < 3 {
		fmt.Printf("Count: %d\n", count)
		count++
	}

	// Infinite loop with break
	fmt.Println("\nInfinite loop with break:")
	counter := 0
	for {
		if counter >= 3 {
			break
		}
		fmt.Printf("Counter: %d\n", counter)
		counter++
	}

	// Continue statement
	fmt.Println("\nUsing continue (skip even numbers):")
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			continue
		}
		fmt.Printf("Odd number: %d\n", i)
	}

	// Range over slice
	fmt.Println("\nRange over slice:")
	fruits := []string{"apple", "banana", "cherry", "date"}
	for index, fruit := range fruits {
		fmt.Printf("Index %d: %s\n", index, fruit)
	}

	// Range with index only
	fmt.Println("\nRange with index only:")
	for index := range fruits {
		fmt.Printf("Index: %d\n", index)
	}

	// Range with value only (using blank identifier)
	fmt.Println("\nRange with value only:")
	for _, fruit := range fruits {
		fmt.Printf("Fruit: %s\n", fruit)
	}

	// Range over map
	fmt.Println("\nRange over map:")
	ages := map[string]int{
		"Alice": 30,
		"Bob":   25,
		"Carol": 35,
	}
	for name, age := range ages {
		fmt.Printf("%s is %d years old\n", name, age)
	}

	// Range over string (gets runes)
	fmt.Println("\nRange over string:")
	text := "Hello"
	for index, char := range text {
		fmt.Printf("Index %d: %c (Unicode: %d)\n", index, char, char)
	}

	// Nested loops
	fmt.Println("\nNested loops (multiplication table):")
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			fmt.Printf("%d×%d=%d ", i, j, i*j)
		}
		fmt.Println()
	}

	// Loop with labeled break
	fmt.Println("\nLabeled break:")
outer:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == 1 && j == 1 {
				break outer
			}
			fmt.Printf("(%d,%d) ", i, j)
		}
	}
	fmt.Println("\nBroke out of nested loop")

	// Range over channel (will be demonstrated in concurrency examples)
	fmt.Println("\nRange over array:")
	numbers := [5]int{10, 20, 30, 40, 50}
	for index, value := range numbers {
		fmt.Printf("numbers[%d] = %d\n", index, value)
	}
}