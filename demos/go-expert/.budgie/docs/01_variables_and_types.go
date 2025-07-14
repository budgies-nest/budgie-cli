/*
METADATA:
Description: Demonstrates Go variable declarations, type inference, and basic data types
Keywords: variables, types, declaration, var, const, string, int, float, bool, type-inference, zero-values
Category: basics
Concepts: variable declaration, type system, constants, zero values
*/

package main

import "fmt"

func main() {
	// Variable declarations with explicit types
	var name string = "Alice"
	var age int = 30
	var salary float64 = 50000.50
	var isEmployed bool = true

	// Short variable declaration (type inference)
	city := "New York"
	zipCode := 10001
	temperature := 72.5
	isWeekend := false

	// Multiple variable declarations
	var (
		firstName string = "John"
		lastName  string = "Doe"
		height    int    = 180
	)

	// Constants
	const PI = 3.14159
	const CompanyName = "Tech Corp"
	const MaxUsers = 1000

	// Zero values (variables declared without initialization)
	var defaultString string    // ""
	var defaultInt int          // 0
	var defaultFloat float64    // 0.0
	var defaultBool bool        // false

	// Type conversion
	var x int = 42
	var y float64 = float64(x)
	var z string = fmt.Sprintf("%d", x)

	// Print all variables
	fmt.Printf("Name: %s, Age: %d, Salary: %.2f, Employed: %t\n", name, age, salary, isEmployed)
	fmt.Printf("City: %s, Zip: %d, Temp: %.1f°F, Weekend: %t\n", city, zipCode, temperature, isWeekend)
	fmt.Printf("Full Name: %s %s, Height: %dcm\n", firstName, lastName, height)
	fmt.Printf("Constants - PI: %.5f, Company: %s, Max Users: %d\n", PI, CompanyName, MaxUsers)
	fmt.Printf("Zero values - String: '%s', Int: %d, Float: %.1f, Bool: %t\n", defaultString, defaultInt, defaultFloat, defaultBool)
	fmt.Printf("Type conversion - Int: %d, Float: %.1f, String: %s\n", x, y, z)
}