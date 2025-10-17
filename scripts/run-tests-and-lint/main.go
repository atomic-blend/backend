package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ServiceResult struct {
	Name         string
	TestResult   *TestResult
	LintResult   *LintResult
	GolintResult *GolintResult
	HasTests     bool
	HasGRPC      bool
	Index        int // For maintaining order
}

type ServiceStatus struct {
	Name         string
	Status       string // "queued", "running", "completed", "failed"
	WorkerID     int
	GolintStatus string // "pending", "running", "passed", "failed"
	TestStatus   string // "pending", "running", "passed", "failed"
	GRPCStatus   string // "pending", "running", "passed", "failed"
	Progress     string // Current operation
	Duration     time.Duration
	Index        int
}

type TestResult struct {
	Passed       bool
	Output       string
	Error        string
	Duration     time.Duration
	TestCount    int
	PassCount    int
	FailCount    int
	SkipCount    int
	SkippedTests []SkippedTest
	FailedTests  []FailedTest
}

type SkippedTest struct {
	Test   string
	File   string
	Line   int
	Reason string
}

type FailedTest struct {
	Test   string
	File   string
	Line   int
	Error  string
	Output string
}

type LintResult struct {
	Passed   bool
	Output   string
	Error    string
	Duration time.Duration
}

type Summary struct {
	TotalServices     int
	PassedTests       int
	FailedTests       int
	PassedLint        int
	FailedLint        int
	PassedGolint      int
	FailedGolint      int
	TotalDuration     time.Duration
	TotalTestCount    int
	TotalPassCount    int
	TotalFailCount    int
	TotalSkipCount    int
	TotalGolintIssues int
}

type GolintResult struct {
	Passed        bool
	Output        string
	Error         string
	Duration      time.Duration
	IssueCount    int
	LintingIssues []LintingIssue
}

type LintingIssue struct {
	File    string
	Line    int
	Column  int
	Message string
	Rule    string
}

// Global status tracking for real-time table updates
var (
	statusMap = make(map[string]*ServiceStatus)
	statusMux sync.RWMutex
)

func main() {
	// Parse command line flags
	var numThreads = flag.Int("threads", 4, "Number of parallel threads to use for running tests and linting")
	var serviceName = flag.String("service", "", "Run tests for a specific service only (e.g., auth, mail, productivity)")
	flag.Parse()

	// Validate service flag
	if *serviceName != "" {
		validServices := []string{"auth", "productivity", "grpc", "mail", "mail-server", "shared"}
		isValid := false
		for _, validService := range validServices {
			if *serviceName == validService {
				isValid = true
				break
			}
		}
		if !isValid {
			fmt.Printf("❌ Error: Invalid service '%s'. Valid services are: %s\n",
				*serviceName, strings.Join(validServices, ", "))
			os.Exit(1)
		}
	}

	if *serviceName != "" {
		fmt.Printf("🚀 Running Tests for Service: %s\n", *serviceName)
		fmt.Println("📍 This script should be run from the backend directory")
		fmt.Println("📋 Full logs will be displayed")
		fmt.Println(strings.Repeat("=", 50))
	} else {
		fmt.Println("🚀 Running Microservice Tests, Golint, and gRPC Linting")
		fmt.Println("📍 This script should be run from the backend directory")
		fmt.Printf("🧵 Using %d parallel threads\n", *numThreads)
		fmt.Println(strings.Repeat("=", 50))
	}

	// Get the workspace root (assuming script is run from backend directory)
	workspaceRoot, err := os.Getwd()
	if err != nil {
		fmt.Printf("❌ Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	startTime := time.Now()

	// Handle single service mode
	if *serviceName != "" {
		servicePath := filepath.Join(workspaceRoot, *serviceName)
		if _, err := os.Stat(servicePath); err != nil {
			fmt.Printf("❌ Error: Service directory '%s' not found at %s\n", *serviceName, servicePath)
			os.Exit(1)
		}

		// Run single service with full logs
		runSingleService(workspaceRoot, *serviceName, startTime)
		return
	}

	// Define microservices to check (from cog.toml)
	services := []string{
		"auth",
		"productivity",
		"grpc",
		"mail",
		"mail-server",
		"shared",
	}

	// Count total services that will be processed
	totalServices := 0
	var validServices []string
	for _, service := range services {
		servicePath := filepath.Join(workspaceRoot, service)
		if _, err := os.Stat(servicePath); err == nil {
			totalServices++
			validServices = append(validServices, service)
		}
	}

	fmt.Printf("📊 Found %d services to process\n", totalServices)
	fmt.Println()

	if totalServices == 0 {
		fmt.Printf("\n❌ Error: No services found. Please run this script from the backend directory.\n")
		fmt.Printf("Current directory: %s\n", workspaceRoot)
		fmt.Printf("Expected to find services like: auth, mail, productivity, etc.\n")
		os.Exit(1)
	}

	// Initialize status map
	for i, service := range validServices {
		statusMux.Lock()
		statusMap[service] = &ServiceStatus{
			Name:         service,
			Status:       "queued",
			GolintStatus: "pending",
			TestStatus:   "pending",
			GRPCStatus:   "pending",
			Progress:     "Waiting...",
			Index:        i,
		}
		statusMux.Unlock()
	}

	// Start table updater
	stopUpdater := make(chan bool)
	go updateTable(totalServices, stopUpdater)

	// Create channels for parallel processing
	serviceChan := make(chan string, len(validServices))
	resultChan := make(chan *ServiceResult, len(validServices))
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < *numThreads; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for service := range serviceChan {
				processService(workspaceRoot, service, resultChan, workerID)
			}
		}(i)
	}

	// Send services to workers
	for _, service := range validServices {
		serviceChan <- service
	}
	close(serviceChan)

	// Start a goroutine to close result channel when all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	results := make([]*ServiceResult, 0, totalServices)
	for result := range resultChan {
		results = append(results, result)
	}

	// Stop table updater
	stopUpdater <- true
	fmt.Println() // Clear the table

	// Sort results by index to maintain original order
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].Index > results[j].Index {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	// Display results
	displayResults(results, startTime)
}

