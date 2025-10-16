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
	} else {
		result.Passed = true
	}

	return result
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

	// Display summary in table format
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📈 SUMMARY")
	fmt.Println(strings.Repeat("=", 80))

	// Summary table
	fmt.Printf("%-20s %-15s %-15s %-15s %-15s\n",
		"CATEGORY", "PASSED", "FAILED", "TOTAL", "ISSUES")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%-20s %-15d %-15d %-15d %-15s\n",
		"Services", summary.TotalServices, 0, summary.TotalServices, "N/A")
	fmt.Printf("%-20s %-15d %-15d %-15d %-15d\n",
		"Golint", summary.PassedGolint, summary.FailedGolint,
		summary.PassedGolint+summary.FailedGolint, summary.TotalGolintIssues)
	fmt.Printf("%-20s %-15d %-15d %-15d %-15s\n",
		"Tests", summary.PassedTests, summary.FailedTests,
		summary.PassedTests+summary.FailedTests, "N/A")
	fmt.Printf("%-20s %-15d %-15d %-15d %-15s\n",
		"gRPC Lint", summary.PassedLint, summary.FailedLint,
		summary.PassedLint+summary.FailedLint, "N/A")
	fmt.Println(strings.Repeat("-", 80))

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
