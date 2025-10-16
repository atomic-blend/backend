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
	Passed    bool
	Output    string
	Error     string
	Duration  time.Duration
	TestCount int
	PassCount int
	FailCount int
	SkipCount int
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
	Passed     bool
	Output     string
	Error      string
	Duration   time.Duration
	IssueCount int
}

// Global status tracking for real-time table updates
var (
	statusMap = make(map[string]*ServiceStatus)
	statusMux sync.RWMutex
)

func main() {
	// Parse command line flags
	var numThreads = flag.Int("threads", 4, "Number of parallel threads to use for running tests and linting")
	flag.Parse()

	fmt.Println("🚀 Running Microservice Tests, Golint, and gRPC Linting")
	fmt.Println("📍 This script should be run from the backend directory")
	fmt.Printf("🧵 Using %d parallel threads\n", *numThreads)
	fmt.Println(strings.Repeat("=", 50))

	// Get the workspace root (assuming script is run from backend directory)
	workspaceRoot, err := os.Getwd()
	if err != nil {
		fmt.Printf("❌ Error getting current directory: %v\n", err)
		os.Exit(1)
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

	startTime := time.Now()

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
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-stopChan:
			return
		case <-ticker.C:
			displayStatusTable(totalServices)
		}
	}
}

// displayStatusTable displays the current status of all services
func displayStatusTable(totalServices int) {
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
		golintIcon := getCheckIcon(service.GolintStatus)
		testIcon := getCheckIcon(service.TestStatus)
		grpcIcon := getCheckIcon(service.GRPCStatus)

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
		return "🔄"
	case "passed":
		return "✅"
	case "failed":
		return "❌"
	default:
		return "❓"
	}
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
func parseTestOutput(output string) (int, int, int, int) {
	var testCount, passCount, failCount, skipCount int

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var testEvent struct {
			Action string `json:"Action"`
			Test   string `json:"Test"`
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
		case "skip":
			skipCount++
		}
	}

	return testCount, passCount, failCount, skipCount
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
	testCount, passCount, failCount, skipCount := parseTestOutput(output)

	result := &TestResult{
		Output:    output,
		Error:     errorOutput,
		Duration:  duration,
		TestCount: testCount,
		PassCount: passCount,
		FailCount: failCount,
		SkipCount: skipCount,
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
		fmt.Printf("  ❌ gRPC lint failed for %s (%.2fs)\n", serviceName, duration.Seconds())
	} else {
		result.Passed = true
		fmt.Printf("  ✅ gRPC lint passed for %s (%.2fs)\n", serviceName, duration.Seconds())
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

	result := &GolintResult{
		Output:   stdout.String(),
		Error:    stderr.String(),
		Duration: duration,
	}

	// Count issues (lines in Vim quickfix format: filename:line:column: message)
	issueCount := 0
	output := stdout.String()
	if output != "" {
		lines := strings.Split(strings.TrimSpace(output), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && strings.Contains(line, ":") && strings.Count(line, ":") >= 3 {
				// Check if it's a valid quickfix format line
				parts := strings.SplitN(line, ":", 4)
				if len(parts) >= 4 {
					issueCount++
				}
			}
		}
	}
	result.IssueCount = issueCount

	// Determine if golint passed based on exit code (1 = issues found, 0 = no issues)
	if err != nil {
		result.Passed = false
		fmt.Printf("  ❌ golint failed for %s (%.2fs) - Found %d issues\n", serviceName, duration.Seconds(), issueCount)
	} else {
		result.Passed = true
		fmt.Printf("  ✅ golint passed for %s (%.2fs) - Found %d issues\n", serviceName, duration.Seconds(), issueCount)
	}

	return result
}

