# Feature Flag Implementation Example

This document shows exactly how to apply feature flags to test new features before releasing to all users.

## Example: Protecting the Mute Button Feature

Here's how to make the mute button only appear for testers during the beta phase.

### Current Code (notifications.go)

```go
// shouldShowMuteButton checks if mute button should be shown for this content
func shouldShowMuteButton(content *NotificationContent) bool {
    // Only show for episodes, not movies
    if content.Type != "Episode" {
        return false
    }

    // Don't show if series name is empty or "Unknown Series"
    if content.SeriesName == "" || content.SeriesName == "Unknown Series" {
        return false
    }

    return true
}
```

### With Feature Flag (Beta Testing)

**Step 1: Update function signature to accept config and chatID**

```go
// shouldShowMuteButton checks if mute button should be shown for this content
func shouldShowMuteButton(content *NotificationContent, cfg *config.Config, chatID int64) bool {
    // Only show for episodes, not movies
    if content.Type != "Episode" {
        return false
    }

    // Don't show if series name is empty or "Unknown Series"
    if content.SeriesName == "" || content.SeriesName == "Unknown Series" {
        return false
    }

    // FEATURE FLAG: Only show to testers during beta
    if !cfg.IsTester(chatID) {
        slog.Debug("Mute button hidden for non-tester", "chat_id", chatID)
        return false
    }

    return true
}
```

**Step 2: Pass config and chatID when calling**

In the `BroadcastNotification` function, you'll need to:

1. Accept config as a parameter:
```go
func (b *Bot) BroadcastNotification(ctx context.Context, content *NotificationContent, cfg *config.Config) error {
```

2. Pass chatID when checking if button should show:
```go
// Inside the subscriber loop
for _, chatID := range filteredSubscribers {
    // ... existing code ...

    if shouldShowMuteButton(content, cfg, chatID) {
        keyboard = createMuteButton(content.SeriesName)
    }

    // ... send message with keyboard ...
}
```

**Step 3: Update caller to pass config**

In `cmd/bot/main.go` or wherever `BroadcastNotification` is called:

```go
// Load config
cfg, err := config.LoadConfig()
if err != nil {
    log.Fatalf("Failed to load config: %v", err)
}

// Pass config to broadcast
err = broadcaster.BroadcastNotification(ctx, content, cfg)
```

### After Testing is Complete (Full Release)

Once you've tested and verified everything works, simply remove the feature flag check:

```go
func shouldShowMuteButton(content *NotificationContent) bool {
    // Only show for episodes, not movies
    if content.Type != "Episode" {
        return false
    }

    // Don't show if series name is empty or "Unknown Series"
    if content.SeriesName == "" || content.SeriesName == "Unknown Series" {
        return false
    }

    // Feature flag check REMOVED - now enabled for everyone!
    return true
}
```

---

## Complete Implementation Pattern

### 1. Add Config to Bot Struct

```go
// In internal/telegram/bot.go
type Bot struct {
    bot            *bot.Bot
    db             SubscriberDB
    jellyfinClient JellyfinClient
    serverURL      string
    config         *config.Config  // Add this
}

// Update NewBot to accept config
func NewBot(token string, db SubscriberDB, jellyfinClient JellyfinClient, serverURL string, cfg *config.Config) (*Bot, error) {
    // ... existing code ...

    return &Bot{
        bot:            b,
        db:             db,
        jellyfinClient: jellyfinClient,
        serverURL:      serverURL,
        config:         cfg,  // Store config
    }, nil
}
```

### 2. Use Config in Methods

```go
func (b *Bot) BroadcastNotification(ctx context.Context, content *NotificationContent) error {
    // ... existing code ...

    for _, chatID := range filteredSubscribers {
        // Check feature flag using stored config
        if shouldShowMuteButton(content, b.config, chatID) {
            keyboard = createMuteButton(content.SeriesName)
        }

        // ... send notification ...
    }
}
```

---

## Alternative: Per-Feature Flags

For more granular control, you can add specific feature flags:

### In config.go

