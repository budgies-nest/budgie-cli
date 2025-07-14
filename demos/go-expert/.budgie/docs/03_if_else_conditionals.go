/*
METADATA:
Description: Demonstrates Go conditional statements including if/else, if with initialization, and nested conditions
Keywords: if, else, conditional, boolean, comparison, logical-operators, nested-if, if-initialization
Category: control-flow
Concepts: conditional logic, boolean expressions, comparison operators, logical operators
*/

package main

import "fmt"

func main() {
	age := 25
	score := 85
	hasLicense := true
	name := "Alice"

	// Basic if statement
	if age >= 18 {
		fmt.Println("You are an adult")
	}

	// If-else statement
	if age >= 65 {
		fmt.Println("You are a senior citizen")
	} else {
		fmt.Println("You are not a senior citizen")
	}

	// If-else if-else chain
	if score >= 90 {
		fmt.Println("Grade: A")
	} else if score >= 80 {
		fmt.Println("Grade: B")
	} else if score >= 70 {
		fmt.Println("Grade: C")
	} else if score >= 60 {
		fmt.Println("Grade: D")
	} else {
		fmt.Println("Grade: F")
	}

	// If with initialization statement
	if length := len(name); length > 5 {
		fmt.Printf("Name '%s' is long (%d characters)\n", name, length)
	} else {
		fmt.Printf("Name '%s' is short (%d characters)\n", name, length)
	}

	// Logical operators
	if age >= 16 && hasLicense {
		fmt.Println("You can drive")
	}

	if age < 13 || age > 65 {
		fmt.Println("You get a discount")
	}

	if !hasLicense {
		fmt.Println("You need to get a license")
	} else {
		fmt.Println("You have a valid license")
	}

	// Nested if statements
	if age >= 18 {
		if hasLicense {
			fmt.Println("You can rent a car")
		} else {
			fmt.Println("You need a license to rent a car")
		}
	} else {
		fmt.Println("You must be 18 or older to rent a car")
	}

	// Multiple conditions
	income := 50000
	creditScore := 720
	if income >= 40000 && creditScore >= 700 && age >= 21 {
		fmt.Println("Loan approved")
	} else {
		fmt.Println("Loan denied")
	}

	// String comparison
	status := "active"
	if status == "active" {
		fmt.Println("Account is active")
	} else if status == "inactive" {
		fmt.Println("Account is inactive")
	} else {
		fmt.Println("Unknown account status")
	}

	// Complex boolean expression
	isWeekend := true
	temperature := 75
	if (temperature > 70 && temperature < 85) && !isWeekend {
		fmt.Println("Perfect weather for work")
	} else if isWeekend && temperature > 60 {
		fmt.Println("Good weather for outdoor activities")
	}
}