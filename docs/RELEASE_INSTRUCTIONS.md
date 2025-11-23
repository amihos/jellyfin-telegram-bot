# Release v1.0.0 Instructions

This document provides detailed instructions for creating the first public release (v1.0.0) of the jellyfin-telegram-bot.

---

## Prerequisites

Before creating the release, ensure:

- [x] All Phase 8 tasks completed
- [x] Pre-launch checklist verified
- [x] All tests passing (`go test ./... -count=1`)
- [x] Docker build successful
- [x] No uncommitted changes (`git status` clean)
- [x] Repository pushed to GitHub
- [x] GitHub repository configured (see `GITHUB_REPOSITORY_SETUP.md`)
- [x] CI/CD workflows passing on main branch

---

## Release Process

### Step 1: Final Pre-Release Verification

Run these commands to verify everything is ready:

```bash
# Navigate to project root
cd /home/huso/jellyfin-telegram-bot

# Verify tests pass
go test ./... -count=1

# Expected output: All tests should PASS
# If any tests fail, DO NOT proceed with release

# Verify Docker build
docker build -t jellyfin-telegram-bot:v1.0.0-test .

# Expected output: Successfully built and tagged
# Image size should be around 34MB

# Test Docker image
docker run --rm jellyfin-telegram-bot:v1.0.0-test --help

# Expected output: Bot help message or version info

# Clean up test image
docker rmi jellyfin-telegram-bot:v1.0.0-test

# Verify no uncommitted changes
git status

# Expected output: "nothing to commit, working tree clean"
```

**If any of these fail, STOP and fix the issues before proceeding.**

---

### Step 2: Review Release Content

Verify these files are ready for release:

```bash
# Check critical files exist
ls -l README.md
ls -l LICENSE
ls -l CONTRIBUTING.md
ls -l CODE_OF_CONDUCT.md
ls -l .env.example
ls -l Dockerfile
ls -l docker-compose.example.yml
ls -l .goreleaser.yaml

# Check GitHub workflows exist
ls -l .github/workflows/test.yml
ls -l .github/workflows/build.yml
ls -l .github/workflows/release.yml
ls -l .github/workflows/docker.yml

# Check documentation exists
ls -l docs/

# Check locales exist
ls -l locales/active.en.toml
ls -l locales/active.fa.toml
```

**All files should be present. If any are missing, STOP and create them.**

---

### Step 3: Create Release Tag

**Important:** This step will trigger the automated release process. Make sure you're ready!

```bash
# Ensure you're on the main branch
git checkout main

# Pull latest changes (if working with a team)
git pull origin main

# Create annotated tag for v1.0.0
git tag -a v1.0.0 -m "Release v1.0.0 - First public release

## Features

### Internationalization (i18n)
- Multi-language support (English and Persian)
- Auto-detection of user language from Telegram settings
- Manual language selection via /language command
- Complete translations for all UI strings and notifications
- Fallback chain: saved preference → Telegram language → English

### Notifications
- Webhook-based notifications from Jellyfin media server
- Support for movies and TV episodes
- Series muting/unmuting for personalized notifications
- Rich notification formatting with metadata
- Poster images included when available
- Inline keyboard buttons for quick actions

### User Interface
- Interactive inline keyboard menus
- Navigation buttons for common actions
- Mute/unmute buttons on episode notifications
- Multi-language menu system
- Help and command reference

### Commands
- /start - Subscribe to notifications
- /recent - View recent content
- /search <query> - Search for content
- /mutedlist - View muted series
- /language - Change language preference

### Database
- SQLite database for data persistence
- Subscriber management
- Language preference storage
- Muted series tracking
- Content notification history (duplicate prevention)

### Deployment
- Docker support with multi-stage builds
- Pre-built binaries for multiple platforms (Linux, Windows, macOS)
- Cross-architecture support (AMD64, ARM64)
- Docker Compose configuration example
- Systemd service file for Linux

### Development
- Comprehensive test suite (184 tests)
- CI/CD pipelines (test, build, release, Docker)
- Code quality checks with golangci-lint
- Automated cross-platform builds with GoReleaser
- Structured logging with rotation

### Documentation
- Comprehensive README
- Quick Start guide (5-minute setup)
- Architecture documentation
- Deployment guides (Docker, binaries, systemd)
- Configuration reference
- Troubleshooting guide
- Contributing guidelines
- API integration documentation

### Security
- Environment variable-based configuration
- No hardcoded secrets
- Webhook secret validation support
- Non-root Docker container
- Clean git history (no leaked secrets)
- Branch protection for main branch

### Community
- MIT License
- Code of Conduct (Contributor Covenant)
- Issue templates (bug report, feature request)
- Pull request template
- Contributing guidelines
- GitHub Discussions for Q&A

## Technical Details

- **Language:** Go 1.23
- **Telegram Library:** go-telegram/bot
- **i18n Library:** nicksnyder/go-i18n v2
- **Database:** SQLite with GORM
- **Logging:** slog with lumberjack rotation
- **Docker:** Multi-stage Alpine-based (34.2MB)
- **Platforms:** Linux, Windows, macOS (AMD64, ARM64)

## Breaking Changes

This is the first public release, so there are no breaking changes from previous versions.

## Migration Notes

If migrating from a private/development version:
1. Database schema is automatically migrated (GORM AutoMigrate)
2. Existing subscribers will use English by default
3. Users can change language via /language command
4. No data loss - all existing data is preserved

## Known Limitations

- Only English and Persian languages supported (community can contribute more)
- Requires manual Jellyfin webhook configuration
- No built-in Jellyfin server (bot connects to existing Jellyfin)

## Contributors

- Initial development and implementation
- Community contributions welcome!

## Acknowledgments

- nicksnyder/go-i18n for internationalization
- go-telegram/bot for Telegram API
- GORM for database ORM
- lumberjack for log rotation
- GoReleaser for multi-platform builds

---

**Ready for production use!** 🎉"

# Verify tag was created
git tag -l v1.0.0

# Should show: v1.0.0

# View tag details
git show v1.0.0

# Should show the tag message and commit details
```

