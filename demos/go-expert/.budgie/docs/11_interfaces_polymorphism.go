/*
METADATA:
Description: Demonstrates Go interfaces, polymorphism, type assertions, and interface composition
Keywords: interface, polymorphism, type-assertion, empty-interface, interface-composition, duck-typing
Category: interfaces
Concepts: interface definition, interface implementation, type assertions, polymorphism, interface composition
*/

package main

import (
	"fmt"
	"strconv"
)

// Basic interface definition
type Writer interface {
	Write(data string) error
}

// Another interface
type Reader interface {
	Read() (string, error)
}

// Interface composition
type ReadWriter interface {
	Reader
	Writer
}

// Interface with multiple methods
type Shape interface {
	Area() float64
	Perimeter() float64
	Name() string
}

// Struct implementing Writer interface
type FileWriter struct {
	filename string
}

func (fw FileWriter) Write(data string) error {
	fmt.Printf("Writing to file %s: %s\n", fw.filename, data)
	return nil
}

// Another struct implementing Writer interface
type ConsoleWriter struct{}

func (cw ConsoleWriter) Write(data string) error {
	fmt.Printf("Console output: %s\n", data)
	return nil
}

// Struct implementing Reader interface
type StringReader struct {
	data string
}

func (sr StringReader) Read() (string, error) {
	return sr.data, nil
}

// Struct implementing both Reader and Writer (ReadWriter)
type MemoryBuffer struct {
	data string
}

func (mb *MemoryBuffer) Read() (string, error) {
	return mb.data, nil
}

func (mb *MemoryBuffer) Write(data string) error {
	mb.data = data
	return nil
}

// Shape implementations
type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func (r Rectangle) Name() string {
	return "Rectangle"
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return 3.14159 * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * 3.14159 * c.Radius
}

func (c Circle) Name() string {
	return "Circle"
}

type Triangle struct {
	Base, Height, Side1, Side2 float64
}

func (t Triangle) Area() float64 {
	return 0.5 * t.Base * t.Height
}

func (t Triangle) Perimeter() float64 {
	return t.Base + t.Side1 + t.Side2
}

func (t Triangle) Name() string {
	return "Triangle"
}

// Interface for different animal behaviors
type Animal interface {
	Speak() string
	Move() string
}

type Dog struct {
	Name string
}

func (d Dog) Speak() string {
	return "Woof!"
}

func (d Dog) Move() string {
	return "Running"
}

type Cat struct {
	Name string
}

func (c Cat) Speak() string {
	return "Meow!"
}

func (c Cat) Move() string {
	return "Stalking"
}

type Bird struct {
	Name string
}

func (b Bird) Speak() string {
	return "Chirp!"
}

func (b Bird) Move() string {
	return "Flying"
}

// Interface for type conversion
type Stringer interface {
	String() string
}

type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	return fmt.Sprintf("%s (%d years old)", p.Name, p.Age)
}

type Product struct {
	Name  string
	Price float64
}

func (p Product) String() string {
	return fmt.Sprintf("%s - $%.2f", p.Name, p.Price)
}

// Function that accepts Writer interface
func writeMessage(w Writer, message string) {
	err := w.Write(message)
	if err != nil {
		fmt.Printf("Error writing: %v\n", err)
	}
}

// Function that accepts ReadWriter interface
func processData(rw ReadWriter, newData string) {
	// Read current data
	if data, err := rw.Read(); err == nil {
		fmt.Printf("Current data: %s\n", data)
	}
	
	// Write new data
	if err := rw.Write(newData); err != nil {
		fmt.Printf("Error writing: %v\n", err)
	}
}

// Function that works with shapes
func describeShape(s Shape) {
	fmt.Printf("Shape: %s\n", s.Name())
	fmt.Printf("  Area: %.2f\n", s.Area())
	fmt.Printf("  Perimeter: %.2f\n", s.Perimeter())
}

// Function demonstrating interface slice
func calculateTotalArea(shapes []Shape) float64 {
	total := 0.0
	for _, shape := range shapes {
		total += shape.Area()
	}
	return total
}

// Function using empty interface
func printValue(value interface{}) {
	fmt.Printf("Value: %v, Type: %T\n", value, value)
}

// Function with type assertion
func processInterface(value interface{}) {
	switch v := value.(type) {
	case string:
		fmt.Printf("String: %s (length: %d)\n", v, len(v))
	case int:
		fmt.Printf("Integer: %d (squared: %d)\n", v, v*v)
	case float64:
		fmt.Printf("Float: %.2f (sqrt: %.2f)\n", v, v*0.5)
	case bool:
		fmt.Printf("Boolean: %t\n", v)
	case Person:
		fmt.Printf("Person: %s\n", v.String())
	case Shape:
		fmt.Printf("Shape: %s with area %.2f\n", v.Name(), v.Area())
	default:
		fmt.Printf("Unknown type: %T\n", v)
	}
}

// Function demonstrating comma ok idiom with type assertion
func safeTypeAssertion(value interface{}) {
	// Safe string assertion
	if str, ok := value.(string); ok {
		fmt.Printf("Successfully converted to string: %s\n", str)
	} else {
		fmt.Printf("Value is not a string: %T\n", value)
	}
	
	// Safe integer assertion
	if num, ok := value.(int); ok {
		fmt.Printf("Successfully converted to int: %d\n", num)
	} else {
		fmt.Printf("Value is not an int: %T\n", value)
	}
}

