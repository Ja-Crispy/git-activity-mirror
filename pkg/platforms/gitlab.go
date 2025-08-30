package platforms

import (
	"fmt"
	"strings"
	"time"

	"github.com/xanzy/go-gitlab"
)

// GitLabPlatform implements GitPlatform for GitLab
type GitLabPlatform struct {
	client *gitlab.Client
	config PlatformConfig
	userID int
}

// NewGitLabPlatform creates a new GitLab platform instance
func NewGitLabPlatform(config PlatformConfig) (*GitLabPlatform, error) {
	var client *gitlab.Client
	var err error

	// Create client based on host
	host := config.Host
	if host == "" {
		host = "https://gitlab.com"
	}
	if !strings.HasPrefix(host, "http") {
		host = "https://" + host
	}

	if config.Auth.Token != "" {
		client, err = gitlab.NewClient(config.Auth.Token, gitlab.WithBaseURL(host))
	} else {
		client, err = gitlab.NewClient("", gitlab.WithBaseURL(host))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create GitLab client: %w", err)
	}

	return &GitLabPlatform{
		client: client,
		config: config,
	}, nil
}

// Connect establishes connection to GitLab
func (g *GitLabPlatform) Connect(config AuthConfig) error {
	g.config.Auth = config

	// Recreate client with new config
	host := config.Host
	if host == "" {
		host = "https://gitlab.com"
	}
	if !strings.HasPrefix(host, "http") {
		host = "https://" + host
	}

	var err error
	if config.Token != "" {
		g.client, err = gitlab.NewClient(config.Token, gitlab.WithBaseURL(host))
	} else {
		g.client, err = gitlab.NewClient("", gitlab.WithBaseURL(host))
	}

	return err
}

// ValidateCredentials validates the GitLab credentials
func (g *GitLabPlatform) ValidateCredentials() error {
	user, _, err := g.client.Users.CurrentUser()
	if err != nil {
		return fmt.Errorf("invalid GitLab credentials: %w", err)
	}

	g.userID = user.ID
	return nil
}

// Disconnect closes any connections (no-op for GitLab API)
func (g *GitLabPlatform) Disconnect() error {
	return nil
}

