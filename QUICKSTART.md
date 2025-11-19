# ⚡ Quick Start Guide - Get Running in 5 Minutes!

## What You Need (Before Starting)

1. **Telegram Account** - Just download Telegram app
2. **5 minutes** - That's all!

---

## 🚀 Step-by-Step Setup

### Step 1: Get Your Telegram Bot Token (2 minutes)

1. **Open Telegram** and search for: `@BotFather`
2. **Send this message**: `/newbot`
3. **Answer the questions**:
   ```
   BotFather: Alright, a new bot. How are we going to call it?
   You: My Jellyfin Bot

   BotFather: Good. Now let's choose a username for your bot.
   You: my_jellyfin_bot
   ```
4. **Copy the token** - it looks like this:
   ```
   123456789:ABCdefGHIjklMNOpqrsTUVwxyz-123456789
   ```

✅ **You now have a bot token!**

---

### Step 2: Configure the Bot (1 minute)

Open the `.env` file and replace these 3 values:

```bash
# 1. Paste your bot token here:
TELEGRAM_BOT_TOKEN=123456789:ABCdefGHIjklMNOpqrsTUVwxyz-123456789

# 2. Your Jellyfin server URL:
JELLYFIN_SERVER_URL=http://192.168.1.100:8096

# 3. Your Jellyfin API key (get from Dashboard → API Keys):
JELLYFIN_API_KEY=abc123def456ghi789jkl012mno345pqr678
```

**Don't have Jellyfin yet?** You can still test the bot! Just use fake values for now.

---

### Step 3: Run the Bot (30 seconds)

```bash
# Create logs folder
mkdir -p logs

# Start the bot
./jellyfin-bot
```

You should see:
```
INFO Connected to database path=./jellyfin_bot.db
INFO Jellyfin client initialized
INFO Telegram bot initialized
INFO Starting webhook server port=8080
INFO Bot is running. Press Ctrl+C to stop.
```

✅ **Bot is running!**

---

### Step 4: Test It! (1 minute)

#### 4.1 Subscribe to Your Bot

1. Open Telegram
2. Search for your bot: `@my_jellyfin_bot` (use the username you chose)
3. Click **START** or send `/start`
4. You'll get a Persian welcome message! 🎉

#### 4.2 Test Commands

Send these commands to your bot:

- `/start` - Subscribe
- `/recent` - See recent content
- `/search interstellar` - Search for a movie

---

### Step 5: Test Notifications (1 minute)

**Open a new terminal** (keep the bot running) and run:

```bash
./test-webhook.sh
```

Choose option `1` to send a test movie notification.

**Check Telegram** - you should receive a notification! 📱

---

## 🎯 That's It! You're Done!

Your bot is now running and ready to receive Jellyfin notifications.

---

## 🔧 What's Next?

### Connect to Real Jellyfin

1. **Get Jellyfin API Key**:
   - Open Jellyfin web interface
   - Go to: Dashboard → API Keys → Add
   - Copy the key

2. **Install Webhook Plugin**:
   - Dashboard → Plugins → Catalog
   - Install "Webhook" plugin
   - Restart Jellyfin

3. **Configure Webhook**:
   - Dashboard → Plugins → Webhook → Add
   - URL: `http://YOUR_SERVER_IP:8080/webhook`
   - Events: Select "Item Added"
   - Types: Movies and Episodes

4. **Add Content**: Add a movie to Jellyfin and watch the magic! ✨

---

## 📊 Monitoring

### View Logs
```bash
tail -f logs/bot.log
```

### Check Database
```bash
sqlite3 jellyfin_bot.db "SELECT * FROM subscribers;"
```

### Test Webhook Endpoint
```bash
curl http://localhost:8080/health
# Should return: {"status":"ok"}
```

---

## 🛑 Stopping the Bot

Press `Ctrl+C` in the terminal where the bot is running.

Or if running in background:
```bash
pkill jellyfin-bot
```

---

## ❓ Troubleshooting

### "Bot is not responding"
- Check bot token is correct
- Make sure bot is running: `ps aux | grep jellyfin-bot`

### "No notifications"
- Send `/start` to the bot first
- Check logs: `tail -f logs/bot.log`
- Test with: `./test-webhook.sh`

### "Can't connect to Jellyfin"
- Verify server URL is reachable
- Check API key is correct
- Try: `curl http://YOUR_JELLYFIN_SERVER:8096`

---

## 📝 Summary

You now have:
- ✅ A working Telegram bot
- ✅ Persian language notifications
- ✅ Commands working (`/start`, `/recent`, `/search`)
- ✅ Webhook endpoint ready for Jellyfin

**Total time**: ~5 minutes
**Difficulty**: Easy 🟢

---

## 🎓 Learn More

- Full documentation: `SETUP_GUIDE.md`
- Architecture details: `docs/architecture.md`
- Production deployment: `docs/deployment.md`

---

**Enjoy your automated Jellyfin notifications! 🍿**
