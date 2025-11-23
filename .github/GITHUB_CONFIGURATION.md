# GitHub Repository Configuration Guide

This document provides step-by-step instructions for configuring the GitHub repository settings for jellyfin-telegram-bot.

## Table of Contents

1. [Repository Description and Topics](#repository-description-and-topics)
2. [Branch Protection Rules](#branch-protection-rules)
3. [Repository Features](#repository-features)
4. [Issue Labels](#issue-labels)
5. [GitHub Discussions](#github-discussions)

---

## Repository Description and Topics

### Setting Repository Description

1. Go to your repository on GitHub: `https://github.com/yourusername/jellyfin-telegram-bot`
2. Click the **Settings** tab (gear icon, top right)
3. In the **About** section (right sidebar on the main repository page), click the gear icon
4. Set the following:
   - **Description**: `Telegram bot for Jellyfin media server notifications with multi-language support (English, Persian). Receive instant updates for new movies, shows, and episodes.`
   - **Website** (optional): Add your documentation URL or leave blank
   - **Topics**: Add the following tags (click in the Topics field and type each one, pressing Enter after each):
     - `golang`
     - `telegram-bot`
     - `jellyfin`
     - `webhook`
     - `notifications`
     - `i18n`
     - `docker`
     - `telegram`
     - `bot`
     - `media-server`
5. Click **Save changes**

---

## Branch Protection Rules

### Protect the Main Branch

1. Go to **Settings** → **Branches** (in the left sidebar)
2. Under **Branch protection rules**, click **Add rule** or **Add branch protection rule**
3. Configure the rule:

   **Branch name pattern**:
   ```
   main
   ```

   **Protect matching branches** - Check the following options:

   - ✅ **Require a pull request before merging**
     - ✅ Require approvals: **1**
     - ✅ Dismiss stale pull request approvals when new commits are pushed
     - ✅ Require review from Code Owners (optional)

   - ✅ **Require status checks to pass before merging**
     - ✅ Require branches to be up to date before merging
     - **Status checks** (these will appear after CI/CD is set up in Phase 7):
       - `test` (from test.yml workflow)
       - `lint` (from test.yml workflow)
       - `build-matrix` (from build.yml workflow)

   - ✅ **Require conversation resolution before merging**

   - ✅ **Require linear history** (optional, keeps git history clean)

   - ✅ **Include administrators** (enforce rules for everyone, including repo owner)

   - ✅ **Restrict who can push to matching branches** (optional)
     - If you want to restrict direct pushes, add specific users/teams

   - ❌ **Allow force pushes** (keep unchecked for safety)

   - ❌ **Allow deletions** (keep unchecked to prevent accidental branch deletion)

4. Click **Create** or **Save changes**

**Note**: If you haven't set up CI/CD workflows yet (Phase 7), you won't see status checks to add yet. You can come back and add them later after workflows are created.

---

## Repository Features

### Configure Repository Features

1. Go to **Settings** → **General** (should be the default page)
2. Scroll down to the **Features** section
3. Configure the following:

   - ✅ **Issues** - Enabled (for bug reports and feature requests)
   - ❌ **Wikis** - Disabled (documentation is in the repository)
   - ✅ **Discussions** - Enabled (for community Q&A, see below for setup)
   - ✅ **Projects** - Optional (enable if you want to use GitHub Projects for planning)
   - ✅ **Preserve this repository** - Optional (for archival)

4. Scroll down and click **Save changes**

### Configure Merge Button Options

Still in **Settings** → **General**, scroll to **Pull Requests**:

- ✅ **Allow merge commits** - Enabled
- ✅ **Allow squash merging** - Enabled (recommended, keeps history clean)
- ✅ **Allow rebase merging** - Enabled
- ✅ **Always suggest updating pull request branches** - Enabled
- ✅ **Allow auto-merge** - Optional
- ✅ **Automatically delete head branches** - Enabled (cleans up merged PR branches)

Click **Save changes**

---

## Issue Labels

### Create Issue Labels

GitHub provides default labels, but we'll add project-specific ones:

1. Go to **Issues** → **Labels** (or direct URL: `https://github.com/yourusername/jellyfin-telegram-bot/labels`)
2. For each label below, click **New label** and fill in the details:

#### Required Labels

| Name | Description | Color |
|------|-------------|-------|
| `bug` | Something isn't working | `#d73a4a` (red) |
| `enhancement` | New feature or request | `#a2eeef` (light blue) |
| `documentation` | Improvements or additions to documentation | `#0075ca` (blue) |
| `good first issue` | Good for newcomers | `#7057ff` (purple) |
| `help wanted` | Extra attention is needed | `#008672` (green) |
| `question` | Further information is requested | `#d876e3` (pink) |
| `translation` | Translation-related issues or new language support | `#f9d0c4` (peach) |
| `wontfix` | This will not be worked on | `#ffffff` (white) |
| `duplicate` | This issue or pull request already exists | `#cfd3d7` (gray) |
| `invalid` | This doesn't seem right | `#e4e669` (yellow) |
| `priority: critical` | Critical priority, needs immediate attention | `#b60205` (dark red) |
| `priority: high` | High priority | `#d93f0b` (orange-red) |
| `priority: medium` | Medium priority | `#fbca04` (yellow) |
| `priority: low` | Low priority | `#0e8a16` (green) |
| `status: in progress` | Currently being worked on | `#c2e0c6` (light green) |
| `status: blocked` | Blocked by another issue or external dependency | `#b60205` (dark red) |
| `docker` | Related to Docker deployment | `#0db7ed` (docker blue) |
| `telegram` | Related to Telegram bot functionality | `#0088cc` (telegram blue) |
| `jellyfin` | Related to Jellyfin integration | `#00a4dc` (jellyfin blue) |
| `i18n` | Internationalization and localization | `#f9d0c4` (peach) |
| `security` | Security-related issue | `#d73a4a` (red) |

**Note**: Default labels like `bug`, `enhancement`, `documentation`, `duplicate`, `invalid`, `wontfix`, `help wanted`, `good first issue`, and `question` may already exist. You can edit their colors and descriptions if needed.

---

## GitHub Discussions

### Enable and Configure Discussions

1. Go to **Settings** → **General**
2. In the **Features** section, check ✅ **Discussions**
3. Click **Set up discussions**
4. GitHub will create a welcome post. You can edit it or create categories.

### Create Discussion Categories

1. Go to **Discussions** tab
2. Click **Categories** (or the pencil icon to edit categories)
3. Create the following categories:

| Category | Description | Format |
|----------|-------------|---------|
| 📣 Announcements | Updates and news about the project | Announcement |
| 💡 Ideas | Share ideas for new features or improvements | Open-ended discussion |
| 🙏 Q&A | Ask questions and get help from the community | Question/Answer |
| 🐛 Troubleshooting | Get help with issues and errors | Question/Answer |
| 📚 Show and Tell | Share your setup, configurations, or use cases | Open-ended discussion |
| 🌍 Translations | Discuss translation updates and new language support | Open-ended discussion |

### Create a Welcome Post (Optional)

After enabling Discussions, create a welcome post:

1. Go to **Discussions**
2. Click **New discussion**
3. Select **Announcements** category
4. Title: `Welcome to Jellyfin Telegram Bot Discussions!`
5. Body:

```markdown
# Welcome to the Jellyfin Telegram Bot Community!

Thank you for using and contributing to this project! This is a space where you can:

- **Ask questions** about installation, configuration, or usage
- **Share ideas** for new features or improvements
- **Get help** troubleshooting issues
- **Show off** your setup and configurations
- **Discuss translations** and help add new languages

## Quick Links

- [Installation Guide](../blob/main/README.md#quick-start)
- [Contributing Guidelines](../blob/main/CONTRIBUTING.md)
- [Report a Bug](../issues/new?template=bug_report.md)
- [Request a Feature](../issues/new?template=feature_request.md)

## Community Guidelines

Please read our [Code of Conduct](../blob/main/CODE_OF_CONDUCT.md) before participating.

**Tips for getting help:**
- Check existing discussions before posting
- Provide details about your setup (OS, bot version, installation method)
- Share relevant logs (with sensitive info removed)
- Be respectful and patient

Happy automating!
```

6. Click **Start discussion**

---

## Verification Checklist

After completing the configuration, verify:

- [ ] Repository has description and all 10 topics
- [ ] Main branch has protection rules enabled
- [ ] Branch protection requires PR reviews (at least 1 approval)
- [ ] Status checks will be required when CI/CD is set up
- [ ] Issues are enabled
- [ ] Wikis are disabled
- [ ] Discussions are enabled with categories
- [ ] All issue labels are created
- [ ] Automatically delete head branches is enabled
- [ ] Force pushes are disabled on main branch

---

## Notes

- **Status checks**: After you set up CI/CD workflows in Phase 7, return to branch protection settings and add the required status checks (`test`, `lint`, `build-matrix`)
- **Collaborators**: If you want to add collaborators, go to **Settings** → **Collaborators** and invite them
- **Secrets**: GitHub Actions secrets will be configured in Phase 7 (CI/CD setup)
- **GitHub Pages**: Not needed for this project, but can be enabled if you want to host documentation

---

## Additional Resources

- [GitHub Branch Protection Documentation](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/about-protected-branches)
- [GitHub Discussions Guide](https://docs.github.com/en/discussions)
- [Managing Labels](https://docs.github.com/en/issues/using-labels-and-milestones-to-track-work/managing-labels)
