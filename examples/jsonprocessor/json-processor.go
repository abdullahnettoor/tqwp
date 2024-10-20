package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/abdullahnettoor/tqwp"
)

// UserData represents the structure of our input JSON
type UserData struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

// ProcessedUserData represents the transformed data
type ProcessedUserData struct {
	ID          int    `json:"id"`
	FullName    string `json:"full_name"`
	EmailDomain string `json:"email_domain"`
	AccountAge  string `json:"account_age"`
	ProcessedAt string `json:"processed_at"`
}

// JSONProcessTask represents a task to process a JSON file
type JSONProcessTask struct {
	tqwp.TaskModel
	InputPath  string
	OutputPath string
}

// Process implements the Task interface
func (t *JSONProcessTask) Process() error {
	// Read input file
	data, err := os.ReadFile(t.InputPath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", t.InputPath, err)
	}

	// Parse input JSON
	var user UserData
	if err := json.Unmarshal(data, &user); err != nil {
		return fmt.Errorf("failed to parse JSON from %s: %v", t.InputPath, err)
	}

	// Process the data
	createdTime, err := time.Parse("2006-01-02", user.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to parse date: %v", err)
	}

	// Calculate account age
	accountAge := time.Since(createdTime).Round(24 * time.Hour)

	// Extract email domain
	emailParts := strings.Split(user.Email, "@")
	domain := emailParts[1]

	// Create processed data
	processed := ProcessedUserData{
		ID:          user.ID,
		FullName:    user.Name,
		EmailDomain: domain,
		AccountAge:  fmt.Sprintf("%d days", int(accountAge.Hours()/24)),
		ProcessedAt: time.Now().Format(time.RFC3339),
	}

	// Convert to JSON
	outputData, err := json.MarshalIndent(processed, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to create output JSON: %v", err)
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(t.OutputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Write output file
	if err := os.WriteFile(t.OutputPath, outputData, 0644); err != nil {
		return fmt.Errorf("failed to write output file %s: %v", t.OutputPath, err)
	}

	return nil
}

func main() {
	// Configure and create worker pool
	wp := tqwp.New(&tqwp.WorkerPoolConfig{
		NumOfWorkers: 4,
		MaxRetries:   2,
		QueueSize:    100,
	})
	defer wp.Summary()
	defer wp.Stop()

	// Start the worker pool
	wp.Start()

	// Create sample input directory and files for testing
	createSampleFiles()

	// Process all JSON files in the input directory
	inputDir := "./examples/jsonprocessor/input"
	outputDir := "./examples/jsonprocessor/output"

	files, err := os.ReadDir(inputDir)
	if err != nil {
		fmt.Printf("Failed to read input directory: %v\n", err)
		return
	}

	// Enqueue a task for each JSON file
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			task := &JSONProcessTask{
				InputPath:  filepath.Join(inputDir, file.Name()),
				OutputPath: filepath.Join(outputDir, "processed_"+file.Name()),
			}
			wp.EnqueueTask(task)
		}
	}
}

// Helper function to create sample input files for testing
func createSampleFiles() {
	inputDir := "./examples/jsonprocessor/input"
	os.MkdirAll(inputDir, 0755)

	// Sample user data
	users := []UserData{
		{
			ID:        1,
			Name:      "John Doe",
			Email:     "john@example.com",
			CreatedAt: "2022-01-15",
		},
		{
			ID:        2,
			Name:      "Jane Smith",
			Email:     "jane@company.com",
			CreatedAt: "2023-06-20",
		},
		{
			ID:        3,
			Name:      "Bob Wilson",
			Email:     "bob@domain.com",
			CreatedAt: "2021-11-30",
		},
	}

	// Create sample JSON files
	for i, user := range users {
		data, _ := json.MarshalIndent(user, "", "  ")
		filename := filepath.Join(inputDir, fmt.Sprintf("user_%d.json", i+1))
		os.WriteFile(filename, data, 0644)
	}
}
