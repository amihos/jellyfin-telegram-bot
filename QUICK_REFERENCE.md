# Quick Reference Card

## 📍 Your Network Setup

```
WSL (Bot Server)           Windows (Jellyfin)
172.31.143.209:8080   →   10.255.255.254:8096
```

## ⚡ Essential Commands

```bash
# Verify setup status
./verify-setup.sh

# Monitor bot
./monitor.sh

# Check logs
tail -f logs/bot.log

# Restart bot
pkill jellyfin-bot && sleep 2 && nohup ./jellyfin-bot > logs/bot.log 2>&1 &

# Test webhook locally
./test-webhook.sh
```

## ✅ Configuration Checklist

### 1. Jellyfin API Key (❌ NOT DONE)

**Get the key:**
1. Open: http://10.255.255.254:8096
2. Go to: **Dashboard** → **API Keys**
3. Click **"+"** button (top right)
4. Name: `Telegram Bot`
5. **Copy the key** (looks like: `a1b2c3d4e5f6...`)

**Add to .env:**
```bash
nano .env
```

Replace line 6:
```
JELLYFIN_API_KEY=PASTE_YOUR_API_KEY_HERE
```

With:
```
JELLYFIN_API_KEY=your-actual-key-here
```

Save: Ctrl+O, Enter, Ctrl+X

### 2. Jellyfin Webhook (❌ NOT DONE)

**Configure in Jellyfin:**
1. Go to: **Dashboard** → **Plugins** → **Webhook**
2. Click **"Add Generic Destination"**
3. Configure:

| Setting | Value |
|---------|-------|
| Webhook Name | `Telegram Bot` |
| Webhook URL | `http://172.31.143.209:8080/webhook` |
| Notification Type | ✅ **Item Added** |
| Item Type | ✅ **Movies** and ✅ **Episodes** |
| Status | ✅ **Enable** |

4. Click **Save**

### 3. Restart Bot (After configuring above)

```bash
pkill jellyfin-bot
sleep 2
nohup ./jellyfin-bot > logs/bot.log 2>&1 &
./monitor.sh
```

### 4. Test in Telegram

1. Find your bot: `@jhuso_jellyfin_bot`
2. Send: `/start`
3. You should see: "سلام! به ربات اطلاع‌رسانی جلیفین خوش آمدید"

## 🔧 Troubleshooting

### Bot shows "Unauthorized" errors

**Check token:**
```bash
grep TELEGRAM_BOT_TOKEN .env
```

If wrong, update `.env` and restart bot.

### Can't reach Jellyfin from WSL

**Test connection:**
```bash
curl http://10.255.255.254:8096
```

**If fails:**
- Check Windows Firewall allows port 8096
- Verify Jellyfin is running on Windows
- Ping the Windows host: `ping 10.255.255.254`

### Jellyfin can't send webhooks to WSL

**Verify bot is listening:**
```bash
curl http://localhost:8080/health
```

Should return: `{"status":"ok"}`

**If fails:**
- Check bot is running: `ps aux | grep jellyfin-bot`
- Check logs: `tail -f logs/bot.log`

### No notifications received

**Checklist:**
1. ✅ API key configured in .env
2. ✅ Bot restarted after .env changes
3. ✅ Webhook configured in Jellyfin with correct WSL IP
4. ✅ Sent /start to bot in Telegram
5. ✅ Added actual content to Jellyfin (not just test webhook)
6. ✅ Content is Movie or Episode (not music, etc.)

## 📱 Bot Commands

| Command | Description | Example |
|---------|-------------|---------|
| `/start` | Subscribe to notifications | `/start` |
| `/recent` | Show 5 latest items | `/recent` |
| `/search` | Search for content | `/search interstellar` |

## 🎬 Testing

### Test 1: Webhook from Windows PowerShell

```powershell
curl -X POST http://172.31.143.209:8080/webhook `
  -H "Content-Type: application/json" `
  -H "X-Webhook-Secret: my-webhook-secret-123" `
  -d '{\"NotificationType\":\"ItemAdded\",\"ItemType\":\"Movie\",\"ItemName\":\"Test Movie\",\"ItemId\":\"test-123\"}'
```

### Test 2: Add Real Content

1. Add a movie/episode to Jellyfin
2. Wait for scan (or force scan in Dashboard)
3. Check Telegram for notification

## 📊 Status Indicators

| Indicator | Meaning |
|-----------|---------|
| ✅ Bot Status: RUNNING | Bot process is active |
| ✅ Webhook Health: OK | Webhook endpoint responding |
| 👥 Active Subscribers: 1+ | Users subscribed |
| ❌ Unauthorized errors | Wrong/expired Telegram token |

## 🆘 Get Help

1. Run verification: `./verify-setup.sh`
2. Check detailed guide: `WSL_SETUP.md`
3. Monitor live logs: `tail -f logs/bot.log`
4. Check Jellyfin logs in Dashboard

---

**Current Status:** Run `./verify-setup.sh` to see what needs to be configured.
