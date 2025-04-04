package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var (
	url         string
	method      string
	requests    int
	concurrency int
	intervalSec int
	body        string
)

type result struct {
	statusCode int
	duration   time.Duration
}

var rootCmd = &cobra.Command{
	Use:   "go-stress-test",
	Short: "A CLI for stress testing APIs",
	Long: `A simple CLI tool for stress testing APIs. 
It allows you to specify the target URL, HTTP method, number of requests, concurrency level, and request body.
This tool is useful for performance testing and benchmarking your APIs.`,
	Run: func(cmd *cobra.Command, args []string) {
		if url == "" {
			fmt.Println("Error: --url is required")
			return
		}

		validMethods := map[string]bool{
			"GET":    true,
			"POST":   true,
			"PUT":    true,
			"DELETE": true,
		}

		method = strings.ToUpper(method)
		if !validMethods[method] {
			fmt.Printf("Invalid HTTP method: %s\nAllowed methods are: GET, POST, PUT, DELETE\n", method)
			return
		}

		stressTest(url, method, requests, concurrency, intervalSec, body)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&url, "url", "u", "", "Target URL (required)")
	rootCmd.PersistentFlags().StringVarP(&method, "method", "X", "GET", "HTTP method: GET, POST, PUT, DELETE")
	rootCmd.PersistentFlags().IntVarP(&requests, "requests", "r", 10, "Number of requests")
	rootCmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "c", 1, "Number of concurrent requests")
	rootCmd.PersistentFlags().IntVarP(&intervalSec, "interval", "i", 0, "Interval between requests in seconds")
	rootCmd.PersistentFlags().StringVarP(&body, "body", "d", "", "Request body (for POST, PUT)")

	rootCmd.RegisterFlagCompletionFunc("method", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"GET", "POST", "PUT", "DELETE"}, cobra.ShellCompDirectiveDefault
	})
}

func stressTest(url, method string, totalRequests, concurrentRequests, intervalSec int, body string) {
	fmt.Println("Starting stress test...")
	fmt.Printf("URL: %s, Method: %s, Requests: %d, Concurrency: %d\n", url, method, totalRequests, concurrentRequests)
	if body != "" {
		fmt.Printf("Body: %s\n", body)
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrentRequests)
	resultsChannel := make(chan result, totalRequests)

	startTime := time.Now()
	statusCodeCounter := make(map[int]int)
	var totalRequestDuration time.Duration
	done := make(chan struct{})

	// goroutine to collect results
	go func() {
		for res := range resultsChannel {
			statusCodeCounter[res.statusCode]++
			totalRequestDuration += res.duration
		}
		close(done) // close the done channel when all results are collected
	}()

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)               // increment the wait group counter
		semaphore <- struct{}{} // acquire a semaphore slot

		// start a goroutine for each request
		go func(i int) {
			defer wg.Done() // decrement the wait group counter when the goroutine completes
			defer func(intervalSec int) {
				if intervalSec > 0 {
					time.Sleep(time.Duration(intervalSec) * time.Second)
				}
				<-semaphore // release the semaphore slot
			}(intervalSec)

			start := time.Now()

			client := &http.Client{}
			reqBody := strings.NewReader(body)
			req, err := http.NewRequest(strings.ToUpper(method), url, reqBody)
			if err != nil {
				fmt.Printf("request %d failed: %v\n", i+1, err)
				resultsChannel <- result{statusCode: -1, duration: time.Since(start)}
				return
			}

			if body != "" {
				req.Header.Set("Content-Type", "application/json")
			}

			resp, err := client.Do(req)
			elapsed := time.Since(start)

			if err != nil {
				fmt.Printf("request %d failed: %v\n", i+1, err)
				resultsChannel <- result{statusCode: -1, duration: elapsed}
				return
			}
			defer resp.Body.Close()

			fmt.Printf("Request %d: Status Code: %d (%v)\n", i+1, resp.StatusCode, elapsed)

			resultsChannel <- result{statusCode: resp.StatusCode, duration: elapsed}
		}(i)
	}

	wg.Wait()             // wait for all requests to finish
	close(resultsChannel) // close the results channel to signal completion
	<-done                // wait for the results collection to finish

	totalDuration := time.Since(startTime)
	avgDuration := time.Duration(0)
	if totalRequests > 0 {
		avgDuration = totalRequestDuration / time.Duration(totalRequests)
	}

	fmt.Println("\n--- Summary Report ---")
	for code, count := range statusCodeCounter {
		if code == -1 {
			fmt.Printf("Failed requests: %d\n", count)
		} else {
			fmt.Printf("Status %d: %d responses\n", code, count)
		}
	}
	fmt.Printf("Total time for %d requests: %v\n", totalRequests, totalDuration)
	fmt.Printf("Average response time: %v\n", avgDuration)
}
