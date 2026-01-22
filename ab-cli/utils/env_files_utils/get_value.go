// Package envfilesutils provides small helpers to read and write
// environment variables inside a dotenv-style file while preserving
// comments and quoting styles. It is used by CLI commands that need to
// programmatically inspect or mutate a project's `.env` file.
package envfilesutils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// GetEnvVarValue reads the provided .env file and returns the value for the
// given variable name. The parser attempts to respect quotes and inline
// comments (i.e. '#' outside of quotes). It also accepts a leading `export`
// token on the left-hand side (e.g. `export FOO=bar`).
func GetEnvVarValue(envPath, varName string) (string, error) {
	f, err := os.Open(envPath)
	if err != nil {
		return "", fmt.Errorf("failed to open env file %s: %w", envPath, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// skip empty and comment-only lines
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		// find first '=' that is not inside quotes
		eqIdx := -1
		inSingle := false
		inDouble := false
		for i, r := range line {
			if r == '\'' {
				if !inDouble {
					inSingle = !inSingle
				}
				continue
			}
			if r == '"' {
				if !inSingle {
					inDouble = !inDouble
				}
				continue
			}
			if r == '=' {
				if !inSingle && !inDouble {
					eqIdx = i
					break
				}
			}
		}

		if eqIdx < 0 {
			continue
		}

		left := strings.TrimSpace(line[:eqIdx])
		right := line[eqIdx+1:]

		// extract key (allow leading export)
		leftParts := strings.Fields(left)
		if len(leftParts) == 0 {
			continue
		}
		key := leftParts[len(leftParts)-1]
		if key != varName {
			continue
		}

		// parse right side to separate value and inline comment, respecting quotes
		valPart := right
		inS := false
		inD := false
		for i, r := range right {
			if r == '\'' {
				if !inD {
					inS = !inS
				}
				continue
			}
			if r == '"' {
				if !inS {
					inD = !inD
				}
				continue
			}
			if r == '#' {
				if !inS && !inD {
					valPart = right[:i]
					break
				}
			}
		}

		vtrim := strings.TrimSpace(valPart)

		// strip surrounding quotes if present
		if len(vtrim) >= 2 {
			if vtrim[0] == '\'' && vtrim[len(vtrim)-1] == '\'' {
				return vtrim[1 : len(vtrim)-1], nil
			}
			if vtrim[0] == '"' && vtrim[len(vtrim)-1] == '"' {
				return vtrim[1 : len(vtrim)-1], nil
			}
		}

		return vtrim, nil
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading env file %s: %w", envPath, err)
	}

	return "", fmt.Errorf("variable %q not found in %s", varName, envPath)
}
