package types

// GithubRemoteFile describes a file stored in a GitHub repo that we want to
// download into a local path.
type GithubRemoteFile struct {
	LocalPath  string
	GitHubPath string
	Repository string
}
