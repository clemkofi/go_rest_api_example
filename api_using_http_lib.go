package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv" // Import godotenv
)

// Define a struct for the input on the /process endpoint
type ProcessInput struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// Define a struct for the response from the /process endpoint
type ProcessResponse struct {
	Message       string `json:"message"`
	ReceivedName  string `json:"received_name"`
	ReceivedValue int    `json:"received_value"`
	SecretFromEnv string `json:"secret_from_env"`
}

// healthHandler responds to /health requests
func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Create a simple JSON response
	response := map[string]string{"status": "ok", "message": "API is healthy!"}
	json.NewEncoder(w).Encode(response)
	log.Println("Health check successful")
}

// processHandler responds to /process requests
func processHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests for this endpoint
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the JSON input
	var input ProcessInput
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		http.Error(w, "Invalid JSON input: "+err.Error(), http.StatusBadRequest)
		log.Printf("Error decoding JSON: %v\n", err)
		return
	}
	defer r.Body.Close() // Good practice to close the request body

	log.Printf("Received input: Name=%s, Value=%d\n", input.Name, input.Value)

	// Get the secret message from environment variables
	secretMessage := os.Getenv("API_SECRET_MESSAGE")
	if secretMessage == "" {
		secretMessage = "Default secret (env var not set)" // Fallback
	}

	// Prepare the response
	response := ProcessResponse{
		Message:       fmt.Sprintf("Successfully processed input for %s.", input.Name),
		ReceivedName:  input.Name,
		ReceivedValue: input.Value,
		SecretFromEnv: secretMessage,
	}

	// Set the content type and write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	log.Printf("Processed input for %s, sent response.\n", input.Name)
}

func main() {
	// Load .env file.
	// If it's not found, godotenv.Load() will return an error, but we can often proceed
	// if environment variables are set externally (e.g., in Docker, Kubernetes).
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Could not load .env file. Using environment variables if set directly.")
	}

	// Get port from environment variable, with a default
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	// Get another config value (just to show you can load more)
	anotherConfig := os.Getenv("ANOTHER_CONFIG_VALUE")
	if anotherConfig != "" {
		log.Printf("Loaded another config: %s\n", anotherConfig)
	} else {
		log.Println("ANOTHER_CONFIG_VALUE not found in environment.")
	}


	// Define routes
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/process", processHandler)

	// Start the server
	log.Printf("Server starting on port %s...\n", port)
	log.Printf("Access health check at http://localhost:%s/health\n", port)
	log.Printf("Send POST requests to http://localhost:%s/process\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
