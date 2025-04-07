package internal

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type result struct {
	statusCode int
	duration   time.Duration
}

func HttpStress(url, method string, totalRequests, concurrentRequests, intervalSec int, body string, headersMap map[string]string) {
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

			setUpHeaders(req, headersMap)

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

	report(totalRequests, totalRequestDuration, totalDuration, statusCodeCounter)
}

func ValidateInputHttpStress(url, method string, totalRequests, concurrentRequests, intervalSec int) error {
	if url == "" {
		return fmt.Errorf("URL is required")
	}

	validMethods := map[string]bool{
		"GET":    true,
		"HEAD":   true,
		"PATCH":  true,
		"POST":   true,
		"PUT":    true,
		"DELETE": true,
	}

	method = strings.ToUpper(method)
	if !validMethods[method] {
		return fmt.Errorf("invalid HTTP method: %s. Allowed methods are: GET, HEAD, PATCH, POST, PUT, DELETE", method)
	}

	if totalRequests <= 0 {
		return fmt.Errorf("total requests must be greater than 0")
	}

	if concurrentRequests <= 0 {
		return fmt.Errorf("concurrent requests must be greater than 0")
	}

	if intervalSec < 0 {
		return fmt.Errorf("interval seconds cannot be negative")
	}

	return nil
}

func setUpHeaders(req *http.Request, headersMap map[string]string) {
	for key, value := range headersMap {
		req.Header.Set(key, value)
	}

	req.Header.Set("User-Agent", "go-stress-test/1.0")
}

func report(totalRequests int, totalRequestDuration, totalDuration time.Duration, statusCodeCounter map[int]int) {
	avgDuration := time.Duration(0)
	if totalRequests > 0 {
		avgDuration = totalRequestDuration / time.Duration(totalRequests)
	}

	fmt.Println("\n--- Summary Report ---")
	fmt.Printf("Total time: %v\n", totalDuration)
	fmt.Printf("Total requests: %d\n", totalRequests)
	fmt.Printf("Average response time: %v\n", avgDuration)
	fmt.Printf("Status code 200: %d response(s)\n", statusCodeCounter[200])
	for code, count := range statusCodeCounter {
		if code == 200 {
			continue
		}

		fmt.Println("Request with other status codes:")
		if code == -1 {
			fmt.Printf("Failed requests: %d\n", count)
		} else {
			fmt.Printf("Status code %d: %d response(s)\n", code, count)
		}
	}
}
