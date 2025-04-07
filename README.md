# Go Stress Test

A simple CLI tool for stress testing APIs, built using the [Cobra CLI](https://github.com/spf13/cobra) framework. This tool allows you to test the performance and reliability of your APIs by sending multiple HTTP requests concurrently and analyzing the results.

## Features
- Specify the target URL, HTTP method, number of requests, concurrency level, and request body.
- Measure response times and collect status code statistics.
- Supports interval delays between requests for controlled testing.

## How It Works
The project is divided into two main components:
1. **Cobra CLI**: The command-line interface is built using Cobra, a powerful library for creating CLI applications in Go. It provides a structured way to define commands, flags, and argument parsing.
2. **HTTP Stress Testing Logic**: The core logic for stress testing is implemented in the `internal/http_stress.go` file.

### Cobra CLI
Cobra is a library for creating modern CLI applications in Go. It simplifies the process of defining commands, flags, and argument parsing. In this project, the `cmd/root.go` file defines the main command (`go-stress-test`) and its flags, such as:
- `--url` (`-u`): The target URL for the stress test. **Required flag**.
- `--method` (`-X`): The HTTP method (e.g., GET, HEAD, PATCH, POST, PUT). Default is GET.
- `--requests` (`-r`): The total number of requests to send. Default is 10.
- `--concurrency` (`-c`): The number of concurrent requests. Default is 2.
- `--interval` (`-i`): The interval (in seconds) between requests. Default is 0.
- `--body` (`-d`): The request body for requests.

### `http_stress.go`
The `internal/http_stress.go` file contains the core logic for performing the stress test. Here's a breakdown of its functionality:
- **Concurrency Management**: Uses Go's `sync.WaitGroup` and a semaphore channel to manage concurrent HTTP requests.
- **HTTP Requests**: Sends HTTP requests using the `net/http` package, with support for custom methods and request bodies.
- **Result Collection**: Collects response times and status codes in a thread-safe manner.
- **Reporting**: Generates a summary report with metrics like total time, average response time, and status code distribution.

## Usage
### Prerequisites
- Docker for build and run the project.

### Build and Run - Makefile
1. Build the project:
   ```bash
   make build
    ```

2. Run the CLI tool:
   ```bash
   make run URL=http://example.com REQUESTS=20 CONCURRENCY=5 HTTP_METHOD=GET
   ```

   For default values, you can run the command without specifying all parameters:
   ```bash
   make run URL=http://example.com
   ```

### Build and Run - Docker
1. Build the Docker image:
   ```bash
   docker build -t go-http-stress .
   ```
2. Run the Docker container:
   ```bash
    docker run --rm go-http-stress:latest -u http://example.com -r 20 -c 5 -X GET
    ```

    Post example:
    ```bash
    docker run --rm go-http-stress:latest -u https://67f4273dcbef97f40d2d8a5b.mockapi.io/users -X POST -d "{\"name\": \"John Doe\"}"
    ```

    For default values, you can run the command without specifying all flags:
    ```bash
    docker run --rm go-http-stress:latest -u http://example.com
    ```

    To see all available flags, run:
    ```bash
    docker run --rm go-http-stress:latest --help
    ```
