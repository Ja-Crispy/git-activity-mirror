---
name: Bug Report
about: Create a report to help us improve
title: '[BUG] '
labels: 'bug'
assignees: ''
---

## 🐛 Bug Description
A clear and concise description of what the bug is.

## 📋 System Information
- **OS**: [e.g., Windows 11, macOS 14.1, Ubuntu 22.04]
- **Architecture**: [e.g., x64, ARM64]
- **Go Version**: [e.g., 1.21.4]
- **git-activity-mirror Version**: [e.g., v0.1.0]

## 🔄 Steps to Reproduce
1. Run command: `git-activity-mirror ...`
2. Configure with: `...`
3. See error: `...`

## ✅ Expected Behavior
A clear description of what you expected to happen.

## ❌ Actual Behavior
A clear description of what actually happened.

## 📝 Configuration (Sanitized)
```yaml
# Remove all tokens and sensitive information
sources:
  - name: work
    platform: gitlab
    repositories:
      - repo1
```

## 🔍 Error Output
```
Paste the full error message or output here.
Use verbose mode: git-activity-mirror --verbose [command]
```

## 🧪 Workaround
If you found a temporary workaround, please describe it.

## 📎 Additional Context
- Screenshots (if applicable)
- Related issues or PRs
- Any other relevant information

## ✅ Checklist
- [ ] I have searched existing issues for this bug
- [ ] I have tested with the latest version
- [ ] I have provided a complete configuration (sanitized)
- [ ] I have included the full error output