---

### Step 4: Push Tag to GitHub

**This step triggers the automated release workflow!**

```bash
# Push the tag to GitHub
git push origin v1.0.0

# Expected output:
# Enumerating objects: 1, done.
# Counting objects: 100% (1/1), done.
# Writing objects: 100% (1/1), 123 bytes | 123.00 KiB/s, done.
# Total 1 (delta 0), reused 0 (delta 0), pack-reused 0
# To github.com:YOUR_USERNAME/jellyfin-telegram-bot.git
#  * [new tag]         v1.0.0 -> v1.0.0
```

**The automated release workflow will now:**
1. Run all tests
2. Build binaries for all platforms
3. Build Docker images
4. Create GitHub Release
5. Upload binary artifacts
6. Push Docker images to GHCR
7. Generate checksums

---

### Step 5: Monitor Release Workflow

```bash
# Open GitHub Actions page
# https://github.com/YOUR_USERNAME/jellyfin-telegram-bot/actions

# Or use GitHub CLI (if installed)
gh run list --workflow=release.yml

# Watch the workflow in real-time
gh run watch
```

**Expected workflow stages:**
1. ✅ Checkout code
2. ✅ Setup Go
3. ✅ Run tests
4. ✅ Run GoReleaser
5. ✅ Build Docker images
6. ✅ Push Docker images
7. ✅ Create GitHub Release
8. ✅ Upload artifacts

**Timeline:** 5-10 minutes for complete workflow

**If any step fails:**
1. Check the workflow logs for errors
2. Fix the issue
3. Delete the tag: `git tag -d v1.0.0` and `git push origin :refs/tags/v1.0.0`
4. Create a new commit with the fix
5. Repeat from Step 3

---

### Step 6: Verify Release Assets

Once the workflow completes, verify the release was created successfully:

**Via Web UI:**
1. Go to `https://github.com/YOUR_USERNAME/jellyfin-telegram-bot/releases`
2. You should see "v1.0.0" release
3. Verify these assets are present:

```
✅ jellyfin-telegram-bot_1.0.0_linux_amd64.tar.gz
✅ jellyfin-telegram-bot_1.0.0_linux_arm64.tar.gz
✅ jellyfin-telegram-bot_1.0.0_windows_amd64.zip
✅ jellyfin-telegram-bot_1.0.0_darwin_amd64.tar.gz
✅ jellyfin-telegram-bot_1.0.0_darwin_arm64.tar.gz
✅ jellyfin-telegram-bot_1.0.0_checksums.txt
✅ Source code (zip)
✅ Source code (tar.gz)
```

