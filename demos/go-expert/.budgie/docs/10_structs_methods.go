/*
METADATA:
Description: Demonstrates Go structs, methods, embedded structs, and struct patterns
Keywords: struct, method, receiver, embedded-struct, composition, constructor, pointer-receiver, value-receiver
Category: data-structures
Concepts: struct definition, methods, method receivers, struct embedding, composition over inheritance
*/

package main

import "fmt"

// Basic struct definition
type Person struct {
	Name string
	Age  int
	City string
}

// Struct with different field types
type Product struct {
	ID          int
	Name        string
	Price       float64
	InStock     bool
	Categories  []string
	Attributes  map[string]string
}

// Method with value receiver
func (p Person) GetInfo() string {
	return fmt.Sprintf("%s is %d years old and lives in %s", p.Name, p.Age, p.City)
}

// Method with pointer receiver (can modify the struct)
func (p *Person) HaveBirthday() {
	p.Age++
}

// Method with pointer receiver for efficiency (large structs)
func (p *Person) MoveTo(newCity string) {
	p.City = newCity
}

// Method that returns multiple values
func (p Person) IsAdult() (bool, string) {
	if p.Age >= 18 {
		return true, "adult"
	}
	return false, "minor"
}

// Constructor function (convention: New + StructName)
func NewPerson(name string, age int, city string) *Person {
	return &Person{
		Name: name,
		Age:  age,
		City: city,
	}
}

// Struct with embedded structs (composition)
type Address struct {
	Street   string
	City     string
	State    string
	ZipCode  string
	Country  string
}

type Employee struct {
	Person                    // Embedded struct
	Address                   // Embedded struct
	ID       int
	Department string
	Salary   float64
}

// Method for embedded struct
func (a Address) FullAddress() string {
	return fmt.Sprintf("%s, %s, %s %s, %s", a.Street, a.City, a.State, a.ZipCode, a.Country)
}

// Method for Employee
func (e Employee) GetDetails() string {
	return fmt.Sprintf("Employee %d: %s, Department: %s, Salary: $%.2f", 
		e.ID, e.Name, e.Department, e.Salary)
}

// Interface and struct methods
type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

type Circle struct {
	Radius float64
}

// Rectangle methods implementing Shape interface
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Circle methods implementing Shape interface
func (c Circle) Area() float64 {
	return 3.14159 * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * 3.14159 * c.Radius
}

// Struct with unexported fields (encapsulation)
type BankAccount struct {
	accountNumber string  // unexported (private)
	balance       float64 // unexported (private)
	Owner         string  // exported (public)
}

// Constructor for BankAccount
func NewBankAccount(owner, accountNumber string, initialBalance float64) *BankAccount {
	return &BankAccount{
		accountNumber: accountNumber,
		balance:       initialBalance,
		Owner:         owner,
	}
}

// Getter methods (accessing private fields)
func (ba *BankAccount) GetBalance() float64 {
	return ba.balance
}

func (ba *BankAccount) GetAccountNumber() string {
	return ba.accountNumber
}

// Methods that modify private fields
func (ba *BankAccount) Deposit(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("deposit amount must be positive")
	}
	ba.balance += amount
	return nil
}

func (ba *BankAccount) Withdraw(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("withdrawal amount must be positive")
	}
	if amount > ba.balance {
		return fmt.Errorf("insufficient funds")
	}
	ba.balance -= amount
	return nil
}

