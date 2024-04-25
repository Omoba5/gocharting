package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// MonitoringData represents the structure of monitoring data
type MonitoringData struct {
	Timestamp time.Time `json:"timestamp"`
	CPU       float64   `json:"cpu"`
	Memory    float64   `json:"memory"`
}

// Function to generate mock monitoring data
func generateMonitoringData() MonitoringData {
	return MonitoringData{
		Timestamp: time.Now(),
		CPU:       20 + rand.Float64()*80, // Random CPU usage between 20% and 100%
		Memory:    40 + rand.Float64()*60, // Random memory usage between 40% and 100%
	}
}

// SSEHandler handles SSE connections and streams monitoring data
func SSEHandler(w http.ResponseWriter, r *http.Request) {
	// Set response headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Create a new ticker to send monitoring data at regular intervals
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Send monitoring data periodically
	for {
		select {
		case <-ticker.C:
			// Generate mock monitoring data
			data := generateMonitoringData()

			// Convert monitoring data to JSON
			jsonData, err := json.Marshal(data)
			if err != nil {
				log.Println("Error encoding JSON:", err)
				continue
			}

			// Write JSON data as SSE event
			fmt.Fprintf(w, "data: %s\n\n", jsonData)
			w.(http.Flusher).Flush() // Flush the response buffer to send data immediately
		case <-r.Context().Done():
			// Client closed connection
			log.Println("Client disconnected")
			return
		}
	}
}

func main() {
	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Handle SSE endpoint
	http.HandleFunc("/events", SSEHandler)

	// Serve HTML page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// Start HTTP server
	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
