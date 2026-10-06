package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	targetURL := flag.String("url", "http://localhost:8080/orders", "Target URL")
	concurrency := flag.Int("c", 10, "Concurrent workers")
	totalRequests := flag.Int("n", 100, "Total requests to send")
	itemID := flag.String("item", "00000000-0000-0000-0000-000000000001", "Item ID to order")
	flag.Parse()

	// Dedicated high-throughput HTTP transport
	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        500,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	var (
		successCount   atomic.Int64
		rateLimitCount atomic.Int64
		failureCount   atomic.Int64
	)

	jobs := make(chan int, *totalRequests)
	for i := 0; i < *totalRequests; i++ {
		jobs <- i
	}
	close(jobs)

	start := time.Now()
	var wg sync.WaitGroup

	for w := 0; w < *concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobs {
				payload := []byte(fmt.Sprintf(`{"item_id": "%s", "amount": 1}`, *itemID))
				resp, err := client.Post(*targetURL, "application/json", bytes.NewReader(payload))
				if err != nil {
					failureCount.Add(1)
					continue
				}
				resp.Body.Close()

				switch resp.StatusCode {
				case http.StatusCreated:
					successCount.Add(1)
				case http.StatusTooManyRequests:
					rateLimitCount.Add(1)
				default:
					failureCount.Add(1)
				}
			}
		}()
	}

	wg.Wait()
	duration := time.Since(start)

	fmt.Printf("\n--- Load Test Summary ---\n")
	fmt.Printf("Total Time:     %v\n", duration)
	fmt.Printf("Success (201):  %d\n", successCount.Load())
	fmt.Printf("Limited (429):  %d\n", rateLimitCount.Load())
	fmt.Printf("Failed:         %d\n", failureCount.Load())
	fmt.Printf("Throughput:     %.2f req/s\n", float64(*totalRequests)/duration.Seconds())
}
