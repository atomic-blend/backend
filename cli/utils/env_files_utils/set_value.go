package envfilesutils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// SetEnvVarValue updates or adds an environment variable in the given .env
// file path. It preserves comments and other lines. If the variable exists,
// its value is replaced in-place while attempting to keep the original
// quoting style and inline comment. If the variable doesn't exist, it is
// appended at the end of the file.
func SetEnvVarValue(envPath, varName, newVal string) error {
	// Read file
	f, err := os.Open(envPath)
	if err != nil {
		return fmt.Errorf("failed to open env file %s: %w", envPath, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var lines []string
	found := false

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// skip empty and comment-only lines
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			lines = append(lines, line)
			continue
		}

		// look for key = value pattern (allow optional export)
		// find first '=' outside quotes
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
			lines = append(lines, line)
			continue
		}

		left := strings.TrimSpace(line[:eqIdx])
		right := line[eqIdx+1:]

		// extract key (allow leading export)
		leftParts := strings.Fields(left)
		key := leftParts[len(leftParts)-1]

		if key != varName {
			lines = append(lines, line)
			continue
		}

		// parse right side to separate value and inline comment, respecting quotes
		valPart := right
		commentPart := ""
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
					commentPart = right[i:]
					break
				}
			}
		}

		// determine original quoting
		vtrim := strings.TrimSpace(valPart)
		quoteChar := byte(0)
		if len(vtrim) > 1 {
			if vtrim[0] == '\'' && vtrim[len(vtrim)-1] == '\'' {
				quoteChar = '\''
			} else if vtrim[0] == '"' && vtrim[len(vtrim)-1] == '"' {
				quoteChar = '"'
			}
		}

		// build new value preserving quote style if present
		newValStr := newVal
		if quoteChar != 0 {
			newValStr = string(quoteChar) + newVal + string(quoteChar)
		}

		// preserve spacing around '=' by reusing left part exactly as in original
		newLine := line[:eqIdx+1] + " " + newValStr
		if commentPart != "" {
			// ensure a space before inline comment if not present
			if !strings.HasPrefix(commentPart, " ") {
				newLine += " " + commentPart
			} else {
				newLine += commentPart
			}
		}

		lines = append(lines, newLine)
		found = true
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading env file %s: %w", envPath, err)
	}

	if !found {
		// append new variable at end
		lines = append(lines, fmt.Sprintf("%s=%s", varName, newVal))
	}

	// write back
	out := strings.Join(lines, "\n")
	// ensure file ends with newline
	if !strings.HasSuffix(out, "\n") {
		out += "\n"
	}

	if err := os.WriteFile(envPath, []byte(out), 0o644); err != nil {
		return fmt.Errorf("failed to write env file %s: %w", envPath, err)
	}

	return nil
}
