/*
METADATA:
Description: Demonstrates Go HTTP server creation including routing, middleware, handlers, and REST API implementation
Keywords: http, server, handler, middleware, routing, mux, rest-api, json-response, middleware-chain
Category: networking
Concepts: HTTP server, request handling, middleware, routing, JSON APIs, server patterns
*/

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// Data structures for our API
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// In-memory storage for demo
var (
	users   = make(map[int]User)
	usersMu = sync.RWMutex{}
	nextID  = 1
)

// Initialize some sample data
func init() {
	users[1] = User{ID: 1, Name: "Alice Johnson", Email: "alice@example.com", Username: "alice"}
	users[2] = User{ID: 2, Name: "Bob Smith", Email: "bob@example.com", Username: "bob"}
	users[3] = User{ID: 3, Name: "Carol Davis", Email: "carol@example.com", Username: "carol"}
	nextID = 4
}

// Basic HTTP handlers
func basicHandlers() {
	fmt.Println("=== BASIC HTTP HANDLERS ===")
	
	// Simple handler function
	helloHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World! Method: %s, URL: %s", r.Method, r.URL.Path)
	}
	
	// Handler with JSON response
	jsonHandler := func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"message":   "Hello, JSON!",
			"timestamp": time.Now().Format(time.RFC3339),
			"method":    r.Method,
			"path":      r.URL.Path,
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
	
	// Handler with URL parameters
	paramsHandler := func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		name := query.Get("name")
		age := query.Get("age")
		
		if name == "" {
			name = "Anonymous"
		}
		
		response := map[string]string{
			"greeting": fmt.Sprintf("Hello, %s!", name),
			"age":      age,
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
	
	// Register handlers
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/json", jsonHandler)
	http.HandleFunc("/params", paramsHandler)
	
	fmt.Println("Basic handlers registered:")
	fmt.Println("  /hello - Simple text response")
	fmt.Println("  /json - JSON response")
	fmt.Println("  /params - URL parameters example")
}

// Middleware functions
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("%s %s %v", r.Method, r.URL.Path, duration)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		
		// Simple token validation (for demo purposes)
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		
		token := strings.TrimPrefix(auth, "Bearer ")
		if token != "valid-token" {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func rateLimitMiddleware(requests int, window time.Duration) func(http.Handler) http.Handler {
	type client struct {
		requests int
		lastSeen time.Time
	}
	
	clients := make(map[string]*client)
	mu := sync.Mutex{}
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			
			mu.Lock()
			defer mu.Unlock()
			
			now := time.Now()
			c, exists := clients[ip]
			
			if !exists {
				clients[ip] = &client{requests: 1, lastSeen: now}
				next.ServeHTTP(w, r)
				return
			}
			
			if now.Sub(c.lastSeen) > window {
				c.requests = 1
				c.lastSeen = now
				next.ServeHTTP(w, r)
				return
			}
			
			if c.requests >= requests {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			
			c.requests++
			next.ServeHTTP(w, r)
		})
	}
}

