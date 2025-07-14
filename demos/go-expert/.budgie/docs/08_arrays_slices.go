/*
METADATA:
Description: Demonstrates Go arrays and slices including declaration, initialization, manipulation, and common operations
Keywords: array, slice, make, append, copy, len, cap, slice-operations, dynamic-array
Category: data-structures
Concepts: fixed arrays, dynamic slices, slice operations, memory management, slice tricks
*/

package main

import "fmt"

func main() {
	// Array declaration and initialization
	fmt.Println("=== ARRAYS ===")
	
	// Fixed-size arrays
	var numbers [5]int
	numbers[0] = 10
	numbers[1] = 20
	numbers[2] = 30
	fmt.Printf("Array: %v\n", numbers)
	
	// Array with literal initialization
	fruits := [3]string{"apple", "banana", "cherry"}
	fmt.Printf("Fruits array: %v\n", fruits)
	
	// Array with inferred size
	colors := [...]string{"red", "green", "blue", "yellow"}
	fmt.Printf("Colors array: %v (length: %d)\n", colors, len(colors))
	
	// Array with specific indices
	grades := [5]int{0: 85, 2: 92, 4: 78}
	fmt.Printf("Grades array: %v\n", grades)

	// Arrays are value types (copied when assigned)
	originalArray := [3]int{1, 2, 3}
	copiedArray := originalArray
	copiedArray[0] = 999
	fmt.Printf("Original: %v, Copied: %v\n", originalArray, copiedArray)

	fmt.Println("\n=== SLICES ===")
	
	// Slice declaration and initialization
	var names []string
	fmt.Printf("Empty slice: %v (len: %d, cap: %d)\n", names, len(names), cap(names))
	
	// Slice literal
	cities := []string{"New York", "London", "Tokyo"}
	fmt.Printf("Cities slice: %v (len: %d, cap: %d)\n", cities, len(cities), cap(cities))
	
	// Creating slice with make
	scores := make([]int, 5)        // length 5, capacity 5
	ages := make([]int, 3, 10)      // length 3, capacity 10
	fmt.Printf("Scores: %v (len: %d, cap: %d)\n", scores, len(scores), cap(scores))
	fmt.Printf("Ages: %v (len: %d, cap: %d)\n", ages, len(ages), cap(ages))
	
	// Slice from array
	numbersSlice := numbers[1:4]    // elements 1, 2, 3 from array
	fmt.Printf("Slice from array: %v\n", numbersSlice)
	
	// Various slicing operations
	allFruits := fruits[:]          // all elements
	firstTwo := fruits[:2]          // first 2 elements
	lastTwo := fruits[1:]           // from index 1 to end
	fmt.Printf("All fruits: %v\n", allFruits)
	fmt.Printf("First two: %v\n", firstTwo)
	fmt.Printf("Last two: %v\n", lastTwo)

	fmt.Println("\n=== SLICE OPERATIONS ===")
	
	// Append operations
	var items []string
	items = append(items, "first")
	items = append(items, "second", "third")
	items = append(items, cities...)  // append another slice
	fmt.Printf("After appends: %v (len: %d, cap: %d)\n", items, len(items), cap(items))
	
	// Append with capacity growth
	numbers2 := []int{1, 2, 3}
	fmt.Printf("Before append: %v (len: %d, cap: %d)\n", numbers2, len(numbers2), cap(numbers2))
	for i := 4; i <= 10; i++ {
		numbers2 = append(numbers2, i)
		fmt.Printf("After append %d: len=%d, cap=%d\n", i, len(numbers2), cap(numbers2))
	}
	
	// Copy operation
	source := []int{1, 2, 3, 4, 5}
	dest := make([]int, len(source))
	copied := copy(dest, source)
	fmt.Printf("Source: %v\n", source)
	fmt.Printf("Destination: %v (copied %d elements)\n", dest, copied)
	
	// Partial copy
	partial := make([]int, 3)
	copy(partial, source)
	fmt.Printf("Partial copy: %v\n", partial)

	fmt.Println("\n=== SLICE MANIPULATION ===")
	
	// Insert element at beginning
	values := []int{2, 3, 4}
	values = append([]int{1}, values...)
	fmt.Printf("Insert at beginning: %v\n", values)
	
	// Insert element at position
	pos := 2
	values = append(values[:pos], append([]int{99}, values[pos:]...)...)
	fmt.Printf("Insert 99 at position %d: %v\n", pos, values)
	
	// Remove element
	indexToRemove := 1
	values = append(values[:indexToRemove], values[indexToRemove+1:]...)
	fmt.Printf("Remove element at index %d: %v\n", indexToRemove, values)
	
	// Remove last element
	if len(values) > 0 {
		values = values[:len(values)-1]
		fmt.Printf("Remove last element: %v\n", values)
	}

	fmt.Println("\n=== MULTI-DIMENSIONAL SLICES ===")
	
	// 2D slice
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	fmt.Printf("Matrix: %v\n", matrix)
	
	// Creating 2D slice with make
	rows, cols := 3, 4
	grid := make([][]int, rows)
	for i := range grid {
		grid[i] = make([]int, cols)
	}
	
	// Fill the grid
	counter := 1
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			grid[i][j] = counter
			counter++
		}
	}
	
	fmt.Printf("Grid:\n")
	for _, row := range grid {
		fmt.Printf("%v\n", row)
	}

	fmt.Println("\n=== SLICE COMPARISON AND SEARCH ===")
	
	// Slices cannot be compared directly, need to iterate
	slice1 := []int{1, 2, 3}
	slice2 := []int{1, 2, 3}
	
	equal := true
	if len(slice1) != len(slice2) {
		equal = false
	} else {
		for i := range slice1 {
			if slice1[i] != slice2[i] {
				equal = false
				break
			}
		}
	}
	fmt.Printf("Slices equal: %t\n", equal)
	
	// Search in slice
	searchFor := 3
	found := false
	foundIndex := -1
	for i, v := range slice1 {
		if v == searchFor {
			found = true
			foundIndex = i
			break
		}
	}
	fmt.Printf("Found %d: %t at index %d\n", searchFor, found, foundIndex)

	fmt.Println("\n=== SLICE MEMORY BEHAVIOR ===")
	
	// Demonstrating shared underlying array
	original := []int{1, 2, 3, 4, 5}
	slice1 = original[1:3]  // [2, 3]
	slice2 = original[2:4]  // [3, 4]
	
	fmt.Printf("Original: %v\n", original)
	fmt.Printf("Slice1: %v\n", slice1)
	fmt.Printf("Slice2: %v\n", slice2)
	
	// Modifying slice1 affects original and slice2
	slice1[1] = 999
	fmt.Printf("After modifying slice1[1]:\n")
	fmt.Printf("Original: %v\n", original)
	fmt.Printf("Slice1: %v\n", slice1)
	fmt.Printf("Slice2: %v\n", slice2)
	
	// Creating independent copy
	independent := make([]int, len(slice1))
	copy(independent, slice1)
	independent[0] = 777
	
	fmt.Printf("Independent copy after modification: %v\n", independent)
	fmt.Printf("Original slice1 unchanged: %v\n", slice1)
}