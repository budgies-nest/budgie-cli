/*
METADATA:
Description: Demonstrates Go JSON handling including marshaling, unmarshaling, custom JSON tags, and JSON streaming
Keywords: json, marshal, unmarshal, encoding, decoding, struct-tags, json-streaming, custom-json
Category: data-serialization
Concepts: JSON encoding/decoding, struct tags, custom JSON handling, streaming JSON, JSON validation
*/

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// Basic struct for JSON demonstration
type Person struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Email    string `json:"email"`
	IsActive bool   `json:"is_active"`
}

// Struct with various JSON tags
type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Description string    `json:"description,omitempty"` // Omit if empty
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Tags        []string  `json:"tags"`
	Metadata    map[string]interface{} `json:"metadata"`
	Secret      string    `json:"-"` // Never include in JSON
}

// Struct with nested objects
type Company struct {
	Name      string    `json:"name"`
	Founded   int       `json:"founded"`
	Employees []Person  `json:"employees"`
	Address   Address   `json:"address"`
	Products  []Product `json:"products"`
}

type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}

// Custom JSON marshaling and unmarshaling
type CustomDate struct {
	time.Time
}

func (cd CustomDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(cd.Time.Format("2006-01-02"))
}

func (cd *CustomDate) UnmarshalJSON(data []byte) error {
	var dateStr string
	if err := json.Unmarshal(data, &dateStr); err != nil {
		return err
	}
	
	parsed, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err
	}
	
	cd.Time = parsed
	return nil
}

type Event struct {
	ID       int        `json:"id"`
	Name     string     `json:"name"`
	Date     CustomDate `json:"date"`
	Location string     `json:"location"`
}

// Interface for custom JSON handling
type JSONCustomizer interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

// Struct implementing custom JSON logic
type Temperature struct {
	Celsius float64
}

func (t Temperature) MarshalJSON() ([]byte, error) {
	// Always output in both Celsius and Fahrenheit
	data := map[string]float64{
		"celsius":    t.Celsius,
		"fahrenheit": t.Celsius*9/5 + 32,
	}
	return json.Marshal(data)
}

func (t *Temperature) UnmarshalJSON(data []byte) error {
	var temp map[string]float64
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	
	if celsius, ok := temp["celsius"]; ok {
		t.Celsius = celsius
	} else if fahrenheit, ok := temp["fahrenheit"]; ok {
		t.Celsius = (fahrenheit - 32) * 5 / 9
	} else {
		return fmt.Errorf("temperature must have either celsius or fahrenheit")
	}
	
	return nil
}

// Function demonstrating basic JSON operations
func basicJSONOperations() {
	fmt.Println("=== BASIC JSON OPERATIONS ===")
	
	// Create a person
	person := Person{
		Name:     "Alice Johnson",
		Age:      30,
		Email:    "alice@example.com",
		IsActive: true,
	}
	
	// Marshal to JSON
	jsonData, err := json.Marshal(person)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}
	
	fmt.Printf("Marshaled JSON: %s\n", string(jsonData))
	
	// Marshal with indentation
	prettyJSON, err := json.MarshalIndent(person, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling pretty JSON: %v\n", err)
		return
	}
	
	fmt.Printf("Pretty JSON:\n%s\n", string(prettyJSON))
	
	// Unmarshal JSON
	jsonString := `{"name":"Bob Smith","age":25,"email":"bob@example.com","is_active":false}`
	var newPerson Person
	
	err = json.Unmarshal([]byte(jsonString), &newPerson)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		return
	}
	
	fmt.Printf("Unmarshaled person: %+v\n", newPerson)
}