// runSingleService runs tests and linting for a single service with full log output
func runSingleService(workspaceRoot, serviceName string, startTime time.Time) {
	servicePath := filepath.Join(workspaceRoot, serviceName)

	fmt.Printf("🔍 Processing service: %s\n", serviceName)
	fmt.Printf("📁 Service path: %s\n", servicePath)
	fmt.Println(strings.Repeat("-", 50))

	allPassed := true

	// Check if service has tests (and can run golint)
	if hasTests(servicePath) {
		fmt.Println("🧪 Running golint...")
		fmt.Println(strings.Repeat("-", 30))

		golintResult := runGolintWithFullOutput(servicePath, serviceName)
		if !golintResult.Passed {
			allPassed = false
			fmt.Printf("❌ Golint failed for %s\n", serviceName)
			fmt.Printf("Issues found: %d\n", golintResult.IssueCount)
			if golintResult.Output != "" {
				fmt.Println("Golint output:")
				fmt.Println(golintResult.Output)
			}
			if golintResult.Error != "" {
				fmt.Println("Golint errors:")
				fmt.Println(golintResult.Error)
			}
		} else {
			fmt.Printf("✅ Golint passed for %s\n", serviceName)
			if golintResult.Output != "" {
				fmt.Println("Golint output:")
				fmt.Println(golintResult.Output)
			}
		}

		fmt.Println(strings.Repeat("-", 30))
		fmt.Println("🧪 Running tests...")
		fmt.Println(strings.Repeat("-", 30))

		testResult := runTestsWithFullOutput(servicePath, serviceName)
		if !testResult.Passed {
			allPassed = false
			fmt.Printf("❌ Tests failed for %s\n", serviceName)
		} else {
			fmt.Printf("✅ Tests passed for %s\n", serviceName)
		}
	} else {
		fmt.Printf("ℹ️  No tests found for %s (no go.mod file)\n", serviceName)
	}

	// Check if service has gRPC
	if hasGRPC(servicePath) {
		fmt.Println(strings.Repeat("-", 30))
		fmt.Println("🔍 Running gRPC lint...")
		fmt.Println(strings.Repeat("-", 30))

		lintResult := runGRPCLintWithFullOutput(servicePath, serviceName)
		if !lintResult.Passed {
			allPassed = false
			fmt.Printf("❌ gRPC lint failed for %s\n", serviceName)
			if lintResult.Output != "" {
				fmt.Println("gRPC lint output:")
				fmt.Println(lintResult.Output)
			}
			if lintResult.Error != "" {
				fmt.Println("gRPC lint errors:")
				fmt.Println(lintResult.Error)
			}
		} else {
			fmt.Printf("✅ gRPC lint passed for %s\n", serviceName)
			if lintResult.Output != "" {
				fmt.Println("gRPC lint output:")
				fmt.Println(lintResult.Output)
			}
		}
	}

	fmt.Println(strings.Repeat("=", 50))
	duration := time.Since(startTime)
	fmt.Printf("⏱️  Total duration: %.2fs\n", duration.Seconds())

	if allPassed {
		fmt.Printf("✅ All checks passed for %s!\n", serviceName)
	} else {
		fmt.Printf("❌ Some checks failed for %s\n", serviceName)
		os.Exit(1)
	}
}

