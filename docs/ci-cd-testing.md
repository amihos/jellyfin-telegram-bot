# CI/CD Testing Guide

This document provides instructions for testing the CI/CD pipeline locally before pushing to GitHub.

## Prerequisites

- Go 1.22 or 1.23 installed
- Docker with buildx support (for multi-platform builds)
- GoReleaser installed (optional, for release testing)
- golangci-lint installed (optional, for linting)

## Install Testing Tools

### Install GoReleaser

**macOS:**
```bash
brew install goreleaser
```

**Linux:**
```bash
go install github.com/goreleaser/goreleaser@latest
```

### Install golangci-lint

**macOS:**
```bash
brew install golangci-lint
```

**Linux:**
```bash
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```

## Test Workflow Steps Locally

### 1. Test Compilation (Test Workflow)

```bash
# Download dependencies
go mod download

# Verify dependencies
go mod verify

# Run tests (will show compilation errors if any)
go test -v -race ./...
```

**Note:** Some tests may fail due to incomplete mocks from Phase 3. This is expected and will be addressed in Phase 8.

### 2. Test Linting (Test Workflow)

```bash
# Run golangci-lint
golangci-lint run --timeout=5m

# Or with verbose output
golangci-lint run --timeout=5m -v
```

### 3. Test Multi-Platform Builds (Build Workflow)

```bash
# Test building for all platforms
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bot-linux-amd64 ./cmd/bot
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bot-linux-arm64 ./cmd/bot
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bot-windows-amd64.exe ./cmd/bot
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bot-darwin-amd64 ./cmd/bot
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bot-darwin-arm64 ./cmd/bot

# Clean up
rm -f bot-*
```

### 4. Test GoReleaser (Release Workflow)

```bash
# Check configuration
goreleaser check

# Test snapshot build (without publishing)
goreleaser build --snapshot --clean --skip=validate

# View built artifacts
ls -lh dist/

# Clean up
rm -rf dist/
```

**Expected output:**
- Binaries for 5 platforms (linux/amd64, linux/arm64, windows/amd64, darwin/amd64, darwin/arm64)
- Each binary should be ~10-15MB (stripped and optimized)

### 5. Test Docker Build (Docker Workflow)

```bash
# Build for current platform
docker build -t jellyfin-telegram-bot:test .

# Test running the container
docker run --rm jellyfin-telegram-bot:test --version

# Build for multiple platforms (requires buildx)
docker buildx create --use
docker buildx build --platform linux/amd64,linux/arm64 -t jellyfin-telegram-bot:test .

# Clean up
docker rmi jellyfin-telegram-bot:test
```

**Expected output:**
- Image size should be ~35MB
- Binary should execute without errors
- Multi-platform build should succeed

## Workflow Triggers

### Test Workflow (.github/workflows/test.yml)

**Triggers:**
- Push to `main` branch
- Pull requests to `main` branch

**Jobs:**
- Run tests on Go 1.22 and 1.23
- Run golangci-lint
- Upload code coverage to Codecov

### Build Workflow (.github/workflows/build.yml)

**Triggers:**
- Push to `main` branch
- Pull requests to `main` branch

**Jobs:**
- Build binaries for all 5 platforms
- Verify Linux/amd64 binary executes

### Docker Workflow (.github/workflows/docker.yml)

**Triggers:**
- Push to `main` branch
- Pull requests to `main` branch

**Jobs:**
- Build multi-platform Docker images
- Push to GHCR on push to main (not on PRs)

### Release Workflow (.github/workflows/release.yml)

**Triggers:**
- Push tags matching `v*.*.*` (e.g., v1.0.0, v1.2.3)

**Jobs:**
- Build and publish release with GoReleaser
- Create GitHub release with binaries
- Build and push Docker images to GHCR with version tags

## Testing on GitHub

### Test with a Pull Request

1. Create a feature branch:
   ```bash
   git checkout -b test-ci-workflows
   ```

2. Make a small change (e.g., update a comment in README)

3. Commit and push:
   ```bash
   git add README.md
   git commit -m "test: Verify CI workflows"
   git push origin test-ci-workflows
   ```

4. Create a pull request on GitHub

