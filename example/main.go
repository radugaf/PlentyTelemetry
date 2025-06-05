package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	_ "github.com/radugaf/PlentyTelemetry/adapters" // trigger init()

	c "github.com/radugaf/PlentyTelemetry/config"
	d "github.com/radugaf/PlentyTelemetry/domain"
	p "github.com/radugaf/PlentyTelemetry/ports"
)

var logger p.LoggingService

type Post struct {
	ID     int    `json:"id,omitempty"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
}

func main() {
	fmt.Println("=== Starting PlentyTelemetry Demo ===")

	// Load configuration
	config, err := c.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	// Create writers based on config
	var writers []p.LogWriter
	fmt.Printf("Processing %d driver configs...\n", len(config.Drivers))

	for i, driverCfg := range config.Drivers {
		fmt.Printf("Processing driver %d: type=%s, enabled=%t\n", i, driverCfg.Type, driverCfg.Enabled)

		if !driverCfg.Enabled {
			fmt.Printf("Skipping disabled driver: %s\n", driverCfg.Type)
			continue
		}

		writer := c.CreateDriver(driverCfg.Type, driverCfg.Settings)
		if writer != nil {
			writers = append(writers, writer)
			fmt.Printf("Added writer for driver: %s\n", driverCfg.Type)
		} else {
			fmt.Printf("Failed to create driver: %s\n", driverCfg.Type)
		}
	}

	// Create logger
	logger = d.NewLogger(writers...)

	fmt.Printf("Logger initialized with %d writers\n", len(writers))

	// Test the logging
	logger.Info("Starting API tests", map[string]string{
		"service": "api-tester",
		"target":  "jsonplaceholder.typicode.com",
	})

	testGetPosts()

	fmt.Println("=== Demo completed ===")
}

func testGetPosts() {
	txID := logger.StartTransaction()

	logger.Info("Testing GET posts", map[string]string{
		"operation": "fetch_posts",
		"endpoint":  "/posts",
	}, txID)

	res, err := http.Get("https://jsonplaceholder.typicode.com/posts?_limit=3")
	if err != nil {
		logger.Error("Failed to fetch posts", map[string]string{
			"error":  err.Error(),
			"reason": "Some error",
		}, txID)
		return
	}
	defer res.Body.Close()

	// Read response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("Failed to read response", map[string]string{
			"error": err.Error(),
		}, txID)
		return
	}

	// Parse JSON
	var posts []Post
	err = json.Unmarshal(body, &posts)
	if err != nil {
		logger.Warning("Failed to parse JSON response", map[string]string{
			"error": err.Error(),
		}, txID)
		return
	}

	logger.Info("Posts fetched successfully", map[string]string{
		"status_code": fmt.Sprintf("%d", res.StatusCode),
		"count":       fmt.Sprintf("%d", len(posts)),
		"first_title": posts[0].Title,
	}, txID)

	logger.Debug("Response details", map[string]string{
		"content_type":  res.Header.Get("Content-Type"),
		"response_size": fmt.Sprintf("%d bytes", len(body)),
	}, txID)
}
