# Testing Guide for Production Bots

This guide explains how to test new features without affecting existing users.

## 🎯 Strategy Overview

We use a **feature flag system** that enables beta features only for specific chat IDs (testers). This allows you to:
- Test new features in production safely
- Keep features hidden from regular users
- Send mock notifications to yourself only
- Verify everything works before releasing to all users

---

## 📋 Quick Start

### Step 1: Get Your Chat ID

Run the helper script:
```bash
./scripts/get-my-chat-id.sh
```

Or send `/start` to your bot and check the database:
```bash
sqlite3 bot.db "SELECT chat_id FROM subscribers;"
```

Your chat ID is a number like `123456789`.

### Step 2: Configure Testing Mode

Edit your `.env` file:
```bash
# Enable beta features
ENABLE_BETA_FEATURES=true

# Add your chat ID (can be multiple, comma-separated)
TESTER_CHAT_IDS=123456789
```

Restart your bot for changes to take effect.

### Step 3: Send Test Notifications

Use the test notification script:

**Test an episode notification:**
```bash
./scripts/send-test-notification.sh episode "Breaking Bad" 5 3
```

**Test a movie notification:**
```bash
./scripts/send-test-notification.sh movie "The Matrix"
```

**Custom examples:**
```bash
# Test series muting feature
./scripts/send-test-notification.sh episode "Stranger Things" 4 1
./scripts/send-test-notification.sh episode "Stranger Things" 4 2

# Test with Persian series names
./scripts/send-test-notification.sh episode "سریال تست" 1 1

# Test movies (should not show mute button)
./scripts/send-test-notification.sh movie "فیلم تست"
```

---

## 🔧 How Feature Flags Work

### In Your Code

The config provides an `IsTester()` method:

```go
// Check if user is a tester
if cfg.IsTester(chatID) {
    // Show beta features to this user
}
```

### Example: Applying to Mute Button

In `internal/telegram/notifications.go`:

```go
func shouldShowMuteButton(content *handlers.NotificationContent, cfg *config.Config, chatID int64) bool {
    // Only show for episodes
    if content.Type != "Episode" {
        return false
    }

    // Only show for valid series names
    if content.SeriesName == "" || content.SeriesName == "Unknown Series" {
        return false
    }

    // FEATURE FLAG: Only show to testers during beta
    if !cfg.IsTester(chatID) {
        return false
    }

    return true
}
```

This way, the mute button only appears for your chat ID, not for regular users!

---

## 🧪 Testing Workflow

### 1. Test Basic Functionality
```bash
# Send test notification
./scripts/send-test-notification.sh episode "Test Series" 1 1

# Check Telegram - you should see:
# - Notification with poster/text
# - "دنبال نکردن" button (only you see it!)
```

### 2. Test Mute Feature
```bash
# Click the mute button in Telegram
# You should see: "✓ شما دیگر اعلان‌های Test Series را دریافت نخواهید کرد"

# Send another episode
./scripts/send-test-notification.sh episode "Test Series" 1 2

# You should NOT receive this notification (you're muted)
```

### 3. Test Unmute Feature
```bash
# Run /mutedlist command in Telegram
# You should see "Test Series" with "رفع مسدودیت" button

# Click the unmute button
# You should see: "✓ Test Series از لیست مسدودی‌ها حذف شد"

# Send another episode
./scripts/send-test-notification.sh episode "Test Series" 1 3

# You SHOULD receive this notification (you're unmuted)
```

### 4. Verify Other Users Unaffected
- Have a friend subscribe to the bot (or use another account)
- Send test notifications
- Verify they don't see the mute button
- Verify they receive all notifications normally

---

## 🚀 Releasing Features to All Users

Once you've tested and verified everything works:

### Option 1: Remove Feature Flag (Full Release)

**In code** (e.g., `notifications.go`), remove the tester check:
```go
func shouldShowMuteButton(content *handlers.NotificationContent) bool {
    // Only show for episodes
    if content.Type != "Episode" {
        return false
    }

    // Only show for valid series names
    if content.SeriesName == "" || content.SeriesName == "Unknown Series" {
        return false
    }

    // REMOVED: Feature flag check - now enabled for everyone!

    return true
}
```

Rebuild and deploy:
```bash
go build -o jellyfin-bot cmd/bot/main.go
# Deploy to production
```

### Option 2: Gradual Rollout

Keep the feature flag but expand the tester list:
```bash
# Add more chat IDs gradually
TESTER_CHAT_IDS=123456789,987654321,555555555
```

Monitor for issues before full release.

### Option 3: Disable in Production

To disable a feature that's causing issues:
```bash
# In .env
ENABLE_BETA_FEATURES=false
```

Restart the bot - feature is hidden from everyone.

---

## 💡 Additional Testing Tips

### Test with Multiple Accounts

