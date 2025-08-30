package platforms

import (
	"fmt"
	"time"
)

// GitPlatform defines the interface that all git hosting platforms must implement
// This enables platform-agnostic mirroring between any git hosts
type GitPlatform interface {
	// Authentication and connection
	Connect(config AuthConfig) error
	ValidateCredentials() error
	Disconnect() error

	// Source operations (reading commits from source platform)
	ListRepositories() ([]Repository, error)
	GetCommits(repo Repository, since time.Time) ([]Commit, error)
	GetCommitCount(repo Repository, since time.Time) (int, error)

	// Target operations (writing mirror commits to target platform)
	InitializeMirror(name string, visibility string) error
	MirrorCommits(commits []Commit) error
	GetMirrorStatus() (MirrorStatus, error)

	// Platform information
	GetPlatformName() string
	GetPlatformType() PlatformType
	SupportsWebhooks() bool
}

// PlatformType represents different git hosting platform types
type PlatformType string

const (
	PlatformGitHub     PlatformType = "github"
	PlatformGitLab     PlatformType = "gitlab"
	PlatformBitbucket  PlatformType = "bitbucket"
	PlatformAzureDevOps PlatformType = "azuredevops"
	PlatformGenericGit PlatformType = "generic"
)

// AuthConfig holds authentication configuration for a platform
type AuthConfig struct {
	Type     AuthType               `yaml:"type"`
	Token    string                 `yaml:"token,omitempty"`
	Username string                 `yaml:"username,omitempty"`
	Password string                 `yaml:"password,omitempty"`
	SSHKey   string                 `yaml:"ssh_key,omitempty"`
	Host     string                 `yaml:"host,omitempty"`
	Extra    map[string]interface{} `yaml:"extra,omitempty"`
}

// AuthType represents different authentication methods
type AuthType string

const (
	AuthToken    AuthType = "token"
	AuthPassword AuthType = "password"
	AuthSSH      AuthType = "ssh"
	AuthOAuth    AuthType = "oauth"
)

// Repository represents a git repository
type Repository struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	FullName    string    `json:"full_name"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	CloneURL    string    `json:"clone_url"`
	Private     bool      `json:"private"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Platform    string    `json:"platform"`
}

// Commit represents a git commit with metadata
type Commit struct {
	SHA       string    `json:"sha"`
	Message   string    `json:"message"`
	Author    Author    `json:"author"`
	Committer Author    `json:"committer"`
	Date      time.Time `json:"date"`
	URL       string    `json:"url"`
	Repo      string    `json:"repo"`
	Platform  string    `json:"platform"`
}

// Author represents a commit author/committer
type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// MirrorStatus represents the status of a mirror repository
type MirrorStatus struct {
	Repository    string    `json:"repository"`
	LastSync      time.Time `json:"last_sync"`
	TotalCommits  int       `json:"total_commits"`
	LastCommitSHA string    `json:"last_commit_sha"`
	Status        string    `json:"status"`
	Error         string    `json:"error,omitempty"`
}

// PlatformConfig holds platform-specific configuration
type PlatformConfig struct {
	Name      string                 `yaml:"name"`
	Platform  PlatformType           `yaml:"platform"`
	Host      string                 `yaml:"host,omitempty"`
	Auth      AuthConfig             `yaml:"auth"`
	Repos     []string               `yaml:"repositories,omitempty"`
	Mirror    MirrorConfig           `yaml:"mirror,omitempty"`
	Extra     map[string]interface{} `yaml:"extra,omitempty"`
}

// MirrorConfig holds mirror-specific configuration
type MirrorConfig struct {
	Repository string `yaml:"repository"`      // Name for mirror repo (e.g., "work-activity-mirror")
	Visibility string `yaml:"visibility"`      // public, private
	Branch     string `yaml:"branch,omitempty"` // Target branch (default: main)
	Strategy   string `yaml:"strategy,omitempty"` // unified, separate, hashed
}

// NewPlatform creates a new platform instance based on the platform type
func NewPlatform(platformType PlatformType, config PlatformConfig) (GitPlatform, error) {
	switch platformType {
	case PlatformGitHub:
		return NewGitHubPlatform(config)
	case PlatformGitLab:
		return NewGitLabPlatform(config)
	case PlatformBitbucket, PlatformAzureDevOps, PlatformGenericGit:
		return nil, fmt.Errorf("platform %s not implemented yet - coming in v0.2.0", platformType)
	default:
		return nil, ErrUnsupportedPlatform
	}
}

// Common errors
var (
	ErrUnsupportedPlatform = fmt.Errorf("unsupported platform")
	ErrInvalidAuth        = fmt.Errorf("invalid authentication")
	ErrRepositoryNotFound = fmt.Errorf("repository not found")
	ErrPermissionDenied   = fmt.Errorf("permission denied")
	ErrRateLimit          = fmt.Errorf("rate limit exceeded")
)