func displayResults(results []*ServiceResult, startTime time.Time) {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("📊 RESULTS SUMMARY")
	fmt.Println(strings.Repeat("=", 50))

	summary := &Summary{
		TotalServices: len(results),
	}

	for _, result := range results {
		fmt.Printf("\n🏗️  %s\n", strings.ToUpper(result.Name))
		fmt.Printf("   %s\n", strings.Repeat("-", len(result.Name)+2))

		if result.HasTests {
			// Show golint results first
			if result.GolintResult != nil {
				if result.GolintResult.Passed {
					fmt.Printf("   ✅ Golint: PASSED (%.2fs) - %d issues found\n",
						result.GolintResult.Duration.Seconds(),
						result.GolintResult.IssueCount)
					summary.PassedGolint++
				} else {
					fmt.Printf("   ❌ Golint: FAILED (%.2fs) - %d issues found\n",
						result.GolintResult.Duration.Seconds(),
						result.GolintResult.IssueCount)
					summary.FailedGolint++
				}
				summary.TotalGolintIssues += result.GolintResult.IssueCount
			}

			// Then show test results
			if result.TestResult.Passed {
				fmt.Printf("   ✅ Tests: PASSED (%.2fs) - %d passed, %d failed, %d skipped\n",
					result.TestResult.Duration.Seconds(),
					result.TestResult.PassCount,
					result.TestResult.FailCount,
					result.TestResult.SkipCount)
				summary.PassedTests++
			} else {
				fmt.Printf("   ❌ Tests: FAILED (%.2fs) - %d passed, %d failed, %d skipped\n",
					result.TestResult.Duration.Seconds(),
					result.TestResult.PassCount,
					result.TestResult.FailCount,
					result.TestResult.SkipCount)
				summary.FailedTests++
			}
		} else {
			fmt.Printf("   ⚠️  Tests: No test files found\n")
		}

		if result.HasGRPC {
			if result.LintResult.Passed {
				fmt.Printf("   ✅ gRPC Lint: PASSED (%.2fs)\n", result.LintResult.Duration.Seconds())
				summary.PassedLint++
			} else {
				fmt.Printf("   ❌ gRPC Lint: FAILED (%.2fs)\n", result.LintResult.Duration.Seconds())
				summary.FailedLint++
			}
		} else {
			fmt.Printf("   ⚠️  gRPC Lint: No gRPC files found\n")
		}

		if result.HasTests && result.TestResult != nil {
			summary.TotalDuration += result.TestResult.Duration
			summary.TotalTestCount += result.TestResult.TestCount
			summary.TotalPassCount += result.TestResult.PassCount
			summary.TotalFailCount += result.TestResult.FailCount
			summary.TotalSkipCount += result.TestResult.SkipCount
		}
		if result.HasTests && result.GolintResult != nil {
			summary.TotalDuration += result.GolintResult.Duration
		}
		if result.HasGRPC && result.LintResult != nil {
			summary.TotalDuration += result.LintResult.Duration
		}
	}

	// Display summary
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("📈 SUMMARY")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Total Services: %d\n", summary.TotalServices)
	fmt.Printf("Golint Passed: %d\n", summary.PassedGolint)
	fmt.Printf("Golint Failed: %d\n", summary.FailedGolint)
	fmt.Printf("Tests Passed: %d\n", summary.PassedTests)
	fmt.Printf("Tests Failed: %d\n", summary.FailedTests)
	fmt.Printf("gRPC Lint Passed: %d\n", summary.PassedLint)
	fmt.Printf("gRPC Lint Failed: %d\n", summary.FailedLint)
	fmt.Printf("Total Duration: %.2fs\n", summary.TotalDuration.Seconds())
	fmt.Printf("Total Golint Issues: %d\n", summary.TotalGolintIssues)
	fmt.Printf("\n📊 Test Statistics:\n")
	fmt.Printf("   Total Tests: %d\n", summary.TotalTestCount)
	fmt.Printf("   Passed: %d\n", summary.TotalPassCount)
	fmt.Printf("   Failed: %d\n", summary.TotalFailCount)
	fmt.Printf("   Skipped: %d\n", summary.TotalSkipCount)

	// Display detailed failures
	if summary.FailedTests > 0 || summary.FailedLint > 0 || summary.FailedGolint > 0 {
		fmt.Println("\n" + strings.Repeat("=", 50))
		fmt.Println("🚨 DETAILED FAILURES")
		fmt.Println(strings.Repeat("=", 50))

		for _, result := range results {
			if result.HasTests && result.GolintResult != nil && !result.GolintResult.Passed {
				fmt.Printf("\n❌ %s - Golint Failure:\n", result.Name)
				fmt.Printf("Issues Found: %d\n", result.GolintResult.IssueCount)
				if result.GolintResult.Output != "" {
					fmt.Printf("Issues:\n")
					displayGolintIssuesTable(result.GolintResult.Output)
				}
			}

			if result.HasTests && !result.TestResult.Passed {
				fmt.Printf("\n❌ %s - Test Failure:\n", result.Name)
				fmt.Printf("Error: %s\n", result.TestResult.Error)
				if result.TestResult.Output != "" {
					fmt.Printf("Output: %s\n", result.TestResult.Output)
				}
			}

			if result.HasGRPC && !result.LintResult.Passed {
				fmt.Printf("\n❌ %s - gRPC Lint Failure:\n", result.Name)
				fmt.Printf("Error: %s\n", result.LintResult.Error)
				if result.LintResult.Output != "" {
					fmt.Printf("Output: %s\n", result.LintResult.Output)
				}
			}
		}
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
