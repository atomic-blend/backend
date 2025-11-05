package ghutils

// GetLatestImageVersion fetches the latest image version from GHCR.
// If rc is true, it fetches the latest release candidate version,
// otherwise it fetches the latest stable version.
func GetLatestImageVersion(imageName string, rc *bool) (string, error) {
	//TODO: get latest stable 
	//TODO: get latest rc
	//TODO: return the most recent between stable and rc if rc is true
	//TODO: return stable if rc is false
	return "latest", nil
}
