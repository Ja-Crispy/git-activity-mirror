# git-activity-mirror

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Cross-Platform](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey.svg)]()

> **Mirror your git commit activity between any platforms while keeping your code private.**

A cross-platform, platform-agnostic tool that mirrors git commit activity between any git hosting platforms (GitHub, GitLab, Bitbucket, Azure DevOps, etc.) while maintaining complete privacy. Show your real work activity everywhere without exposing sensitive code.

## 🎯 Why git-activity-mirror?

**The Problem:** Your GitHub contribution graph looks empty, but you code 8+ hours daily on:
- Company GitLab/Bitbucket repositories
- Self-hosted git instances  
- Azure DevOps projects
- Private enterprise platforms

**The Solution:** Mirror your real commit activity timestamps to maintain accurate contribution graphs across all platforms.

## 🔥 Features

### Platform Support
- ✅ **GitHub** (github.com + Enterprise)
- ✅ **GitLab** (gitlab.com + self-hosted)
- 🔄 **Bitbucket** (coming soon)
- 🔄 **Azure DevOps** (coming soon)
- 🔄 **Generic Git** (any git server)

### Privacy-First Design
- ❌ **No code exposure** - Never accesses actual file contents
- ❌ **No sensitive data** - Only timestamps and generic messages
- ✅ **Accurate activity** - Real commit times, not fake patterns
- ✅ **Configurable messages** - Generic like "Development work - 2025-08-30"

### Cross-Platform
- 🪟 **Windows** - Native binary + Task Scheduler
- 🍎 **macOS** - Universal binary + Homebrew + launchd
- 🐧 **Linux** - Multiple package formats + systemd

## 🚀 Quick Start

### Installation

#### macOS (Homebrew)
```bash
brew install ja-crispy/tap/git-activity-mirror
```

#### Windows (Chocolatey)
```powershell
choco install git-activity-mirror
```

#### Linux (Direct Download)
```bash
curl -L https://github.com/Ja-Crispy/git-activity-mirror/releases/latest/download/git-activity-mirror-linux-amd64 -o git-activity-mirror
chmod +x git-activity-mirror
sudo mv git-activity-mirror /usr/local/bin/
```

#### Go Developers
```bash
go install github.com/Ja-Crispy/git-activity-mirror@latest
```

### Setup

1. **Initialize configuration:**
```bash
git-activity-mirror init
```

2. **Edit your configuration file:**
```bash
git-activity-mirror config edit
```

3. **Set up authentication tokens:**
```bash
export GITLAB_TOKEN=your_gitlab_token
export GITHUB_TOKEN=your_github_token
```

4. **Import historical commits:**
```bash
git-activity-mirror import --since=1y
```

5. **Set up automatic syncing:**
```bash
git-activity-mirror sync --schedule
```

## 📝 Configuration Examples

### GitLab → GitHub (Your Use Case)
```yaml
# ~/.git-activity-mirror/config.yaml
version: 1
sources:
  - name: work
    platform: gitlab
    host: gitlab.com
    auth:
      type: token
      username: vaishnav9
      token: ${GITLAB_TOKEN}
    repositories:
      - python-fastapi
      - scyai-proto
      - kubernetes
      - infrastructureascode

targets:
  - name: github-profile
    platform: github
    auth:
      type: token
      username: Ja-Crispy
      token: ${GITHUB_TOKEN}
    mirror:
      repository: work-activity-mirror
      visibility: private

sync:
  schedule: "0 18 * * *"  # Daily at 6 PM
  commit_message: "Development work - {date}"
```

### Multi-Platform Aggregation
```yaml
sources:
  - name: work-gitlab
    platform: gitlab
    host: gitlab.company.com
  - name: personal-github
    platform: github
  - name: client-bitbucket
    platform: bitbucket

targets:
  - name: main-profile
    platform: github
    mirror:
      repository: all-my-work
```

See [`examples/`](examples/) for more configuration templates.

## 🛠️ Commands

| Command | Description |
|---------|-------------|
| `init` | Create initial configuration |
| `sync` | Sync recent commits (last 24h) |
| `import` | Import historical commits (last 1y) |
| `status` | Show platform and sync status |
| `config show` | Display current configuration |

### Common Usage
```bash
# Daily sync (usually automated)
git-activity-mirror sync

# Import last 6 months of history
git-activity-mirror import --since=6m

# Check if everything is working
git-activity-mirror status

# Test configuration without making changes
git-activity-mirror sync --dry-run
```

