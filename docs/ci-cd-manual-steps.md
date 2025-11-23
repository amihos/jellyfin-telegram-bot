# Phase 7: Manual Steps Required

This document outlines the manual steps required to complete the CI/CD setup after pushing the code to GitHub.

## Before Pushing to GitHub

### 1. Update Repository References

All workflows and documentation use `yourusername` as a placeholder. Before pushing, update:

**Files to update:**
- `.goreleaser.yaml` - Lines 93-94 (owner and name)
- `README.md` - All badge URLs
- `docs/ci-cd-setup.md` - Example commands and URLs

**Search and replace:**
```bash
# In the repository root
find . -type f \( -name "*.yml" -o -name "*.yaml" -o -name "*.md" \) -exec sed -i 's/yourusername/ACTUAL_GITHUB_USERNAME/g' {} +
```

Replace `ACTUAL_GITHUB_USERNAME` with your actual GitHub username.

## After Pushing to GitHub

### Step 1: Configure Actions Permissions

1. Go to your repository on GitHub
2. Click **Settings** → **Actions** → **General**
3. Scroll to "Workflow permissions"
4. Select **"Read and write permissions"**
5. Check **"Allow GitHub Actions to create and approve pull requests"**
6. Click **Save**

**Why:** This allows workflows to create releases and push Docker images to GHCR.

### Step 2: Set Up Branch Protection (Optional but Recommended)

1. Go to **Settings** → **Branches**
2. Click **Add rule**
3. Branch name pattern: `main`
4. Enable:
   - ✅ Require a pull request before merging
     - Require approvals: 1
   - ✅ Require status checks to pass before merging
     - Add: `Test (Go 1.22)`, `Test (Go 1.23)`, `Lint`
     - Add: `Build (linux/amd64)`, etc. (all 5 platforms)
   - ✅ Require conversation resolution before merging
5. Click **Create**

**Why:** Ensures all tests pass before merging to main.

### Step 3: Configure Codecov (Optional)

1. Go to https://codecov.io
2. Sign in with GitHub
3. Add your repository
4. Copy the upload token
5. In GitHub: Settings → Secrets and variables → Actions
6. Click **New repository secret**
7. Name: `CODECOV_TOKEN`
8. Value: (paste the token)
9. Click **Add secret**

**Why:** Enables code coverage reporting and tracking.

### Step 4: Make Docker Images Public (Optional)

1. Go to your GitHub profile → **Packages**
2. Find `jellyfin-telegram-bot` package
3. Click **Package settings**
4. Scroll to "Danger Zone"
5. Click **Change visibility**
6. Select **Public**
7. Confirm

**Why:** Allows users to pull Docker images without authentication:
```bash
docker pull ghcr.io/yourusername/jellyfin-telegram-bot:latest
```

### Step 5: Test Workflows with a Pull Request

1. Create a test branch:
   ```bash
   git checkout -b test-ci-workflows
   ```

2. Make a small change (e.g., update a comment):
   ```bash
   echo "# CI/CD Test" >> README.md
   git add README.md
   git commit -m "test: Verify CI workflows"
   git push origin test-ci-workflows
   ```

3. Create a pull request on GitHub

4. Verify workflows run:
   - Go to **Actions** tab
   - You should see workflows running
   - Click on each to view logs

5. Expected results:
   - ✅ Test workflow: Runs (may have test failures from Phase 3)
   - ✅ Build workflow: All builds succeed
   - ✅ Docker workflow: Build succeeds (no push on PR)

6. Close the PR without merging

### Step 6: Create First Release (When Ready)

**IMPORTANT:** Only do this after Phase 8 is complete and all tests pass.

1. Ensure main branch is clean and up to date:
   ```bash
   git checkout main
   git pull
   git status  # Should show "nothing to commit, working tree clean"
   ```

2. Create and push the v1.0.0 tag:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0: Initial open source release"
   git push origin v1.0.0
   ```

3. Monitor the release workflow:
   - Go to **Actions** tab
   - Click on the "Release" workflow run
   - Wait for completion (~8-10 minutes)

4. Verify the release:
   - Go to **Releases** page
   - Check that v1.0.0 release exists
   - Verify all platform binaries are attached
   - Check checksums.txt is present

5. Verify Docker images:
   - Go to **Packages** tab
   - Check `jellyfin-telegram-bot` package
   - Verify tags: `latest`, `1.0.0`, `1.0`, `1`

6. Test downloading and running a binary:
   ```bash
   # Example for Linux
   wget https://github.com/yourusername/jellyfin-telegram-bot/releases/download/v1.0.0/jellyfin-telegram-bot_1.0.0_linux_x86_64.tar.gz
   tar -xzf jellyfin-telegram-bot_1.0.0_linux_x86_64.tar.gz
   ./jellyfin-telegram-bot --version
   ```

7. Test Docker image:
   ```bash
   docker pull ghcr.io/yourusername/jellyfin-telegram-bot:1.0.0
   docker run --rm ghcr.io/yourusername/jellyfin-telegram-bot:1.0.0 --version
   ```

## Verification Checklist

After completing all steps:

- [ ] Actions permissions set to "Read and write"
- [ ] Branch protection rules configured (optional)
- [ ] Codecov token added (optional)
- [ ] Docker package visibility set to public (optional)
- [ ] Test PR created and workflows verified
- [ ] Workflows show correct status badges in README
- [ ] Ready to create v1.0.0 release (after Phase 8)

## Troubleshooting

### Workflow Fails with "Permission denied"

**Issue:** Docker push fails with permission error

**Solution:** Check Actions permissions (Step 1 above)

### Badge Shows "No status"

**Issue:** Workflow badges in README don't update

**Solution:** 
1. Ensure workflows have run at least once
2. Check workflow file names match badge URLs
3. Wait a few minutes for GitHub cache to update

### Release Workflow Doesn't Trigger

**Issue:** Pushing a tag doesn't trigger the release workflow

**Solution:**
1. Check tag format is `v*.*.*` (e.g., v1.0.0)
2. Verify workflow file exists in `.github/workflows/release.yml`
3. Check Actions tab for any errors

### Docker Image Push Fails

**Issue:** "denied: permission_denied" when pushing to GHCR

**Solution:**
1. Verify Actions have "Read and write" permissions
2. Check that GITHUB_TOKEN is available in workflow
3. Ensure workflow has `packages: write` permission

### Can't Pull Docker Image

**Issue:** "unauthorized: unauthenticated" when pulling

**Solution:** Make the package public (Step 4 above)

## Next Steps

1. ✅ Complete Phase 7 manual steps (this document)
2. ⏭️ Proceed to Phase 8: Testing, Polish & Pre-Launch Verification
3. ⏭️ Fix any remaining test failures
4. ⏭️ Final documentation review
5. ⏭️ Create v1.0.0 release
6. ⏭️ Make repository public

## Support

If you encounter issues:

1. Check workflow logs in Actions tab
2. Review documentation:
   - `docs/ci-cd-setup.md`
   - `docs/ci-cd-testing.md`
3. Verify configuration against this checklist
4. Check GitHub Actions documentation

---

**After completing these steps, your CI/CD pipeline will be fully operational!**
