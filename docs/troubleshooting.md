# Troubleshooting Guide

This guide helps you diagnose and fix common issues with the Jellyfin Telegram Bot.

## Table of Contents

- [Quick Diagnostics](#quick-diagnostics)
- [Bot Not Starting](#bot-not-starting)
- [Bot Not Responding](#bot-not-responding)
- [Webhook Issues](#webhook-issues)
- [Notification Issues](#notification-issues)
- [Language and Translation Issues](#language-and-translation-issues)
- [Database Issues](#database-issues)
- [Docker-Specific Issues](#docker-specific-issues)
- [Performance Issues](#performance-issues)
- [Getting Help](#getting-help)

## Quick Diagnostics

Before diving into specific issues, run these quick checks:

### 1. Check Bot Status

```bash
# Docker
docker-compose ps

# Systemd
sudo systemctl status jellyfin-telegram-bot

# Manual process
ps aux | grep jellyfin-telegram-bot
```

### 2. Check Logs

```bash
# Docker
docker-compose logs -f --tail=100

# Systemd
sudo journalctl -u jellyfin-telegram-bot -f -n 100

# Manual (if LOG_FILE is set)
tail -f ./logs/bot.log
```

### 3. Test Connectivity

```bash
# Test webhook endpoint
curl http://localhost:8080/health
# Should return: {"status":"healthy"}

# Test Jellyfin API
curl http://your-jellyfin:8096/System/Info
```

### 4. Verify Configuration

```bash
# Check .env file exists and is readable
cat .env

# Verify required variables are set
grep -E "(TELEGRAM_BOT_TOKEN|JELLYFIN_SERVER_URL|JELLYFIN_API_KEY)" .env
```

---

## Bot Not Starting

### Issue: Bot exits immediately after starting

**Symptoms**:
- Bot process exits right after starting
- Error messages in logs
- Container restarts continuously (Docker)

**Common Causes**:

#### 1. Missing Required Environment Variables

**Error Message**:
```
FATAL: TELEGRAM_BOT_TOKEN is required
```

**Solution**:
```bash
# Check .env file has all required variables
cat .env

# Required variables:
TELEGRAM_BOT_TOKEN=your_token
JELLYFIN_SERVER_URL=http://your-server:8096
JELLYFIN_API_KEY=your_api_key

# Ensure no extra spaces around =
# Correct:   VARIABLE=value
# Incorrect: VARIABLE = value
```

#### 2. Invalid Telegram Bot Token

**Error Message**:
```
ERROR: Failed to initialize Telegram bot: Unauthorized
ERROR: Invalid bot token
```

**Solution**:
1. Verify token is correct (check for typos)
2. Get new token from [@BotFather](https://t.me/BotFather):
   - Send `/mybots`
   - Select your bot
   - Go to "API Token"
   - Copy the token

3. Update .env with correct token
4. Restart bot

#### 3. Cannot Connect to Jellyfin Server

**Error Message**:
```
ERROR: Failed to connect to Jellyfin server
ERROR: dial tcp: lookup jellyfin: no such host
```

**Solution**:
```bash
# Test Jellyfin URL is accessible
curl http://your-jellyfin:8096/System/Info

# If using Docker, ensure correct network:
# - Use 'jellyfin' if Jellyfin container is on same network
# - Use 'host.docker.internal:8096' for Jellyfin on host machine
# - Use actual IP address 'http://192.168.1.100:8096'

# Update .env with correct URL
JELLYFIN_SERVER_URL=http://correct-url:8096
```

#### 4. Invalid Jellyfin API Key

**Error Message**:
```
ERROR: Jellyfin API authentication failed: 401 Unauthorized
```

**Solution**:
1. Generate new API key in Jellyfin:
   - Dashboard → API Keys → Click **+**
   - Copy the new key

2. Update .env:
```env
JELLYFIN_API_KEY=new_api_key_here
```

3. Restart bot

#### 5. Port Already in Use

**Error Message**:
```
ERROR: Failed to start HTTP server: bind: address already in use
```

**Solution**:
```bash
# Check what's using port 8080
sudo lsof -i :8080
# or
sudo netstat -tlnp | grep 8080

# Option 1: Use different port
# Edit .env:
PORT=8081

# Option 2: Stop the conflicting service
sudo systemctl stop other-service
```

#### 6. Database Permission Issues

**Error Message**:
```
ERROR: Failed to open database: unable to open database file: permission denied
```

**Solution**:
```bash
# Check database directory exists
mkdir -p $(dirname ./bot.db)

# Fix permissions
chmod 755 $(dirname ./bot.db)
chmod 644 ./bot.db

# If running as systemd service, ensure user has access
sudo chown botuser:botuser ./bot.db
```

---

## Bot Not Responding

### Issue: Bot is running but doesn't respond to commands

**Symptoms**:
- Bot process is running
- Messages sent to bot are not answered
- No errors in logs

**Diagnostics**:

```bash
# 1. Check bot is actually running
ps aux | grep jellyfin-telegram-bot

# 2. Check logs for errors
tail -f ./logs/bot.log

# 3. Test bot directly in Telegram
# Send /start command
```

**Common Causes**:

#### 1. Bot is Not Connected to Telegram

**Check Logs For**:
```
INFO: Connected to Telegram successfully
INFO: Bot started successfully
```

**If Missing**:
- Check TELEGRAM_BOT_TOKEN is correct
- Check internet connection
- Check firewall allows outbound HTTPS (port 443)

**Solution**:
```bash
# Test internet connectivity
curl https://api.telegram.org/

# Check firewall
sudo ufw status
sudo ufw allow out 443/tcp

# Restart bot
sudo systemctl restart jellyfin-telegram-bot
```

#### 2. User Not Subscribed

**Issue**: Bot only responds to subscribed users in some configurations

**Solution**:
```bash
# Send /start to the bot to subscribe
# Check logs confirm subscription:
# INFO: User subscribed chat_id=123456789
```

#### 3. Bot Was Blocked by User

**Check**: Have you previously blocked this bot?

**Solution**:
1. Unblock the bot in Telegram
2. Send `/start` again

#### 4. Wrong Bot Token (Different Bot)

**Verify**: Are you messaging the correct bot?

**Solution**:
1. Check bot username matches your bot
2. Search for bot by exact username from @BotFather
3. Verify TELEGRAM_BOT_TOKEN in .env is for this bot

---

## Webhook Issues

### Issue: Jellyfin webhooks not being received

**Symptoms**:
- No notifications when adding content to Jellyfin
- Webhook endpoint returns errors
- Jellyfin webhook test fails

**Diagnostics**:

```bash
# 1. Test webhook endpoint
curl http://localhost:8080/health
# Should return: {"status":"healthy"}

# 2. Test webhook accepts POST
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -d '{"NotificationType":"ItemAdded","ItemType":"Movie","Name":"Test"}'
# Should return HTTP 200

# 3. Check bot is listening on correct port
sudo netstat -tlnp | grep 8080
```

**Common Causes**:

#### 1. Webhook URL Incorrect in Jellyfin

**Check Jellyfin Configuration**:
1. Dashboard → Plugins → Webhook
2. Verify URL format: `http://bot-server-ip:8080/webhook`

**Common Mistakes**:
- Using `https://` when bot uses `http://`
- Using `localhost` when bot is on different machine
- Wrong port number
- Missing `/webhook` path

**Solution**:
```bash
# If bot on same machine as Jellyfin:
http://localhost:8080/webhook

# If bot on different machine:
http://192.168.1.x:8080/webhook

# If bot in Docker on same host:
http://host.docker.internal:8080/webhook

# If bot in Docker on same network:
http://jellyfin-telegram-bot:8080/webhook
```

#### 2. Firewall Blocking Webhook Port

**Test from Jellyfin Server**:
```bash
# On the Jellyfin server, test bot endpoint
curl http://bot-ip:8080/health

# If this fails, firewall is blocking it
```

**Solution**:
```bash
# On bot server, allow incoming connections
sudo ufw allow 8080/tcp

# Verify firewall rule
sudo ufw status
```

#### 3. Webhook Secret Mismatch

**Error in Logs**:
```
WARN: Webhook authentication failed: invalid secret
```

**Solution**:
```bash
# Option 1: Remove webhook secret
# In .env:
WEBHOOK_SECRET=

# Option 2: Add secret to Jellyfin webhook
# In Jellyfin webhook config, add header:
# X-Webhook-Secret: same_value_as_in_env
```

#### 4. Jellyfin Webhook Plugin Not Configured

**Check**:
1. Dashboard → Plugins → Installed Plugins
2. Verify "Webhook" is installed
3. Go to Webhook settings
4. Check configuration:
   - ✓ Item Added notification type
   - ✓ Movies item type
   - ✓ Episodes item type

**Solution**:
```bash
# Install webhook plugin:
# 1. Dashboard → Plugins → Catalog
# 2. Search "Webhook"
# 3. Install
# 4. Restart Jellyfin
# 5. Configure webhook destination
```

#### 5. Bot Behind Reverse Proxy

**Issue**: Reverse proxy configuration incorrect

**Solution**:
```nginx
# nginx example
location /webhook {
    proxy_pass http://localhost:8080/webhook;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
}
```

#### 6. Docker Network Issues

**Issue**: Bot container not accessible from Jellyfin

**Solution**:
```yaml
# docker-compose.yml
version: '3.8'
services:
  bot:
    networks:
      - jellyfin-network
    ports:
      - "8080:8080"

networks:
  jellyfin-network:
    external: true  # If Jellyfin on same network
```

---

## Notification Issues

### Issue: Not receiving notifications for new content

**Symptoms**:
- Webhook is received (check logs)
- No notifications sent to Telegram
- No errors in logs

**Diagnostics**:

```bash
# 1. Enable debug logging
# Edit .env:
LOG_LEVEL=DEBUG

# 2. Restart bot
sudo systemctl restart jellyfin-telegram-bot

# 3. Add test content to Jellyfin
# Watch logs for detailed processing

# 4. Check subscription status
# Send /start to bot again
```

**Common Causes**:

#### 1. Not Subscribed to Bot

**Solution**:
```bash
# Send /start to the bot in Telegram
# Check logs for confirmation:
# INFO: User subscribed chat_id=123456789
```

#### 2. Series is Muted

**Issue**: You previously muted notifications for this series

**Solution**:
```bash
# Check muted series in database
sqlite3 bot.db "SELECT * FROM muted_series;"

# Unmute via bot:
# Send /recent
# Navigate to the series
# Click "Unmute" button
```

#### 3. Duplicate Notification Prevention

**Issue**: Content already notified (bot tracks this)

**Check Logs**:
```
DEBUG: Content already notified, skipping: Movie Title (2024)
```

**Solution**:
```bash
# This is normal behavior to prevent spam
# If you want to test, delete the content cache:
sqlite3 bot.db "DELETE FROM notified_content WHERE title='Movie Title';"

# Or delete entire cache (testing only):
sqlite3 bot.db "DELETE FROM notified_content;"
```

#### 4. Item Type Not Supported

**Issue**: Trying to notify for unsupported content type

**Supported Types**:
- Movies
- Episodes
- (Audio books and music not supported)

**Check Logs**:
```
DEBUG: Ignoring notification for unsupported type: Audio
```

**Solution**: This is expected behavior. Bot only notifies for video content.

#### 5. User Blocked Bot

**Error in Logs**:
```
ERROR: Failed to send notification: Forbidden: bot was blocked by the user
```

**Solution**:
1. Unblock bot in Telegram
2. Send `/start` to resubscribe

#### 6. Network Issues to Telegram API

**Error in Logs**:
```
ERROR: Failed to send message: dial tcp: lookup api.telegram.org: no such host
```

**Solution**:
```bash
# Check internet connectivity
curl https://api.telegram.org/

# Check DNS resolution
nslookup api.telegram.org

# Check firewall allows outbound HTTPS
sudo ufw allow out 443/tcp
```

---

## Language and Translation Issues

### Issue: Language not changing or wrong language

**Symptoms**:
- Bot doesn't switch language after `/language` command
- Always shows English despite preference
- Missing translations

**Diagnostics**:

```bash
# 1. Check translation files exist
ls -l locales/
# Should show: active.en.toml, active.fa.toml

# 2. Check database for language preference
sqlite3 bot.db "SELECT chat_id, language_code FROM subscribers;"

# 3. Enable debug logging
LOG_LEVEL=DEBUG
```

**Common Causes**:

#### 1. Translation Files Missing

**Error in Logs**:
```
ERROR: Failed to load translation file: locales/active.fa.toml
```

**Solution**:
```bash
# Verify files exist
ls -l locales/

# If using Docker, ensure locales are copied
docker exec -it jellyfin-telegram-bot ls -l /app/locales/

# Rebuild Docker image if necessary
docker-compose build --no-cache
```

#### 2. Language Preference Not Saved

**Issue**: Language resets after restart

**Check Database**:
```bash
sqlite3 bot.db "SELECT * FROM subscribers;"
# language_code column should show your preference
```

**Solution**:
```bash
# If column is missing, database needs migration
# Stop bot
# Delete old database (backup first!)
cp bot.db bot.db.backup
rm bot.db

# Start bot (it will recreate with correct schema)
```

#### 3. Invalid Language Code

**Issue**: Trying to use unsupported language

**Supported Languages**:
- `en` - English
- `fa` - Persian

**Solution**: Use `/language` command and select from available options

#### 4. Telegram Language Auto-detection Not Working

**Issue**: Bot doesn't detect your Telegram language

**Explanation**: This is expected if:
- Your Telegram language is not supported (only en/fa currently)
- You haven't set a language in Telegram settings

**Solution**:
```bash
# Manually select language:
# Send /language to bot
# Choose your preferred language

# Or change Telegram app language:
# Telegram Settings → Language
```

---

## Database Issues

### Issue: Database errors or corruption

**Symptoms**:
- Errors about database being locked
- Can't subscribe or unsubscribe
- Data not persisting

**Common Causes**:

#### 1. Database Locked

**Error**:
```
ERROR: database is locked
```

**Solution**:
```bash
# Check for multiple bot instances
ps aux | grep jellyfin-telegram-bot
# Should only show one instance

# If multiple, kill extras
kill <pid>

# Check for stale lock files
ls -la bot.db*
rm bot.db-shm bot.db-wal  # If exist and bot is stopped
```

#### 2. Database Corrupted

**Error**:
```
ERROR: database disk image is malformed
```

**Solution**:
```bash
# Backup database
cp bot.db bot.db.corrupted

# Try to recover
sqlite3 bot.db ".recover" | sqlite3 bot.db.recovered

# If recovery works
mv bot.db.recovered bot.db

# If recovery fails, start fresh (you'll lose subscriptions)
rm bot.db
# Bot will recreate on next start
```

#### 3. Disk Full

**Error**:
```
ERROR: disk I/O error
ERROR: database or disk is full
```

**Solution**:
```bash
# Check disk space
df -h

# Clean up if needed
# Delete old logs
find ./logs -name "*.log.*" -mtime +7 -delete

# Vacuum database to reclaim space
sqlite3 bot.db "VACUUM;"
```

#### 4. Permission Issues

**Error**:
```
ERROR: unable to open database file: permission denied
```

**Solution**:
```bash
# Fix permissions
chmod 644 bot.db
chmod 755 $(dirname bot.db)

# If running as systemd service
sudo chown botuser:botuser bot.db

# Verify
ls -l bot.db
```

---

## Docker-Specific Issues

### Issue: Docker container won't start or keeps restarting

**Diagnostics**:

```bash
# Check container status
docker-compose ps

# View container logs
docker-compose logs -f

# Check Docker daemon logs
sudo journalctl -u docker -f

# Inspect container
docker inspect jellyfin-telegram-bot
```

**Common Causes**:

#### 1. .env File Not Loaded

**Solution**:
```bash
# Verify .env file exists
ls -la .env

# Check docker-compose.yml references it
grep "env_file" docker-compose.yml

# Recreate container
docker-compose down
docker-compose up -d
```

#### 2. Volume Mount Issues

**Error**:
```
ERROR: Cannot create directory: Read-only file system
```

**Solution**:
```bash
# Check volume mounts in docker-compose.yml
docker-compose config

# Ensure host directories exist
mkdir -p ./data ./logs

# Fix permissions
chmod 755 ./data ./logs

# Recreate container
docker-compose down
docker-compose up -d
```

#### 3. Network Issues

**Issue**: Container can't reach Jellyfin

**Solution**:
```bash
# Test from inside container
docker-compose exec jellyfin-telegram-bot sh
wget -O- http://jellyfin:8096/System/Info

# If fails, check network configuration
docker network ls
docker network inspect network-name

# Use host network mode (less secure)
# In docker-compose.yml:
network_mode: host
```

#### 4. Image Build Failed

**Error during `docker-compose build`**:

**Solution**:
```bash
# Clear build cache
docker-compose build --no-cache

# Check Dockerfile syntax
docker-compose config

# Pull base image manually
docker pull golang:1.22-alpine
docker pull alpine:latest

# Retry build
docker-compose build
```

---

## Performance Issues

### Issue: Bot is slow or using excessive resources

**Symptoms**:
- High CPU usage
- High memory usage
- Slow response to commands
- Delayed notifications

**Diagnostics**:

```bash
# Check resource usage
top
htop

# Docker stats
docker stats jellyfin-telegram-bot

# Check log file size
du -h logs/bot.log

# Check database size
du -h bot.db
```

**Solutions**:

#### 1. Log File Too Large

```bash
# Check log file size
ls -lh logs/

# Logs should auto-rotate at 10MB
# If not, manually clean up
# Edit .env:
LOG_LEVEL=INFO  # Reduce verbosity

# Delete old logs
rm logs/bot.log.*

# Restart bot
```

#### 2. Database Growing Too Large

```bash
# Check database size
sqlite3 bot.db "SELECT COUNT(*) FROM notified_content;"

# Clean old entries (older than 30 days)
sqlite3 bot.db "DELETE FROM notified_content WHERE created_at < datetime('now', '-30 days');"

# Vacuum to reclaim space
sqlite3 bot.db "VACUUM;"
```

#### 3. Too Many API Calls

**Issue**: Polling Jellyfin too frequently

**Check**: Bot uses webhooks, not polling, so this shouldn't happen

**If it does**:
- Check logs for API call patterns
- Report issue on GitHub

#### 4. Memory Leak

**Symptoms**: Memory usage grows over time

**Temporary Solution**:
```bash
# Restart bot daily with cron
# Add to crontab: sudo crontab -e
0 3 * * * systemctl restart jellyfin-telegram-bot

# Or with Docker
0 3 * * * docker-compose restart
```

**Permanent Solution**: Report issue on GitHub with:
- Memory usage graph
- How long bot was running
- Number of subscribers
- Version information

---

## Getting Help

If this guide doesn't solve your issue:

### 1. Enable Debug Logging

```bash
# Edit .env
LOG_LEVEL=DEBUG

# Restart bot
sudo systemctl restart jellyfin-telegram-bot

# Reproduce the issue

# Collect logs
tail -n 200 logs/bot.log > debug.log
```

### 2. Gather Information

Collect:
- Bot version (`git describe --tags`)
- Operating system and version
- Installation method (Docker, binary, source)
- Go version (if building from source)
- Relevant configuration (sanitize secrets!)
- Error messages and logs (sanitize secrets!)
- Steps to reproduce the issue

### 3. Search Existing Issues

[Search GitHub Issues](https://github.com/yourusername/jellyfin-telegram-bot/issues) to see if someone else reported this.

### 4. Ask for Help

**GitHub Issues**: [Report a Bug](https://github.com/yourusername/jellyfin-telegram-bot/issues/new?template=bug_report.md)

**GitHub Discussions**: [Ask a Question](https://github.com/yourusername/jellyfin-telegram-bot/discussions)

### 5. Include Debug Information

When asking for help, include:

```markdown
## Environment
- OS: Ubuntu 22.04
- Bot Version: v1.2.3
- Installation: Docker
- Go Version: 1.22 (if building from source)

## Configuration (sanitized)
```env
JELLYFIN_SERVER_URL=http://192.168.1.100:8096
PORT=8080
LOG_LEVEL=DEBUG
```

## Steps to Reproduce
1. Start bot with configuration above
2. Add new movie to Jellyfin
3. Observe error in logs

## Expected Behavior
Should send notification to Telegram

## Actual Behavior
No notification sent, error in logs:
```
ERROR: Failed to send notification: some error
```

## Additional Context
- Worked fine until yesterday
- Changed nothing in configuration
- Jellyfin webhook test succeeds
```

---

**Still stuck?** Don't hesitate to ask for help. We're here to help you get the bot working!
