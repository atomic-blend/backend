package config

// Config holds CLI configuration that can be populated via env / unmarshal.
// Fields must be exported (capitalized) so reflection-based unmarshalers can set them.
type Config struct {
	Debug          bool   `mapstructure:"debug" json:"debug"`
	Channel        string `mapstructure:"channel" json:"channel"`
	Directory      string `mapstructure:"directory" json:"directory"`
	DevelopperMode bool   `mapstructure:"developper-mode" json:"developper_mode"`
	GithubToken    string `mapstructure:"github-token" json:"github_token"`
}

// CliConfig is the package-level configuration instance used by the CLI.
var CliConfig Config

// PlatformComponents is the list of all platform components.
// These correspond to services defined in docker-compose.yml and
// components that have versions defined in the .env file.
var PlatformComponents = []string{
	"auth", "mail", "mail-server", "calendar", "productivity", "task-app", "notes-app", "mail-app",
}

// EnvComponentVersionMapping maps component names to their corresponding environment variable names for versioning.
var EnvComponentVersionMapping = map[string]string{
	"auth":         "AUTH_SERVICE_VERSION",
	"mail":         "MAIL_SERVICE_VERSION",
	"mail-server":  "MAIL_SERVER_SERVICE_VERSION",
	"calendar":     "CALENDAR_SERVICE_VERSION",
	"productivity": "PRODUCTIVITY_SERVICE_VERSION",
	"task-app":     "TASK_APP_VERSION",
	"notes-app":    "NOTES_APP_VERSION",
	"mail-app":     "MAIL_APP_VERSION",
}

var BackendServices = []string{
	"auth",
	"productivity",
	"calendar",
	"grpc",
	"mail",
	"mail-server",
	"shared",
}
