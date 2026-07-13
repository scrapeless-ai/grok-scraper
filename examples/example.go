// Grok Scraper — Scrapeless LLM Chat Scraper (Go example)
//
// Docs:  https://docs.scrapeless.com/en/llm-chat-scraper/quickstart/introduction/
// Token: https://app.scrapeless.com/passport/login?redirect=/quick-start
//
// Run:
//
//	export SCRAPELESS_API_TOKEN="your_api_token"
//	go run example.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const apiURL = "https://api.scrapeless.com/api/v2/scraper/execute"

func main() {
	apiToken := os.Getenv("SCRAPELESS_API_TOKEN")
	if apiToken == "" {
		apiToken = "YOUR_API_TOKEN"
	}

	payload := map[string]any{
		"actor": "scraper.grok",
		"input": map[string]any{
			"prompt":  "Most reliable proxy service for data extraction",
			"country": "US",
			"mode":    "MODEL_MODE_AUTO",
		},
		// Optional: receive the result via webhook instead of the sync response.
		// "webhook": map[string]any{"url": "https://www.your-webhook.com"},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-token", apiToken)

	client := &http.Client{Timeout: 180 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode >= 300 {
		panic(fmt.Sprintf("request failed: %d %s", resp.StatusCode, raw))
	}

	var data struct {
		Status     string `json:"status"`
		TaskID     string `json:"task_id"`
		TaskResult struct {
			UserModel        string `json:"user_model"`
			FullResponse     string `json:"full_response"`
			WebSearchResults []struct {
				Title string `json:"title"`
				URL   string `json:"url"`
			} `json:"web_search_results"`
		} `json:"task_result"`
	}
	if err := json.Unmarshal(raw, &data); err != nil {
		panic(err)
	}

	fmt.Println("Status: ", data.Status)
	fmt.Println("Task ID:", data.TaskID)
	fmt.Println("Model:  ", data.TaskResult.UserModel)
	fmt.Println("\nAnswer:\n", data.TaskResult.FullResponse)

	for _, link := range data.TaskResult.WebSearchResults {
		fmt.Printf("- %s -> %s\n", link.Title, link.URL)
	}

	var pretty bytes.Buffer
	_ = json.Indent(&pretty, raw, "", "  ")
	fmt.Println("\nRaw response:\n", pretty.String())
}