func main() {
	fmt.Println("=== BASIC INTERFACE USAGE ===")
	
	// Different implementations of Writer interface
	fileWriter := FileWriter{filename: "data.txt"}
	consoleWriter := ConsoleWriter{}
	
	writeMessage(fileWriter, "Hello, File!")
	writeMessage(consoleWriter, "Hello, Console!")

	fmt.Println("\n=== INTERFACE COMPOSITION ===")
	
	// Using ReadWriter interface
	buffer := &MemoryBuffer{}
	processData(buffer, "Initial data")
	processData(buffer, "Updated data")

	fmt.Println("\n=== POLYMORPHISM WITH SHAPES ===")
	
	// Different shape implementations
	rectangle := Rectangle{Width: 5, Height: 3}
	circle := Circle{Radius: 4}
	triangle := Triangle{Base: 6, Height: 4, Side1: 5, Side2: 5}
	
	// Using polymorphism
	shapes := []Shape{rectangle, circle, triangle}
	
	for _, shape := range shapes {
		describeShape(shape)
		fmt.Println()
	}
	
	totalArea := calculateTotalArea(shapes)
	fmt.Printf("Total area of all shapes: %.2f\n", totalArea)

	fmt.Println("\n=== ANIMAL INTERFACE EXAMPLE ===")
	
	// Different animal implementations
	dog := Dog{Name: "Buddy"}
	cat := Cat{Name: "Whiskers"}
	bird := Bird{Name: "Tweety"}
	
	animals := []Animal{dog, cat, bird}
	
	for _, animal := range animals {
		fmt.Printf("%T says: %s and is %s\n", 
			animal, animal.Speak(), animal.Move())
	}

	fmt.Println("\n=== EMPTY INTERFACE ===")
	
	// Empty interface can hold any type
	var anything interface{}
	
	anything = 42
	printValue(anything)
	
	anything = "Hello, World!"
	printValue(anything)
	
	anything = 3.14159
	printValue(anything)
	
	anything = true
	printValue(anything)
	
	anything = Person{Name: "Alice", Age: 30}
	printValue(anything)

	fmt.Println("\n=== TYPE ASSERTIONS ===")
	
	// Type assertions with switch
	values := []interface{}{
		"Hello",
		42,
		3.14,
		true,
		Person{Name: "Bob", Age: 25},
		Rectangle{Width: 3, Height: 4},
		[]int{1, 2, 3},
	}
	
	for _, value := range values {
		processInterface(value)
	}

	fmt.Println("\n=== SAFE TYPE ASSERTIONS ===")
	
	testValues := []interface{}{
		"test string",
		123,
		3.14,
		true,
	}
	
	for _, value := range testValues {
		fmt.Printf("Testing value: %v\n", value)
		safeTypeAssertion(value)
		fmt.Println()
	}

	fmt.Println("\n=== INTERFACE WITH CUSTOM STRING METHOD ===")
	
	// Types implementing Stringer interface
	person := Person{Name: "Alice", Age: 30}
	product := Product{Name: "Laptop", Price: 999.99}
	
	stringers := []Stringer{person, product}
	
	for _, s := range stringers {
		fmt.Printf("String representation: %s\n", s.String())
	}

	fmt.Println("\n=== INTERFACE NIL VALUES ===")
	
	var writer Writer
	fmt.Printf("Nil interface: %v\n", writer)
	fmt.Printf("Is nil: %t\n", writer == nil)
	
	// Assign nil pointer of concrete type
	var fw *FileWriter
	writer = fw
	fmt.Printf("Interface with nil pointer: %v\n", writer)
	fmt.Printf("Is nil: %t\n", writer == nil) // false! interface is not nil
	
	// Check for nil pointer inside interface
	if writer != nil {
		// This would panic if we tried to call methods on nil pointer
		fmt.Printf("Interface is not nil, but pointer inside might be nil\n")
	}

	fmt.Println("\n=== INTERFACE CONVERSION ===")
	
	// Converting between interfaces
	var shape Shape = Rectangle{Width: 5, Height: 3}
	
	// Check if shape also implements Stringer (it doesn't in this example)
	if stringer, ok := shape.(Stringer); ok {
		fmt.Printf("Shape as string: %s\n", stringer.String())
	} else {
		fmt.Printf("Shape does not implement Stringer interface\n")
	}
	
	// Type assertion to concrete type
	if rect, ok := shape.(Rectangle); ok {
		fmt.Printf("Width: %.2f, Height: %.2f\n", rect.Width, rect.Height)
	}

	fmt.Println("\n=== INTERFACE EMBEDDING ===")
	
	// Example of interface with embedded interfaces
	type Closer interface {
		Close() error
	}
	
	type ReadWriteCloser interface {
		Reader
		Writer
		Closer
	}
	
	// This would be implemented by types that satisfy all three interfaces
	fmt.Println("ReadWriteCloser interface embeds Reader, Writer, and Closer")

	fmt.Println("\n=== FUNCTIONAL INTERFACES ===")
	
	// Interface for function types
	type Processor interface {
		Process(string) string
	}
	
	// Function type implementing interface
	type ProcessorFunc func(string) string
	
	func (f ProcessorFunc) Process(s string) string {
		return f(s)
	}
	
	// Using function as interface
	upperProcessor := ProcessorFunc(func(s string) string {
		return fmt.Sprintf("PROCESSED: %s", s)
	})
	
	result := upperProcessor.Process("hello world")
	fmt.Printf("Processed result: %s\n", result)
}