**Via Command Line:**
```bash
# Using GitHub CLI
gh release view v1.0.0

# Check release assets
gh release view v1.0.0 --json assets -q '.assets[].name'
```

---

### Step 7: Verify Docker Images

Check that Docker images were published to GitHub Container Registry:

```bash
# Pull the latest tag
docker pull ghcr.io/YOUR_USERNAME/jellyfin-telegram-bot:latest

# Pull the specific version tag
docker pull ghcr.io/YOUR_USERNAME/jellyfin-telegram-bot:v1.0.0

# Verify image size (should be around 34MB)
docker images ghcr.io/YOUR_USERNAME/jellyfin-telegram-bot

# Test the image
docker run --rm ghcr.io/YOUR_USERNAME/jellyfin-telegram-bot:v1.0.0 --help

# Should display help information or version
```

**Image tags available:**
- `ghcr.io/YOUR_USERNAME/jellyfin-telegram-bot:v1.0.0` - Specific version
- `ghcr.io/YOUR_USERNAME/jellyfin-telegram-bot:latest` - Latest release

---

### Step 8: Test Download and Installation

Verify users can download and use the release:

**Test Linux AMD64 binary:**
```bash
# Download from GitHub Release
wget https://github.com/YOUR_USERNAME/jellyfin-telegram-bot/releases/download/v1.0.0/jellyfin-telegram-bot_1.0.0_linux_amd64.tar.gz

# Extract
tar -xzf jellyfin-telegram-bot_1.0.0_linux_amd64.tar.gz

# Run
./jellyfin-bot --version

# Clean up
rm -rf jellyfin-bot jellyfin-telegram-bot_1.0.0_linux_amd64.tar.gz
```

**Test Docker installation:**
```bash
# Pull and run
docker run --rm \
  -e TELEGRAM_BOT_TOKEN=test \
  -e JELLYFIN_SERVER_URL=http://localhost:8096 \
  -e JELLYFIN_API_KEY=test \
  ghcr.io/YOUR_USERNAME/jellyfin-telegram-bot:v1.0.0 \
  --help

# Should show help output without errors
```

---

### Step 9: Update README Badge (Optional)

Add a release badge to README.md:

```bash
# Edit README.md and add after existing badges:
```

```markdown
[![Release](https://img.shields.io/github/v/release/YOUR_USERNAME/jellyfin-telegram-bot)](https://github.com/YOUR_USERNAME/jellyfin-telegram-bot/releases)
[![Docker Image](https://img.shields.io/badge/docker-ghcr.io-blue)](https://github.com/YOUR_USERNAME/jellyfin-telegram-bot/pkgs/container/jellyfin-telegram-bot)
```

```bash
# Commit and push
git add README.md
git commit -m "docs: Add release and Docker badges to README"
git push origin main
```

---

### Step 10: Announce Release (Optional)

Consider announcing the release on:

1. **GitHub Discussions**
   - Create an announcement in the Announcements category
   - Link to the release notes
   - Highlight key features

2. **Social Media** (if desired)
   - Reddit: r/jellyfin, r/golang, r/selfhosted
   - Twitter/X
   - LinkedIn
   - Dev.to

3. **Jellyfin Community**
   - Jellyfin forum
   - Jellyfin Discord

**Example announcement:**
```markdown
# Jellyfin Telegram Bot v1.0.0 Released! 🎉

I'm excited to announce the first public release of jellyfin-telegram-bot!

## What is it?

A multi-language Telegram bot that sends notifications when new content is added to your Jellyfin media server.

## Key Features

- 🌍 Multi-language support (English & Persian) with auto-detection
- 🔕 Series muting for personalized notifications
- 🎨 Interactive inline keyboard menus
- 🐳 Docker support
- 📦 Pre-built binaries for all major platforms
- 🧪 Fully tested (184 tests!)

## Quick Start

1. Download the binary for your platform
2. Configure with your Telegram bot token and Jellyfin API
3. Run and start receiving notifications!

Full documentation: [link to repository]

## Contributing

Contributions welcome! Especially:
- Translations for additional languages
- Bug reports and feature requests
- Documentation improvements

MIT Licensed | Open Source | Community Driven

[Link to Release]
```

---

## Post-Release Tasks

After successful release:

### Day 1
- [x] Monitor GitHub Issues for bug reports
- [x] Monitor GitHub Discussions for questions
- [x] Respond to any issues or questions
- [x] Check download statistics

