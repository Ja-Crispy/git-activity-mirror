## 📝 What This PR Does
<!-- Brief description of the changes in this PR -->

## 🎯 Type of Change
<!-- Mark the relevant option with an "x" -->
- [ ] 🐛 Bug fix (non-breaking change that fixes an issue)
- [ ] ✨ New feature (non-breaking change that adds functionality)
- [ ] 💥 Breaking change (fix or feature that would cause existing functionality to change)
- [ ] 📚 Documentation update
- [ ] 🔧 Code refactoring (no functional changes)
- [ ] 🧪 Test improvements
- [ ] 🔒 Security fix

## 🧪 Testing Done
<!-- Describe how you tested your changes -->
- [ ] Unit tests pass locally (`go test ./...`)
- [ ] Code formatted (`gofmt -s -l .`)
- [ ] Linting passes (`go vet ./...`)
- [ ] Manual testing completed
- [ ] Integration testing (if applicable)

### Manual Testing Details
<!-- If you did manual testing, describe what you tested -->
```bash
# Example commands you ran
git-activity-mirror sync --dry-run
git-activity-mirror import --since=1y --verbose
```

## 🔧 Platform Testing (if applicable)
<!-- For platform integrations, confirm you tested with real accounts -->
- [ ] GitHub: Tested with personal/test account
- [ ] GitLab: Tested with personal/test account  
- [ ] Bitbucket: Tested with personal/test account
- [ ] Azure DevOps: Tested with personal/test account
- [ ] Other: ___________

## 💥 Breaking Changes
<!-- List any breaking changes and migration instructions -->
- None

OR

- **Configuration format changed**: Users need to update their `config.yaml`
- **CLI flag renamed**: `--old-flag` is now `--new-flag`
- **API method signature changed**: Update custom platform implementations

## 📚 Documentation
<!-- Check all that apply -->
- [ ] README.md updated
- [ ] CONTRIBUTING.md updated (if process changed)
- [ ] Help text/CLI documentation updated
- [ ] Code comments added/updated
- [ ] Configuration examples provided

## 🔗 Related Issues
<!-- Link to related issues -->
- Closes #123
- Relates to #456
- Fixes #789

## 📋 Checklist
<!-- Check all items before requesting review -->
- [ ] My code follows the project's style guidelines
- [ ] I have performed a self-review of my code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
- [ ] Any dependent changes have been merged and published

## 🖼️ Screenshots (if applicable)
<!-- Add screenshots to help reviewers understand UI/CLI changes -->

## 📎 Additional Context
<!-- Add any other context about the PR here -->