func main() {
	fmt.Println("=== BASIC STRUCT USAGE ===")
	
	// Creating structs
	var person1 Person
	person1.Name = "Alice"
	person1.Age = 30
	person1.City = "New York"
	
	// Struct literal
	person2 := Person{
		Name: "Bob",
		Age:  25,
		City: "London",
	}
	
	// Struct literal with positional values (not recommended)
	person3 := Person{"Carol", 35, "Tokyo"}
	
	fmt.Printf("Person1: %+v\n", person1)
	fmt.Printf("Person2: %+v\n", person2)
	fmt.Printf("Person3: %+v\n", person3)
	
	// Using constructor
	person4 := NewPerson("Dave", 28, "Sydney")
	fmt.Printf("Person4: %+v\n", *person4)

	fmt.Println("\n=== STRUCT METHODS ===")
	
	// Calling methods
	fmt.Println(person1.GetInfo())
	fmt.Println(person2.GetInfo())
	
	// Method with pointer receiver
	fmt.Printf("Before birthday: %s is %d\n", person1.Name, person1.Age)
	person1.HaveBirthday()
	fmt.Printf("After birthday: %s is %d\n", person1.Name, person1.Age)
	
	// Moving person
	fmt.Printf("Before move: %s\n", person2.GetInfo())
	person2.MoveTo("Paris")
	fmt.Printf("After move: %s\n", person2.GetInfo())
	
	// Method with multiple returns
	isAdult, category := person1.IsAdult()
	fmt.Printf("%s is an %s: %t\n", person1.Name, category, isAdult)

	fmt.Println("\n=== EMBEDDED STRUCTS ===")
	
	// Creating Employee with embedded structs
	employee := Employee{
		Person: Person{
			Name: "John Doe",
			Age:  30,
			City: "Seattle",
		},
		Address: Address{
			Street:  "123 Main St",
			City:    "Seattle",
			State:   "WA",
			ZipCode: "98101",
			Country: "USA",
		},
		ID:         1001,
		Department: "Engineering",
		Salary:     75000.00,
	}
	
	// Accessing embedded fields directly
	fmt.Printf("Employee name: %s\n", employee.Name) // From embedded Person
	fmt.Printf("Employee age: %d\n", employee.Age)   // From embedded Person
	fmt.Printf("Employee street: %s\n", employee.Street) // From embedded Address
	
	// Calling methods from embedded structs
	fmt.Printf("Info: %s\n", employee.GetInfo()) // Person method
	fmt.Printf("Address: %s\n", employee.FullAddress()) // Address method
	fmt.Printf("Details: %s\n", employee.GetDetails()) // Employee method
	
	// Method on embedded struct can be called directly
	employee.HaveBirthday() // Person method
	fmt.Printf("After birthday: %d years old\n", employee.Age)

	fmt.Println("\n=== INTERFACES AND POLYMORPHISM ===")
	
	// Creating shapes
	rectangle := Rectangle{Width: 5, Height: 3}
	circle := Circle{Radius: 4}
	
	// Using shapes through interface
	shapes := []Shape{rectangle, circle}
	
	for i, shape := range shapes {
		fmt.Printf("Shape %d:\n", i+1)
		fmt.Printf("  Area: %.2f\n", shape.Area())
		fmt.Printf("  Perimeter: %.2f\n", shape.Perimeter())
	}

	fmt.Println("\n=== ENCAPSULATION ===")
	
	// Using BankAccount with private fields
	account := NewBankAccount("Alice Johnson", "ACC-001", 1000.00)
	
	fmt.Printf("Account owner: %s\n", account.Owner)
	fmt.Printf("Account number: %s\n", account.GetAccountNumber())
	fmt.Printf("Initial balance: $%.2f\n", account.GetBalance())
	
	// Deposit money
	err := account.Deposit(500.00)
	if err != nil {
		fmt.Printf("Deposit error: %v\n", err)
	} else {
		fmt.Printf("After deposit: $%.2f\n", account.GetBalance())
	}
	
	// Withdraw money
	err = account.Withdraw(200.00)
	if err != nil {
		fmt.Printf("Withdrawal error: %v\n", err)
	} else {
		fmt.Printf("After withdrawal: $%.2f\n", account.GetBalance())
	}
	
	// Try to withdraw more than balance
	err = account.Withdraw(2000.00)
	if err != nil {
		fmt.Printf("Withdrawal error: %v\n", err)
	}

	fmt.Println("\n=== STRUCT COMPARISON AND COPYING ===")
	
	// Structs are comparable if all fields are comparable
	p1 := Person{Name: "Alice", Age: 30, City: "NYC"}
	p2 := Person{Name: "Alice", Age: 30, City: "NYC"}
	p3 := Person{Name: "Bob", Age: 25, City: "LA"}
	
	fmt.Printf("p1 == p2: %t\n", p1 == p2) // true
	fmt.Printf("p1 == p3: %t\n", p1 == p3) // false
	
	// Struct copying (value semantics)
	p4 := p1  // Copy
	p4.Name = "Alice Smith"
	fmt.Printf("Original: %s\n", p1.Name)
	fmt.Printf("Copy: %s\n", p4.Name)

	fmt.Println("\n=== ANONYMOUS STRUCTS ===")
	
	// Anonymous struct for temporary use
	config := struct {
		Host string
		Port int
		SSL  bool
	}{
		Host: "localhost",
		Port: 8080,
		SSL:  false,
	}
	
	fmt.Printf("Config: %+v\n", config)
	
	// Slice of anonymous structs
	users := []struct {
		Name  string
		Email string
	}{
		{Name: "Alice", Email: "alice@example.com"},
		{Name: "Bob", Email: "bob@example.com"},
	}
	
	fmt.Println("Users:")
	for _, user := range users {
		fmt.Printf("  %s: %s\n", user.Name, user.Email)
	}

	fmt.Println("\n=== STRUCT TAGS ===")
	
	// Struct with tags (used by reflection, JSON, etc.)
	type User struct {
		ID       int    `json:"id" db:"user_id"`
		Username string `json:"username" db:"username"`
		Email    string `json:"email" db:"email_address"`
		Password string `json:"-" db:"password_hash"` // - means exclude from JSON
	}
	
	user := User{
		ID:       1,
		Username: "alice",
		Email:    "alice@example.com",
		Password: "secret",
	}
	
	fmt.Printf("User struct: %+v\n", user)
	// Note: JSON marshaling would be demonstrated in JSON examples
}