// Function demonstrating JSON tags
func jsonTags() {
	fmt.Println("\n=== JSON TAGS ===")
	
	product := Product{
		ID:          1,
		Name:        "Laptop",
		Price:       999.99,
		Description: "", // This will be omitted due to omitempty tag
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Tags:        []string{"electronics", "computer", "portable"},
		Metadata: map[string]interface{}{
			"warranty_years": 2,
			"manufacturer":   "TechCorp",
			"weight_kg":      1.5,
		},
		Secret: "This won't appear in JSON", // Will be ignored due to "-" tag
	}
	
	jsonData, err := json.MarshalIndent(product, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling product: %v\n", err)
		return
	}
	
	fmt.Printf("Product JSON:\n%s\n", string(jsonData))
	
	// Demonstrate omitempty with empty description
	product.Description = "High-performance laptop"
	jsonData, _ = json.MarshalIndent(product, "", "  ")
	fmt.Printf("Product with description:\n%s\n", string(jsonData))
}

// Function demonstrating nested JSON structures
func nestedJSON() {
	fmt.Println("\n=== NESTED JSON STRUCTURES ===")
	
	company := Company{
		Name:    "TechCorp Inc.",
		Founded: 2010,
		Employees: []Person{
			{Name: "Alice Johnson", Age: 30, Email: "alice@techcorp.com", IsActive: true},
			{Name: "Bob Smith", Age: 25, Email: "bob@techcorp.com", IsActive: true},
			{Name: "Carol Davis", Age: 35, Email: "carol@techcorp.com", IsActive: false},
		},
		Address: Address{
			Street:  "123 Tech Street",
			City:    "San Francisco",
			State:   "CA",
			ZipCode: "94105",
			Country: "USA",
		},
		Products: []Product{
			{
				ID:        1,
				Name:      "Super Laptop",
				Price:     1299.99,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Tags:      []string{"premium", "laptop"},
				Metadata:  map[string]interface{}{"warranty": "3 years"},
			},
		},
	}
	
	jsonData, err := json.MarshalIndent(company, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling company: %v\n", err)
		return
	}
	
	fmt.Printf("Company JSON:\n%s\n", string(jsonData))
	
	// Unmarshal the JSON back
	var unmarshaledCompany Company
	err = json.Unmarshal(jsonData, &unmarshaledCompany)
	if err != nil {
		fmt.Printf("Error unmarshaling company: %v\n", err)
		return
	}
	
	fmt.Printf("Unmarshaled company name: %s\n", unmarshaledCompany.Name)
	fmt.Printf("Number of employees: %d\n", len(unmarshaledCompany.Employees))
	fmt.Printf("First employee: %s\n", unmarshaledCompany.Employees[0].Name)
}

// Function demonstrating custom JSON marshaling
func customJSONMarshaling() {
	fmt.Println("\n=== CUSTOM JSON MARSHALING ===")
	
	// Custom date handling
	event := Event{
		ID:       1,
		Name:     "Go Conference",
		Date:     CustomDate{time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)},
		Location: "San Francisco",
	}
	
	jsonData, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling event: %v\n", err)
		return
	}
	
	fmt.Printf("Event with custom date:\n%s\n", string(jsonData))
	
	// Unmarshal custom date
	jsonString := `{"id":2,"name":"Tech Meetup","date":"2024-04-20","location":"New York"}`
	var newEvent Event
	
	err = json.Unmarshal([]byte(jsonString), &newEvent)
	if err != nil {
		fmt.Printf("Error unmarshaling event: %v\n", err)
		return
	}
	
	fmt.Printf("Unmarshaled event: ID=%d, Name=%s, Date=%s, Location=%s\n",
		newEvent.ID, newEvent.Name, newEvent.Date.Format("2006-01-02"), newEvent.Location)
	
	// Custom temperature handling
	temp := Temperature{Celsius: 25.0}
	
	tempJSON, err := json.MarshalIndent(temp, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling temperature: %v\n", err)
		return
	}
	
	fmt.Printf("Temperature JSON:\n%s\n", string(tempJSON))
	
	// Unmarshal temperature from Fahrenheit
	fahrenheitJSON := `{"fahrenheit": 77}`
	var newTemp Temperature
	
	err = json.Unmarshal([]byte(fahrenheitJSON), &newTemp)
	if err != nil {
		fmt.Printf("Error unmarshaling temperature: %v\n", err)
		return
	}
	
	fmt.Printf("Temperature from Fahrenheit: %.2f°C\n", newTemp.Celsius)
}