```go
type TestingConfig struct {
    TesterChatIDs          []int64
    EnableBetaFeatures     bool
    EnableMuteButton       bool  // Specific feature flag
    EnableSearchFilters    bool  // Another feature flag
    EnableAdminCommands    bool  // Yet another feature flag
}

func LoadConfig() (*Config, error) {
    // ... existing code ...

    Testing: TestingConfig{
        TesterChatIDs:       getEnvInt64Slice("TESTER_CHAT_IDS", []int64{}),
        EnableBetaFeatures:  getEnvBool("ENABLE_BETA_FEATURES", false),
        EnableMuteButton:    getEnvBool("ENABLE_MUTE_BUTTON", false),
        EnableSearchFilters: getEnvBool("ENABLE_SEARCH_FILTERS", false),
        EnableAdminCommands: getEnvBool("ENABLE_ADMIN_COMMANDS", false),
    },
}

// Helper method for specific feature
func (c *Config) CanUseMuteButton(chatID int64) bool {
    if !c.Testing.EnableMuteButton {
        return false
    }
    return c.IsTester(chatID)
}
```

### In .env

```bash
ENABLE_BETA_FEATURES=true
ENABLE_MUTE_BUTTON=true
ENABLE_SEARCH_FILTERS=false
ENABLE_ADMIN_COMMANDS=false
TESTER_CHAT_IDS=123456789
```

### In Code

```go
if b.config.CanUseMuteButton(chatID) {
    keyboard = createMuteButton(content.SeriesName)
}
```

---

## Testing Scenarios

### Scenario 1: Brand New Feature (Full Protection)

```bash
# .env
ENABLE_BETA_FEATURES=true
TESTER_CHAT_IDS=YOUR_CHAT_ID
```

- Only you see the feature
- Test thoroughly
- Other users are completely unaffected

### Scenario 2: Gradual Rollout

```bash
# .env - Week 1
TESTER_CHAT_IDS=YOUR_CHAT_ID

# .env - Week 2 (add trusted users)
TESTER_CHAT_IDS=YOUR_CHAT_ID,FRIEND_1,FRIEND_2

# .env - Week 3 (add more users)
TESTER_CHAT_IDS=YOUR_CHAT_ID,FRIEND_1,FRIEND_2,USER_3,USER_4,USER_5
```

- Gradually expand tester base
- Monitor for issues
- Collect feedback

### Scenario 3: Emergency Disable

```bash
# .env - Disable immediately if there's a bug
ENABLE_BETA_FEATURES=false
```

- Restart bot
- Feature is hidden from everyone
- Fix the bug
- Re-enable when ready

---

## Best Practices

### DO:
- ✅ Test with feature flags for significant new features
- ✅ Start with just your chat ID
- ✅ Monitor logs during beta testing
- ✅ Keep tester list small initially
- ✅ Remove feature flags after full release
- ✅ Document which features have flags

### DON'T:
- ❌ Forget to remove flags after release (code bloat)
- ❌ Add flags to every tiny change (over-engineering)
- ❌ Skip testing even with flags
- ❌ Leave broken features enabled for testers
- ❌ Forget to update tests when adding flags

---

## Quick Reference

### Get Your Chat ID
```bash
./scripts/get-my-chat-id.sh
```

### Send Test Notification
```bash
./scripts/send-test-notification.sh episode "Test Series" 1 1
```

### Check Tester Status
```go
if cfg.IsTester(chatID) {
    // This user is a tester
}
```

### Enable Beta for Testing
```bash
# .env
ENABLE_BETA_FEATURES=true
TESTER_CHAT_IDS=YOUR_CHAT_ID
```

### Disable All Beta Features
```bash
# .env
ENABLE_BETA_FEATURES=false
```

---

## Example: Adding Flag to Any New Feature

**Generic Template:**

```go
func shouldEnableNewFeature(cfg *config.Config, chatID int64) bool {
    // Your feature logic here
    if !meetsRequirements() {
        return false
    }

    // FEATURE FLAG: Only show to testers during beta
    if !cfg.IsTester(chatID) {
        return false
    }

    return true
}
```

**After beta testing:**

```go
func shouldEnableNewFeature(cfg *config.Config, chatID int64) bool {
    // Your feature logic here
    if !meetsRequirements() {
        return false
    }

    // REMOVED: Feature flag - now enabled for everyone!

    return true
}
```

---

## Summary

1. **Setup**: Add `TESTER_CHAT_IDS` and `ENABLE_BETA_FEATURES` to `.env`
2. **Code**: Add `cfg.IsTester(chatID)` check before showing feature
3. **Test**: Send mock notifications, verify only you see feature
4. **Release**: Remove feature flag check when ready
5. **Monitor**: Watch logs and user feedback after release

That's it! You now have a production-safe way to test new features. 🎉
