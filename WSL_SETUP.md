# WSL + Windows Jellyfin Setup Guide

## 🌐 Your Network Configuration

```
┌─────────────────────────────────────────────────────────┐
│                                                         │
│   Windows (Jellyfin Server)                            │
│   IP: 10.255.255.254                                   │
│   Port: 8096                                            │
│                                                         │
│         ↓ sends webhooks to ↓                          │
│                                                         │
│   WSL (Telegram Bot)                                   │
│   IP: 172.31.143.209                                   │
│   Port: 8080                                            │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

---

## ✅ Step 1: Get Jellyfin API Key

### In Windows (Jellyfin Dashboard):

1. Open Jellyfin in browser: `http://localhost:8096`
2. Login to your Jellyfin account
3. Go to: **Dashboard** → **API Keys** (left sidebar)
4. Click the **"+"** button (top right)
5. Enter a name: `Telegram Bot`
6. Click **OK**
7. **COPY the API key** - it looks like this:
   ```
   a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6
   ```

### Save the API Key:

In WSL, edit your `.env` file:

```bash
nano .env
```

Replace this line:
```
JELLYFIN_API_KEY=PASTE_YOUR_API_KEY_HERE
```

With your actual key:
```
JELLYFIN_API_KEY=a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6
```

Save (Ctrl+O, Enter, Ctrl+X)

---

## ✅ Step 2: Configure Jellyfin Webhook

### In Jellyfin Dashboard → Plugins → Webhook:

Click **"Add Generic Destination"** and configure:

### Required Settings:

| Field | Value |
|-------|-------|
| **Webhook Name** | `Telegram Bot` |
| **Webhook URL** | `http://172.31.143.209:8080/webhook` |
| **Notification Type** | ✅ **Item Added** (check this box) |
| **Item Type** | ✅ **Movies** and ✅ **Episodes** |
| **Status** | ✅ **Enable** (check this box) |

### Optional Settings:

- **Request Content Type**: `application/json` (default)
- **Webhook Secret**: `my-webhook-secret-123` (add as custom header if available)
  - Header name: `X-Webhook-Secret`
  - Header value: `my-webhook-secret-123`

### Important Notes:

⚠️  **Use WSL IP**: `172.31.143.209` (not localhost!)
⚠️  **Check "Item Added"** - This is the notification type
⚠️  **Check Movies and Episodes** - These are the content types

---

## ✅ Step 3: Restart the Bot

The bot is already running with the old config. Restart it:

```bash
# Stop the bot
pkill jellyfin-bot

# Wait a moment
sleep 2

# Start it again with new config
nohup ./jellyfin-bot > logs/bot.log 2>&1 &

# Check it's running
./monitor.sh
```

---

## ✅ Step 4: Test the Connection

### Test 1: Check if Jellyfin is reachable from WSL

```bash
# This should return Jellyfin's web page
curl -s http://10.255.255.254:8096 | head -20
```

If you see HTML output, ✅ connection works!

### Test 2: Test webhook from Windows

Open PowerShell or Command Prompt in Windows and run:

```powershell
curl -X POST http://172.31.143.209:8080/webhook `
  -H "Content-Type: application/json" `
  -H "X-Webhook-Secret: my-webhook-secret-123" `
  -d '{\"NotificationType\":\"ItemAdded\",\"ItemType\":\"Movie\",\"ItemName\":\"Test Movie\",\"ItemId\":\"test-123\"}'
```

Check your Telegram bot - you should get a notification!

---

## 🔥 Common Issues and Fixes

### Issue 1: "Connection refused" to Jellyfin

**Problem**: Bot can't reach Windows Jellyfin

**Fix**: Check Windows Firewall

1. Open Windows Firewall settings
2. Allow incoming connections on port 8096
3. Or temporarily disable firewall to test

### Issue 2: Jellyfin can't reach WSL webhook

**Problem**: Webhook returns error in Jellyfin

**Fix**: Check WSL firewall and bot status

```bash
# Check bot is running
ps aux | grep jellyfin-bot

# Check webhook is accessible
curl http://localhost:8080/health
```

### Issue 3: API Key doesn't work

**Problem**: Bot gets "Unauthorized" from Jellyfin

**Solutions**:
1. Verify API key is copied correctly (no extra spaces)
2. Make sure you created the API key in Jellyfin Dashboard
3. Try regenerating the API key

---

## 📊 Verify Everything is Working

Run this command in WSL:

```bash
./monitor.sh
```

You should see:
- ✅ Bot Status: RUNNING
- ✅ Webhook Health: OK
- 👥 Active Subscribers: 1+ (after you send /start to bot)

---

## 🎬 Final Test: Add Real Content

1. Add a movie or episode to your Jellyfin library
2. Wait for Jellyfin to scan (or force scan in Dashboard)
3. Check Telegram - you should get a notification!

---

## 📝 Quick Reference

### Your Configuration:

```bash
# WSL Bot IP (for Jellyfin webhook)
WSL_IP=172.31.143.209

# Windows Jellyfin IP (for bot to call)
JELLYFIN_IP=10.255.255.254

# Webhook URL (use in Jellyfin)
http://172.31.143.209:8080/webhook

# Jellyfin URL (already in .env)
http://10.255.255.254:8096
```

### Useful Commands:

```bash
# Restart bot
pkill jellyfin-bot && sleep 2 && nohup ./jellyfin-bot > logs/bot.log 2>&1 &

# Check status
./monitor.sh

# Watch logs
tail -f logs/bot.log

# Test from WSL
./test-webhook.sh

# Check Jellyfin connection
curl http://10.255.255.254:8096
```

---

## 🆘 Need Help?

1. Check logs: `tail -f logs/bot.log`
2. Test health: `curl http://localhost:8080/health`
3. Verify network: `ping 10.255.255.254`
4. Check Jellyfin logs in Dashboard → Logs

---

## ✅ Checklist

- [ ] Got Jellyfin API key from Dashboard
- [ ] Updated `.env` with API key
- [ ] Configured webhook in Jellyfin with WSL IP
- [ ] Restarted the bot
- [ ] Sent /start to bot in Telegram
- [ ] Tested webhook with `./test-webhook.sh`
- [ ] Added content to Jellyfin and got notification

---

**Once all checked, you're ready! 🎉**