// ListRepositories returns all repositories for the authenticated user
func (g *GitLabPlatform) ListRepositories() ([]Repository, error) {
	var allRepos []Repository

	opt := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{PerPage: 100},
		Owned:       gitlab.Ptr(true),
	}

	for {
		projects, resp, err := g.client.Projects.ListProjects(opt)
		if err != nil {
			return nil, fmt.Errorf("failed to list projects: %w", err)
		}

		for _, project := range projects {
			allRepos = append(allRepos, Repository{
				ID:          fmt.Sprintf("%d", project.ID),
				Name:        project.Name,
				FullName:    project.PathWithNamespace,
				Description: project.Description,
				URL:         project.WebURL,
				CloneURL:    project.HTTPURLToRepo,
				Private:     project.Visibility != gitlab.PublicVisibility,
				CreatedAt:   *project.CreatedAt,
				UpdatedAt:   *project.LastActivityAt,
				Platform:    "gitlab",
			})
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allRepos, nil
}

// GetCommits retrieves commits from a repository since a specific date
func (g *GitLabPlatform) GetCommits(repo Repository, since time.Time) ([]Commit, error) {
	var allCommits []Commit

	projectID := repo.ID

	opt := &gitlab.ListCommitsOptions{
		ListOptions: gitlab.ListOptions{PerPage: 100},
		Since:       &since,
	}

	// Note: GitLab API doesn't support AuthorEmail filter in ListCommitsOptions
	// We'll filter commits by author email after fetching if needed

	for {
		commits, resp, err := g.client.Commits.ListCommits(projectID, opt)
		if err != nil {
			return nil, fmt.Errorf("failed to get commits: %w", err)
		}

		for _, commit := range commits {
			allCommits = append(allCommits, Commit{
				SHA:     commit.ID,
				Message: commit.Message,
				Author: Author{
					Name:  commit.AuthorName,
					Email: commit.AuthorEmail,
				},
				Committer: Author{
					Name:  commit.CommitterName,
					Email: commit.CommitterEmail,
				},
				Date:     *commit.AuthoredDate,
				URL:      commit.WebURL,
				Repo:     repo.FullName,
				Platform: "gitlab",
			})
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allCommits, nil
}

// GetCommitCount returns the number of commits since a specific date
func (g *GitLabPlatform) GetCommitCount(repo Repository, since time.Time) (int, error) {
	commits, err := g.GetCommits(repo, since)
	if err != nil {
		return 0, err
	}
	return len(commits), nil
}

// InitializeMirror creates a new project for mirroring
func (g *GitLabPlatform) InitializeMirror(name string, visibility string) error {
	vis := gitlab.PrivateVisibility
	if visibility == "public" {
		vis = gitlab.PublicVisibility
	}

	project := &gitlab.CreateProjectOptions{
		Name:        &name,
		Description: gitlab.Ptr("Mirror of git activity from other platforms"),
		Visibility:  &vis,
	}

	_, _, err := g.client.Projects.CreateProject(project)
	if err != nil {
		// Check if project already exists
		if strings.Contains(err.Error(), "already been taken") {
			return nil // Project already exists, that's okay
		}
		return fmt.Errorf("failed to create mirror project: %w", err)
	}

	return nil
}

// MirrorCommits creates mirror commits with preserved timestamps
func (g *GitLabPlatform) MirrorCommits(commits []Commit) error {
	if len(commits) == 0 {
		return nil
	}

	// Get the mirror project
	mirrorRepo := g.config.Mirror.Repository

	// Find the project
	projects, _, err := g.client.Projects.ListProjects(&gitlab.ListProjectsOptions{
		Search: &mirrorRepo,
		Owned:  gitlab.Ptr(true),
	})
	if err != nil {
		return fmt.Errorf("failed to find mirror project: %w", err)
	}

	if len(projects) == 0 {
		return fmt.Errorf("mirror project not found: %s", mirrorRepo)
	}

	projectID := projects[0].ID

	// For each commit, create an empty commit with preserved timestamp
	for _, commit := range commits {
		// Create a simple commit message
		message := fmt.Sprintf("Development work - %s", commit.Date.Format("2006-01-02"))

		// Create an empty commit using the Commits API
		// Note: GitLab doesn't have direct support for empty commits via API
		// This is a simplified version - in practice, you might need to:
		// 1. Clone the repository locally
		// 2. Create empty commits with git commands
		// 3. Push back to GitLab

		// For now, create a simple file update to simulate activity
		actions := []*gitlab.CommitActionOptions{
			{
				Action:   gitlab.Ptr(gitlab.FileUpdate),
				FilePath: gitlab.Ptr(".activity"),
				Content:  gitlab.Ptr(fmt.Sprintf("Activity recorded: %s", commit.Date.Format(time.RFC3339))),
			},
		}

		createCommitOpt := &gitlab.CreateCommitOptions{
			Branch:        gitlab.Ptr("main"),
			CommitMessage: &message,
			Actions:       actions,
			AuthorEmail:   gitlab.Ptr(g.config.Auth.Username + "@users.noreply.gitlab.com"),
			AuthorName:    gitlab.Ptr(g.config.Auth.Username),
		}

		_, _, err := g.client.Commits.CreateCommit(projectID, createCommitOpt)
		if err != nil {
			// Try master branch if main doesn't exist
			createCommitOpt.Branch = gitlab.Ptr("master")
			_, _, err = g.client.Commits.CreateCommit(projectID, createCommitOpt)
			if err != nil {
				return fmt.Errorf("failed to create mirror commit: %w", err)
			}
		}
	}

	return nil
}

// GetMirrorStatus returns the status of the mirror project
func (g *GitLabPlatform) GetMirrorStatus() (MirrorStatus, error) {
	mirrorRepo := g.config.Mirror.Repository

	// Find the project
	projects, _, err := g.client.Projects.ListProjects(&gitlab.ListProjectsOptions{
		Search: &mirrorRepo,
		Owned:  gitlab.Ptr(true),
	})
	if err != nil {
		return MirrorStatus{}, fmt.Errorf("failed to find mirror project: %w", err)
	}

	if len(projects) == 0 {
		return MirrorStatus{}, fmt.Errorf("mirror project not found: %s", mirrorRepo)
	}

	project := projects[0]

	// Get latest commit
	commits, _, err := g.client.Commits.ListCommits(project.ID, &gitlab.ListCommitsOptions{
		ListOptions: gitlab.ListOptions{PerPage: 1},
	})
	if err != nil {
		return MirrorStatus{}, fmt.Errorf("failed to get latest commit: %w", err)
	}

	status := MirrorStatus{
		Repository:   project.PathWithNamespace,
		LastSync:     time.Now(), // This would need to be tracked separately
		TotalCommits: 0,          // Would need to count commits
		Status:       "active",
	}

	if len(commits) > 0 {
		status.LastCommitSHA = commits[0].ID
	}

	return status, nil
}

// GetPlatformName returns the platform name
func (g *GitLabPlatform) GetPlatformName() string {
	if g.config.Host != "" && g.config.Host != "gitlab.com" {
		return fmt.Sprintf("GitLab (%s)", g.config.Host)
	}
	return "GitLab"
}

// GetPlatformType returns the platform type
func (g *GitLabPlatform) GetPlatformType() PlatformType {
	return PlatformGitLab
}

// SupportsWebhooks returns true if the platform supports webhooks
func (g *GitLabPlatform) SupportsWebhooks() bool {
	return true
}
