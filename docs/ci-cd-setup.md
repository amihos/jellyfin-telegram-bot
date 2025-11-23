# CI/CD Configuration Guide

This document explains how to configure your GitHub repository for the CI/CD pipeline.

## Repository Secrets

The following secrets need to be configured in your GitHub repository settings.

### Required Secrets

Navigate to **Settings** → **Secrets and variables** → **Actions** → **New repository secret**

1. **CODECOV_TOKEN** (Optional but recommended)
   - Sign up at https://codecov.io
   - Add your repository
   - Copy the upload token
   - Add as a repository secret named `CODECOV_TOKEN`
   - This enables code coverage reporting

### Built-in Secrets

The following secrets are automatically provided by GitHub Actions and don't need configuration:

- **GITHUB_TOKEN** - Used for:
  - Creating GitHub releases
  - Uploading release assets
  - Pushing Docker images to GHCR
  - Commenting on pull requests

## Repository Permissions

### 1. Actions Permissions

Go to **Settings** → **Actions** → **General**

Configure the following:

#### Workflow permissions
- Select: **Read and write permissions**
- Check: **Allow GitHub Actions to create and approve pull requests**

This allows workflows to:
- Create releases
- Push Docker images to GitHub Container Registry
- Upload release assets

### 2. Branch Protection Rules

Go to **Settings** → **Branches** → **Add rule**

**Branch name pattern:** `main`

Enable the following:

- ✅ **Require a pull request before merging**
  - ✅ Require approvals: 1
  - ✅ Dismiss stale pull request approvals when new commits are pushed
  
- ✅ **Require status checks to pass before merging**
  - ✅ Require branches to be up to date before merging
  - Add required status checks:
    - `Test (Go 1.22)`
    - `Test (Go 1.23)`
    - `Lint`
    - `Build (linux/amd64)`
    - `Build (linux/arm64)`
    - `Build (windows/amd64)`
    - `Build (darwin/amd64)`
    - `Build (darwin/arm64)`

- ✅ **Require conversation resolution before merging**

- ✅ **Do not allow bypassing the above settings**

### 3. GitHub Container Registry

The GitHub Container Registry (GHCR) is automatically available for your repository.

To make your Docker images public (optional):

1. Go to your GitHub profile → **Packages**
2. Find `jellyfin-telegram-bot`
3. Click **Package settings**
4. Scroll to **Danger Zone**
5. Click **Change visibility**
6. Select **Public**

This allows anyone to pull your Docker images without authentication:

```bash
docker pull ghcr.io/yourusername/jellyfin-telegram-bot:latest
```

## Testing the CI/CD Pipeline

### Test Workflows Locally

Before pushing, you can test some workflows locally:

#### Test Build Locally

```bash
# Test building for all platforms
GOOS=linux GOARCH=amd64 go build -o bot-linux-amd64 ./cmd/bot
GOOS=linux GOARCH=arm64 go build -o bot-linux-arm64 ./cmd/bot
GOOS=windows GOARCH=amd64 go build -o bot-windows-amd64.exe ./cmd/bot
GOOS=darwin GOARCH=amd64 go build -o bot-darwin-amd64 ./cmd/bot
GOOS=darwin GOARCH=arm64 go build -o bot-darwin-arm64 ./cmd/bot
```

#### Test GoReleaser Locally

Install GoReleaser:
```bash
# macOS
brew install goreleaser

# Linux
go install github.com/goreleaser/goreleaser@latest
```

Test release (without publishing):
```bash
goreleaser release --snapshot --clean --skip=publish
```

This will create builds in the `dist/` directory.

#### Test Docker Build Locally

```bash
# Build for current platform
docker build -t jellyfin-telegram-bot:test .

# Build for multiple platforms (requires buildx)
docker buildx build --platform linux/amd64,linux/arm64 -t jellyfin-telegram-bot:test .
```

### Test on GitHub

#### 1. Test Test and Build Workflows

Create a feature branch and push:

```bash
git checkout -b test-ci
git push origin test-ci
```

Create a pull request. The test and build workflows will run automatically.

#### 2. Test Release Workflow

**WARNING: This creates a real release!**

Only do this when you're ready to release version 1.0.0:

```bash
# Ensure main branch is up to date
git checkout main
git pull

# Create and push a version tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

This will:
- Trigger the release workflow
- Build binaries for all platforms
- Create a GitHub release
- Upload all binaries as release assets
- Build and push Docker images to GHCR

#### 3. Verify Release

After the release workflow completes:

1. Go to your repository's **Releases** page
2. You should see release v1.0.0 with:
   - Release notes
   - Binary assets for all platforms
   - Checksums file

3. Go to your repository's **Packages**
4. You should see Docker images with tags:
   - `latest`
   - `1.0.0`
   - `1.0`
   - `1`

## Monitoring Workflows

### View Workflow Runs

Go to **Actions** tab in your repository to see all workflow runs.

### View Workflow Logs

Click on any workflow run to see detailed logs for each job and step.

### Workflow Badges

Add these badges to your README to show workflow status:

```markdown
[![Test](https://github.com/yourusername/jellyfin-telegram-bot/workflows/Test/badge.svg)](https://github.com/yourusername/jellyfin-telegram-bot/actions/workflows/test.yml)
[![Build](https://github.com/yourusername/jellyfin-telegram-bot/workflows/Build/badge.svg)](https://github.com/yourusername/jellyfin-telegram-bot/actions/workflows/build.yml)
[![Release](https://github.com/yourusername/jellyfin-telegram-bot/workflows/Release/badge.svg)](https://github.com/yourusername/jellyfin-telegram-bot/actions/workflows/release.yml)
[![Docker](https://github.com/yourusername/jellyfin-telegram-bot/workflows/Docker/badge.svg)](https://github.com/yourusername/jellyfin-telegram-bot/actions/workflows/docker.yml)
```

## Troubleshooting

### Workflow Fails with Permission Denied

Check that Actions have write permissions:
- Settings → Actions → General → Workflow permissions
- Select "Read and write permissions"

### Docker Push Fails

Ensure GITHUB_TOKEN has package write permissions (it should by default).

### Release Not Created

- Verify tag matches pattern `v*.*.*` (e.g., v1.0.0)
- Check that release workflow has `contents: write` permission
- View workflow logs for specific errors

### Tests Fail on CI but Pass Locally

- Check Go version matches (1.22 or 1.23)
- Ensure all test dependencies are in go.mod
- Check for race conditions with `-race` flag locally:
  ```bash
  go test -race ./...
  ```

### Coverage Upload Fails

- CODECOV_TOKEN may be missing or invalid
- This is non-blocking (fail_ci_if_error: false)
- Coverage upload only runs on Go 1.23

## Next Steps

After configuration:

1. ✅ Configure repository secrets (CODECOV_TOKEN)
2. ✅ Set Actions permissions to "Read and write"
3. ✅ Configure branch protection on `main`
4. ✅ Make GHCR package public (optional)
5. ✅ Test workflows with a pull request
6. ✅ Create v1.0.0 release when ready
7. ✅ Add workflow badges to README
8. ✅ Monitor first few workflow runs

## Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [GoReleaser Documentation](https://goreleaser.com/)
- [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [Codecov Documentation](https://docs.codecov.com/)