// Function demonstrating JSON with interfaces
func jsonWithInterfaces() {
	fmt.Println("\n=== JSON WITH INTERFACES ===")
	
	// Using interface{} for dynamic JSON
	data := map[string]interface{}{
		"string_field":  "Hello, World!",
		"number_field":  42,
		"float_field":   3.14159,
		"bool_field":    true,
		"array_field":   []interface{}{1, 2, 3, "four", 5.0},
		"object_field": map[string]interface{}{
			"nested_string": "nested value",
			"nested_number": 123,
		},
		"null_field": nil,
	}
	
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling interface data: %v\n", err)
		return
	}
	
	fmt.Printf("Interface JSON:\n%s\n", string(jsonData))
	
	// Unmarshal into interface{}
	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		fmt.Printf("Error unmarshaling to interface: %v\n", err)
		return
	}
	
	// Type assertions to access values
	fmt.Println("Accessing unmarshaled interface values:")
	if str, ok := result["string_field"].(string); ok {
		fmt.Printf("  String field: %s\n", str)
	}
	
	if num, ok := result["number_field"].(float64); ok { // JSON numbers are float64
		fmt.Printf("  Number field: %.0f\n", num)
	}
	
	if arr, ok := result["array_field"].([]interface{}); ok {
		fmt.Printf("  Array field: %v\n", arr)
	}
	
	if obj, ok := result["object_field"].(map[string]interface{}); ok {
		fmt.Printf("  Object field: %v\n", obj)
		if nestedStr, ok := obj["nested_string"].(string); ok {
			fmt.Printf("    Nested string: %s\n", nestedStr)
		}
	}
}

// Function demonstrating JSON streaming
func jsonStreaming() {
	fmt.Println("\n=== JSON STREAMING ===")
	
	// Create a JSON file with multiple objects
	filename := "people.json"
	people := []Person{
		{Name: "Alice", Age: 30, Email: "alice@example.com", IsActive: true},
		{Name: "Bob", Age: 25, Email: "bob@example.com", IsActive: false},
		{Name: "Carol", Age: 35, Email: "carol@example.com", IsActive: true},
	}
	
	// Write JSON array to file
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()
	defer os.Remove(filename)
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	
	err = encoder.Encode(people)
	if err != nil {
		fmt.Printf("Error encoding JSON to file: %v\n", err)
		return
	}
	
	fmt.Printf("Wrote %d people to %s\n", len(people), filename)
	
	// Read JSON from file using decoder
	file, err = os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()
	
	decoder := json.NewDecoder(file)
	
	var readPeople []Person
	err = decoder.Decode(&readPeople)
	if err != nil {
		fmt.Printf("Error decoding JSON from file: %v\n", err)
		return
	}
	
	fmt.Printf("Read %d people from file:\n", len(readPeople))
	for i, person := range readPeople {
		fmt.Printf("  %d: %s (age %d)\n", i+1, person.Name, person.Age)
	}
	
	// Streaming multiple JSON objects
	fmt.Println("\nStreaming individual JSON objects:")
	
	// Create file with individual JSON objects (one per line)
	streamFile := "stream.json"
	file, err = os.Create(streamFile)
	if err != nil {
		fmt.Printf("Error creating stream file: %v\n", err)
		return
	}
	defer os.Remove(streamFile)
	
	encoder = json.NewEncoder(file)
	for _, person := range people {
		err = encoder.Encode(person)
		if err != nil {
			fmt.Printf("Error encoding person: %v\n", err)
			continue
		}
	}
	file.Close()
	
	// Read streaming JSON objects
	file, err = os.Open(streamFile)
	if err != nil {
		fmt.Printf("Error opening stream file: %v\n", err)
		return
	}
	defer file.Close()
	
	decoder = json.NewDecoder(file)
	count := 0
	
	for decoder.More() {
		var person Person
		err = decoder.Decode(&person)
		if err != nil {
			fmt.Printf("Error decoding person: %v\n", err)
			break
		}
		count++
		fmt.Printf("  Streamed person %d: %s\n", count, person.Name)
	}
}

