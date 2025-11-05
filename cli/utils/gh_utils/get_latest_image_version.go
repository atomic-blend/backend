package ghutils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/atomic-blend/backend/cli/config"
	"github.com/google/go-github/v77/github"
	"golang.org/x/oauth2"
)

// GetLatestImageVersion fetches the latest image version from GHCR.
// If rc is true, it fetches the latest release candidate version,
// otherwise it fetches the latest stable version.
//
// Expected input: imageName must be a GHCR image, e.g. "ghcr.io/atomic-blend/auth".
// Tags in GHCR are expected in the form: "auth/v0.12.0-rc-47833be" or "auth/0.12.0".
// For stable versions we remove the "-rc-xxxxx" suffix if present and return the
// version part (e.g. "v0.12.0" or "0.12.0"). For RC mode we return the RC
// tag (without removing the rc suffix) if it is the most recent.
func GetLatestImageVersion(imageName string, rc *bool) (string, error) {
	if !strings.HasPrefix(imageName, "ghcr.io/") {
		return "", fmt.Errorf("unsupported image registry: %s", imageName)
	}

	// parse org and repo
	parts := strings.Split(strings.TrimPrefix(imageName, "ghcr.io/"), "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid ghcr image name: %s", imageName)
	}
	owner := parts[0]
	pkgName := parts[1]

	ctx := context.Background()

	// build HTTP client with optional token
	var httpClient *http.Client
	if token := config.CliConfig.GithubToken; token != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		httpClient = oauth2.NewClient(ctx, ts)
	}
	client := github.NewClient(httpClient)

	// Use the Packages API endpoint for versions: /orgs/{owner}/packages/container/{packageName}/versions
	url := fmt.Sprintf("/orgs/%s/packages/container/%s/versions", owner, pkgName)

	req, err := client.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// We will decode the response into a lightweight struct matching the fields we need
	type containerMeta struct {
		Tags []string `json:"tags"`
	}
	type meta struct {
		Container containerMeta `json:"container"`
	}
	type versionItem struct {
		CreatedAt time.Time `json:"created_at"`
		Metadata  meta      `json:"metadata"`
	}

	var items []versionItem
	_, err = client.Do(ctx, req, &items)
	if err != nil {
		return "", err
	}

	if len(items) == 0 {
		return "", errors.New("no versions found for package")
	}

	// Helpers
	stripServicePrefix := func(tag string) string {
		if idx := strings.Index(tag, "/"); idx != -1 {
			return tag[idx+1:]
		}
		return tag
	}

	isRC := func(ver string) bool {
		return strings.Contains(ver, "-rc-")
	}

	// normalize semver for comparison (simple numeric compare)
	parseSemver := func(v string) ([]int, error) {
		v = strings.TrimPrefix(v, "v")
		parts := strings.Split(v, ".")
		nums := make([]int, 3)
		for i := 0; i < 3; i++ {
			if i < len(parts) {
				// remove any suffix like -rc-... if present
				p := parts[i]
				if idx := strings.Index(p, "-rc-"); idx != -1 {
					p = p[:idx]
				}
				n, err := strconv.Atoi(p)
				if err != nil {
					return nil, err
				}
				nums[i] = n
			} else {
				nums[i] = 0
			}
		}
		return nums, nil
	}

	semverRegexp := regexp.MustCompile(`^v?\d+\.\d+(?:\.\d+)?$`)

	// Track best stable (by semver) and best RC (by created_at)
	var bestStable string
	var bestStableTime time.Time
	var bestStableNums []int

	var bestRC string
	var bestRCTime time.Time

	// Also keep a fallback most recent tag
	var mostRecentTag string
	var mostRecentTime time.Time

	for _, it := range items {
		created := it.CreatedAt
		for _, rawTag := range it.Metadata.Container.Tags {
			tag := stripServicePrefix(rawTag)

			// update most recent
			if mostRecentTag == "" || created.After(mostRecentTime) {
				mostRecentTag = tag
				mostRecentTime = created
			}

			// RC handling
			if isRC(tag) {
				if bestRC == "" || created.After(bestRCTime) {
					bestRC = tag
					bestRCTime = created
				}
				continue
			}

			// Check stable semver
			if semverRegexp.MatchString(tag) {
				nums, err := parseSemver(tag)
				if err != nil {
					// skip unparsable
					continue
				}
				if bestStable == "" {
					bestStable = tag
					bestStableNums = nums
					bestStableTime = created
					continue
				}
				// compare nums
				greater := false
				for i := 0; i < 3; i++ {
					if nums[i] > bestStableNums[i] {
						greater = true
						break
					} else if nums[i] < bestStableNums[i] {
						break
					}
				}
				if greater {
					bestStable = tag
					bestStableNums = nums
					bestStableTime = created
				}
			}
		}
	}

	// If rc flag is provided and true -> prefer RC when it's more recent than stable
	wantRC := false
	if rc != nil && *rc {
		wantRC = true
	}

	if wantRC {
		if bestRC != "" {
			// if we have a stable, compare dates (prefer RC only if newer)
			if bestStable == "" {
				return bestRC, nil
			}
			if bestRCTime.After(bestStableTime) {
				return bestRC, nil
			}
			// otherwise fall back to stable
			return bestStable, nil
		}
		// no RC found -> fall back to stable if present
		if bestStable != "" {
			return bestStable, nil
		}
	} else {
		if bestStable != "" {
			return bestStable, nil
		}
	}

	// last resort: return the most recent tag found
	if mostRecentTag != "" {
		// If the most recent tag contains -rc- and stable was requested, strip the rc suffix
		if !wantRC && isRC(mostRecentTag) {
			// strip -rc-xxxxx
			if idx := strings.Index(mostRecentTag, "-rc-"); idx != -1 {
				return mostRecentTag[:idx], nil
			}
		}
		return mostRecentTag, nil
	}

	return "", errors.New("no suitable version found")
}
