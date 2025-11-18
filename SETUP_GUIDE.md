# Complete Setup & Testing Guide

## Prerequisites Checklist

- ✅ Bot binary built (`jellyfin-bot`)
- ⬜ Telegram Bot Token (get from @BotFather)
- ⬜ Jellyfin server with API access
- ⬜ Jellyfin Webhook plugin installed

---

## Step 1: Get Telegram Bot Token

1. Open Telegram and search for [@BotFather](https://t.me/BotFather)
2. Send `/newbot` command
3. Follow the prompts:
   - Choose a name (e.g., "My Jellyfin Bot")
   - Choose a username (must end in 'bot', e.g., "my_jellyfin_bot")
4. Copy the token that looks like: `123456789:ABCdefGHIjklMNOpqrsTUVwxyz`

**Save this token - you'll need it!**

---

## Step 2: Get Jellyfin API Key

1. Log into your Jellyfin server web interface
2. Go to **Dashboard** → **API Keys**
3. Click **+ (Add)** button
4. Enter a name like "Telegram Bot"
5. Copy the generated API key

**Save this key!**

---

## Step 3: Install Jellyfin Webhook Plugin

1. In Jellyfin Dashboard, go to **Plugins** → **Catalog**
2. Find and install **Webhook** plugin
3. Restart Jellyfin server
4. Go to **Dashboard** → **Plugins** → **Webhook**
5. Click **Add Generic Destination**
6. Configure:
   - **Webhook Name**: Telegram Bot
   - **Webhook Url**: `http://YOUR_BOT_SERVER:8080/webhook`
   - **Notification Type**: Select **Item Added**
   - **Item Type**: Select **Movies** and **Episodes**
   - **Request Content Type**: `application/json`
7. Save configuration

---

## Step 4: Configure the Bot

Create your `.env` file:

```bash
cat > .env <<'EOF'
# Telegram Configuration
TELEGRAM_BOT_TOKEN=YOUR_BOT_TOKEN_HERE

# Jellyfin Configuration
JELLYFIN_SERVER_URL=http://YOUR_JELLYFIN_SERVER:8096
JELLYFIN_API_KEY=YOUR_API_KEY_HERE

# Webhook Configuration
WEBHOOK_PORT=8080
WEBHOOK_SECRET=my-secret-key-123

# Database Configuration
DATABASE_PATH=./jellyfin_bot.db

# Logging Configuration (optional)
LOG_LEVEL=INFO
LOG_FILE=./logs/bot.log
EOF
```

**Replace these values:**
- `YOUR_BOT_TOKEN_HERE` - Token from @BotFather
- `YOUR_JELLYFIN_SERVER` - Your Jellyfin server IP/hostname
- `YOUR_API_KEY_HERE` - API key from Jellyfin
- `my-secret-key-123` - Optional security secret

---

## Step 5: Run the Bot

### Option A: Direct Run (for testing)

```bash
# Create logs directory
mkdir -p logs

# Run the bot
./jellyfin-bot
```

You should see output like:
```
INFO Connected to database path=./jellyfin_bot.db
INFO Jellyfin client initialized server_url=http://...
INFO Telegram bot initialized
INFO Webhook handler initialized
INFO Starting webhook server port=8080
INFO Bot is running. Press Ctrl+C to stop.
```

### Option B: Run in background

```bash
# Run in background and save logs
nohup ./jellyfin-bot > logs/bot.log 2>&1 &

# Check if it's running
ps aux | grep jellyfin-bot

# View logs
tail -f logs/bot.log
```

---

## Step 6: Test the Bot

### Test 1: Subscribe to Bot

1. Open Telegram
2. Search for your bot username (e.g., @my_jellyfin_bot)
3. Start a chat and send: `/start`
4. You should receive a Persian welcome message:

```
سلام! به ربات اطلاع‌رسانی جلیفین خوش آمدید.

شما از این پس اطلاعیه‌های محتوای جدید را دریافت خواهید کرد.

دستورات موجود:
/start - عضویت در ربات
/recent - مشاهده محتوای اخیر
/search - جستجوی محتوا
```

### Test 2: Check Recent Content

Send command: `/recent`

You should see a list of recently added movies/episodes with:
- Poster images
- Titles
- Descriptions
- Ratings

### Test 3: Search Content

Send command: `/search interstellar`

You should see search results for "Interstellar" (if it exists in your library).

### Test 4: Test Webhook (Manual)

You can test the webhook endpoint manually:

```bash
# Test webhook with a movie payload
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Secret: my-secret-key-123" \
  -d '{
    "NotificationType": "ItemAdded",
    "ItemType": "Movie",
    "ItemName": "Test Movie",
    "ItemId": "test-123",
    "Year": 2024,
    "Overview": "A test movie notification"
  }'
```

You should receive a notification in Telegram!

### Test 5: Add Real Content to Jellyfin

1. Add a new movie or episode to your Jellyfin library
2. Wait for Jellyfin to scan the library
3. Check Telegram - you should receive a notification!

---

## Troubleshooting

### Bot doesn't start

**Check logs:**
```bash
cat logs/bot.log
```

**Common issues:**
- Invalid Telegram token → Check .env file
- Jellyfin server unreachable → Verify JELLYFIN_SERVER_URL
- Database permission error → Check directory permissions

### Not receiving notifications

**Check webhook is configured:**
```bash
# Test webhook endpoint is accessible
curl http://localhost:8080/health
```

Should return: `{"status":"ok"}`

**Check Jellyfin webhook logs:**
- Go to Jellyfin Dashboard → Logs
- Look for webhook delivery attempts

### Commands don't work

**Check bot is subscribed:**
```bash
# Check database
sqlite3 jellyfin_bot.db "SELECT * FROM subscribers;"
```

---

## Production Deployment

For production use with systemd:

```bash
# Install as service
sudo cp deployments/systemd/jellyfin-bot.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable jellyfin-bot
sudo systemctl start jellyfin-bot

# Check status
sudo systemctl status jellyfin-bot

# View logs
sudo journalctl -u jellyfin-bot -f
```

---

## Security Notes

1. **Keep your bot token secret** - Never commit it to git
2. **Use WEBHOOK_SECRET** - Protects against unauthorized webhook calls
3. **Firewall rules** - Only allow webhook calls from your Jellyfin server
4. **HTTPS recommended** - For production, use nginx/caddy with SSL

---

## Quick Reference

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| TELEGRAM_BOT_TOKEN | Yes | - | Bot token from @BotFather |
| JELLYFIN_SERVER_URL | Yes | - | Jellyfin server URL |
| JELLYFIN_API_KEY | Yes | - | Jellyfin API key |
| WEBHOOK_PORT | No | 8080 | Webhook server port |
| WEBHOOK_SECRET | No | - | Optional webhook security |
| DATABASE_PATH | No | ./bot.db | SQLite database path |
| LOG_LEVEL | No | INFO | Log level (DEBUG/INFO/WARN/ERROR) |
| LOG_FILE | No | ./logs/bot.log | Log file path |

### Bot Commands

| Command | Description |
|---------|-------------|
| `/start` | Subscribe to notifications |
| `/recent` | View recently added content |
| `/search <query>` | Search for content |

### Useful Commands

```bash
# Check bot status
ps aux | grep jellyfin-bot

# Stop bot
pkill jellyfin-bot

# View real-time logs
tail -f logs/bot.log

# Check database
sqlite3 jellyfin_bot.db ".tables"
sqlite3 jellyfin_bot.db "SELECT COUNT(*) FROM subscribers;"

# Test webhook
curl http://localhost:8080/health
```

---

## Need Help?

- Check logs in `logs/bot.log`
- Review documentation in `docs/` folder
- Open an issue on GitHub
- Check Jellyfin webhook logs

---

**Enjoy your Persian-language Jellyfin notifications! 🎬📺**
