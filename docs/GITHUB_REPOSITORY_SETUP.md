# GitHub Repository Configuration Guide

This document provides step-by-step instructions for configuring the GitHub repository settings for optimal discoverability and community engagement.

---

## Prerequisites

- Repository must be pushed to GitHub
- You must have admin access to the repository
- Repository should be made public before following these steps

---

## Step-by-Step Configuration

### 1. Repository General Settings

**Navigate to:** Repository → Settings → General

#### Description & Website

1. **Description** (160 character limit):
   ```
   Multi-language Telegram bot for Jellyfin media server notifications. Supports English and Persian with auto-detection, series muting, and interactive inline keyboards.
   ```

2. **Website** (optional):
   - If you have a deployment URL or documentation site, add it here
   - Otherwise, leave blank

#### Topics/Tags

Click "Add topics" and add these tags (one at a time):

1. `golang`
2. `telegram-bot`
3. `jellyfin`
4. `webhook`
5. `notifications`
6. `i18n`
7. `docker`
8. `media-server`
9. `persian`
10. `localization`

These tags help users discover the project when searching GitHub.

#### Features

Under "Features" section:

- ✅ **Issues:** Keep ENABLED
- ❌ **Wikis:** DISABLE (we use repository docs instead)
- ✅ **Discussions:** ENABLE (for Q&A and community support)
- ❌ **Projects:** DISABLE (not needed for this project)
- ✅ **Sponsorships:** ENABLE if you want to accept sponsorships

#### Pull Requests

Under "Pull Requests" section:

- ✅ **Allow merge commits:** ENABLE
- ✅ **Allow squash merging:** ENABLE
- ✅ **Allow rebase merging:** ENABLE
- ✅ **Automatically delete head branches:** ENABLE
- ❌ **Allow auto-merge:** DISABLE
- ✅ **Always suggest updating pull request branches:** ENABLE

#### Archives

