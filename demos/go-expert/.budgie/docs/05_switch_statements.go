/*
METADATA:
Description: Demonstrates Go switch statements including basic switch, switch with expressions, type switch, and fallthrough
Keywords: switch, case, default, fallthrough, type-switch, expression-switch, multiple-cases
Category: control-flow
Concepts: switch statements, pattern matching, type assertion, fallthrough behavior
*/

package main

import "fmt"

func main() {
	// Basic switch statement
	day := 3
	fmt.Println("Basic switch:")
	switch day {
	case 1:
		fmt.Println("Monday")
	case 2:
		fmt.Println("Tuesday")
	case 3:
		fmt.Println("Wednesday")
	case 4:
		fmt.Println("Thursday")
	case 5:
		fmt.Println("Friday")
	case 6, 7:
		fmt.Println("Weekend")
	default:
		fmt.Println("Invalid day")
	}

	// Switch with expressions
	score := 85
	fmt.Println("\nSwitch with expressions:")
	switch {
	case score >= 90:
		fmt.Println("Grade: A")
	case score >= 80:
		fmt.Println("Grade: B")
	case score >= 70:
		fmt.Println("Grade: C")
	case score >= 60:
		fmt.Println("Grade: D")
	default:
		fmt.Println("Grade: F")
	}

	// Switch with initialization
	fmt.Println("\nSwitch with initialization:")
	switch grade := score / 10; grade {
	case 10, 9:
		fmt.Println("Excellent")
	case 8:
		fmt.Println("Good")
	case 7:
		fmt.Println("Average")
	case 6:
		fmt.Println("Below Average")
	default:
		fmt.Println("Poor")
	}

	// Switch with multiple values in case
	month := "December"
	fmt.Println("\nSwitch with multiple values:")
	switch month {
	case "December", "January", "February":
		fmt.Println("Winter")
	case "March", "April", "May":
		fmt.Println("Spring")
	case "June", "July", "August":
		fmt.Println("Summer")
	case "September", "October", "November":
		fmt.Println("Fall")
	default:
		fmt.Println("Invalid month")
	}

	// Switch with fallthrough
	fmt.Println("\nSwitch with fallthrough:")
	number := 2
	switch number {
	case 1:
		fmt.Println("One")
		fallthrough
	case 2:
		fmt.Println("Two or fell through from One")
		fallthrough
	case 3:
		fmt.Println("Three or fell through")
	default:
		fmt.Println("Other number")
	}

	// Type switch
	fmt.Println("\nType switch:")
	var value interface{} = 42
	switch v := value.(type) {
	case int:
		fmt.Printf("Integer: %d\n", v)
	case string:
		fmt.Printf("String: %s\n", v)
	case bool:
		fmt.Printf("Boolean: %t\n", v)
	case float64:
		fmt.Printf("Float: %.2f\n", v)
	default:
		fmt.Printf("Unknown type: %T\n", v)
	}

	// Type switch with multiple types
	fmt.Println("\nType switch with multiple types:")
	values := []interface{}{42, "hello", 3.14, true, []int{1, 2, 3}}
	for i, val := range values {
		switch v := val.(type) {
		case int, int64:
			fmt.Printf("values[%d] is an integer: %v\n", i, v)
		case string:
			fmt.Printf("values[%d] is a string: %s\n", i, v)
		case float64, float32:
			fmt.Printf("values[%d] is a float: %v\n", i, v)
		case bool:
			fmt.Printf("values[%d] is a boolean: %t\n", i, v)
		default:
			fmt.Printf("values[%d] is of type %T: %v\n", i, v, v)
		}
	}

	// Switch without expression (same as switch true)
	temperature := 75
	fmt.Println("\nSwitch without expression:")
	switch {
	case temperature < 32:
		fmt.Println("Freezing")
	case temperature < 60:
		fmt.Println("Cold")
	case temperature < 80:
		fmt.Println("Comfortable")
	case temperature < 90:
		fmt.Println("Warm")
	default:
		fmt.Println("Hot")
	}

	// Switch with function calls
	fmt.Println("\nSwitch with function calls:")
	operation := "add"
	a, b := 10, 5
	switch operation {
	case "add":
		fmt.Printf("%d + %d = %d\n", a, b, a+b)
	case "subtract":
		fmt.Printf("%d - %d = %d\n", a, b, a-b)
	case "multiply":
		fmt.Printf("%d * %d = %d\n", a, b, a*b)
	case "divide":
		if b != 0 {
			fmt.Printf("%d / %d = %d\n", a, b, a/b)
		} else {
			fmt.Println("Cannot divide by zero")
		}
	default:
		fmt.Println("Unknown operation")
	}
}