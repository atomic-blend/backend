package main

import (
	"bytes"
	"context"
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
	Passed   bool
	Output   string
	Error    string
	Duration time.Duration
}

type LintResult struct {
	Passed   bool
	Output   string
	Error    string
	Duration time.Duration
}

type Summary struct {
	TotalServices int
	PassedTests   int
	FailedTests   int
	PassedLint    int
	FailedLint    int
	TotalDuration time.Duration
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

	// Run tests and linting for each service
	foundServices := 0
	for _, service := range services {
		servicePath := filepath.Join(workspaceRoot, service)

		// Check if service directory exists
		if _, err := os.Stat(servicePath); os.IsNotExist(err) {
			fmt.Printf("⚠️  Service %s not found, skipping...\n", service)
			continue
		}
		foundServices++

		fmt.Printf("\n🔍 Processing %s...\n", service)

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

func runTests(servicePath, serviceName string) *TestResult {
	fmt.Printf("  🧪 Running tests for %s...\n", serviceName)

	startTime := time.Now()

	// Create command with context and set the working directory
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "test", "./...", "-v")
	cmd.Dir = servicePath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	duration := time.Since(startTime)

	result := &TestResult{
		Output:   stdout.String(),
		Error:    stderr.String(),
		Duration: duration,
	}

	if err != nil {
		result.Passed = false
		fmt.Printf("  ❌ Tests failed for %s (%.2fs)\n", serviceName, duration.Seconds())
	} else {
		result.Passed = true
		fmt.Printf("  ✅ Tests passed for %s (%.2fs)\n", serviceName, duration.Seconds())
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
				fmt.Printf("   ✅ Tests: PASSED (%.2fs)\n", result.TestResult.Duration.Seconds())
				summary.PassedTests++
			} else {
				fmt.Printf("   ❌ Tests: FAILED (%.2fs)\n", result.TestResult.Duration.Seconds())
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