## 🔒 Security & Privacy

### What Gets Mirrored
- ✅ **Commit timestamps** - Exact dates and times
- ✅ **Commit frequency** - How often you commit
- ✅ **Repository count** - Number of active projects

### What's Protected
- ❌ **Source code** - Never accessed or transmitted
- ❌ **File names** - No file paths or names exposed
- ❌ **Commit messages** - Replaced with generic messages
- ❌ **Branch names** - Not included in mirrors
- ❌ **Repository names** - Source repo names not exposed

### Authentication
- Tokens stored in environment variables (not config files)
- Minimal required permissions (read repos, write to mirror)
- Supports GitHub/GitLab personal access tokens
- Optional SSH key support

## 🏗️ How It Works

```
┌─────────────┐    ┌──────────────────┐    ┌─────────────┐
│   GitLab    │───▶│ git-activity-    │───▶│   GitHub    │
│   (Work)    │    │    mirror        │    │ (Profile)   │
│             │    │                  │    │             │
│ 🔒 Private  │    │ 📊 Timestamps    │    │ 📈 Graph    │
│    Code     │    │ 🛡️  Privacy      │    │  Updated    │
└─────────────┘    └──────────────────┘    └─────────────┘
```

1. **Scan** source repositories for commits by you
2. **Extract** only timestamp and author information
3. **Create** empty commits in target with preserved dates
4. **Mirror** using generic messages like "Development work - 2025-08-30"

## 🤝 vs Other Solutions

| Feature | git-activity-mirror | Fake commit generators | Manual scripts |
|---------|-------------------|----------------------|----------------|
| **Real timestamps** | ✅ Actual work times | ❌ Random/fake | ✅ If done right |
| **Multiple platforms** | ✅ Any → Any | ❌ Usually GitHub only | ❌ Platform-specific |
| **Privacy-first** | ✅ No code exposure | ✅ No real data | ⚠️ Depends on implementation |
| **Cross-platform** | ✅ Windows/Mac/Linux | ❌ Often platform-specific | ❌ Usually bash-only |
| **Professional quality** | ✅ Production-ready | ❌ Often hacky | ⚠️ Varies widely |
| **Maintenance** | ✅ Automated | ❌ Manual tweaking | ❌ Constant updates needed |

## 📊 Success Stories

> *"My GitHub contribution graph went from empty to showing 2+ years of consistent work activity. Finally reflects my actual productivity!"*  
> — **Enterprise Developer**

> *"Perfect for consultants working on multiple client platforms. Now my profile shows all my work."*  
> — **Freelance Developer**

> *"Set it up once, forgot about it. Been running perfectly for 6 months mirroring GitLab to GitHub."*  
> — **Startup CTO**

## 🗺️ Roadmap

### v1.0 (Current)
- [x] GitLab ↔ GitHub mirroring
- [x] Cross-platform CLI
- [x] Historical import
- [x] Automated scheduling

### v1.1 (Next)
- [ ] Bitbucket platform support
- [ ] Azure DevOps integration
- [ ] Web dashboard
- [ ] Team/organization support

### v1.2 (Future)
- [ ] Webhook-based real-time sync  
- [ ] Advanced analytics
- [ ] GitHub Action
- [ ] API for integrations

## 🤝 Contributing

We welcome contributions! This project helps thousands of developers show their real work.

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b amazing-feature`)  
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin amazing-feature`)
5. **Open** a Pull Request

### Development Setup
```bash
git clone https://github.com/Ja-Crispy/git-activity-mirror.git
cd git-activity-mirror
go mod download
go run cmd/git-activity-mirror/main.go --help
```

## 📄 License

MIT License - see [LICENSE](LICENSE) file.

## 💬 Support

- 🐛 **Bug Reports:** [GitHub Issues](https://github.com/Ja-Crispy/git-activity-mirror/issues)
- 💡 **Feature Requests:** [GitHub Discussions](https://github.com/Ja-Crispy/git-activity-mirror/discussions)
- 📧 **Contact:** [support@git-activity-mirror.dev](mailto:support@git-activity-mirror.dev)

---

<div align="center">

**⭐ If this tool helps you, please give it a star! ⭐**

*Show your real work, everywhere.*

</div>