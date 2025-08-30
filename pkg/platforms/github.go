package platforms

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v56/github"
	"golang.org/x/oauth2"
)

// GitHubPlatform implements GitPlatform for GitHub
type GitHubPlatform struct {
	client *github.Client
	ctx    context.Context
	config PlatformConfig
	owner  string
}

// NewGitHubPlatform creates a new GitHub platform instance
func NewGitHubPlatform(config PlatformConfig) (*GitHubPlatform, error) {
	ctx := context.Background()
	
	// Set up authentication
	var client *github.Client
	if config.Auth.Token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: config.Auth.Token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	} else {
		client = github.NewClient(nil) // Public access only
	}

	// For GitHub Enterprise, set base URL
	if config.Host != "" && config.Host != "github.com" {
		var err error
		client, err = client.WithEnterpriseURLs(
			fmt.Sprintf("https://%s/api/v3/", config.Host),
			fmt.Sprintf("https://%s/api/uploads/", config.Host),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to configure GitHub Enterprise: %w", err)
		}
	}

	return &GitHubPlatform{
		client: client,
		ctx:    ctx,
		config: config,
		owner:  config.Auth.Username,
	}, nil
}

// Connect establishes connection to GitHub
func (g *GitHubPlatform) Connect(config AuthConfig) error {
	g.config.Auth = config
	g.owner = config.Username
	
	// Recreate client with new config
	if config.Token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: config.Token},
		)
		tc := oauth2.NewClient(g.ctx, ts)
		g.client = github.NewClient(tc)
	}
	
	return nil
}

// ValidateCredentials validates the GitHub credentials
func (g *GitHubPlatform) ValidateCredentials() error {
	_, _, err := g.client.Users.Get(g.ctx, "")
	if err != nil {
		return fmt.Errorf("invalid GitHub credentials: %w", err)
	}
	return nil
}

// Disconnect closes any connections (no-op for GitHub API)
func (g *GitHubPlatform) Disconnect() error {
	return nil
}