// processService processes a single service in a worker goroutine
func processService(workspaceRoot, service string, resultChan chan<- *ServiceResult, workerID int) {
	servicePath := filepath.Join(workspaceRoot, service)

	// Find the index of this service in the original list for ordering
	services := []string{"auth", "productivity", "grpc", "mail", "mail-server", "shared"}
	index := -1
	for i, s := range services {
		if s == service {
			index = i
			break
		}
	}

	result := &ServiceResult{
		Name:  service,
		Index: index,
	}

	// Update status to running
	updateServiceStatus(service, "running", workerID, "Starting...")

	// Check if service has tests (and can run golint)
	if hasTests(servicePath) {
		result.HasTests = true

		// Run golint first
		updateServiceStatus(service, "running", workerID, "Running golint...")
		result.GolintResult = runGolint(servicePath, service)
		updateGolintStatus(service, result.GolintResult.Passed)

		// Then run tests
		updateServiceStatus(service, "running", workerID, "Running tests...")
		result.TestResult = runTests(servicePath, service)
		updateTestStatus(service, result.TestResult.Passed)
	}

	// Check if service has gRPC
	if hasGRPC(servicePath) {
		result.HasGRPC = true
		updateServiceStatus(service, "running", workerID, "Running gRPC lint...")
		result.LintResult = runGRPCLint(servicePath, service)
		updateGRPCStatus(service, result.LintResult.Passed)
	}

	// Mark as completed
	updateServiceStatus(service, "completed", workerID, "Completed")
	resultChan <- result
}

func hasTests(servicePath string) bool {
	// Check for go.mod file - if it exists, the service can run tests
	_, err := os.Stat(filepath.Join(servicePath, "go.mod"))
	return err == nil
}

// updateServiceStatus updates the status of a service
func updateServiceStatus(service, status string, workerID int, progress string) {
	statusMux.Lock()
	defer statusMux.Unlock()
	if s, exists := statusMap[service]; exists {
		s.Status = status
		s.WorkerID = workerID
		s.Progress = progress
	}
}

// updateGolintStatus updates the golint status
func updateGolintStatus(service string, passed bool) {
	statusMux.Lock()
	defer statusMux.Unlock()
	if s, exists := statusMap[service]; exists {
		if passed {
			s.GolintStatus = "passed"
		} else {
			s.GolintStatus = "failed"
		}
	}
}

// updateTestStatus updates the test status
func updateTestStatus(service string, passed bool) {
	statusMux.Lock()
	defer statusMux.Unlock()
	if s, exists := statusMap[service]; exists {
		if passed {
			s.TestStatus = "passed"
		} else {
			s.TestStatus = "failed"
		}
	}
}

// updateGRPCStatus updates the gRPC lint status
func updateGRPCStatus(service string, passed bool) {
	statusMux.Lock()
	defer statusMux.Unlock()
	if s, exists := statusMap[service]; exists {
		if passed {
			s.GRPCStatus = "passed"
		} else {
			s.GRPCStatus = "failed"
		}
	}
}

func hasGRPC(servicePath string) bool {
	// Only run gRPC lint for the grpc directory
	serviceName := filepath.Base(servicePath)
	return serviceName == "grpc"
}

// updateTable continuously updates and displays the status table
func updateTable(totalServices int, stopChan <-chan bool) {
	ticker := time.NewTicker(200 * time.Millisecond) // Faster refresh for animation
	defer ticker.Stop()
	frame := 0

	for {
		select {
		case <-stopChan:
			return
		case <-ticker.C:
			displayStatusTable(totalServices, frame)
			frame++
		}
	}
}

// displayStatusTable displays the current status of all services
func displayStatusTable(totalServices int, frame int) {
	statusMux.RLock()
	defer statusMux.RUnlock()

	// Clear screen and move cursor to top
	fmt.Print("\033[2J\033[H")

	fmt.Println("🚀 Running Microservice Tests, Golint, and gRPC Linting")
	fmt.Println("📍 Real-time Status Table")
	fmt.Println(strings.Repeat("=", 80))

	// Table header
	fmt.Printf("%-15s %-8s %-8s %-8s %-8s %-20s\n",
		"SERVICE", "STATUS", "GOLINT", "TESTS", "GRPC", "PROGRESS")
	fmt.Println(strings.Repeat("-", 80))

	// Sort services by index
	var services []*ServiceStatus
	for _, status := range statusMap {
		services = append(services, status)
	}

	// Simple bubble sort by index
	for i := 0; i < len(services); i++ {
		for j := i + 1; j < len(services); j++ {
			if services[i].Index > services[j].Index {
				services[i], services[j] = services[j], services[i]
			}
		}
	}

	// Display each service
	for _, service := range services {
		statusIcon := getStatusIcon(service.Status)
		golintIcon := getAnimatedIcon(service.GolintStatus, frame)
		testIcon := getAnimatedIcon(service.TestStatus, frame)
		grpcIcon := getAnimatedIcon(service.GRPCStatus, frame)

		fmt.Printf("%-15s %-8s %-8s %-8s %-8s %-20s\n",
			service.Name,
			statusIcon,
			golintIcon,
			testIcon,
			grpcIcon,
			service.Progress)
	}

	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("Total Services: %d | Use Ctrl+C to stop\n", totalServices)
}

// getStatusIcon returns an icon for the service status
func getStatusIcon(status string) string {
	switch status {
	case "queued":
		return "⏳"
	case "running":
		return "🔄"
	case "completed":
		return "✅"
	case "failed":
		return "❌"
	default:
		return "❓"
	}
}

