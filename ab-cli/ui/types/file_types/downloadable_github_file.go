package filetypes

// DownloadableGithubFile describes a file stored in a GitHub repo that we want to
// download into a local path.
type DownloadableGithubFile struct {
	LocalPath  string
	GitHubPath string
	Repository string
}
