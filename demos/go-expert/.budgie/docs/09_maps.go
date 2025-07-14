/*
METADATA:
Description: Demonstrates Go maps including declaration, initialization, operations, and common patterns
Keywords: map, hash-map, key-value, make, delete, iteration, map-operations, zero-value
Category: data-structures
Concepts: hash maps, key-value pairs, map operations, map iteration, nested maps
*/

package main

import "fmt"

func main() {
	fmt.Println("=== MAP DECLARATION AND INITIALIZATION ===")
	
	// Map declaration with var (zero value is nil)
	var ages map[string]int
	fmt.Printf("Nil map: %v (length: %d)\n", ages, len(ages))
	
	// Initialize with make
	ages = make(map[string]int)
	ages["Alice"] = 30
	ages["Bob"] = 25
	ages["Carol"] = 35
	fmt.Printf("Ages map: %v\n", ages)
	
	// Map literal initialization
	scores := map[string]int{
		"Math":    95,
		"Science": 87,
		"History": 92,
		"English": 88,
	}
	fmt.Printf("Scores map: %v\n", scores)
	
	// Empty map literal
	countries := map[string]string{}
	countries["US"] = "United States"
	countries["UK"] = "United Kingdom"
	countries["JP"] = "Japan"
	fmt.Printf("Countries: %v\n", countries)
	
	// Map with different key types
	idToName := map[int]string{
		101: "John",
		102: "Jane",
		103: "Bob",
	}
	fmt.Printf("ID to Name: %v\n", idToName)

	fmt.Println("\n=== MAP OPERATIONS ===")
	
	// Reading values
	mathScore := scores["Math"]
	fmt.Printf("Math score: %d\n", mathScore)
	
	// Check if key exists (comma ok idiom)
	englishScore, exists := scores["English"]
	if exists {
		fmt.Printf("English score: %d\n", englishScore)
	}
	
	artScore, exists := scores["Art"]
	if !exists {
		fmt.Printf("Art score not found (got zero value: %d)\n", artScore)
	}
	
	// Adding new key-value pairs
	scores["Art"] = 90
	scores["Music"] = 85
	fmt.Printf("Updated scores: %v\n", scores)
	
	// Updating existing values
	scores["Math"] = 98
	fmt.Printf("Updated Math score: %d\n", scores["Math"])
	
	// Deleting keys
	delete(scores, "History")
	fmt.Printf("After deleting History: %v\n", scores)
	
	// Deleting non-existent key (no error)
	delete(scores, "NonExistent")
	fmt.Printf("After deleting non-existent key: %v\n", scores)

	fmt.Println("\n=== MAP ITERATION ===")
	
	// Iterate over key-value pairs
	fmt.Println("Iterating over scores:")
	for subject, score := range scores {
		fmt.Printf("  %s: %d\n", subject, score)
	}
	
	// Iterate over keys only
	fmt.Println("Subject names:")
	for subject := range scores {
		fmt.Printf("  %s\n", subject)
	}
	
	// Iterate over values only
	fmt.Println("Scores only:")
	for _, score := range scores {
		fmt.Printf("  %d\n", score)
	}

	fmt.Println("\n=== NESTED MAPS ===")
	
	// Map of maps
	studentGrades := map[string]map[string]int{
		"Alice": {
			"Math":    95,
			"Science": 92,
			"English": 88,
		},
		"Bob": {
			"Math":    78,
			"Science": 85,
			"English": 90,
		},
	}
	
	fmt.Println("Student grades:")
	for student, grades := range studentGrades {
		fmt.Printf("  %s:\n", student)
		for subject, grade := range grades {
			fmt.Printf("    %s: %d\n", subject, grade)
		}
	}
	
	// Adding new student
	studentGrades["Carol"] = make(map[string]int)
	studentGrades["Carol"]["Math"] = 88
	studentGrades["Carol"]["Science"] = 91
	
	// Access nested map values
	aliceMath := studentGrades["Alice"]["Math"]
	fmt.Printf("Alice's Math grade: %d\n", aliceMath)

	fmt.Println("\n=== MAP WITH COMPLEX VALUE TYPES ===")
	
	// Map with slice values
	teamMembers := map[string][]string{
		"Development": {"Alice", "Bob", "Charlie"},
		"Testing":     {"Dave", "Eve"},
		"Design":      {"Frank", "Grace", "Henry", "Ivy"},
	}
	
	fmt.Println("Team members:")
	for team, members := range teamMembers {
		fmt.Printf("  %s: %v\n", team, members)
	}
	
	// Add member to team
	teamMembers["Testing"] = append(teamMembers["Testing"], "Jack")
	fmt.Printf("Updated Testing team: %v\n", teamMembers["Testing"])

	fmt.Println("\n=== MAP WITH STRUCT VALUES ===")
	
	type Person struct {
		Name string
		Age  int
		City string
	}
	
	people := map[int]Person{
		1: {Name: "Alice", Age: 30, City: "New York"},
		2: {Name: "Bob", Age: 25, City: "London"},
		3: {Name: "Carol", Age: 35, City: "Tokyo"},
	}
	
	fmt.Println("People map:")
	for id, person := range people {
		fmt.Printf("  ID %d: %s, %d years old, lives in %s\n", 
			id, person.Name, person.Age, person.City)
	}

	fmt.Println("\n=== MAP FUNCTIONS AND PATTERNS ===")
	
	// Count occurrences
	text := "hello world hello go"
	wordCount := make(map[string]int)
	for _, char := range text {
		if char != ' ' {
			word := string(char)
			wordCount[word]++
		}
	}
	fmt.Printf("Character count: %v\n", wordCount)
	
	// Group by property
	ages2 := map[string]int{
		"Alice": 30,
		"Bob":   25,
		"Carol": 30,
		"Dave":  25,
		"Eve":   35,
	}
	
	ageGroups := make(map[int][]string)
	for name, age := range ages2 {
		ageGroups[age] = append(ageGroups[age], name)
	}
	
	fmt.Println("Age groups:")
	for age, names := range ageGroups {
		fmt.Printf("  Age %d: %v\n", age, names)
	}
	
	// Map as set (using map[type]bool)
	fruits := map[string]bool{
		"apple":  true,
		"banana": true,
		"cherry": true,
	}
	
	// Check if item is in set
	if fruits["apple"] {
		fmt.Println("Apple is in the fruit set")
	}
	
	if !fruits["orange"] {
		fmt.Println("Orange is not in the fruit set")
	}
	
	// Add to set
	fruits["orange"] = true
	
	// Remove from set
	delete(fruits, "banana")
	
	fmt.Printf("Fruit set: %v\n", fruits)

	fmt.Println("\n=== MAP COMPARISON AND COPYING ===")
	
	// Maps cannot be compared directly (except to nil)
	map1 := map[string]int{"a": 1, "b": 2}
	map2 := map[string]int{"a": 1, "b": 2}
	
	// Function to compare maps
	mapsEqual := func(m1, m2 map[string]int) bool {
		if len(m1) != len(m2) {
			return false
		}
		for k, v1 := range m1 {
			if v2, exists := m2[k]; !exists || v1 != v2 {
				return false
			}
		}
		return true
	}
	
	fmt.Printf("Maps equal: %t\n", mapsEqual(map1, map2))
	
	// Copy map
	mapCopy := make(map[string]int)
	for k, v := range map1 {
		mapCopy[k] = v
	}
	fmt.Printf("Original: %v\n", map1)
	fmt.Printf("Copy: %v\n", mapCopy)
	
	// Modify copy to verify independence
	mapCopy["c"] = 3
	fmt.Printf("After modifying copy:\n")
	fmt.Printf("Original: %v\n", map1)
	fmt.Printf("Copy: %v\n", mapCopy)
}