// getCheckIcon returns an icon for check status
func getCheckIcon(status string) string {
	switch status {
	case "pending":
		return "⏳"
	case "running":
		return "⏳" // Will be animated with spinner
	case "passed":
		return "✅"
	case "failed":
		return "❌"
	default:
		return "❓"
	}
}

// getAnimatedIcon returns an animated icon for running status
func getAnimatedIcon(status string, frame int) string {
	if status == "running" {
		// Different hourglass states for animation
		spinners := []string{"⏳", "⏳", "⏳", "⏳", "⏳", "⏳", "⏳", "⏳"}
		return spinners[frame%len(spinners)]
	}
	return getCheckIcon(status)
}

// printProgressBar prints a simple progress bar
func printProgressBar(current, total int, prefix string) {
	const barLength = 30
	filled := int(float64(current) / float64(total) * barLength)
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barLength-filled)
	percentage := int(float64(current) / float64(total) * 100)
	fmt.Printf("\r%s [%s] %d%% (%d/%d)", prefix, bar, percentage, current, total)
}

// printSpinner prints a spinning animation
func printSpinner(done chan bool, prefix string) {
	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	for {
		select {
		case <-done:
			fmt.Printf("\r%s ✅ Complete\n", prefix)
			return
		default:
			fmt.Printf("\r%s %s Running...", prefix, spinner[i])
			time.Sleep(100 * time.Millisecond)
			i = (i + 1) % len(spinner)
		}
	}
}

// parseTestOutput parses Go test JSON output
func parseTestOutput(output string) (int, int, int, int, []SkippedTest, []FailedTest) {
	var testCount, passCount, failCount, skipCount int
	var skippedTests []SkippedTest
	var failedTests []FailedTest

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var testEvent struct {
			Action string `json:"Action"`
			Test   string `json:"Test"`
			Output string `json:"Output"`
		}

		if err := json.Unmarshal([]byte(line), &testEvent); err != nil {
			continue
		}

		switch testEvent.Action {
		case "run":
			testCount++
		case "pass":
			passCount++
		case "fail":
			failCount++
			// Extract failed test details
			failedTest := parseFailedTest(testEvent.Test, testEvent.Output)
			failedTests = append(failedTests, failedTest)
		case "skip":
			skipCount++
			// Extract skipped test details
			skippedTest := parseSkippedTest(testEvent.Test, testEvent.Output)
			skippedTests = append(skippedTests, skippedTest)
		}
	}

	return testCount, passCount, failCount, skipCount, skippedTests, failedTests
}

// parseSkippedTest extracts details from a skipped test
func parseSkippedTest(testName, output string) SkippedTest {
	skippedTest := SkippedTest{
		Test:   testName,
		Reason: "No reason provided",
		File:   "Unknown",
		Line:   0,
	}

	// Try to extract file and line information from test name
	// Test names are typically in format: package.TestName or package/file.TestName
	// For Go tests, we need to look for patterns like: package/file.TestName
	parts := strings.Split(testName, "/")
	if len(parts) > 1 {
		// Check if the last part contains a dot (indicating TestName)
		lastPart := parts[len(parts)-1]
		if strings.Contains(lastPart, ".") {
			// This is package/file.TestName format
			filePart := parts[len(parts)-2]
			skippedTest.File = filePart + ".go"
		} else {
			// This might be package/file format
			filePart := parts[len(parts)-1]
			skippedTest.File = filePart + ".go"
		}
	} else if strings.Contains(testName, ".") {
		// This might be package.TestName format, try to extract package name
		dotIndex := strings.LastIndex(testName, ".")
		if dotIndex > 0 {
			packageName := testName[:dotIndex]
			// Try to convert package name to file name
			skippedTest.File = packageName + ".go"
		}
	}

	// Try to extract file path, line, and column from output
	if output != "" {
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)

			// Look for file:line:column patterns in the output
			if strings.Contains(line, ":") && strings.Count(line, ":") >= 2 {
				// Try to parse as file:line:column: message
				parts := strings.SplitN(line, ":", 4)
				if len(parts) >= 3 {
					// Check if it looks like a file path with line/column
					if lineNum, err := strconv.Atoi(parts[1]); err == nil {
						skippedTest.File = parts[0]
						skippedTest.Line = lineNum
						if len(parts) >= 3 {
							if _, err := strconv.Atoi(parts[2]); err == nil {
								// This is file:line:column format
								skippedTest.Line = lineNum
								if len(parts) >= 4 {
									skippedTest.Reason = strings.TrimSpace(parts[3])
								}
								break
							}
						}
					}
				}
			}

			// Look for skip reason patterns
			if strings.Contains(strings.ToLower(line), "skip") {
				// Extract the reason after "skip" or similar keywords
				if idx := strings.Index(strings.ToLower(line), "skip"); idx != -1 {
					reason := strings.TrimSpace(line[idx+4:])
					if reason != "" {
						skippedTest.Reason = reason
					}
				}
			}
		}
	}

	return skippedTest
}