- ❌ **Include Git LFS objects in archives:** DISABLE (we don't use LFS)

Click **Save** at the bottom.

---

### 2. Branch Protection Rules

**Navigate to:** Repository → Settings → Branches

Click "Add branch protection rule" and configure:

#### Branch Name Pattern
```
main
```

#### Protection Settings

**Protect matching branches:**
- ✅ Require a pull request before merging
  - ✅ Require approvals: **1** approval required
  - ✅ Dismiss stale pull request approvals when new commits are pushed
  - ❌ Require review from Code Owners (we don't have CODEOWNERS file)

- ✅ Require status checks to pass before merging
  - ✅ Require branches to be up to date before merging
  - Add required status checks (click "Search for status checks"):
    - Type: `test` (from test.yml workflow) - press Enter
    - Type: `build` (from build.yml workflow) - press Enter
  - Note: These will only appear after the workflows have run at least once

- ❌ Require conversation resolution before merging (optional, can enable if desired)

- ❌ Require signed commits (optional, advanced feature)

- ❌ Require linear history (we allow merge commits)

- ❌ Require deployments to succeed before merging (we don't have deployments)

**Do not allow bypassing the above settings:**
- ❌ Allow force pushes (keep DISABLED for main branch)
- ❌ Allow deletions (keep DISABLED for main branch)

Click **Create** to save the branch protection rule.

---

### 3. GitHub Actions Permissions

**Navigate to:** Repository → Settings → Actions → General

#### Actions permissions

Select: **Allow all actions and reusable workflows**

#### Workflow permissions

Select: **Read and write permissions**

- ✅ Allow GitHub Actions to create and approve pull requests

This is required for the release workflow to create releases and upload assets.

Click **Save**.

---

### 4. GitHub Pages (Optional)

**Navigate to:** Repository → Settings → Pages

If you want to host documentation via GitHub Pages:

#### Source
- Branch: `main`
- Folder: `/docs`

Click **Save**.

Your documentation will be available at:
```
https://YOUR_USERNAME.github.io/jellyfin-telegram-bot/
```

**Note:** This is optional. The README.md already contains comprehensive documentation.

---

### 5. Secrets and Variables

**Navigate to:** Repository → Settings → Secrets and variables → Actions

#### Repository Secrets

The following secrets are **automatically provided** by GitHub Actions:
- `GITHUB_TOKEN` - No action needed

#### Optional Secrets (only if needed)

If you plan to publish to Docker Hub (in addition to GHCR):

1. Click "New repository secret"
2. Name: `DOCKER_USERNAME`
3. Value: Your Docker Hub username
4. Click "Add secret"

5. Click "New repository secret"
6. Name: `DOCKER_PASSWORD`
7. Value: Your Docker Hub access token (not password!)
8. Click "Add secret"

**Note:** The current configuration uses GitHub Container Registry (GHCR), so Docker Hub secrets are not required.

---

### 6. Security Settings

**Navigate to:** Repository → Settings → Security

#### Dependabot

- ✅ **Dependabot alerts:** ENABLE
- ✅ **Dependabot security updates:** ENABLE

This will automatically create PRs to update vulnerable dependencies.

#### Code scanning

GitHub will automatically suggest enabling CodeQL for Go projects.

- ✅ **Enable CodeQL analysis:** ENABLE (click "Set up" button)
- Use default configuration

This adds an additional security scanning workflow.

---

### 7. General Repository Settings

**Navigate to:** Repository → Settings → General

#### Social Preview

Upload a social preview image (optional):
- Recommended size: 1280x640 pixels
- Format: PNG or JPEG
- Shows when sharing the repository on social media

If you don't have a logo, you can:
- Skip this step
- Create one later
- Use a screenshot of the bot in action

---

### 8. Repository Visibility

**IMPORTANT:** Only do this when you're ready to launch!

**Navigate to:** Repository → Settings → General → Danger Zone

1. Click "Change visibility"
2. Select "Make public"
3. Type the repository name to confirm
4. Click "I understand, change repository visibility"

**Note:** Once public, the repository will be:
- Visible to everyone
- Indexed by search engines
- Cloneable by anyone
- Subject to the MIT license terms

---

## Verification Checklist

After completing the configuration, verify:

- [ ] Repository description is set and visible on main page
- [ ] Topics/tags are visible on main page
- [ ] Discussions tab is visible
- [ ] Wiki tab is NOT visible
- [ ] Branch protection is active on `main` branch
- [ ] GitHub Actions workflows can run successfully
- [ ] Dependabot is enabled and scanning
- [ ] Repository is public (if ready to launch)

---

## Post-Configuration Tasks

### Initialize Discussions

1. Go to **Discussions** tab
2. Click "New discussion"
3. Create a welcome discussion:
   - Title: "Welcome to Jellyfin Telegram Bot!"
   - Category: Announcements
   - Content:
     ```markdown
     # Welcome! 👋

     Thank you for your interest in the Jellyfin Telegram Bot!

     This bot provides multi-language Telegram notifications for your Jellyfin media server, with support for English and Persian.

     ## Quick Links

     - [Quick Start Guide](../QUICKSTART.md)
     - [Full Documentation](../README.md)
     - [Contributing Guidelines](../CONTRIBUTING.md)

     ## How to Get Help

     - 🐛 **Found a bug?** [Open an issue](../issues/new/choose)
     - 💡 **Have a feature idea?** [Open an issue](../issues/new/choose)
     - ❓ **Have a question?** Start a discussion in Q&A

     ## Contributing

     We welcome contributions! Check out the [Contributing Guide](../CONTRIBUTING.md) to get started.

     Happy streaming! 🍿
     ```

### Enable Issue Labels

GitHub automatically provides default labels. Consider adding these custom labels:

**Navigate to:** Repository → Issues → Labels

Add these labels (click "New label"):

1. **priority: high**
   - Color: `#d73a4a` (red)
   - Description: High priority issue

2. **priority: low**
   - Color: `#0e8a16` (green)
   - Description: Low priority issue

3. **good first issue**
   - Color: `#7057ff` (purple)
   - Description: Good for newcomers

4. **help wanted**
   - Color: `#008672` (teal)
   - Description: Extra attention is needed

5. **translation**
   - Color: `#fbca04` (yellow)
   - Description: Related to i18n/translations

6. **docker**
   - Color: `#0969da` (blue)
   - Description: Related to Docker deployment

---

## Troubleshooting

### Branch Protection Not Showing Status Checks

**Problem:** When adding branch protection, the status checks don't appear in the search.

**Solution:**
1. Push a commit to trigger the workflows
2. Wait for workflows to complete
3. The status checks will now appear in the search
4. Add them to the branch protection rule

### Workflows Not Running

**Problem:** GitHub Actions workflows don't run after push.

**Solution:**
1. Check Repository → Settings → Actions → General
2. Ensure "Allow all actions" is selected
3. Ensure workflows have "Read and write permissions"
4. Check if there's a `.github/workflows/` directory
5. Check if YAML files are valid (no syntax errors)

### Dependabot Not Creating PRs

**Problem:** Dependabot is enabled but not creating PRs.

**Solution:**
1. It can take up to 24 hours for first scan
2. Check Repository → Security → Dependabot alerts
3. Ensure "Dependabot security updates" is enabled
4. Check if `go.mod` is in the repository root

---

## Security Best Practices

1. **Never commit secrets**
   - Always use GitHub Secrets for sensitive data
   - Use `.env.example` for configuration templates
   - Add `.env` to `.gitignore`

2. **Review PRs carefully**
   - Always require at least 1 approval
   - Enable branch protection
   - Run all tests before merging

3. **Keep dependencies updated**
   - Enable Dependabot
   - Review and merge security updates promptly
   - Test updates before merging

4. **Monitor security alerts**
   - Check Security tab regularly
   - Enable email notifications for security alerts
   - Respond to alerts within 24 hours

---

## Community Management Tips

1. **Respond promptly**
   - Try to respond to issues within 24-48 hours
   - Even if you can't fix immediately, acknowledge the report

2. **Be welcoming**
   - Thank contributors for their time
   - Be patient with first-time contributors
   - Follow the Code of Conduct

3. **Keep documentation updated**
   - Update docs when features change
   - Fix documentation issues reported by users
   - Accept documentation PRs

4. **Release regularly**
   - Use semantic versioning
   - Write clear release notes
   - Tag releases with useful information

---

**Configuration Complete!** 🎉

Your repository is now configured for optimal community engagement and discoverability.

**Next Step:** Create the v1.0.0 release following the [Pre-Launch Checklist](./PRE_LAUNCH_CHECKLIST.md).