// ListRepositories returns all repositories for the authenticated user
func (g *GitHubPlatform) ListRepositories() ([]Repository, error) {
	var allRepos []Repository
	
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	
	for {
		repos, resp, err := g.client.Repositories.List(g.ctx, g.owner, opt)
		if err != nil {
			return nil, fmt.Errorf("failed to list repositories: %w", err)
		}
		
		for _, repo := range repos {
			allRepos = append(allRepos, Repository{
				ID:          fmt.Sprintf("%d", repo.GetID()),
				Name:        repo.GetName(),
				FullName:    repo.GetFullName(),
				Description: repo.GetDescription(),
				URL:         repo.GetHTMLURL(),
				CloneURL:    repo.GetCloneURL(),
				Private:     repo.GetPrivate(),
				CreatedAt:   repo.GetCreatedAt().Time,
				UpdatedAt:   repo.GetUpdatedAt().Time,
				Platform:    "github",
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
func (g *GitHubPlatform) GetCommits(repo Repository, since time.Time) ([]Commit, error) {
	var allCommits []Commit
	
	// Parse owner/repo from full name
	parts := strings.Split(repo.FullName, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid repository full name: %s", repo.FullName)
	}
	owner, repoName := parts[0], parts[1]
	
	opt := &github.CommitsListOptions{
		Since: since,
		ListOptions: github.ListOptions{PerPage: 100},
	}
	
	// Filter by author email if specified
	if g.config.Auth.Username != "" {
		opt.Author = g.config.Auth.Username
	}
	
	for {
		commits, resp, err := g.client.Repositories.ListCommits(g.ctx, owner, repoName, opt)
		if err != nil {
			return nil, fmt.Errorf("failed to get commits: %w", err)
		}
		
		for _, commit := range commits {
			if commit.Commit == nil || commit.Commit.Author == nil {
				continue
			}
			
			allCommits = append(allCommits, Commit{
				SHA:     commit.GetSHA(),
				Message: commit.Commit.GetMessage(),
				Author: Author{
					Name:  commit.Commit.Author.GetName(),
					Email: commit.Commit.Author.GetEmail(),
				},
				Committer: Author{
					Name:  commit.Commit.Committer.GetName(),
					Email: commit.Commit.Committer.GetEmail(),
				},
				Date:     commit.Commit.Author.GetDate().Time,
				URL:      commit.GetHTMLURL(),
				Repo:     repo.FullName,
				Platform: "github",
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
func (g *GitHubPlatform) GetCommitCount(repo Repository, since time.Time) (int, error) {
	commits, err := g.GetCommits(repo, since)
	if err != nil {
		return 0, err
	}
	return len(commits), nil
}

// InitializeMirror creates a new repository for mirroring
func (g *GitHubPlatform) InitializeMirror(name string, visibility string) error {
	repo := &github.Repository{
		Name:        github.String(name),
		Description: github.String("Mirror of git activity from other platforms"),
		Private:     github.Bool(visibility == "private"),
	}
	
	_, _, err := g.client.Repositories.Create(g.ctx, "", repo)
	if err != nil {
		// Check if repo already exists
		if strings.Contains(err.Error(), "already exists") {
			return nil // Repo already exists, that's okay
		}
		return fmt.Errorf("failed to create mirror repository: %w", err)
	}
	
	return nil
}

// MirrorCommits creates mirror commits with preserved timestamps
func (g *GitHubPlatform) MirrorCommits(commits []Commit) error {
	if len(commits) == 0 {
		return nil
	}
	
	// Get the mirror repository
	mirrorRepo := g.config.Mirror.Repository
	parts := strings.Split(mirrorRepo, "/")
	var owner, repoName string
	
	if len(parts) == 2 {
		owner, repoName = parts[0], parts[1]
	} else {
		owner = g.owner
		repoName = mirrorRepo
	}
	
	// For each commit, create an empty commit with preserved timestamp
	for _, commit := range commits {
		// Create a simple commit message
		message := fmt.Sprintf("Development work - %s", commit.Date.Format("2006-01-02"))
		
		// Get the current HEAD to create a new commit
		ref, _, err := g.client.Git.GetRef(g.ctx, owner, repoName, "refs/heads/main")
		if err != nil {
			// Try master branch if main doesn't exist
			ref, _, err = g.client.Git.GetRef(g.ctx, owner, repoName, "refs/heads/master")
			if err != nil {
				return fmt.Errorf("failed to get repository head: %w", err)
			}
		}
		
		// Create a tree (empty change)
		baseTree, _, err := g.client.Git.GetTree(g.ctx, owner, repoName, ref.Object.GetSHA(), false)
		if err != nil {
			return fmt.Errorf("failed to get base tree: %w", err)
		}
		
		// Create commit parameters
		author := &github.CommitAuthor{
			Date:  &github.Timestamp{Time: commit.Date},
			Name:  github.String(g.config.Auth.Username),
			Email: github.String(g.config.Auth.Username + "@users.noreply.github.com"),
		}
		
		committer := &github.CommitAuthor{
			Date:  &github.Timestamp{Time: commit.Date},
			Name:  github.String(g.config.Auth.Username),
			Email: github.String(g.config.Auth.Username + "@users.noreply.github.com"),
		}
		
		parents := []*github.Commit{
			{SHA: ref.Object.SHA},
		}
		
		tree := &github.Tree{SHA: baseTree.SHA}
		
		// Create the commit
		createdCommit, _, err := g.client.Git.CreateCommit(g.ctx, owner, repoName, &github.Commit{
			Message: github.String(message),
			Tree:    tree,
			Parents: parents,
			Author:  author,
			Committer: committer,
		})
		if err != nil {
			return fmt.Errorf("failed to create mirror commit: %w", err)
		}
		
		// Update the reference
		ref.Object.SHA = createdCommit.SHA
		_, _, err = g.client.Git.UpdateRef(g.ctx, owner, repoName, ref, false)
		if err != nil {
			return fmt.Errorf("failed to update ref: %w", err)
		}
	}
	
	return nil
}

// GetMirrorStatus returns the status of the mirror repository
func (g *GitHubPlatform) GetMirrorStatus() (MirrorStatus, error) {
	mirrorRepo := g.config.Mirror.Repository
	parts := strings.Split(mirrorRepo, "/")
	var owner, repoName string
	
	if len(parts) == 2 {
		owner, repoName = parts[0], parts[1]
	} else {
		owner = g.owner
		repoName = mirrorRepo
	}
	
	// Get repository info
	repo, _, err := g.client.Repositories.Get(g.ctx, owner, repoName)
	if err != nil {
		return MirrorStatus{}, fmt.Errorf("failed to get mirror repository: %w", err)
	}
	
	// Get latest commit
	commits, _, err := g.client.Repositories.ListCommits(g.ctx, owner, repoName, &github.CommitsListOptions{
		ListOptions: github.ListOptions{PerPage: 1},
	})
	if err != nil {
		return MirrorStatus{}, fmt.Errorf("failed to get latest commit: %w", err)
	}
	
	status := MirrorStatus{
		Repository:   repo.GetFullName(),
		LastSync:     time.Now(), // This would need to be tracked separately
		TotalCommits: 0,          // Would need to count commits
		Status:       "active",
	}
	
	if len(commits) > 0 {
		status.LastCommitSHA = commits[0].GetSHA()
	}
	
	return status, nil
}

// GetPlatformName returns the platform name
func (g *GitHubPlatform) GetPlatformName() string {
	if g.config.Host != "" && g.config.Host != "github.com" {
		return fmt.Sprintf("GitHub Enterprise (%s)", g.config.Host)
	}
	return "GitHub"
}

// GetPlatformType returns the platform type
func (g *GitHubPlatform) GetPlatformType() PlatformType {
	return PlatformGitHub
}

// SupportsWebhooks returns true if the platform supports webhooks
func (g *GitHubPlatform) SupportsWebhooks() bool {
	return true
}