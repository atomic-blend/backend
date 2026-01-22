# run with docker

## `ab-manager` (helper script)

A helper utility is provided to check the latest versions of `ghcr.io/atomic-blend/*` images and help update the version variables in the `.env` file.

Usage (run from the `backend/docker` folder or provide the script path):

```bash
./ab-manager.sh [OPTIONS]
```

Available options:
- `-c, --compose-file FILE`    : Docker Compose file to parse (default: `docker-compose.yaml`).
- `-e, --env-file FILE`        : `.env` file to update (default: `.env`).
- `-b, --branch BRANCH`        : GitHub branch to use when downloading files if needed (default: `main`).
- `--rc`                      : Prefer release-candidate (RC) images when available and newer than the stable release.
- `-h, --help`                : Show help.

Behavior of the `--rc` flag:
- By default the script selects the latest stable version using semantic versioning.
- If `--rc` is provided, the script will also look for RC-style tags (for example `0.12.0-rc-d5fa8f8`) and will propose the most recent RC only if it is newer than the latest stable release. The comparison is based on the release date. If the stable release is newer or there is no valid RC, the stable version is returned.

Requirements:
- `curl` and `jq` must be installed.
- The script will prompt for `GITHUB_TOKEN` if the environment variable is not set (recommended for authenticated GitHub API access).

Examples:

```bash
# Check and propose the latest stable versions (default behavior)
./ab-manager.sh

# Check and propose RCs when they are more recent than stable
./ab-manager.sh --rc
```