// REST API handlers
func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	usersMu.RLock()
	defer usersMu.RUnlock()
	
	userList := make([]User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}
	
	response := APIResponse{
		Success: true,
		Data:    userList,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	id, err := strconv.Atoi(path)
	if err != nil {
		response := APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	usersMu.RLock()
	user, exists := users[id]
	usersMu.RUnlock()
	
	if !exists {
		response := APIResponse{
			Success: false,
			Error:   "User not found",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	response := APIResponse{
		Success: true,
		Data:    user,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		response := APIResponse{
			Success: false,
			Error:   "Invalid JSON data",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	// Validate required fields
	if newUser.Name == "" || newUser.Email == "" {
		response := APIResponse{
			Success: false,
			Error:   "Name and email are required",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	usersMu.Lock()
	newUser.ID = nextID
	nextID++
	users[newUser.ID] = newUser
	usersMu.Unlock()
	
	response := APIResponse{
		Success: true,
		Message: "User created successfully",
		Data:    newUser,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	id, err := strconv.Atoi(path)
	if err != nil {
		response := APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	var updatedUser User
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		response := APIResponse{
			Success: false,
			Error:   "Invalid JSON data",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	usersMu.Lock()
	defer usersMu.Unlock()
	
	if _, exists := users[id]; !exists {
		response := APIResponse{
			Success: false,
			Error:   "User not found",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	updatedUser.ID = id
	users[id] = updatedUser
	
	response := APIResponse{
		Success: true,
		Message: "User updated successfully",
		Data:    updatedUser,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/users/")
	id, err := strconv.Atoi(path)
	if err != nil {
		response := APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	usersMu.Lock()
	defer usersMu.Unlock()
	
	if _, exists := users[id]; !exists {
		response := APIResponse{
			Success: false,
			Error:   "User not found",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	delete(users, id)
	
	response := APIResponse{
		Success: true,
		Message: "User deleted successfully",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// API router function
func usersAPIHandler(w http.ResponseWriter, r *http.Request) {
	// Route based on method and path
	switch {
	case r.Method == "GET" && r.URL.Path == "/api/users":
		getUsersHandler(w, r)
	case r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/api/users/"):
		getUserHandler(w, r)
	case r.Method == "POST" && r.URL.Path == "/api/users":
		createUserHandler(w, r)
	case r.Method == "PUT" && strings.HasPrefix(r.URL.Path, "/api/users/"):
		updateUserHandler(w, r)
	case r.Method == "DELETE" && strings.HasPrefix(r.URL.Path, "/api/users/"):
		deleteUserHandler(w, r)
	default:
		response := APIResponse{
			Success: false,
			Error:   "Endpoint not found",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
	}
}

// File server handler
func staticFileHandler() {
	// Create a simple HTML file for demo
	htmlContent := `<!DOCTYPE html>
<html>
<head>
    <title>Go HTTP Server Demo</title>
</head>
<body>
    <h1>Welcome to Go HTTP Server</h1>
    <p>This is a static file served by Go HTTP server.</p>
    <h2>API Endpoints:</h2>
    <ul>
        <li>GET /api/users - Get all users</li>
        <li>GET /api/users/:id - Get user by ID</li>
        <li>POST /api/users - Create new user</li>
        <li>PUT /api/users/:id - Update user</li>
        <li>DELETE /api/users/:id - Delete user</li>
    </ul>
</body>
</html>`
	
	// Serve static content
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, htmlContent)
			return
		}
		
		// For other paths, return 404
		http.NotFound(w, r)
	})
}

// Health check handler
func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    time.Since(time.Now().Add(-time.Hour)).String(), // Demo uptime
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Metrics handler
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	usersMu.RLock()
	userCount := len(users)
	usersMu.RUnlock()
	
	metrics := map[string]interface{}{
		"users_count":     userCount,
		"memory_usage":    "N/A", // Would use runtime.MemStats in real app
		"goroutines":      "N/A", // Would use runtime.NumGoroutine()
		"requests_total":  "N/A", // Would track in middleware
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// Custom router implementation
type Router struct {
	routes map[string]map[string]http.HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]map[string]http.HandlerFunc),
	}
}

func (r *Router) AddRoute(method, path string, handler http.HandlerFunc) {
	if r.routes[method] == nil {
		r.routes[method] = make(map[string]http.HandlerFunc)
	}
	r.routes[method][path] = handler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if methodRoutes, exists := r.routes[req.Method]; exists {
		if handler, exists := methodRoutes[req.URL.Path]; exists {
			handler(w, req)
			return
		}
	}
	
	http.NotFound(w, req)
}

// Server configuration and startup
func setupServer() *http.Server {
	// Register handlers
	basicHandlers()
	staticFileHandler()
	
	// API routes with middleware
	apiHandler := corsMiddleware(loggingMiddleware(http.HandlerFunc(usersAPIHandler)))
	http.Handle("/api/users", apiHandler)
	http.Handle("/api/users/", apiHandler)
	
	// Protected routes
	protectedHandler := authMiddleware(http.HandlerFunc(metricsHandler))
	http.Handle("/metrics", protectedHandler)
	
	// Rate limited routes
	rateLimitedHandler := rateLimitMiddleware(10, time.Minute)(http.HandlerFunc(healthHandler))
	http.Handle("/health", rateLimitedHandler)
	
	// Create server with timeouts
	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	return server
}

// Graceful shutdown
func gracefulShutdown(server *http.Server) {
	// Channel to listen for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	// Block until signal is received
	<-c
	fmt.Println("\nShutting down server...")
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	
	fmt.Println("Server stopped gracefully")
}

func main() {
	fmt.Println("=== HTTP SERVER EXAMPLES ===")
	
	// For demonstration, we'll show the server setup without actually starting it
	server := setupServer()
	
	fmt.Printf("HTTP Server configured on %s\n", server.Addr)
	fmt.Println("\nAvailable endpoints:")
	fmt.Println("  GET  / - Static HTML page")
	fmt.Println("  GET  /hello - Simple greeting")
	fmt.Println("  GET  /json - JSON response")
	fmt.Println("  GET  /params?name=value - URL parameters")
	fmt.Println("  GET  /health - Health check (rate limited)")
	fmt.Println("  GET  /metrics - Metrics (requires auth)")
	fmt.Println("  GET  /api/users - Get all users")
	fmt.Println("  GET  /api/users/:id - Get user by ID")
	fmt.Println("  POST /api/users - Create user")
	fmt.Println("  PUT  /api/users/:id - Update user")
	fmt.Println("  DELETE /api/users/:id - Delete user")
	
	fmt.Println("\nMiddleware stack:")
	fmt.Println("  - CORS middleware (allows cross-origin requests)")
	fmt.Println("  - Logging middleware (logs requests)")
	fmt.Println("  - Auth middleware (requires Bearer token)")
	fmt.Println("  - Rate limiting middleware (10 requests per minute)")
	
	fmt.Println("\nServer features:")
	fmt.Println("  - JSON API responses")
	fmt.Println("  - Error handling")
	fmt.Println("  - Request timeouts")
	fmt.Println("  - Graceful shutdown")
	fmt.Println("  - Static file serving")
	fmt.Println("  - Custom routing")
	
	// In a real application, you would uncomment these lines:
	// go gracefulShutdown(server)
	// fmt.Printf("Starting server on %s...\n", server.Addr)
	// log.Fatal(server.ListenAndServe())
	
	fmt.Println("\n=== HTTP SERVER BEST PRACTICES ===")
	fmt.Println("1. Set appropriate timeouts")
	fmt.Println("2. Use middleware for cross-cutting concerns")
	fmt.Println("3. Implement graceful shutdown")
	fmt.Println("4. Use structured logging")
	fmt.Println("5. Handle errors consistently")
	fmt.Println("6. Validate input data")
	fmt.Println("7. Use HTTPS in production")
	fmt.Println("8. Implement rate limiting")
	fmt.Println("9. Add health check endpoints")
	fmt.Println("10. Monitor and expose metrics")
	
	fmt.Println("\nHTTP server examples completed!")
}