// Package envvareditor provides a tiny interactive helper to inspect and
// optionally update a single environment variable in the project's .env
// file. It uses the env_files_utils helpers to read and write while
// preserving comments and quoting where possible.
package envvareditor

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/atomic-blend/backend/cli/config"
	env_files_utils "github.com/atomic-blend/backend/cli/utils/env_files_utils"
)

// EditEnvVar reads the current value for `envKey` from the project's
// `.env` (determined from `config.CliConfig.Directory`) and prompts the
// user whether they want to update it. The default answer is No.
// If the user chooses to update, they are prompted for the new value and
// the helper updates the file using env_files_utils.SetEnvVarValue.
// If skipIfSet is true and the variable is already defined in the .env,
// the function returns immediately with no output or prompt.
// Optionally, a `validator` function may be provided to validate the
// entered value. The validator should return nil for a valid value or an
// error describing the problem. If provided, the helper will keep asking
// until the user supplies a valid value or cancels. The `generationDoc`
// string, if non-empty, will be displayed as instructions to the user on
// how to generate/obtain a correct value. If `forceSet` is true the user
// will be required to provide a (non-empty, valid) value and the initial
// "Do you want to update?" prompt will be skipped.
func EditEnvVar(envKey string, skipIfSet bool, forceSet bool, validator func(string) error, generationDoc string) error {
	envPath := path.Join(config.CliConfig.Directory, ".env")

	curr, err := env_files_utils.GetEnvVarValue(envPath, envKey)
	if err != nil {
		// failed to read: treat as missing so we can offer to create it
		curr = ""
	} else {
		// If requested, skip entirely when the variable is already set.
		if skipIfSet && curr != "" {
			return nil
		}
		fmt.Printf("Current value for %s: %s\n", envKey, curr)
	}

	reader := bufio.NewReader(os.Stdin)

	// Decide whether to prompt the user or force a value.
	// `forceSet` applies only when the variable is currently missing.
	if forceSet && curr == "" {
		// skip initial confirmation and require a non-empty value
		fmt.Printf("%s must be set. You will be prompted to provide a value.\n", envKey)
	} else {
		// normal behaviour: ask whether to update (default No)
		fmt.Printf("Do you want to update %s? [y/N]: ", envKey)
		ans, _ := reader.ReadString('\n')
		ans = strings.TrimSpace(strings.ToLower(ans))
		if ans != "y" && ans != "yes" {
			fmt.Println("No changes made.")
			return nil
		}
	}

	// Show generation instructions if provided
	if generationDoc != "" {
		fmt.Println("Instructions:")
		fmt.Println(generationDoc)
	}

	// Prompt for new value and validate if a validator was provided.
	for {
		if curr != "" {
			fmt.Printf("Enter new value for %s (current: %s): ", envKey, curr)
		} else {
			fmt.Printf("Enter new value for %s: ", envKey)
		}
		newVal, _ := reader.ReadString('\n')
		newVal = strings.TrimSpace(newVal)
		if newVal == "" {
			if forceSet {
				fmt.Println("Value is required. Please provide a non-empty value.")
				continue
			}
			fmt.Println("Empty value entered, aborting.")
			return nil
		}

		if validator != nil {
			if verr := validator(newVal); verr != nil {
				fmt.Printf("Provided value is invalid: %v\n", verr)
				// Ask whether user wants to retry
				fmt.Printf("Retry entering value? [y/N]: ")
				retryAns, _ := reader.ReadString('\n')
				retryAns = strings.TrimSpace(strings.ToLower(retryAns))
				if retryAns == "y" || retryAns == "yes" {
					// loop again
					continue
				}
				fmt.Println("Aborting without changes.")
				return nil
			}
		}

		if err := env_files_utils.SetEnvVarValue(envPath, envKey, newVal); err != nil {
			return fmt.Errorf("failed to set %s in %s: %w", envKey, envPath, err)
		}

		fmt.Printf("Updated %s -> %s in %s\n", envKey, newVal, envPath)
		break
	}

	return nil
}
