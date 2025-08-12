package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type ServiceResult struct {
	Name       string
	TestResult *TestResult
	LintResult *LintResult
	HasTests   bool
	HasGRPC    bool
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
	TotalServices  int
	PassedTests    int
	FailedTests    int
	PassedLint     int
	FailedLint     int
	TotalDuration  time.Duration
	TotalTestCount int
	TotalPassCount int
	TotalFailCount int
	TotalSkipCount int
}

func main() {
	fmt.Println("🚀 Running Microservice Tests and gRPC Linting")
	fmt.Println("📍 This script should be run from the backend directory")
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
	results := make([]*ServiceResult, 0, len(services))

	// Count total services that will be processed
	totalServices := 0
	for _, service := range services {
		servicePath := filepath.Join(workspaceRoot, service)
		if _, err := os.Stat(servicePath); err == nil {
			totalServices++
		}
	}

	fmt.Printf("📊 Found %d services to process\n", totalServices)
	fmt.Println()

	// Run tests and linting for each service
	foundServices := 0
	currentService := 0
	for _, service := range services {
		servicePath := filepath.Join(workspaceRoot, service)

		// Check if service directory exists
		if _, err := os.Stat(servicePath); os.IsNotExist(err) {
			fmt.Printf("⚠️  Service %s not found, skipping...\n", service)
			continue
		}
		foundServices++
		currentService++

		fmt.Printf("\n🔍 Processing %s... (%d/%d)\n", service, currentService, totalServices)
		printProgressBar(currentService, totalServices, "📈 Overall Progress")

		result := &ServiceResult{
			Name: service,
		}

		// Check if service has tests
		if hasTests(servicePath) {
			result.HasTests = true
			result.TestResult = runTests(servicePath, service)
		}

		// Check if service has gRPC
		if hasGRPC(servicePath) {
			result.HasGRPC = true
			result.LintResult = runGRPCLint(servicePath, service)
		}

		results = append(results, result)
	}

	// Complete the progress bar
	fmt.Println() // New line after progress bar

	// Check if we found any services
	if foundServices == 0 {
		fmt.Printf("\n❌ Error: No services found. Please run this script from the backend directory.\n")
		fmt.Printf("Current directory: %s\n", workspaceRoot)
		fmt.Printf("Expected to find services like: auth, mail, productivity, etc.\n")
		os.Exit(1)
	}

	// Display results
	displayResults(results, startTime)
}

func hasTests(servicePath string) bool {
	// Check for go.mod file - if it exists, the service can run tests
	_, err := os.Stat(filepath.Join(servicePath, "go.mod"))
	return err == nil
}

func hasGRPC(servicePath string) bool {
	// Only run gRPC lint for the grpc directory
	serviceName := filepath.Base(servicePath)
	return serviceName == "grpc"
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
	fmt.Printf("  🧪 Running tests for %s...\n", serviceName)

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
		fmt.Printf("  ❌ Tests failed for %s (%.2fs) - %d passed, %d failed, %d skipped\n",
			serviceName, duration.Seconds(), passCount, failCount, skipCount)
	} else {
		result.Passed = true
		fmt.Printf("  ✅ Tests passed for %s (%.2fs) - %d passed, %d failed, %d skipped\n",
			serviceName, duration.Seconds(), passCount, failCount, skipCount)
	}

	return result
}

func runGRPCLint(servicePath, serviceName string) *LintResult {
	fmt.Printf("  🔍 Running gRPC lint for %s...\n", serviceName)

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
		if result.HasGRPC && result.LintResult != nil {
			summary.TotalDuration += result.LintResult.Duration
		}
	}

	// Display summary
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("📈 SUMMARY")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Total Services: %d\n", summary.TotalServices)
	fmt.Printf("Tests Passed: %d\n", summary.PassedTests)
	fmt.Printf("Tests Failed: %d\n", summary.FailedTests)
	fmt.Printf("gRPC Lint Passed: %d\n", summary.PassedLint)
	fmt.Printf("gRPC Lint Failed: %d\n", summary.FailedLint)
	fmt.Printf("Total Duration: %.2fs\n", summary.TotalDuration.Seconds())
	fmt.Printf("\n📊 Test Statistics:\n")
	fmt.Printf("   Total Tests: %d\n", summary.TotalTestCount)
	fmt.Printf("   Passed: %d\n", summary.TotalPassCount)
	fmt.Printf("   Failed: %d\n", summary.TotalFailCount)
	fmt.Printf("   Skipped: %d\n", summary.TotalSkipCount)

	// Display detailed failures
	if summary.FailedTests > 0 || summary.FailedLint > 0 {
		fmt.Println("\n" + strings.Repeat("=", 50))
		fmt.Println("🚨 DETAILED FAILURES")
		fmt.Println(strings.Repeat("=", 50))

		for _, result := range results {
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
	if summary.FailedTests > 0 || summary.FailedLint > 0 {
		fmt.Printf("\n❌ Some checks failed. Exiting with code 1.\n")
		os.Exit(1)
	} else {
		fmt.Printf("\n✅ All checks passed successfully!\n")
	}
}
