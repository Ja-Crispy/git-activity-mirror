# Contributing

## Development Setup

### Prerequisites
- Go 1.21 or higher
- Git
- Platform tokens for testing (GitHub, GitLab, etc.)

### Development Setup
```bash
# Clone the repository
git clone https://github.com/Ja-Crispy/git-activity-mirror.git
cd git-activity-mirror

# Install dependencies
go mod download

# Build the project
go build -o git-activity-mirror ./cmd/git-activity-mirror

# Run tests
go test ./...
```

## 🎯 Ways to Contribute

### 1. Platform Support
We're actively looking for platform integrations:
- **Bitbucket** (Cloud & Server)
- **Azure DevOps** (Azure Repos)
- **GitKraken Glo** 
- **Gitea/Forgejo**
- **SourceForge**
- **CodeCommit** (AWS)

### 2. Features & Enhancements
- Web dashboard for monitoring sync status
- Progress bars and better CLI UX
- Webhook support for real-time syncing
- Team/organization support
- Activity analytics and insights

### 3. Bug Reports & Fixes
- Cross-platform compatibility issues
- Authentication edge cases
- Rate limiting and API handling
- Configuration validation

## 📋 Code Style Guidelines

### Go Code Standards
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting (automatic in CI)
- Run `go vet` to check for issues
- Add tests for new functionality
- Keep functions focused and well-documented

### Commit Message Format
```
<type>(<scope>): <description>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix  
- `docs`: Documentation
- `style`: Code style/formatting
- `refactor`: Code refactoring
- `test`: Tests
- `chore`: Maintenance

**Examples:**
```
feat(github): add GitHub Enterprise support

Add support for GitHub Enterprise instances with custom hostnames.
Includes authentication handling and API endpoint configuration.

Closes #42
```

```
fix(gitlab): handle rate limiting correctly

Fix 429 responses by implementing exponential backoff retry logic
for GitLab API calls.

Fixes #38
```

## 🧪 Testing Requirements

### Unit Tests
All new code should include unit tests:
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Integration Tests
For platform integrations, include integration tests that:
- Test against real API endpoints (when possible)
- Use mock servers for reliable CI testing
- Cover authentication edge cases
- Test rate limiting behavior

### Manual Testing Checklist
Before submitting platform integrations:
- [ ] Authentication works with tokens/OAuth
- [ ] Repository listing works correctly
- [ ] Commit fetching handles pagination
- [ ] Rate limiting is handled gracefully
- [ ] Error messages are helpful
- [ ] Configuration validation works
- [ ] Dry-run mode works correctly

## 🔄 Pull Request Process

### Before Submitting
1. **Run the full test suite**: `go test ./...`
2. **Check formatting**: `gofmt -s -l .` (should return nothing)
3. **Run linting**: `go vet ./...`
4. **Test your changes**: Build and test manually
5. **Update documentation**: Update README, add examples

### PR Requirements
- [ ] **Clear description**: What does this PR do and why?
- [ ] **Tests included**: Unit tests for new functionality
- [ ] **Documentation updated**: README, examples, help text
- [ ] **No breaking changes**: Or clearly documented migration path
- [ ] **Platform tested**: For platform integrations, test with real accounts

### PR Template
When creating a PR, please include:

```markdown
## What This PR Does
Brief description of the change.

## Testing Done
- [ ] Unit tests pass
- [ ] Manual testing completed
- [ ] Integration tests (if applicable)

## Breaking Changes
None / List any breaking changes

## Related Issues
Closes #123, Relates to #456
```

## 🏗️ Architecture Guidelines

### Platform Implementation
When adding a new platform, implement the `GitPlatform` interface:
```go
type GitPlatform interface {
    Connect(config AuthConfig) error
    ValidateCredentials() error
    ListRepositories() ([]Repository, error)
    GetCommits(repo Repository, since time.Time) ([]Commit, error)
    InitializeMirror(name string, visibility string) error
    MirrorCommits(commits []Commit) error
    GetMirrorStatus() (MirrorStatus, error)
    // ... more methods
}
```

### Configuration
- Use YAML for configuration files
- Support environment variable substitution
- Validate configuration on load
- Provide helpful error messages

### Privacy & Security
- **Never log tokens or credentials**
- Use generic repository names for mirrors
- Don't expose source project information
- Implement proper rate limiting
- Handle API errors gracefully

## 🐛 Bug Report Guidelines

When reporting bugs, please include:

### System Information
- OS and version (Windows 11, macOS 14.1, Ubuntu 22.04)
- Go version used for building
- git-activity-mirror version

### Configuration
Provide a **sanitized** version of your config (remove tokens):
```yaml
sources:
  - name: work
    platform: gitlab
    host: gitlab.com
    repositories:
      - repo1
      - repo2
```

### Steps to Reproduce
1. Clear, numbered steps
2. Expected behavior
3. Actual behavior
4. Error messages (full stack trace if possible)

### Logs
Run with verbose logging:
```bash
git-activity-mirror sync --verbose --dry-run 2>&1 | tee debug.log
```

## 💡 Feature Request Guidelines

For feature requests, please provide:

### Use Case
- What problem does this solve?
- Who would benefit from this feature?
- How would you use it?

### Proposed Solution
- High-level description of the feature
- UI/CLI mockups if applicable
- Configuration examples

### Alternatives Considered
- Other ways to solve this problem
- Why this approach is preferred

## 📞 Getting Help

- **GitHub Issues**: For bugs and feature requests
- **GitHub Discussions**: For questions and general discussion
- **Discord** (coming soon): For real-time chat with contributors

## 🎉 Recognition

Contributors are recognized in:
- `CONTRIBUTORS.md` file
- Release notes for their contributions
- GitHub contributor statistics

Thank you for helping make git-activity-mirror better for developers everywhere! 🚀