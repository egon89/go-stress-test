package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	url         string
	method      string
	requests    int
	concurrency int
	body        string
)

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

		stressTest(url, method, requests, concurrency, body)
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
	rootCmd.PersistentFlags().StringVarP(&body, "body", "d", "", "Request body (for POST, PUT)")

	rootCmd.RegisterFlagCompletionFunc("method", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"GET", "POST", "PUT", "DELETE"}, cobra.ShellCompDirectiveDefault
	})
}

func stressTest(url, method string, totalRequests, concurrentRequests int, body string) {
	fmt.Println("Starting stress test...")
	fmt.Printf("URL: %s, Method: %s, Requests: %d, Concurrency: %d\n", url, method, totalRequests, concurrentRequests)
	if body != "" {
		fmt.Printf("Body: %s\n", body)
	}
}