// parseFailedTest extracts details from a failed test
func parseFailedTest(testName, output string) FailedTest {
	failedTest := FailedTest{
		Test:  testName,
		Error: "No error details provided",
		File:  "Unknown",
		Line:  0,
	}

	// Try to extract file and line information from test name
	// Test names are typically in format: package.TestName or package/file.TestName
	// For Go tests, we need to look for patterns like: package/file.TestName
	parts := strings.Split(testName, "/")
	if len(parts) > 1 {
		// Check if the last part contains a dot (indicating TestName)
		lastPart := parts[len(parts)-1]
		if strings.Contains(lastPart, ".") {
			// This is package/file.TestName format
			filePart := parts[len(parts)-2]
			failedTest.File = filePart + ".go"
		} else {
			// This might be package/file format
			filePart := parts[len(parts)-1]
			failedTest.File = filePart + ".go"
		}
	} else if strings.Contains(testName, ".") {
		// This might be package.TestName format, try to extract package name
		dotIndex := strings.LastIndex(testName, ".")
		if dotIndex > 0 {
			packageName := testName[:dotIndex]
			// Try to convert package name to file name
			failedTest.File = packageName + ".go"
		}
	}

	// Extract error details and file location from output
	if output != "" {
		lines := strings.Split(output, "\n")
		var errorLines []string

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// Look for file:line:column patterns in the output
			if strings.Contains(line, ":") && strings.Count(line, ":") >= 2 {
				// Try to parse as file:line:column: message
				parts := strings.SplitN(line, ":", 4)
				if len(parts) >= 3 {
					// Check if it looks like a file path with line/column
					if lineNum, err := strconv.Atoi(parts[1]); err == nil {
						failedTest.File = parts[0]
						failedTest.Line = lineNum
						if len(parts) >= 4 {
							// This is file:line:column: message format
							failedTest.Error = strings.TrimSpace(parts[3])
							continue
						}
					}
				}
			}

			// Collect other error lines
			errorLines = append(errorLines, line)
		}

		// If we didn't find a specific file:line:column format, use the collected error lines
		if failedTest.Error == "No error details provided" && len(errorLines) > 0 {
			failedTest.Error = strings.Join(errorLines, " ")
		}
	}

	return failedTest
}

// parseLintingIssues extracts linting issues from golint output
func parseLintingIssues(output string) []LintingIssue {
	var issues []LintingIssue

	if output == "" {
		return issues
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Parse golint output format: filename:line:column: message
		parts := strings.SplitN(line, ":", 4)
		if len(parts) >= 4 {
			file := parts[0]
			lineStr := parts[1]
			columnStr := parts[2]
			message := strings.TrimSpace(parts[3])

			// Convert line and column to integers
			lineNum := 0
			columnNum := 0
			if l, err := strconv.Atoi(lineStr); err == nil {
				lineNum = l
			}
			if c, err := strconv.Atoi(columnStr); err == nil {
				columnNum = c
			}

			// Extract rule from message (if present)
			rule := "golint"
			if strings.Contains(message, ":") {
				ruleParts := strings.SplitN(message, ":", 2)
				if len(ruleParts) == 2 {
					rule = strings.TrimSpace(ruleParts[0])
					message = strings.TrimSpace(ruleParts[1])
				}
			}

			issues = append(issues, LintingIssue{
				File:    file,
				Line:    lineNum,
				Column:  columnNum,
				Message: message,
				Rule:    rule,
			})
		}
	}

	return issues
}

func runTests(servicePath, serviceName string) *TestResult {
	startTime := time.Now()

	// Create command with context and set the working directory
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Use JSON output format for better parsing
	cmd := exec.CommandContext(ctx, "go", "test", "./...", "-json")
	cmd.Dir = servicePath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Start the command
	if err := cmd.Start(); err != nil {
		duration := time.Since(startTime)
		return &TestResult{
			Passed:   false,
			Output:   "",
			Error:    err.Error(),
			Duration: duration,
		}
	}

	// Create a channel to signal when tests are done
	done := make(chan bool)

	// Start spinner in a goroutine
	go printSpinner(done, "  ⏳")

	// Wait for completion
	err := cmd.Wait()

	// Signal spinner to stop
	done <- true

	duration := time.Since(startTime)

	output := stdout.String()
	errorOutput := stderr.String()

	// Parse the JSON output to get test counts
	testCount, passCount, failCount, skipCount, skippedTests, failedTests := parseTestOutput(output)

	result := &TestResult{
		Output:       output,
		Error:        errorOutput,
		Duration:     duration,
		TestCount:    testCount,
		PassCount:    passCount,
		FailCount:    failCount,
		SkipCount:    skipCount,
		SkippedTests: skippedTests,
		FailedTests:  failedTests,
	}

	if err != nil {
		result.Passed = false
	} else {
		result.Passed = true
	}

	return result
}

