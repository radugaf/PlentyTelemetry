package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	a "github.com/radugaf/PlentyTelemetry/adapters"
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
	config, err := c.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	var writers []p.LogWriter
	for _, driverCfg := range config.Drivers {
		if !driverCfg.Enabled {
			continue
		}

		switch driverCfg.Type {
		case "cli":
			writers = append(writers, a.NewCLIDriver())
		case "json":
			filename := driverCfg.Settings["filename"]
			if filename == "" {
				filename = "logs.json"
			}
			writers = append(writers, a.NewJSONDriver(filename))
		}
	}

	logger = d.NewLogger(writers...)

	logger.Info("Starting API tests", map[string]string{
		"service": "api-tester",
		"target":  "jsonplaceholder.typicode.com",
	})

	testGetPosts()
}

func testGetPosts() {
	txID := logger.StartTransaction()

	logger.Info("Testing GET posts", map[string]string{
		"operation": "fetch_posts",
		"endpoint":  "/posts",
	}, txID)

	res, err := http.Get("https://jsonplaceholder.typicode.com/posts?_limit=1")
	if err != nil {
		logger.Error("Failed to fetch posts", map[string]string{
			"error": err.Error(),
		}, txID)
		return
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	var posts []Post
	json.Unmarshal(body, &posts)

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