// Function demonstrating JSON validation and error handling
func jsonValidation() {
	fmt.Println("\n=== JSON VALIDATION ===")
	
	// Valid JSON
	validJSON := `{"name":"Alice","age":30,"email":"alice@example.com","is_active":true}`
	var person Person
	
	err := json.Unmarshal([]byte(validJSON), &person)
	if err != nil {
		fmt.Printf("Error with valid JSON: %v\n", err)
	} else {
		fmt.Printf("Valid JSON parsed successfully: %+v\n", person)
	}
	
	// Invalid JSON syntax
	invalidJSON := `{"name":"Alice","age":30,"email":"alice@example.com","is_active":true` // Missing closing brace
	err = json.Unmarshal([]byte(invalidJSON), &person)
	if err != nil {
		fmt.Printf("Invalid JSON syntax error: %v\n", err)
	}
	
	// JSON with wrong types
	wrongTypeJSON := `{"name":"Alice","age":"thirty","email":"alice@example.com","is_active":true}`
	err = json.Unmarshal([]byte(wrongTypeJSON), &person)
	if err != nil {
		fmt.Printf("Wrong type error: %v\n", err)
	}
	
	// JSON with missing fields (will use zero values)
	incompleteJSON := `{"name":"Bob"}`
	err = json.Unmarshal([]byte(incompleteJSON), &person)
	if err != nil {
		fmt.Printf("Error with incomplete JSON: %v\n", err)
	} else {
		fmt.Printf("Incomplete JSON parsed: %+v\n", person)
	}
	
	// JSON with extra fields (will be ignored)
	extraFieldsJSON := `{"name":"Carol","age":35,"email":"carol@example.com","is_active":true,"extra_field":"ignored"}`
	err = json.Unmarshal([]byte(extraFieldsJSON), &person)
	if err != nil {
		fmt.Printf("Error with extra fields: %v\n", err)
	} else {
		fmt.Printf("JSON with extra fields parsed: %+v\n", person)
	}
}

// Function demonstrating raw JSON messages
func rawJSONMessages() {
	fmt.Println("\n=== RAW JSON MESSAGES ===")
	
	// Using json.RawMessage for delayed parsing
	type Response struct {
		Status string          `json:"status"`
		Data   json.RawMessage `json:"data"`
	}
	
	jsonString := `{
		"status": "success",
		"data": {
			"name": "Alice",
			"age": 30,
			"email": "alice@example.com"
		}
	}`
	
	var response Response
	err := json.Unmarshal([]byte(jsonString), &response)
	if err != nil {
		fmt.Printf("Error unmarshaling response: %v\n", err)
		return
	}
	
	fmt.Printf("Response status: %s\n", response.Status)
	fmt.Printf("Raw data: %s\n", string(response.Data))
	
	// Parse the raw data based on status
	if response.Status == "success" {
		var person Person
		err = json.Unmarshal(response.Data, &person)
		if err != nil {
			fmt.Printf("Error unmarshaling person data: %v\n", err)
			return
		}
		fmt.Printf("Parsed person from raw data: %+v\n", person)
	}
}

func main() {
	basicJSONOperations()
	jsonTags()
	nestedJSON()
	customJSONMarshaling()
	jsonWithInterfaces()
	jsonStreaming()
	jsonValidation()
	rawJSONMessages()
	
	fmt.Println("\n=== JSON BEST PRACTICES ===")
	fmt.Println("1. Use struct tags to control JSON field names")
	fmt.Println("2. Use omitempty to exclude empty fields")
	fmt.Println("3. Use json:\"-\" to exclude fields entirely")
	fmt.Println("4. Handle JSON errors properly")
	fmt.Println("5. Use json.RawMessage for delayed parsing")
	fmt.Println("6. Use streaming for large JSON data")
	fmt.Println("7. Implement custom marshaling for complex types")
	fmt.Println("8. Validate JSON structure and types")
	fmt.Println("9. Use interface{} sparingly, prefer typed structs")
	fmt.Println("10. Consider performance implications of reflection")
	
	fmt.Println("\nJSON examples completed!")
}