func runGRPCLint(servicePath, serviceName string) *LintResult {
	startTime := time.Now()

	// Try to find proto files
	protoFiles, err := filepath.Glob(filepath.Join(servicePath, "**/*.proto"))
	if err != nil || len(protoFiles) == 0 {
		return &LintResult{
			Passed:   true,
			Output:   "No proto files found",
			Duration: time.Since(startTime),
		}
	}

	// Run buf lint if available
	cmd := exec.Command("buf", "lint")
	cmd.Dir = servicePath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd = exec.CommandContext(ctx, "buf", "lint")

	err = cmd.Run()
	duration := time.Since(startTime)

	result := &LintResult{
		Output:   stdout.String(),
		Error:    stderr.String(),
		Duration: duration,
	}

	if err != nil {
		result.Passed = false
	} else {
		result.Passed = true
	}

	return result
}

func runGolint(servicePath, serviceName string) *GolintResult {
	startTime := time.Now()

	// Check if golint is available
	if _, err := exec.LookPath("golint"); err != nil {
		return &GolintResult{
			Passed:     true,
			Output:     "golint not available",
			Error:      "golint command not found",
			Duration:   time.Since(startTime),
			IssueCount: 0,
		}
	}

	// Try to find go files
	goFiles, err := filepath.Glob(filepath.Join(servicePath, "**/*.go"))
	if err != nil || len(goFiles) == 0 {
		return &GolintResult{
			Passed:     true,
			Output:     "No go files found",
			Duration:   time.Since(startTime),
			IssueCount: 0,
		}
	}

	// Run golint if available
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "golint", "-set_exit_status", "./...")
	cmd.Dir = servicePath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	duration := time.Since(startTime)

	output := stdout.String()

	// Parse linting issues
	lintingIssues := parseLintingIssues(output)

	result := &GolintResult{
		Output:        output,
		Error:         stderr.String(),
		Duration:      duration,
		IssueCount:    len(lintingIssues),
		LintingIssues: lintingIssues,
	}

	// Determine if golint passed based on exit code (1 = issues found, 0 = no issues)
	if err != nil {
		result.Passed = false
	} else {
		result.Passed = true
	}

	return result
}

// runTestsWithFullOutput runs tests with streaming output directly to stdout
func runTestsWithFullOutput(servicePath, serviceName string) *TestResult {
	startTime := time.Now()

	// Create command with context and set the working directory
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Use classic go test (no JSON flag) for direct stdout output
	cmd := exec.CommandContext(ctx, "go", "test", "./...")
	cmd.Dir = servicePath

	// Stream output directly to stdout and stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the command
	if err := cmd.Start(); err != nil {
		duration := time.Since(startTime)
		return &TestResult{
			Passed:   false,
			Output:   "",
			Error:    err.Error(),
			Duration: duration,
		}
	}

	// Wait for completion
	err := cmd.Wait()
	duration := time.Since(startTime)

	// Since we're streaming directly, we don't capture output
	result := &TestResult{
		Output:       "", // Output is streamed directly to stdout
		Error:        "", // Errors are streamed directly to stderr
		Duration:     duration,
		TestCount:    0, // Can't parse from classic output
		PassCount:    0, // Can't parse from classic output
		FailCount:    0, // Can't parse from classic output
		SkipCount:    0, // Can't parse from classic output
		SkippedTests: []SkippedTest{},
		FailedTests:  []FailedTest{},
	}

	if err != nil {
		result.Passed = false
	} else {
		result.Passed = true
	}

	return result
}

// filterTestOutput extracts only the "Output" fields from Go test JSON output
func filterTestOutput(jsonOutput string) string {
	var filteredLines []string

	lines := strings.Split(jsonOutput, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var testEvent struct {
			Action string `json:"Action"`
			Output string `json:"Output"`
		}

		if err := json.Unmarshal([]byte(line), &testEvent); err != nil {
			continue
		}

		// Only include lines that have an "Output" field
		if testEvent.Output != "" {
			filteredLines = append(filteredLines, testEvent.Output)
		}
	}

	return strings.Join(filteredLines, "\n")
}

// runGolintWithFullOutput runs golint with full output display
func runGolintWithFullOutput(servicePath, serviceName string) *GolintResult {
	startTime := time.Now()

	// Check if golint is available
	if _, err := exec.LookPath("golint"); err != nil {
		return &GolintResult{
			Passed:     true,
			Output:     "golint not available",
			Error:      "golint command not found",
			Duration:   time.Since(startTime),
			IssueCount: 0,
		}
	}

	// Try to find go files
	goFiles, err := filepath.Glob(filepath.Join(servicePath, "**/*.go"))
	if err != nil || len(goFiles) == 0 {
		return &GolintResult{
			Passed:     true,
			Output:     "No go files found",
			Duration:   time.Since(startTime),
			IssueCount: 0,
		}
	}

	// Run golint if available
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "golint", "-set_exit_status", "./...")
	cmd.Dir = servicePath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	duration := time.Since(startTime)

	output := stdout.String()

	// Parse linting issues
	lintingIssues := parseLintingIssues(output)

	result := &GolintResult{
		Output:        output,
		Error:         stderr.String(),
		Duration:      duration,
		IssueCount:    len(lintingIssues),
		LintingIssues: lintingIssues,
	}

	// Determine if golint passed based on exit code (1 = issues found, 0 = no issues)
	if err != nil {
		result.Passed = false
	} else {
		result.Passed = true
	}

	return result
}

