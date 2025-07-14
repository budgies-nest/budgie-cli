/*
METADATA:
Description: Demonstrates Go HTTP client operations including GET, POST, PUT, DELETE requests, headers, timeouts, and response handling
Keywords: http, client, request, response, GET, POST, PUT, DELETE, headers, timeout, json-api
Category: networking
Concepts: HTTP client, REST API calls, request/response handling, HTTP methods, error handling
*/

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Structs for API responses
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Website  string `json:"website"`
}

type Post struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type Comment struct {
	PostID int    `json:"postId"`
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

// Function demonstrating basic GET requests
func basicGETRequests() {
	fmt.Println("=== BASIC GET REQUESTS ===")
	
	// Simple GET request
	resp, err := http.Get("https://jsonplaceholder.typicode.com/users/1")
	if err != nil {
		fmt.Printf("Error making GET request: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Content-Type: %s\n", resp.Header.Get("Content-Type"))
	
	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}
	
	// Parse JSON response
	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}
	
	fmt.Printf("User: ID=%d, Name=%s, Email=%s\n", user.ID, user.Name, user.Email)
}

// Function demonstrating custom HTTP client with timeout
func customHTTPClient() {
	fmt.Println("\n=== CUSTOM HTTP CLIENT ===")
	
	// Create custom client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	// Make request with custom client
	resp, err := client.Get("https://jsonplaceholder.typicode.com/posts")
	if err != nil {
		fmt.Printf("Error with custom client: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	// Check status code
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Non-OK status: %s\n", resp.Status)
		return
	}
	
	// Parse JSON array
	var posts []Post
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&posts)
	if err != nil {
		fmt.Printf("Error decoding JSON: %v\n", err)
		return
	}
	
	fmt.Printf("Retrieved %d posts\n", len(posts))
	
	// Show first few posts
	for i, post := range posts[:3] {
		fmt.Printf("Post %d: %s (User ID: %d)\n", i+1, post.Title, post.UserID)
	}
}

// Function demonstrating custom requests with headers
func customRequestsWithHeaders() {
	fmt.Println("\n=== CUSTOM REQUESTS WITH HEADERS ===")
	
	// Create request
	req, err := http.NewRequest("GET", "https://jsonplaceholder.typicode.com/posts/1", nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}
	
	// Add custom headers
	req.Header.Set("User-Agent", "Go-HTTP-Client/1.0")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer fake-token-for-demo")
	
	// Create client and make request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	fmt.Printf("Request headers sent:\n")
	for key, values := range req.Header {
		fmt.Printf("  %s: %s\n", key, strings.Join(values, ", "))
	}
	
	fmt.Printf("\nResponse headers received:\n")
	for key, values := range resp.Header {
		fmt.Printf("  %s: %s\n", key, strings.Join(values, ", "))
	}
	
	// Read and parse response
	var post Post
	err = json.NewDecoder(resp.Body).Decode(&post)
	if err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		return
	}
	
	fmt.Printf("\nPost details: %s\n", post.Title)
}

