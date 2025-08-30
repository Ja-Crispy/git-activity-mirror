# git-activity-mirror

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Cross-Platform](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey.svg)]()

Mirror git commit activity between platforms while preserving code privacy.

## Overview

Cross-platform tool that synchronizes commit timestamps and metadata between git hosting platforms without exposing source code. Maintains accurate contribution graphs across GitHub, GitLab, and other platforms.

## Features

**Supported Platforms**
- GitHub (github.com + Enterprise)
- GitLab (gitlab.com + self-hosted)
- Bitbucket (planned)
- Azure DevOps (planned)

**Privacy Design**
- Mirrors timestamps and commit metadata only
- No source code access or exposure
- Uses unified generic repository names
- Preserves commit frequency and patterns

**Cross-Platform**
- Windows, macOS, Linux binaries
- Native system scheduling integration

## Installation

Download the binary for your platform from [releases](https://github.com/Ja-Crispy/git-activity-mirror/releases/latest).

```bash
# Linux/macOS
curl -L https://github.com/Ja-Crispy/git-activity-mirror/releases/latest/download/git-activity-mirror-linux-amd64 -o git-activity-mirror
chmod +x git-activity-mirror

# Or build from source
go install github.com/Ja-Crispy/git-activity-mirror@latest
```

## Usage

```bash
# Initialize configuration
git-activity-mirror init

# Set authentication tokens
export GITLAB_TOKEN=your_gitlab_token
export GITHUB_TOKEN=your_github_token

# Import historical commits
git-activity-mirror import --since=1y --dry-run
git-activity-mirror import --since=1y

# Sync recent activity
git-activity-mirror sync --since=24h
```

## Configuration
```yaml
# ~/.git-activity-mirror/config.yaml
sources:
  - name: work
    platform: gitlab
    auth:
      type: token
      username: your-username
      token: ${GITLAB_TOKEN}
    repositories:
      - repo1
      - repo2

targets:
  - name: github
    platform: github
    auth:
      type: token
      username: your-github-username
      token: ${GITHUB_TOKEN}
    mirror:
      repository: work-activity-mirror
      visibility: private

sync:
  schedule: "0 18 * * *"
  commit_message: "Development work - {date}"
```

## Commands

| Command | Description |
|---------|-------------|
| `init` | Initialize configuration |
| `sync` | Synchronize recent commits |
| `import` | Import historical commits |
| `status` | Show sync status |
| `config` | Manage configuration |

## Architecture

Built in Go with platform-agnostic interfaces. Supports authentication via tokens and environment variables. Uses unified mirror repositories to preserve privacy while maintaining accurate contribution patterns.

## License

MIT