// runGRPCLintWithFullOutput runs gRPC lint with full output display
func runGRPCLintWithFullOutput(servicePath, serviceName string) *LintResult {
	startTime := time.Now()

	// Try to find proto files
	protoFiles, err := filepath.Glob(filepath.Join(servicePath, "**/*.proto"))
	if err != nil || len(protoFiles) == 0 {
		return &LintResult{
			Passed:   true,
			Output:   "No proto files found",
			Duration: time.Since(startTime),
		}
	}

	// Run buf lint if available
	cmd := exec.Command("buf", "lint")
	cmd.Dir = servicePath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd = exec.CommandContext(ctx, "buf", "lint")

	err = cmd.Run()
	duration := time.Since(startTime)

	result := &LintResult{
		Output:   stdout.String(),
		Error:    stderr.String(),
		Duration: duration,
	}

	if err != nil {
		result.Passed = false
	} else {
		result.Passed = true
	}

	return result
}

// displayDetailedIssues displays linting issues
func displayDetailedIssues(results []*ServiceResult) {
	hasIssues := false

	// Check if there are any linting issues to display
	for _, result := range results {
		if result.GolintResult != nil && len(result.GolintResult.LintingIssues) > 0 {
			hasIssues = true
			break
		}
	}

	if !hasIssues {
		return
	}

	fmt.Println("\n" + strings.Repeat("=", 100))
	fmt.Println("🔍 DETAILED ISSUES")
	fmt.Println(strings.Repeat("=", 100))

	// Display linting issues
	for _, result := range results {
		if result.GolintResult != nil && len(result.GolintResult.LintingIssues) > 0 {
			fmt.Printf("\n🔍 LINTING ISSUES - %s\n", strings.ToUpper(result.Name))
			fmt.Println(strings.Repeat("-", 50))

			for _, issue := range result.GolintResult.LintingIssues {
				fmt.Printf("%s:%d:%d: %s - %s\n", issue.File, issue.Line, issue.Column, issue.Rule, issue.Message)
			}
			fmt.Println()
		}
	}
}

