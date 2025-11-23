# Configuration Reference

This document provides comprehensive information about all configuration options for the Jellyfin Telegram Bot.

## Table of Contents

- [Configuration Methods](#configuration-methods)
- [Required Variables](#required-variables)
- [Optional Variables](#optional-variables)
- [Environment Variables by Category](#environment-variables-by-category)
- [Configuration Examples](#configuration-examples)
- [Security Best Practices](#security-best-practices)
- [Advanced Configuration](#advanced-configuration)

## Configuration Methods

The bot can be configured using environment variables, which can be set in several ways:

### 1. .env File (Recommended for Local Development)

Create a `.env` file in the project root directory:

```bash
cp .env.example .env
```

Edit `.env` with your values:

```env
TELEGRAM_BOT_TOKEN=123456789:ABCdefGHIjklMNOpqrsTUVwxyz
JELLYFIN_SERVER_URL=http://192.168.1.100:8096
JELLYFIN_API_KEY=a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6
```

### 2. Docker Environment Variables

When using Docker, set variables in your `docker-compose.yml`:

```yaml
environment:
  - TELEGRAM_BOT_TOKEN=your_token
  - JELLYFIN_SERVER_URL=http://jellyfin:8096
  - JELLYFIN_API_KEY=your_api_key
```

Or use an env_file:

```yaml
env_file:
  - .env
```

### 3. System Environment Variables

For systemd services or direct execution:

```bash
export TELEGRAM_BOT_TOKEN=your_token
export JELLYFIN_SERVER_URL=http://localhost:8096
export JELLYFIN_API_KEY=your_api_key
./jellyfin-telegram-bot
```

### 4. Systemd EnvironmentFile

In your systemd service file:

```ini
[Service]
EnvironmentFile=/etc/jellyfin-telegram-bot/.env
```

## Required Variables

These variables **must** be set for the bot to function:

### TELEGRAM_BOT_TOKEN

**Purpose**: Authentication token for your Telegram bot

**Required**: Yes

**Format**: `<bot_id>:<token_string>`

**How to Get**:
1. Open Telegram and message [@BotFather](https://t.me/BotFather)
2. Send `/newbot` command
3. Follow the prompts to create your bot
4. Copy the token provided

**Example**: `123456789:ABCdefGHIjklMNOpqrsTUVwxyz-123456789`

**Validation**: The bot will fail to start if this token is invalid or missing

**Security Note**: Never commit this token to version control or share it publicly

---

### JELLYFIN_SERVER_URL

**Purpose**: URL of your Jellyfin media server

**Required**: Yes

**Format**: `http(s)://hostname:port`

**Examples**:
- Local server: `http://localhost:8096`
- Network server: `http://192.168.1.100:8096`
- Remote server: `https://jellyfin.yourdomain.com`
- Docker internal: `http://jellyfin:8096`

**Default Port**: Jellyfin typically uses port 8096

**Notes**:
- Use `http://` for local/unencrypted connections
- Use `https://` for remote/encrypted connections
- If Jellyfin runs in Docker on same network, use container name
- URL must be accessible from where the bot runs

**Validation**: Bot will log error and fail to fetch content if unreachable

---

### JELLYFIN_API_KEY

**Purpose**: API key for authenticating with Jellyfin server

**Required**: Yes

**Format**: 32-character hexadecimal string

**How to Get**:
1. Log in to Jellyfin web interface
2. Go to Dashboard → Advanced → API Keys
3. Click the **+** button
4. Give it a descriptive name (e.g., "Telegram Bot")
5. Copy the generated key

**Example**: `a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6`

**Permissions**: The API key has full access to the Jellyfin server

**Security Note**: Treat this like a password. Never share it or commit it to version control

**Validation**: Bot will log 401 errors if key is invalid

---

## Optional Variables

These variables have sensible defaults and can be customized as needed:

### PORT

**Purpose**: Port number for the webhook HTTP server

**Required**: No

**Format**: Integer (1-65535)

**Default**: `8080`

**Example**: `8080`

**Usage**: The bot listens for Jellyfin webhook notifications on this port

**Jellyfin Configuration**: Set webhook URL to `http://bot-server-ip:PORT/webhook`

**Notes**:
- Must be available (not used by another service)
- Must be accessible from Jellyfin server
- Firewall rules must allow incoming connections
- Common alternatives: 3000, 5000, 8000, 9000

**Docker**: When using Docker, map this port in docker-compose.yml:
```yaml
ports:
  - "8080:8080"
```

---

### WEBHOOK_SECRET

**Purpose**: Secret token for validating webhook requests from Jellyfin

**Required**: No

**Format**: Any string (alphanumeric recommended)

**Default**: Empty (no validation)

**Example**: `my-super-secret-webhook-token-12345`

**Security Level**: Medium - prevents unauthorized webhook submissions

**How to Use**:
1. Set this variable in bot configuration
2. Add same value to Jellyfin webhook configuration as header:
   - Header Name: `X-Webhook-Secret`
   - Header Value: Same secret value

**Recommended**: Yes for production deployments

**Generation**: Use a random string generator:
```bash
# Linux/macOS
openssl rand -hex 32

# Or
cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1
```

**Validation**: If set, bot will reject webhooks without matching secret (HTTP 401)

---

### DATABASE_PATH

**Purpose**: Path to SQLite database file

**Required**: No

**Format**: File path (relative or absolute)

**Default**: `./bot.db`

**Examples**:
- Development: `./bot.db`
- Production: `/var/lib/jellyfin-telegram-bot/bot.db`
- Docker: `/app/data/bot.db`

**Notes**:
- Bot will create file if it doesn't exist
- Directory must exist and be writable
- Database stores subscriber information and content cache
- Backup this file to preserve user subscriptions

**Permissions**: File and directory must be readable/writable by bot user

**Migration**: If you move the database, copy the file to new location and update path

**Docker Volume**: Mount a volume to persist data:
```yaml
volumes:
  - ./data:/app/data
environment:
  - DATABASE_PATH=/app/data/bot.db
```

---

### LOG_LEVEL

**Purpose**: Controls verbosity of application logs

**Required**: No

**Format**: String (case-insensitive)

**Default**: `INFO`

**Options**:

- **DEBUG**: Very verbose, shows all operations
  - Use for: Development, troubleshooting issues
  - Logs: Every API call, database query, message sent
  - Performance: Slightly slower due to extra logging

- **INFO**: Normal operation logs
  - Use for: Production
  - Logs: Important events, errors, warnings
  - Recommended: Yes

- **WARN**: Only warnings and errors
  - Use for: Production with minimal logging
  - Logs: Problems that don't stop execution

- **ERROR**: Only errors
  - Use for: When you only care about failures
  - Logs: Critical errors that prevent functionality

**Example**: `LOG_LEVEL=DEBUG`

**When to Change**:
- Set to `DEBUG` when troubleshooting issues
- Set to `INFO` for production
- Set to `WARN` or `ERROR` to reduce log file size

**Performance Impact**: DEBUG level can generate large log files

---

### LOG_FILE

**Purpose**: Path to log file

**Required**: No

**Format**: File path (relative or absolute)

**Default**: `./logs/bot.log`

**Examples**:
- Development: `./logs/bot.log`
- Production (Linux): `/var/log/jellyfin-telegram-bot/bot.log`
- Docker: `/app/logs/bot.log`

**Features**:
- Automatic log rotation (prevents disk space issues)
- Maximum file size: 10 MB
- Number of backups: 3
- Old logs compressed automatically

**File Names**:
- Current: `bot.log`
- Rotated: `bot.log.1`, `bot.log.2`, `bot.log.3`

**Notes**:
- Directory must exist before starting bot
- Directory must be writable by bot user
- Logs also output to stdout/stderr (visible in `docker logs`)

**Creating Log Directory**:
```bash
# Create directory
mkdir -p logs

# Set permissions (Linux)
chmod 755 logs
```

**Docker Volume**:
```yaml
volumes:
  - ./logs:/app/logs
environment:
  - LOG_FILE=/app/logs/bot.log
```

---

### ENABLE_BETA_FEATURES

**Purpose**: Enable features that are still being tested

**Required**: No

**Format**: Boolean (`true` or `false`)

**Default**: `false`

**Example**: `ENABLE_BETA_FEATURES=true`

**Usage**: Developers use this to test new features before releasing to all users

**User Impact**: May expose unstable features or UI changes

**Recommendation**: Keep disabled unless you're testing specific features

**Related**: Works with `TESTER_CHAT_IDS` to limit beta features to specific users

---

### TESTER_CHAT_IDS

**Purpose**: Comma-separated list of Telegram chat IDs that can access beta features

**Required**: No (only relevant when `ENABLE_BETA_FEATURES=true`)

**Format**: Comma-separated integers (no spaces)

**Default**: Empty (no testers)

**Example**: `TESTER_CHAT_IDS=123456789,987654321`

**How to Get Your Chat ID**:
1. Send `/start` to the bot
2. Check bot logs - it will show your chat ID:
   ```
   INFO User started bot chat_id=123456789
   ```

**Usage**: Only users with chat IDs in this list will see beta features

**Multiple Users**: Separate IDs with commas (no spaces)

**Testing**: Useful for A/B testing or gradual feature rollouts

---

## Environment Variables by Category

### Telegram Bot Settings

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `TELEGRAM_BOT_TOKEN` | Yes | - | Bot authentication token |

### Jellyfin Integration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `JELLYFIN_SERVER_URL` | Yes | - | Jellyfin server URL |
| `JELLYFIN_API_KEY` | Yes | - | Jellyfin API key |

### Webhook Server

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `PORT` | No | `8080` | Webhook listener port |
| `WEBHOOK_SECRET` | No | (empty) | Webhook validation secret |

### Data Storage

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DATABASE_PATH` | No | `./bot.db` | SQLite database file path |

### Logging

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `LOG_LEVEL` | No | `INFO` | Log verbosity level |
| `LOG_FILE` | No | `./logs/bot.log` | Log file path |

### Feature Flags

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `ENABLE_BETA_FEATURES` | No | `false` | Enable experimental features |
| `TESTER_CHAT_IDS` | No | (empty) | Beta tester chat IDs |

---

## Configuration Examples

### Minimal Configuration (Development)

```env
# Required only
TELEGRAM_BOT_TOKEN=123456789:ABCdefGHIjklMNOpqrsTUVwxyz
JELLYFIN_SERVER_URL=http://localhost:8096
JELLYFIN_API_KEY=a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6
```

All other settings will use defaults.

---

### Production Configuration (Recommended)

```env
# Required
TELEGRAM_BOT_TOKEN=123456789:ABCdefGHIjklMNOpqrsTUVwxyz
JELLYFIN_SERVER_URL=https://jellyfin.yourdomain.com
JELLYFIN_API_KEY=a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6

# Security
WEBHOOK_SECRET=randomly-generated-secret-token-here

# Server
PORT=8080

# Data persistence
DATABASE_PATH=/opt/jellyfin-telegram-bot/data/bot.db

# Logging
LOG_LEVEL=INFO
LOG_FILE=/var/log/jellyfin-telegram-bot/bot.log
```

---

### Docker Configuration

```env
# Required
TELEGRAM_BOT_TOKEN=123456789:ABCdefGHIjklMNOpqrsTUVwxyz
JELLYFIN_SERVER_URL=http://jellyfin:8096
JELLYFIN_API_KEY=a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6

# Docker-friendly paths
DATABASE_PATH=/app/data/bot.db
LOG_FILE=/app/logs/bot.log

# Security
WEBHOOK_SECRET=my-webhook-secret
PORT=8080
```

With docker-compose.yml volumes:
```yaml
volumes:
  - ./data:/app/data
  - ./logs:/app/logs
```

---

### Testing/Debug Configuration

```env
# Required
TELEGRAM_BOT_TOKEN=123456789:ABCdefGHIjklMNOpqrsTUVwxyz
JELLYFIN_SERVER_URL=http://localhost:8096
JELLYFIN_API_KEY=a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6

# Verbose logging for debugging
LOG_LEVEL=DEBUG

# Enable beta features for testing
ENABLE_BETA_FEATURES=true
TESTER_CHAT_IDS=123456789

# Local development
PORT=8080
DATABASE_PATH=./test.db
LOG_FILE=./logs/debug.log
```

---

## Security Best Practices

### 1. Protect Your .env File

```bash
# Never commit .env to git
echo ".env" >> .gitignore

# Set restrictive permissions (Linux/macOS)
chmod 600 .env

# Verify it's in .gitignore
git check-ignore .env
# Should output: .env
```

### 2. Use Strong Secrets

Generate secure random secrets:

```bash
# For WEBHOOK_SECRET
openssl rand -hex 32
```

### 3. Secure File Permissions

```bash
# Database file (Linux)
chmod 600 bot.db

# Config file (Linux)
chmod 600 .env

# Log directory (Linux)
chmod 755 logs
chmod 644 logs/bot.log
```

### 4. Environment Variable Security

**DO**:
- Use `.env` files for local development
- Use systemd `EnvironmentFile` for services
- Use Docker secrets for production containers
- Rotate credentials periodically

**DON'T**:
- Hardcode secrets in code
- Commit `.env` to version control
- Share credentials in chat/email
- Use default/weak secrets

### 5. Network Security

```env
# Use HTTPS when possible
JELLYFIN_SERVER_URL=https://jellyfin.yourdomain.com

# Use webhook secret
WEBHOOK_SECRET=strong-random-secret

# Consider restricting port access via firewall
# Only allow Jellyfin server IP to access webhook port
```

### 6. Backup Configuration

```bash
# Backup your .env file securely
cp .env .env.backup
chmod 600 .env.backup

# Store backup in secure location (encrypted)
# Never store in git repository
```

---

## Advanced Configuration

### Running Multiple Bots

To run multiple bot instances (e.g., for different Jellyfin servers):

**Instance 1** (.env1):
```env
TELEGRAM_BOT_TOKEN=bot1_token
JELLYFIN_SERVER_URL=http://server1:8096
JELLYFIN_API_KEY=server1_api_key
PORT=8081
DATABASE_PATH=./bot1.db
LOG_FILE=./logs/bot1.log
```

**Instance 2** (.env2):
```env
TELEGRAM_BOT_TOKEN=bot2_token
JELLYFIN_SERVER_URL=http://server2:8096
JELLYFIN_API_KEY=server2_api_key
PORT=8082
DATABASE_PATH=./bot2.db
LOG_FILE=./logs/bot2.log
```

Run with different env files:
```bash
# Instance 1
env $(cat .env1 | xargs) ./jellyfin-telegram-bot &

# Instance 2
env $(cat .env2 | xargs) ./jellyfin-telegram-bot &
```

### Custom Paths for Docker

```yaml
version: '3.8'
services:
  bot:
    image: jellyfin-telegram-bot
    environment:
      - TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN}
      - JELLYFIN_SERVER_URL=${JELLYFIN_SERVER_URL}
      - JELLYFIN_API_KEY=${JELLYFIN_API_KEY}
      - DATABASE_PATH=/custom/path/bot.db
      - LOG_FILE=/custom/logs/bot.log
    volumes:
      - /host/custom/path:/custom/path
      - /host/custom/logs:/custom/logs
```

### Environment Variable Validation

The bot validates configuration on startup:

**Required Variables**: Bot exits with error if missing
**Optional Variables**: Uses defaults if not set
**Invalid Values**: Logs warning and uses default or exits

Example error messages:
```
FATAL: TELEGRAM_BOT_TOKEN is required
FATAL: JELLYFIN_SERVER_URL is required
FATAL: JELLYFIN_API_KEY is required
WARN: Invalid LOG_LEVEL 'VERBOSE', using default 'INFO'
WARN: PORT must be between 1-65535, using default 8080
```

---

## Troubleshooting Configuration Issues

### Bot Won't Start

**Check required variables are set**:
```bash
# Print all environment variables
env | grep -E "(TELEGRAM|JELLYFIN)"

# Check specific variable
echo $TELEGRAM_BOT_TOKEN
```

**Validate .env file is loaded**:
```bash
# Check .env exists
ls -la .env

# Check .env is not empty
cat .env

# Check for syntax errors (no spaces around =)
# Correct: VARIABLE=value
# Wrong:   VARIABLE = value
```

### Database Errors

**Check path and permissions**:
```bash
# Create directory
mkdir -p $(dirname ./bot.db)

# Check permissions
ls -l bot.db
chmod 600 bot.db

# Check directory is writable
touch test.db && rm test.db
```

### Webhook Not Receiving Requests

**Check PORT is correct**:
```bash
# Verify bot is listening
netstat -tlnp | grep 8080

# Test webhook endpoint
curl http://localhost:8080/health
```

**Check firewall**:
```bash
# Linux: Allow port
sudo ufw allow 8080/tcp

# Check if port is accessible from Jellyfin server
curl http://bot-server-ip:8080/health
```

### Log File Issues

**Create log directory**:
```bash
# Create directory
mkdir -p logs

# Check permissions
chmod 755 logs

# Test writing
touch logs/test.log && rm logs/test.log
```

---

## Getting Help

If you're having configuration issues:

1. **Check logs** for error messages
2. **Verify all required variables** are set
3. **Test connectivity** to Jellyfin server
4. **Review this documentation** for correct formats
5. **Open an issue** on GitHub with:
   - Error messages (sanitize secrets!)
   - Configuration (sanitize secrets!)
   - Environment (OS, Docker version, etc.)

**Resources**:
- [Troubleshooting Guide](troubleshooting.md)
- [GitHub Issues](https://github.com/yourusername/jellyfin-telegram-bot/issues)
- [GitHub Discussions](https://github.com/yourusername/jellyfin-telegram-bot/discussions)
