package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// NewInitCommand creates the init command
func NewInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a new configuration",
		Long: `Initialize a new git-activity-mirror configuration file.

This command will guide you through setting up source and target platforms,
authentication, and mirroring preferences.`,
		RunE: runInit,
	}
}

func runInit(cmd *cobra.Command, args []string) error {
	fmt.Println("🚀 Welcome to git-activity-mirror!")
	fmt.Println()
	fmt.Println("This wizard will help you set up mirroring between git platforms.")
	fmt.Println("You can mirror activity from any platform (GitLab, GitHub, Bitbucket, etc.)")
	fmt.Println("to any other platform while keeping your code private.")
	fmt.Println()

	// Create basic configuration structure
	config := Config{
		Version: 1,
		Sources: []SourceConfig{
			{
				Name:     "work",
				Platform: "gitlab",
				Host:     "gitlab.com",
				Auth: AuthConfig{
					Type:     "token",
					Username: "your-username",
					Token:    "${GITLAB_TOKEN}",
				},
				Repositories: []string{
					"project1",
					"project2",
				},
			},
		},
		Targets: []TargetConfig{
			{
				Name:     "github-profile",
				Platform: "github",
				Auth: AuthConfig{
					Type:     "token",
					Username: "your-github-username",
					Token:    "${GITHUB_TOKEN}",
				},
				Mirror: MirrorConfig{
					Repository: "work-activity-mirror",
					Visibility: "private",
					Branch:     "main",
				},
			},
		},
		Sync: SyncConfig{
			Schedule:      "0 18 * * *", // 6 PM daily
			Timezone:      "local",
			CommitMessage: "Development work - {date}",
		},
	}

	// Get config directory
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(home, ".git-activity-mirror")
	configFile := filepath.Join(configDir, "config.yaml")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal configuration to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	// Write configuration file
	if err := os.WriteFile(configFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write configuration file: %w", err)
	}

	fmt.Printf("✅ Configuration file created: %s\n", configFile)
	fmt.Println()
	fmt.Println("📝 Next steps:")
	fmt.Println("1. Edit the configuration file with your platform details")
	fmt.Println("2. Set environment variables for your tokens:")
	fmt.Println("   export GITLAB_TOKEN=your_gitlab_token")
	fmt.Println("   export GITHUB_TOKEN=your_github_token")
	fmt.Println("3. Run 'git-activity-mirror import' to import historical commits")
	fmt.Println("4. Run 'git-activity-mirror sync' to start syncing")
	fmt.Println()

	return nil
}

// Configuration structures
type Config struct {
	Version int            `yaml:"version"`
	Sources []SourceConfig `yaml:"sources"`
	Targets []TargetConfig `yaml:"targets"`
	Sync    SyncConfig     `yaml:"sync"`
}

type SourceConfig struct {
	Name         string     `yaml:"name"`
	Platform     string     `yaml:"platform"`
	Host         string     `yaml:"host,omitempty"`
	Auth         AuthConfig `yaml:"auth"`
	Repositories []string   `yaml:"repositories,omitempty"`
}

type TargetConfig struct {
	Name     string       `yaml:"name"`
	Platform string       `yaml:"platform"`
	Host     string       `yaml:"host,omitempty"`
	Auth     AuthConfig   `yaml:"auth"`
	Mirror   MirrorConfig `yaml:"mirror"`
}

type AuthConfig struct {
	Type     string `yaml:"type"`
	Username string `yaml:"username,omitempty"`
	Token    string `yaml:"token,omitempty"`
	Password string `yaml:"password,omitempty"`
	SSHKey   string `yaml:"ssh_key,omitempty"`
}

type MirrorConfig struct {
	Repository string `yaml:"repository"`
	Visibility string `yaml:"visibility"`
	Branch     string `yaml:"branch,omitempty"`
}

type SyncConfig struct {
	Schedule      string `yaml:"schedule"`
	Timezone      string `yaml:"timezone"`
	CommitMessage string `yaml:"commit_message"`
}
