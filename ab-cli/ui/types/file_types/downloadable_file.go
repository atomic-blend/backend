package filetypes

// DownloadableFile represents a file that can be downloaded from a URL and
// written to a local path. LocalPath is the final absolute or relative path
// where the file should be stored (the caller is responsible for joining
// config.CliConfig.Directory if desired).
type DownloadableFile struct {
	URL       string
	LocalPath string
}