// Function demonstrating POST requests
func postRequests() {
	fmt.Println("\n=== POST REQUESTS ===")
	
	// Create new post data
	newPost := Post{
		UserID: 1,
		Title:  "My New Post",
		Body:   "This is the content of my new post created via HTTP POST.",
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(newPost)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}
	
	// Create POST request
	req, err := http.NewRequest("POST", "https://jsonplaceholder.typicode.com/posts", 
		bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating POST request: %v\n", err)
		return
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	
	// Make request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making POST request: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	fmt.Printf("POST Status: %s\n", resp.Status)
	
	// Parse response
	var createdPost Post
	err = json.NewDecoder(resp.Body).Decode(&createdPost)
	if err != nil {
		fmt.Printf("Error parsing POST response: %v\n", err)
		return
	}
	
	fmt.Printf("Created post: ID=%d, Title=%s\n", createdPost.ID, createdPost.Title)
}

// Function demonstrating PUT requests
func putRequests() {
	fmt.Println("\n=== PUT REQUESTS ===")
	
	// Update existing post
	updatedPost := Post{
		UserID: 1,
		ID:     1,
		Title:  "Updated Post Title",
		Body:   "This post has been updated via HTTP PUT request.",
	}
	
	// Convert to JSON
	jsonData, err := json.Marshal(updatedPost)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}
	
	// Create PUT request
	req, err := http.NewRequest("PUT", "https://jsonplaceholder.typicode.com/posts/1", 
		bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating PUT request: %v\n", err)
		return
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	
	// Make request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making PUT request: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	fmt.Printf("PUT Status: %s\n", resp.Status)
	
	// Parse response
	var responsePost Post
	err = json.NewDecoder(resp.Body).Decode(&responsePost)
	if err != nil {
		fmt.Printf("Error parsing PUT response: %v\n", err)
		return
	}
	
	fmt.Printf("Updated post: ID=%d, Title=%s\n", responsePost.ID, responsePost.Title)
}

// Function demonstrating DELETE requests
func deleteRequests() {
	fmt.Println("\n=== DELETE REQUESTS ===")
	
	// Create DELETE request
	req, err := http.NewRequest("DELETE", "https://jsonplaceholder.typicode.com/posts/1", nil)
	if err != nil {
		fmt.Printf("Error creating DELETE request: %v\n", err)
		return
	}
	
	// Make request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making DELETE request: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	fmt.Printf("DELETE Status: %s\n", resp.Status)
	
	// For DELETE, often there's no response body or empty response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading DELETE response: %v\n", err)
		return
	}
	
	if len(body) > 0 {
		fmt.Printf("DELETE Response body: %s\n", string(body))
	} else {
		fmt.Println("DELETE completed successfully (empty response)")
	}
}

// Function demonstrating URL parameters
func urlParameters() {
	fmt.Println("\n=== URL PARAMETERS ===")
	
	// Build URL with parameters
	baseURL := "https://jsonplaceholder.typicode.com/comments"
	params := url.Values{}
	params.Add("postId", "1")
	params.Add("_limit", "3")
	
	fullURL := baseURL + "?" + params.Encode()
	fmt.Printf("URL with parameters: %s\n", fullURL)
	
	// Make request
	resp, err := http.Get(fullURL)
	if err != nil {
		fmt.Printf("Error making request with parameters: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	// Parse response
	var comments []Comment
	err = json.NewDecoder(resp.Body).Decode(&comments)
	if err != nil {
		fmt.Printf("Error parsing comments: %v\n", err)
		return
	}
	
	fmt.Printf("Retrieved %d comments for post 1:\n", len(comments))
	for i, comment := range comments {
		fmt.Printf("  Comment %d: %s (by %s)\n", i+1, comment.Name, comment.Email)
	}
}

// Function demonstrating form data submission
func formDataSubmission() {
	fmt.Println("\n=== FORM DATA SUBMISSION ===")
	
	// Prepare form data
	formData := url.Values{}
	formData.Set("title", "Form Post Title")
	formData.Set("body", "This post was created using form data")
	formData.Set("userId", "1")
	
	// Create POST request with form data
	resp, err := http.PostForm("https://jsonplaceholder.typicode.com/posts", formData)
	if err != nil {
		fmt.Printf("Error posting form data: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	fmt.Printf("Form POST Status: %s\n", resp.Status)
	
	// Parse response
	var post Post
	err = json.NewDecoder(resp.Body).Decode(&post)
	if err != nil {
		fmt.Printf("Error parsing form response: %v\n", err)
		return
	}
	
	fmt.Printf("Created post from form: ID=%d, Title=%s\n", post.ID, post.Title)
}

// Function demonstrating request context and cancellation
func requestContext() {
	fmt.Println("\n=== REQUEST CONTEXT ===")
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", 
		"https://jsonplaceholder.typicode.com/posts", nil)
	if err != nil {
		fmt.Printf("Error creating request with context: %v\n", err)
		return
	}
	
	// Make request
	client := &http.Client{}
	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)
	
	if err != nil {
		fmt.Printf("Request failed after %v: %v\n", duration, err)
		return
	}
	defer resp.Body.Close()
	
	fmt.Printf("Request completed in %v\n", duration)
	fmt.Printf("Status: %s\n", resp.Status)
	
	// Demonstrate cancellation
	fmt.Println("\nDemonstrating request cancellation:")
	ctx2, cancel2 := context.WithCancel(context.Background())
	
	// Cancel after 100ms
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel2()
	}()
	
	req2, err := http.NewRequestWithContext(ctx2, "GET", 
		"https://jsonplaceholder.typicode.com/posts", nil)
	if err != nil {
		fmt.Printf("Error creating cancellable request: %v\n", err)
		return
	}
	
	start = time.Now()
	_, err = client.Do(req2)
	duration = time.Since(start)
	
	if err != nil {
		fmt.Printf("Request cancelled after %v: %v\n", duration, err)
	}
}

// Function demonstrating error handling
func errorHandling() {
	fmt.Println("\n=== ERROR HANDLING ===")
	
	// Test different error scenarios
	testCases := []struct {
		name string
		url  string
	}{
		{"Valid URL", "https://jsonplaceholder.typicode.com/users/1"},
		{"Invalid URL", "not-a-valid-url"},
		{"Non-existent domain", "https://this-domain-does-not-exist-123456.com"},
		{"404 Not Found", "https://jsonplaceholder.typicode.com/users/999999"},
		{"500 Server Error", "https://httpstat.us/500"},
	}
	
	client := &http.Client{Timeout: 5 * time.Second}
	
	for _, tc := range testCases {
		fmt.Printf("\nTesting %s:\n", tc.name)
		
		resp, err := client.Get(tc.url)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
			continue
		}
		defer resp.Body.Close()
		
		fmt.Printf("  Status: %s\n", resp.Status)
		
		// Handle different status codes
		switch {
		case resp.StatusCode >= 200 && resp.StatusCode < 300:
			fmt.Printf("  Success: %s\n", resp.Status)
		case resp.StatusCode >= 400 && resp.StatusCode < 500:
			fmt.Printf("  Client error: %s\n", resp.Status)
		case resp.StatusCode >= 500:
			fmt.Printf("  Server error: %s\n", resp.Status)
		default:
			fmt.Printf("  Unexpected status: %s\n", resp.Status)
		}
		
		// Read response body for errors (often contains error details)
		if resp.StatusCode >= 400 {
			body, err := io.ReadAll(resp.Body)
			if err == nil && len(body) > 0 {
				fmt.Printf("  Error body: %s\n", string(body)[:min(100, len(body))])
			}
		}
	}
}