### Week 1
- [x] Address any critical bugs with hotfix release (v1.0.1)
- [x] Review user feedback
- [x] Update documentation based on questions
- [x] Consider implementing high-demand features

### Ongoing
- [x] Keep dependencies updated (Dependabot PRs)
- [x] Respond to PRs and issues
- [x] Plan for v1.1.0 with new features
- [x] Add community-requested language translations

---

## Hotfix Release Process

If a critical bug is found after release:

```bash
# Create hotfix branch from v1.0.0 tag
git checkout -b hotfix/v1.0.1 v1.0.0

# Make fixes
# ... edit files ...
# ... test thoroughly ...

# Commit fixes
git add .
git commit -m "fix: Critical bug in notification handling"

# Merge to main
git checkout main
git merge hotfix/v1.0.1

# Create new tag
git tag -a v1.0.1 -m "Release v1.0.1 - Hotfix

Fixes:
- Critical bug in notification handling
- Updated documentation

Full changelog: https://github.com/YOUR_USERNAME/jellyfin-telegram-bot/compare/v1.0.0...v1.0.1"

# Push
git push origin main
git push origin v1.0.1

# Clean up hotfix branch
git branch -d hotfix/v1.0.1
```

---

## Rollback Procedure

If you need to retract a release:

```bash
# Delete the release on GitHub
gh release delete v1.0.0 --yes

# Delete the tag
git tag -d v1.0.0
git push origin :refs/tags/v1.0.0

# Fix issues and create new release
# Follow the release process again with fixes
```

**Note:** Rollback should only be done in extreme cases (security issues, data loss bugs, etc.)

---

## Troubleshooting

### Release Workflow Failed

**Problem:** The release workflow fails during execution.

**Solutions:**

1. **Tests Failed:**
   - Check the test output in workflow logs
   - Run tests locally: `go test ./... -v`
   - Fix failing tests
   - Delete tag and recreate after fix

2. **GoReleaser Failed:**
   - Check `.goreleaser.yaml` syntax
   - Verify all platforms are valid
   - Check for file permission issues
   - Review GoReleaser logs in workflow

3. **Docker Build Failed:**
   - Check `Dockerfile` syntax
   - Verify base images are available
   - Check for network issues
   - Review Docker build logs

4. **Asset Upload Failed:**
   - Verify GitHub token has correct permissions
   - Check repository settings → Actions → General
   - Ensure "Read and write permissions" is enabled

### Binary Doesn't Work on Target Platform

**Problem:** Downloaded binary doesn't run on user's system.

**Solutions:**

1. **Linux:**
   - Ensure binary has execute permission: `chmod +x jellyfin-bot`
   - Check architecture: `uname -m` (should be x86_64 or aarch64)
   - Try static build if dynamic linking issues

2. **macOS:**
   - Remove quarantine: `xattr -d com.apple.quarantine jellyfin-bot`
   - Ensure correct architecture (Intel vs ARM)
   - Check macOS version compatibility

3. **Windows:**
   - Ensure antivirus isn't blocking
   - Run from Command Prompt/PowerShell
   - Check architecture (64-bit vs 32-bit)

### Docker Image Not Found

**Problem:** Users can't pull the Docker image.

**Solutions:**

1. Check image visibility in GHCR (should be public)
2. Verify image was pushed successfully
3. Check package permissions in GitHub
4. Try pulling with full path: `ghcr.io/username/repo:tag`

---

## Success Metrics

Track these metrics post-release:

1. **GitHub Metrics:**
   - Stars
   - Forks
   - Issues opened/closed
   - Pull requests
   - Download count

2. **Docker Metrics:**
   - Pull count
   - Unique users

3. **Community Metrics:**
   - Discussion participation
   - Contributors
   - Language translations added

4. **Quality Metrics:**
   - Bug report rate
   - Time to fix bugs
   - Test coverage
   - Code quality scores

---

**Release Complete!** 🚀

Your project is now publicly available and ready for the community to use and contribute to.

**Remember:**
- Respond to issues and PRs promptly
- Keep documentation updated
- Thank contributors
- Release updates regularly
- Have fun maintaining your open source project!

---

**Generated:** 2025-11-22
**For:** jellyfin-telegram-bot v1.0.0
**Status:** Ready to execute
