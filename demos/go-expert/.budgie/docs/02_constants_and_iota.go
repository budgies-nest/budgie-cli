/*
METADATA:
Description: Demonstrates Go constants, iota enumeration, and typed constants
Keywords: constants, const, iota, enumeration, typed-constants, untyped-constants
Category: basics
Concepts: constant declaration, iota counter, enumerated values, constant expressions
*/

package main

import "fmt"

func main() {
	// Basic constants
	const message = "Hello, World!"
	const pi = 3.14159
	const maxRetries = 3

	// Typed constants
	const typedString string = "Typed constant"
	const typedInt int = 100
	const typedFloat float64 = 99.99

	// Multiple constants
	const (
		StatusOK       = 200
		StatusNotFound = 404
		StatusError    = 500
	)

	// iota enumeration
	const (
		Sunday = iota    // 0
		Monday           // 1
		Tuesday          // 2
		Wednesday        // 3
		Thursday         // 4
		Friday           // 5
		Saturday         // 6
	)

	// iota with expressions
	const (
		_  = iota             // ignore first value
		KB = 1 << (10 * iota) // 1024
		MB                    // 1048576
		GB                    // 1073741824
		TB                    // 1099511627776
	)

	// iota with custom values
	const (
		Low = iota * 10    // 0
		Medium             // 10
		High               // 20
		Critical           // 30
	)

	// String constants for enum-like behavior
	const (
		ModeRead  = "read"
		ModeWrite = "write"
		ModeAdmin = "admin"
	)

	// Constant expressions
	const hoursPerDay = 24
	const minutesPerHour = 60
	const minutesPerDay = hoursPerDay * minutesPerHour

	fmt.Printf("Basic constants: %s, %.5f, %d\n", message, pi, maxRetries)
	fmt.Printf("Typed constants: %s, %d, %.2f\n", typedString, typedInt, typedFloat)
	fmt.Printf("HTTP Status codes: OK=%d, NotFound=%d, Error=%d\n", StatusOK, StatusNotFound, StatusError)
	fmt.Printf("Days: Sunday=%d, Wednesday=%d, Saturday=%d\n", Sunday, Wednesday, Saturday)
	fmt.Printf("Storage sizes: KB=%d, MB=%d, GB=%d\n", KB, MB, GB)
	fmt.Printf("Priority levels: Low=%d, Medium=%d, High=%d, Critical=%d\n", Low, Medium, High, Critical)
	fmt.Printf("Modes: %s, %s, %s\n", ModeRead, ModeWrite, ModeAdmin)
	fmt.Printf("Time calculation: %d minutes per day\n", minutesPerDay)
}