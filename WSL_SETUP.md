# WSL + Windows Jellyfin Setup Guide

## 🌐 Your Network Configuration

```
┌─────────────────────────────────────────────────────────┐
│                                                         │
│   Windows (Jellyfin Server)                            │
│   IP: <WINDOWS_IP>                                     │
│   Port: 8096                                            │
│                                                         │
│         ↓ sends webhooks to ↓                          │
│                                                         │
│   WSL (Telegram Bot)                                   │
│   IP: <WSL_IP>                                         │
│   Port: 8080                                            │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

**To find your IPs:**
- Windows IP: Run `ipconfig` in Windows and look for "vEthernet (WSL)"
- WSL IP: Run `hostname -I` in WSL terminal

---

## ✅ Step 1: Get Jellyfin API Key

### In Windows (Jellyfin Dashboard):

1. Open Jellyfin in browser: `http://localhost:8096`
2. Login to your Jellyfin account
3. Go to: **Dashboard** → **API Keys** (left sidebar)
4. Click the **"+"** button (top right)
5. Enter a name: `Telegram Bot`
6. Click **OK**
7. **COPY the API key** - it will be a long string of random characters

### Save the API Key:

In WSL, edit your `.env` file:

```bash
nano .env
```

Replace this line:
```
JELLYFIN_API_KEY=YOUR_JELLYFIN_API_KEY_HERE
```

With your actual key:
```
JELLYFIN_API_KEY=<paste_your_actual_api_key_here>
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
| **Webhook URL** | `http://<WSL_IP>:8080/webhook` |
| **Notification Type** | ✅ **Item Added** (check this box) |
| **Item Type** | ✅ **Movies** and ✅ **Episodes** |
| **Status** | ✅ **Enable** (check this box) |

**Replace `<WSL_IP>` with your actual WSL IP address** (get it by running `hostname -I` in WSL)

### Optional Settings:

- **Request Content Type**: `application/json` (default)
- **Webhook Secret**: `your-secret-token` (add as custom header if available)
  - Header name: `X-Webhook-Secret`
  - Header value: `your-secret-token`

### Important Notes:

⚠️  **Use WSL IP**: Not localhost! Get it with `hostname -I` in WSL
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
ps aux | grep jellyfin-bot
```

---

## ✅ Step 4: Test the Connection

### Test 1: Check if Jellyfin is reachable from WSL

```bash
# Replace <WINDOWS_IP> with your actual Windows IP
curl -s http://<WINDOWS_IP>:8096 | head -20
```

If you see HTML output, ✅ connection works!

### Test 2: Test webhook from Windows

Open PowerShell or Command Prompt in Windows and run:

```powershell
# Replace <WSL_IP> with your actual WSL IP
curl -X POST http://<WSL_IP>:8080/webhook `
  -H "Content-Type: application/json" `
  -H "X-Webhook-Secret: your-secret-token" `
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

Check the bot logs:

```bash
tail -f logs/bot.log
```

You should see:
- ✅ Bot connected to Telegram
- ✅ Webhook server running on port 8080
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
# Get with: hostname -I
WSL_IP=<your_wsl_ip>

# Windows Jellyfin IP (for bot to call)
# Get with: ipconfig in Windows, look for vEthernet (WSL)
WINDOWS_IP=<your_windows_ip>

# Webhook URL (use in Jellyfin)
http://<WSL_IP>:8080/webhook

# Jellyfin URL (already in .env)
http://<WINDOWS_IP>:8096
```

### Useful Commands:

```bash
# Find your WSL IP
hostname -I

# Restart bot
pkill jellyfin-bot && sleep 2 && nohup ./jellyfin-bot > logs/bot.log 2>&1 &

# Check status
ps aux | grep jellyfin-bot

# Watch logs
tail -f logs/bot.log

# Test webhook locally
./scripts/send-test-notification.sh

# Check Jellyfin connection
curl http://<WINDOWS_IP>:8096
```

---

## 🆘 Need Help?

1. Check logs: `tail -f logs/bot.log`
2. Test health: `curl http://localhost:8080/health`
3. Verify network: `ping <WINDOWS_IP>`
4. Check Jellyfin logs in Dashboard → Logs

---

## ✅ Checklist

- [ ] Got Jellyfin API key from Dashboard
- [ ] Updated `.env` with API key
- [ ] Found WSL IP with `hostname -I`
- [ ] Configured webhook in Jellyfin with WSL IP
- [ ] Restarted the bot
- [ ] Sent /start to bot in Telegram
- [ ] Tested webhook with `./scripts/send-test-notification.sh`
- [ ] Added content to Jellyfin and got notification

---

**Once all checked, you're ready! 🎉**