5. Verify workflows run:
   - Go to Actions tab
   - You should see "Test", "Build", and "Docker" workflows running
   - All should complete successfully (tests may fail if not all fixed yet)

6. Close the PR (don't merge yet)

### Test Release Workflow (CAUTION)

**WARNING:** This creates a real release! Only do this when ready for v1.0.0.

1. Ensure main branch is up to date:
   ```bash
   git checkout main
   git pull
   ```

2. Create and push a version tag:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0: Initial open source release"
   git push origin v1.0.0
   ```

3. Monitor the release workflow:
   - Go to Actions tab → Release workflow
   - Wait for completion (~5-10 minutes)

4. Verify release:
   - Go to Releases page
   - Check all binaries are uploaded
   - Check Docker images at Packages tab

## Troubleshooting

### Tests Fail Locally

**Issue:** `go test ./...` fails with compilation errors

**Solution:** This is expected if Phase 3 mocks are incomplete. The workflows will still build binaries successfully.

### golangci-lint Not Found

**Issue:** `golangci-lint: command not found`

**Solution:** Install golangci-lint:
```bash
# macOS
brew install golangci-lint

# Linux
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Ensure GOPATH/bin is in PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

### GoReleaser Version Issues

**Issue:** `version: 2` not supported error

**Solution:** The .goreleaser.yaml uses version 1 format for compatibility. Update GoReleaser if needed:
```bash
go install github.com/goreleaser/goreleaser@latest
```

### Docker Build Fails

**Issue:** Docker build fails with CGO errors

**Solution:** The Dockerfile uses CGO for SQLite. Ensure you have build tools:
```bash
# Ensure Docker is using the correct Dockerfile
docker build -f Dockerfile -t jellyfin-telegram-bot:test .
```

### Multi-Platform Docker Build Fails

**Issue:** `multiple platforms feature is currently not supported`

**Solution:** Enable Docker buildx:
```bash
docker buildx create --use
docker buildx inspect --bootstrap
```

### Workflow Permissions Error

**Issue:** Workflow fails with "Permission denied" when pushing Docker images

**Solution:** 
1. Go to Settings → Actions → General
2. Set Workflow permissions to "Read and write permissions"
3. Re-run the workflow

## Next Steps

After local testing succeeds:

1. ✅ Push workflows to GitHub
2. ✅ Test with a pull request
3. ✅ Configure repository settings (see docs/ci-cd-setup.md)
4. ✅ Set up branch protection rules
5. ✅ Add Codecov token (optional)
6. ✅ Create v1.0.0 release when ready

## CI/CD Pipeline Summary

| Workflow | Trigger | Purpose | Duration |
|----------|---------|---------|----------|
| Test | PR, push to main | Run tests and linting | ~2-3 min |
| Build | PR, push to main | Verify compilation for all platforms | ~3-4 min |
| Docker | PR, push to main | Build multi-platform images | ~4-5 min |
| Release | Tag v*.*.* | Create release with binaries and Docker images | ~8-10 min |

**Total pipeline time for PR:** ~10 minutes
**Total pipeline time for release:** ~10 minutes

## Files Created in Phase 7

```
.github/
  workflows/
    test.yml       - Test and lint workflow
    build.yml      - Build verification workflow
    docker.yml     - Docker image workflow
    release.yml    - Release workflow
.goreleaser.yaml   - GoReleaser configuration
.golangci.yml      - golangci-lint configuration
docs/
  ci-cd-setup.md   - CI/CD configuration guide
  ci-cd-testing.md - This file
```

## Expected Workflow Status

After Phase 7 completion:

- ✅ Test workflow - Ready (may have test failures from Phase 3)
- ✅ Build workflow - Ready
- ✅ Docker workflow - Ready
- ✅ Release workflow - Ready (waiting for v1.0.0 tag)
- ✅ GoReleaser - Configured and tested
- ✅ Multi-platform binaries - Tested locally
- ✅ Docker multi-platform - Tested locally
- ⏳ Codecov integration - Needs token configuration
- ⏳ Branch protection - Needs GitHub configuration
- ⏳ Release v1.0.0 - Waiting for Phase 8 completion