1. Create a second Telegram account
2. Subscribe to bot from both accounts
3. Only add YOUR account to `TESTER_CHAT_IDS`
4. Verify features work for you but not for the second account

### Test Persian/RTL Text

```bash
# Test with Persian series names
./scripts/send-test-notification.sh episode "سریال بازی تاج و تخت" 8 6

# Verify:
# - Text displays correctly (right-to-left)
# - Buttons work with Persian names
# - Database stores Persian characters correctly
```

### Test Edge Cases

```bash
# Empty/unknown series
./scripts/send-test-notification.sh episode "" 1 1
./scripts/send-test-notification.sh episode "Unknown Series" 1 1

# Special characters
./scripts/send-test-notification.sh episode "Series's \"Name\" & More!" 1 1

# Very long names
./scripts/send-test-notification.sh episode "This Is A Very Long Series Name That Might Cause Display Issues" 1 1
```

### Monitor Logs

While testing, watch the bot logs:
```bash
tail -f logs/bot.log
```

Look for:
- Notification sent/filtered messages
- Database operations
- Error messages
- Feature flag checks

---

## 📊 Monitoring Production

Even with testing, monitor production after release:

### Key Metrics to Watch

1. **Notification delivery rate** - Are notifications being sent?
2. **Database growth** - Are muted series being stored correctly?
3. **Error rates** - Any increase in errors?
4. **User engagement** - Are users using the new features?

### Log Queries

```bash
# Check muted series count
sqlite3 bot.db "SELECT COUNT(*) FROM muted_series;"

# Check which series are most muted
sqlite3 bot.db "SELECT series_name, COUNT(*) as mute_count FROM muted_series GROUP BY series_name ORDER BY mute_count DESC LIMIT 10;"

# Check subscribers with most muted series
sqlite3 bot.db "SELECT chat_id, COUNT(*) as muted_count FROM muted_series GROUP BY chat_id ORDER BY muted_count DESC LIMIT 10;"
```

---

## 🆘 Troubleshooting

### Feature Not Appearing for Me

1. Check `.env` configuration:
   ```bash
   cat .env | grep TESTER
   ```

2. Verify chat ID is correct:
   ```bash
   sqlite3 bot.db "SELECT chat_id FROM subscribers;"
   ```

3. Check bot logs for feature flag checks

4. Restart bot after changing `.env`

### Test Notification Not Received

1. Verify bot is running:
   ```bash
   ps aux | grep jellyfin-bot
   ```

2. Check webhook URL is correct:
   ```bash
   curl http://localhost:8080/health  # If you have health endpoint
   ```

3. Check bot logs for errors

4. Verify you're subscribed:
   ```bash
   sqlite3 bot.db "SELECT * FROM subscribers WHERE chat_id=YOUR_CHAT_ID;"
   ```

### Feature Appearing for All Users (Oops!)

1. Immediately disable:
   ```bash
   # In .env
   ENABLE_BETA_FEATURES=false
   ```

2. Restart bot

3. Fix code to add feature flag check

4. Re-test with feature flag enabled

---

## 📚 Alternative Testing Approaches

### Separate Test Bot (Recommended for Large Changes)

For major features, consider a completely separate test instance:

1. **Create a new bot** with @BotFather (separate token)
2. **Use test database**: `DATABASE_PATH=./test.db`
3. **Test freely** without affecting production
4. **Merge to production** after thorough testing

Example `.env.test`:
```bash
TELEGRAM_BOT_TOKEN=your_test_bot_token_here
DATABASE_PATH=./test.db
PORT=8081
ENABLE_BETA_FEATURES=true
```

Run test instance:
```bash
DATABASE_PATH=./test.db PORT=8081 ./jellyfin-bot
```

### Docker Compose for Staging

Use Docker Compose to run staging environment:

```yaml
# docker-compose.test.yml
services:
  bot-test:
    build: .
    environment:
      - TELEGRAM_BOT_TOKEN=${TEST_BOT_TOKEN}
      - DATABASE_PATH=/data/test.db
      - PORT=8081
    ports:
      - "8081:8081"
    volumes:
      - ./test-data:/data
```

---

## 🎓 Best Practices

1. **Always test before releasing** - Even small changes can have unexpected effects
2. **Keep tester list small** - Start with just yourself, expand gradually
3. **Monitor logs during testing** - Catch errors early
4. **Test edge cases** - Empty strings, special characters, long names
5. **Test with real data** - Use realistic series names and scenarios
6. **Document issues** - Keep notes on bugs found during testing
7. **Have a rollback plan** - Know how to disable features quickly

---

## 📞 Need Help?

If you encounter issues:

1. Check this guide first
2. Review bot logs for errors
3. Verify configuration in `.env`
4. Test with mock notifications
5. Check database state with SQLite queries

Happy testing! 🎉