func displayResults(results []*ServiceResult, startTime time.Time) {
	fmt.Println("\n" + strings.Repeat("=", 100))
	fmt.Println("📊 RESULTS SUMMARY")
	fmt.Println(strings.Repeat("=", 100))

	summary := &Summary{
		TotalServices: len(results),
	}

	// Calculate summary data
	for _, result := range results {
		if result.HasTests {
			if result.GolintResult != nil {
				if result.GolintResult.Passed {
					summary.PassedGolint++
				} else {
					summary.FailedGolint++
				}
				summary.TotalGolintIssues += result.GolintResult.IssueCount
			}

			if result.TestResult.Passed {
				summary.PassedTests++
			} else {
				summary.FailedTests++
			}

			if result.TestResult != nil {
				summary.TotalDuration += result.TestResult.Duration
				summary.TotalTestCount += result.TestResult.TestCount
				summary.TotalPassCount += result.TestResult.PassCount
				summary.TotalFailCount += result.TestResult.FailCount
				summary.TotalSkipCount += result.TestResult.SkipCount
			}
			if result.GolintResult != nil {
				summary.TotalDuration += result.GolintResult.Duration
			}
		}

		if result.HasGRPC {
			if result.LintResult.Passed {
				summary.PassedLint++
			} else {
				summary.FailedLint++
			}
			if result.LintResult != nil {
				summary.TotalDuration += result.LintResult.Duration
			}
		}
	}

	// Display results in table format
	fmt.Printf("%-15s %-8s %-8s %-8s %-8s %-12s %-8s\n",
		"SERVICE", "GOLINT", "TESTS", "GRPC", "DURATION", "TEST STATS", "ISSUES")
	fmt.Println(strings.Repeat("-", 100))

	for _, result := range results {
		// Get status icons
		golintStatus := "N/A"
		if result.HasTests && result.GolintResult != nil {
			if result.GolintResult.Passed {
				golintStatus = "✅"
			} else {
				golintStatus = "❌"
			}
		}

		testStatus := "N/A"
		testStats := ""
		if result.HasTests && result.TestResult != nil {
			if result.TestResult.Passed {
				testStatus = "✅"
			} else {
				testStatus = "❌"
			}
			testStats = fmt.Sprintf("%d/%d/%d", result.TestResult.PassCount, result.TestResult.FailCount, result.TestResult.SkipCount)
		}

		grpcStatus := "N/A"
		if result.HasGRPC && result.LintResult != nil {
			if result.LintResult.Passed {
				grpcStatus = "✅"
			} else {
				grpcStatus = "❌"
			}
		}

		// Calculate total duration for this service
		totalDuration := time.Duration(0)
		if result.TestResult != nil {
			totalDuration += result.TestResult.Duration
		}
		if result.GolintResult != nil {
			totalDuration += result.GolintResult.Duration
		}
		if result.LintResult != nil {
			totalDuration += result.LintResult.Duration
		}

		// Get issue count
		issueCount := 0
		if result.GolintResult != nil {
			issueCount = result.GolintResult.IssueCount
		}

		fmt.Printf("%-15s %-8s %-8s %-8s %-8s %-12s %-8d\n",
			result.Name,
			golintStatus,
			testStatus,
			grpcStatus,
			fmt.Sprintf("%.2fs", totalDuration.Seconds()),
			testStats,
			issueCount)
	}

	fmt.Println(strings.Repeat("-", 100))

	// Test statistics table
	fmt.Printf("\n📊 Test Statistics:\n")
	fmt.Printf("%-20s %-15s %-15s %-15s %-15s\n",
		"STATISTIC", "COUNT", "PERCENTAGE", "DURATION", "RATE")
	fmt.Println(strings.Repeat("-", 80))

	totalTests := summary.TotalTestCount
	if totalTests > 0 {
		passRate := float64(summary.TotalPassCount) / float64(totalTests) * 100
		failRate := float64(summary.TotalFailCount) / float64(totalTests) * 100
		skipRate := float64(summary.TotalSkipCount) / float64(totalTests) * 100

		fmt.Printf("%-20s %-15d %-15.1f%% %-15s %-15s\n",
			"Total Tests", totalTests, 100.0,
			fmt.Sprintf("%.2fs", summary.TotalDuration.Seconds()), "N/A")
		fmt.Printf("%-20s %-15d %-15.1f%% %-15s %-15.1f/s\n",
			"Passed", summary.TotalPassCount, passRate,
			fmt.Sprintf("%.2fs", summary.TotalDuration.Seconds()),
			float64(summary.TotalPassCount)/summary.TotalDuration.Seconds())
		fmt.Printf("%-20s %-15d %-15.1f%% %-15s %-15s\n",
			"Failed", summary.TotalFailCount, failRate, "N/A", "N/A")
		fmt.Printf("%-20s %-15d %-15.1f%% %-15s %-15s\n",
			"Skipped", summary.TotalSkipCount, skipRate, "N/A", "N/A")
	} else {
		fmt.Printf("%-20s %-15s %-15s %-15s %-15s\n",
			"No tests found", "N/A", "N/A", "N/A", "N/A")
	}

	fmt.Println(strings.Repeat("-", 80))

	// Display brief failure summary
	if summary.FailedTests > 0 || summary.FailedLint > 0 || summary.FailedGolint > 0 {
		fmt.Println("\n" + strings.Repeat("=", 50))
		fmt.Println("🚨 FAILURE SUMMARY")
		fmt.Println(strings.Repeat("=", 50))

		var failedServices []string
		for _, result := range results {
			hasFailures := false
			failureTypes := []string{}

			if result.HasTests && result.GolintResult != nil && !result.GolintResult.Passed {
				hasFailures = true
				failureTypes = append(failureTypes, fmt.Sprintf("golint (%d issues)", result.GolintResult.IssueCount))
			}

			if result.HasTests && result.TestResult != nil && !result.TestResult.Passed {
				hasFailures = true
				failureTypes = append(failureTypes, fmt.Sprintf("tests (%d failed)", result.TestResult.FailCount))
			}

			if result.HasGRPC && result.LintResult != nil && !result.LintResult.Passed {
				hasFailures = true
				failureTypes = append(failureTypes, "gRPC lint")
			}

			if hasFailures {
				failedServices = append(failedServices, fmt.Sprintf("  • %s: %s", result.Name, strings.Join(failureTypes, ", ")))
			}
		}

		for _, service := range failedServices {
			fmt.Println(service)
		}

		fmt.Println("\n💡 To see detailed logs for a specific service, run:")
		fmt.Println("   go run scripts/run-tests-and-lint/main.go --service <service-name>")
		fmt.Println("\n   Available services: auth, productivity, grpc, mail, mail-server, shared")
	}

	// Exit with appropriate code
	if summary.FailedTests > 0 || summary.FailedLint > 0 || summary.FailedGolint > 0 {
		fmt.Printf("\n❌ Some checks failed. Exiting with code 1.\n")
		os.Exit(1)
	} else {
		fmt.Printf("\n✅ All checks passed successfully!\n")
	}
}

// displayGolintIssuesTable displays golint output in a simple, clean format
func displayGolintIssuesTable(output string) {
	lines := strings.Split(strings.TrimSpace(output), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 4) // filename:line:column:message
		if len(parts) >= 4 {
			filename := parts[0]
			lineNum := parts[1]
			message := strings.TrimSpace(parts[3])

			// Simple format: filename:line: message
			fmt.Printf("  %s:%s: %s\n", filename, lineNum, message)
		}
	}
}
