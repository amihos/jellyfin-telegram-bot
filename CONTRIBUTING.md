# Contributing to Jellyfin Telegram Bot

Thank you for your interest in contributing to the Jellyfin Telegram Bot! We welcome contributions of all kinds: bug fixes, new features, documentation improvements, and translations.

This guide will help you get started with contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Code Style Guidelines](#code-style-guidelines)
- [Testing Requirements](#testing-requirements)
- [Pull Request Process](#pull-request-process)
- [Adding Translations](#adding-translations)
- [Reporting Bugs](#reporting-bugs)
- [Suggesting Features](#suggesting-features)

## Code of Conduct

This project follows a Code of Conduct to ensure a welcoming environment for everyone. By participating, you agree to uphold this code. Please report unacceptable behavior to the project maintainers.

**Expected Behavior**:
- Be respectful and inclusive
- Welcome newcomers and help them get started
- Accept constructive criticism gracefully
- Focus on what's best for the community
- Show empathy towards other contributors

**Unacceptable Behavior**:
- Harassment, discrimination, or offensive comments
- Personal attacks or trolling
- Publishing others' private information
- Any conduct that would be inappropriate in a professional setting

## How Can I Contribute?

### Reporting Bugs

Found a bug? Please [open an issue](https://github.com/yourusername/jellyfin-telegram-bot/issues/new?template=bug_report.md) with:
- Clear description of the bug
- Steps to reproduce
- Expected vs actual behavior
- Your environment (OS, bot version, Docker/binary)
- Relevant log output (with sensitive data removed)

### Suggesting Features

Have an idea for a new feature? Please [open a feature request](https://github.com/yourusername/jellyfin-telegram-bot/issues/new?template=feature_request.md) with:
- Description of the problem you're trying to solve
- Your proposed solution
- Alternative approaches you've considered
- Why this would be valuable to other users

### Improving Documentation

Documentation improvements are always welcome! You can:
- Fix typos or unclear instructions
- Add examples or use cases
- Improve troubleshooting guides
- Translate documentation to other languages

### Adding Features

Want to add a new feature? Great! Please:
1. Check if an issue already exists for this feature
2. If not, open a feature request to discuss it first
3. Wait for maintainer feedback before starting work
4. Follow the development setup and PR process below

### Adding Translations

We welcome translations to new languages! See [Adding Translations](#adding-translations) below.

## Development Setup

### Prerequisites

- **Go 1.22 or higher** - [Download](https://go.dev/dl/)
- **Git** - [Download](https://git-scm.com/downloads)
- **A Telegram bot token** - Get from [@BotFather](https://t.me/BotFather)
- **A Jellyfin server** - For testing (can be local)

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork:

```bash
git clone https://github.com/yourusername/jellyfin-telegram-bot.git
cd jellyfin-telegram-bot
```

3. Add the upstream remote:

```bash
git remote add upstream https://github.com/originalowner/jellyfin-telegram-bot.git
```

### Install Dependencies

```bash
# Download Go dependencies
go mod download

# Verify everything downloaded correctly
go mod verify
```

### Configure for Local Testing

1. Copy the example environment file:

```bash
cp .env.example .env
```

2. Edit `.env` with your test credentials:

```env
TELEGRAM_BOT_TOKEN=your_test_bot_token
JELLYFIN_SERVER_URL=http://localhost:8096
JELLYFIN_API_KEY=your_test_api_key
PORT=8080
LOG_LEVEL=DEBUG
LOG_FILE=./logs/bot.log
```

**Important**: Use a separate test bot token, not your production bot!

3. Create the logs directory:

```bash
mkdir -p logs
```

### Run the Bot

```bash
# Run directly
go run cmd/bot/main.go

# Or build and run
go build -o jellyfin-telegram-bot cmd/bot/main.go
./jellyfin-telegram-bot
```

The bot should start and connect to Telegram. Check the logs for any errors.

### Verify Setup

1. Send `/start` to your test bot in Telegram
2. Add a test movie to Jellyfin
3. Verify you receive a notification

If everything works, you're ready to start developing!

## Making Changes

### Create a Feature Branch

Always create a new branch for your changes:

```bash
# Update your main branch
git checkout main
git pull upstream main

# Create a feature branch
git checkout -b feature/your-feature-name

# Or for bug fixes
git checkout -b fix/bug-description
```

**Branch Naming Conventions**:
- Features: `feature/feature-name`
- Bug fixes: `fix/bug-description`
- Documentation: `docs/what-you-changed`
- Translations: `i18n/language-code`

### Make Your Changes

- Write clear, readable code
- Follow Go conventions and idioms
- Add comments for complex logic
- Keep commits focused and atomic
- Write meaningful commit messages

**Good Commit Messages**:
```
Add support for video notifications

- Detect video item type in webhook handler
- Add video formatting to notification message
- Include video duration and codec in details
- Update tests for video notifications
```

**Bad Commit Messages**:
```
fixed stuff
Update file.go
WIP
```

### Keep Your Branch Updated

Regularly sync with upstream to avoid conflicts:

```bash
# Fetch upstream changes
git fetch upstream

# Rebase your branch on upstream/main
git rebase upstream/main

# If conflicts occur, resolve them and continue
git rebase --continue
```

## Code Style Guidelines

### Go Code Style

We follow standard Go conventions:

1. **Use `gofmt`** - Format all code before committing:

```bash
# Format all Go files
go fmt ./...
```

2. **Use `go vet`** - Check for common errors:

```bash
go vet ./...
```

3. **Use `golangci-lint`** (if installed) - Advanced linting:

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run
```

### Code Organization

- **Keep functions small** - Each function should do one thing well
- **Use meaningful names** - Variables, functions, and types should be self-documenting
- **Add comments** - Explain *why*, not *what* (code shows what)
- **Handle errors properly** - Always check and handle errors, don't ignore them
- **Use structured logging** - Use `slog.Info`, `slog.Error`, etc. with context

**Example of Good Code**:

```go
// GetUserLanguage retrieves the user's preferred language from the database.
// It returns the language code (e.g., "en", "fa") or an error if the query fails.
// If no preference is found, it returns "en" as the default.
func (db *Database) GetUserLanguage(chatID int64) (string, error) {
    var subscriber Subscriber

    result := db.db.Where("chat_id = ?", chatID).First(&subscriber)
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            // No preference set, return default
            return "en", nil
        }
        // Database error
        return "", fmt.Errorf("failed to query user language: %w", result.Error)
    }

    return subscriber.LanguageCode, nil
}
```

### Project Structure

When adding new files, follow the existing structure:

```
jellyfin-telegram-bot/
├── cmd/bot/              # Application entry point only
├── internal/
│   ├── config/           # Configuration and logging
│   ├── database/         # Database models and queries
│   ├── handlers/         # HTTP request handlers
│   ├── telegram/         # Telegram bot logic
│   ├── jellyfin/         # Jellyfin API client
│   └── i18n/             # Internationalization
├── locales/              # Translation files (.toml)
├── docs/                 # Documentation (.md files)
└── test/                 # Integration tests
```

**Rules**:
- `internal/` code can only be imported by this project
- Each package should have a clear, single responsibility
- Avoid circular dependencies
- Put tests next to the code they test (`_test.go` files)

## Testing Requirements

### Writing Tests

All new features and bug fixes should include tests:

1. **Unit tests** - Test individual functions:

```go
// internal/database/subscriber_test.go
func TestGetUserLanguage(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    // Test default when no preference set
    lang, err := db.GetUserLanguage(12345)
    assert.NoError(t, err)
    assert.Equal(t, "en", lang)

    // Test returning saved preference
    db.SetLanguage(12345, "fa")
    lang, err = db.GetUserLanguage(12345)
    assert.NoError(t, err)
    assert.Equal(t, "fa", lang)
}
```

2. **Integration tests** - Test components working together:

```go
// test/integration/notification_test.go
func TestNotificationWithI18n(t *testing.T) {
    // Test full notification flow with different languages
}
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test ./... -v

# Run tests with race detection
go test ./... -race

# Run tests with coverage
go test ./... -cover

# Run specific package tests
go test ./internal/database/...

# Run specific test
go test ./internal/database -run TestGetUserLanguage
```

### Test Coverage

We aim for reasonable test coverage, focusing on:
- Critical business logic
- Error handling paths
- Edge cases and boundary conditions

**Note**: We don't enforce 100% coverage. Focus on testing what matters.

### Before Submitting PR

Ensure all tests pass:

```bash
# Run tests
go test ./...

# Format code
go fmt ./...

# Check for issues
go vet ./...
```

## Pull Request Process

### 1. Prepare Your Changes

- Ensure all tests pass
- Format your code with `go fmt`
- Update documentation if needed
- Add or update tests for your changes

### 2. Push Your Branch

```bash
git push origin feature/your-feature-name
```

### 3. Create Pull Request

1. Go to your fork on GitHub
2. Click "New Pull Request"
3. Select your feature branch
4. Fill out the PR template:
   - Describe what your PR does
   - Reference any related issues
   - Check all items in the checklist
   - Note any breaking changes

### 4. PR Review Process

- A maintainer will review your PR
- They may request changes or ask questions
- Address feedback by pushing new commits
- Once approved, a maintainer will merge your PR

### PR Checklist

Before submitting, verify:

- [ ] Code follows Go style guidelines
- [ ] All tests pass (`go test ./...`)
- [ ] Code is formatted (`go fmt ./...`)
- [ ] No linting errors (`go vet ./...`)
- [ ] Documentation updated if needed
- [ ] Commit messages are clear and descriptive
- [ ] PR description explains the changes
- [ ] Tests added for new functionality
- [ ] No breaking changes (or clearly documented)

### After Your PR is Merged

1. Delete your feature branch:

```bash
git checkout main
git branch -d feature/your-feature-name
git push origin --delete feature/your-feature-name
```

2. Update your main branch:

```bash
git pull upstream main
```

3. Celebrate! 🎉 You've contributed to the project!

## Adding Translations

We welcome translations to new languages!

### 1. Choose Your Language

Check `locales/` directory for existing languages:
- `active.en.toml` - English
- `active.fa.toml` - Persian

If your language doesn't exist, you can add it!

### 2. Create Translation File

Copy the English template:

```bash
cp locales/active.en.toml locales/active.XX.toml
```

Replace `XX` with your language code:
- Spanish: `es`
- French: `fr`
- German: `de`
- Arabic: `ar`
- Chinese (Simplified): `zh-CN`
- Japanese: `ja`
- Russian: `ru`

See [ISO 639-1 codes](https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes) for your language.

### 3. Translate Strings

Edit your new file and translate all strings:

```toml
[welcome.message]
description = "Welcome message when user subscribes"
other = "Welcome to Jellyfin Notifications! You'll receive updates when new content is added."

# Translate to your language
other = "Your translated message here"
```

**Translation Guidelines**:
- Keep the tone friendly and welcoming
- Maintain formatting placeholders: `{{.Title}}`, `{{.Year}}`, etc.
- Preserve emoji if they make sense in your culture
- Test your translations with the bot

### 4. Test Your Translation

1. Build the bot with your translation
2. Change your language preference: send `/language` to bot
3. Verify all messages appear correctly
4. Test all commands: `/start`, `/recent`, `/search`
5. Test notifications with your language

### 5. Submit Your Translation

Create a PR with:
- The new translation file
- Updated README mentioning your language
- A note about testing you've done

**Example PR**:
```
Add Spanish translation

- Add active.es.toml with full Spanish translation
- Update README to list Spanish as supported language
- Tested all commands and notifications in Spanish
```

### Translation Help

If you need help with translations:
- Check existing translations for context
- Ask in GitHub Discussions
- Reference the English version for meaning
- Don't hesitate to ask questions!

## Reporting Bugs

### Before Reporting

1. **Check existing issues** - Your bug may already be reported
2. **Update to latest version** - Bug might be fixed
3. **Check documentation** - Might be configuration issue
4. **Enable debug logging** - Set `LOG_LEVEL=DEBUG`

### Bug Report Template

When reporting a bug, include:

**Description**: Clear description of the issue

**Steps to Reproduce**:
1. Configure bot with...
2. Send command...
3. Observe error...

**Expected Behavior**: What should happen

**Actual Behavior**: What actually happens

**Environment**:
- OS: Linux / Windows / macOS
- Bot Version: v1.2.3
- Installation Method: Docker / Binary / Source
- Go Version (if building from source): 1.22

**Logs**: Relevant log output (sanitize sensitive data!)

**Additional Context**: Screenshots, configuration (sanitize secrets!)

## Suggesting Features

### Good Feature Requests Include:

1. **Problem Description**: What problem are you trying to solve?
2. **Proposed Solution**: How should the feature work?
3. **Alternatives**: What other approaches did you consider?
4. **Benefits**: Why would this be valuable to users?
5. **Use Cases**: When would you use this feature?

### Feature Evaluation Criteria

Features are evaluated based on:
- **User value** - Does it solve a real problem?
- **Scope** - Does it fit the project's goals?
- **Maintenance** - Can we maintain it long-term?
- **Complexity** - Is it worth the implementation cost?
- **Breaking changes** - Does it break existing functionality?

## Getting Help

Need help contributing?

- **Questions**: [GitHub Discussions](https://github.com/yourusername/jellyfin-telegram-bot/discussions)
- **Bug Reports**: [GitHub Issues](https://github.com/yourusername/jellyfin-telegram-bot/issues)
- **Chat**: Join our community (if we have one)

## Recognition

Contributors are recognized in:
- GitHub contributors list
- Release notes (for significant contributions)
- Special thanks in README (for major features)

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

**Thank you for contributing to Jellyfin Telegram Bot!**

Your contributions help make this project better for everyone in the Jellyfin community. We appreciate your time and effort!