// Function demonstrating response streaming
func responseStreaming() {
	fmt.Println("\n=== RESPONSE STREAMING ===")
	
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts")
	if err != nil {
		fmt.Printf("Error making streaming request: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	fmt.Println("Streaming response parsing:")
	
	// Use decoder for streaming JSON parsing
	decoder := json.NewDecoder(resp.Body)
	
	// Read opening bracket of array
	token, err := decoder.Token()
	if err != nil {
		fmt.Printf("Error reading opening token: %v\n", err)
		return
	}
	fmt.Printf("Opening token: %v\n", token)
	
	count := 0
	// Read array elements one by one
	for decoder.More() {
		var post Post
		err := decoder.Decode(&post)
		if err != nil {
			fmt.Printf("Error decoding post: %v\n", err)
			break
		}
		
		count++
		if count <= 3 { // Only show first 3
			fmt.Printf("  Streamed post %d: %s\n", count, post.Title)
		}
	}
	
	// Read closing bracket
	token, err = decoder.Token()
	if err != nil {
		fmt.Printf("Error reading closing token: %v\n", err)
		return
	}
	
	fmt.Printf("Closing token: %v\n", token)
	fmt.Printf("Total posts streamed: %d\n", count)
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	basicGETRequests()
	customHTTPClient()
	customRequestsWithHeaders()
	postRequests()
	putRequests()
	deleteRequests()
	urlParameters()
	formDataSubmission()
	requestContext()
	errorHandling()
	responseStreaming()
	
	fmt.Println("\n=== HTTP CLIENT BEST PRACTICES ===")
	fmt.Println("1. Always set timeouts for HTTP clients")
	fmt.Println("2. Always close response bodies")
	fmt.Println("3. Check status codes before processing responses")
	fmt.Println("4. Use context for request cancellation")
	fmt.Println("5. Handle network errors gracefully")
	fmt.Println("6. Set appropriate headers (User-Agent, Content-Type)")
	fmt.Println("7. Use streaming for large responses")
	fmt.Println("8. Reuse HTTP clients when possible")
	fmt.Println("9. Handle redirects appropriately")
	fmt.Println("10. Validate and sanitize URLs")
	
	fmt.Println("\nHTTP client examples completed!")
}