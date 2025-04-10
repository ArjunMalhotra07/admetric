package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

// Click represents a click event
type Click struct {
	ID           string    `json:"id"`
	AdID         string    `json:"ad_id"`
	IP           string    `json:"ip"`
	PlaybackTime int       `json:"playback_time"`
	Timestamp    time.Time `json:"timestamp"`
}

// ClickRequest represents the request body for recording a click
type ClickRequest struct {
	AdID         string `json:"ad_id"`
	IP           string `json:"ip"`
	PlaybackTime int    `json:"playback_time"`
}

func main() {
	// Configuration
	serverURL := "http://localhost:8888/ads/click"
	totalClicks := 450 // Target total clicks
	batchSize := 50    // Clicks per batch
	pauseDuration := 1 * time.Minute

	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Track progress
	clicksSent := 0
	batchCount := 0

	// Track clicks per ID
	clicksPerID := make(map[string]int)

	fmt.Printf("Starting click simulation. Target: %d clicks in batches of %d\n", totalClicks, batchSize)

	// Continue until we reach the target number of clicks
	for clicksSent < totalClicks {
		batchCount++
		fmt.Printf("\nSending batch #%d (%d clicks)...\n", batchCount, batchSize)

		// Send a batch of clicks
		for i := 0; i < batchSize && clicksSent < totalClicks; i++ {
			// Generate a random ad ID between 1 and 10
			adID := fmt.Sprintf("%d", rand.Intn(10)+1)

			// Create a click request
			clickReq := ClickRequest{
				AdID:         adID,
				IP:           "127.0.0.1",
				PlaybackTime: rand.Intn(300) + 1, // Random playback time between 1-300 seconds
			}

			// Send the request
			err := sendClick(serverURL, clickReq)
			if err != nil {
				fmt.Printf("Error sending click for ad %s: %v\n", adID, err)
			} else {
				clicksSent++
				clicksPerID[adID]++
				fmt.Printf("Click sent for ad %s (%d/%d)\n", adID, clicksSent, totalClicks)
			}

			// Small delay between requests to avoid overwhelming the server
			time.Sleep(100 * time.Millisecond)
		}

		// Pause between batches
		if clicksSent < totalClicks {
			fmt.Printf("Pausing for %v before next batch...\n", pauseDuration)
			time.Sleep(pauseDuration)
		}
	}

	// Print summary of clicks per ID
	fmt.Printf("\nClick simulation completed. Sent %d clicks.\n", clicksSent)
	fmt.Println("\nSummary of clicks per ID:")
	for id := 1; id <= 10; id++ {
		idStr := fmt.Sprintf("%d", id)
		fmt.Printf("Ad ID %s: %d clicks\n", idStr, clicksPerID[idStr])
	}
}

// sendClick sends a click request to the server
func sendClick(serverURL string, click ClickRequest) error {
	// Convert the request to JSON
	jsonData, err := json.Marshal(click)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", serverURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusAccepted {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
