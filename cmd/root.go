package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/egon89/go-stress-test/internal"
	"github.com/spf13/cobra"
)

var (
	url         string
	method      string
	requests    int
	concurrency int
	intervalSec int
	body        string
	headers     string
	headersMap  map[string]string
)

var rootCmd = &cobra.Command{
	Use:   "go-stress-test",
	Short: "A CLI for stress testing APIs",
	Long: `A simple CLI tool for stress testing APIs. 
It allows you to specify the target URL, HTTP method, number of requests, concurrency level, and request body.
This tool is useful for performance testing and benchmarking your APIs.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := internal.ValidateInputHttpStress(url, method, requests, concurrency, intervalSec)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		internal.HttpStress(url, method, requests, concurrency, intervalSec, body, headersMap)
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
	rootCmd.PersistentFlags().StringVarP(&method, "method", "X", "GET", "HTTP method: GET, HEAD, PATCH, POST, PUT, DELETE")
	rootCmd.PersistentFlags().IntVarP(&requests, "requests", "r", 10, "Number of requests")
	rootCmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "c", 2, "Number of concurrent requests")
	rootCmd.PersistentFlags().IntVarP(&intervalSec, "interval", "i", 0, "Interval between requests in seconds")
	rootCmd.PersistentFlags().StringVarP(&body, "body", "d", "", "Request body for requests")
	rootCmd.PersistentFlags().StringVarP(&headers, "headers", "H", "", "Request headers in key:value format, separated by commas")
	rootCmd.MarkPersistentFlagRequired("url")

	rootCmd.RegisterFlagCompletionFunc("method", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"GET", "POST", "PUT", "DELETE"}, cobra.ShellCompDirectiveDefault
	})

	cobra.OnInitialize(parseHeaders)
}

func parseHeaders() {
	headersMap = make(map[string]string)
	if headers != "" {
		pairs := strings.Split(headers, ",")
		for _, pair := range pairs {
			kv := strings.SplitN(pair, "=", 2)
			if len(kv) == 2 {
				headersMap[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
			}
		}
	}
}
