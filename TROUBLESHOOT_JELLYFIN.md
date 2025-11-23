# Jellyfin Connection Troubleshooting

## Current Issue

✅ **Ping works**: Windows host (<WINDOWS_IP>) is reachable
❌ **Port 8096 blocked**: Jellyfin not accepting connections from WSL

```
ping <WINDOWS_IP>        → Success
curl <WINDOWS_IP>:8096   → Connection refused
```

## Why This Happens

Jellyfin might be configured to listen only on `localhost` (127.0.0.1) instead of all network interfaces. This means it only accepts connections from the local Windows machine, not from WSL.

---

## Solution 1: Check Jellyfin Network Settings (Recommended)

### On Windows, in Jellyfin:

1. Open Jellyfin: http://localhost:8096
2. Go to: **Dashboard** → **Networking**
3. Look for: **"Bind to local network address"** or **"LAN Networks"**
4. Make sure it's set to:
   - **Bind to**: `0.0.0.0:8096` (all interfaces)
   - **OR** Leave empty to bind to all interfaces
5. **Save** and **Restart Jellyfin**

### Alternative: Check config file

1. Open: `C:\ProgramData\Jellyfin\Server\config\network.xml`
2. Find the `<PublicPort>` and `<LocalNetworkAddresses>` sections
3. Make sure it includes your network range

Example:
```xml
<NetworkConfiguration>
  <PublicPort>8096</PublicPort>
  <LocalNetworkAddresses>
    <string>10.0.0.0/8</string>
    <string>172.16.0.0/12</string>
  </LocalNetworkAddresses>
</NetworkConfiguration>
```

---

## Solution 2: Verify Firewall Rule

### Check the firewall rule you created:

1. Open: **Windows Defender Firewall with Advanced Security**
2. Click: **Inbound Rules**
3. Find: Your Jellyfin rule (or create new)
4. Make sure:
   - ✅ Protocol: **TCP**
   - ✅ Port: **8096**
   - ✅ Action: **Allow the connection**
   - ✅ Profile: **All** (Domain, Private, Public)
   - ✅ Enabled: **Yes**

### Create new rule if needed:

```powershell
# Run in PowerShell as Administrator
New-NetFirewallRule -DisplayName "Jellyfin WSL Access" -Direction Inbound -Protocol TCP -LocalPort 8096 -Action Allow
```

---

## Solution 3: Check Jellyfin is Running

### In Windows PowerShell:

```powershell
# Check if Jellyfin is running
Get-Process | Where-Object {$_.ProcessName -like "*jellyfin*"}

# Check if port 8096 is listening
netstat -ano | findstr :8096
```

You should see:
```
TCP    0.0.0.0:8096           0.0.0.0:0              LISTENING       <PID>
```

If you see `127.0.0.1:8096` instead of `0.0.0.0:8096`, Jellyfin is only listening on localhost!

---

## Solution 4: Temporarily Disable Firewall (Testing Only)

**⚠️ ONLY for testing - not recommended for production!**

1. Open: **Windows Defender Firewall**
2. Click: **Turn Windows Defender Firewall on or off**
3. Select: **Turn off** (for Private network only)
4. Test connection from WSL: `curl http://<WINDOWS_IP>:8096`
5. **Turn firewall back on** after testing

---

## After Making Changes

### 1. Restart Jellyfin

- Windows Services → Jellyfin Server → Restart
- **OR** Restart from Jellyfin Dashboard → Server → Restart

### 2. Test from WSL

```bash
# Test connection
curl -I http://<WINDOWS_IP>:8096

# Should see HTTP/1.1 302 Found or similar
```

### 3. Check bot can connect

```bash
# Send test webhook with new movie ID
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Secret: my-webhook-secret-123" \
  -d '{
    "NotificationType": "ItemAdded",
    "ItemType": "Movie",
    "ItemName": "تست فیلم",
    "ItemId": "test-new-'$(date +%s)'",
    "Year": 2024,
    "Overview": "تست اتصال"
  }'

# Check logs - should fetch poster successfully now
tail -10 logs/bot.log
```

---

## Quick Test Commands

### From Windows PowerShell:

```powershell
# Test if Jellyfin responds locally
curl http://localhost:8096

# Check what IP Jellyfin is bound to
netstat -ano | findstr :8096

# Check firewall rules
Get-NetFirewallRule | Where-Object {$_.DisplayName -like "*Jellyfin*"}
```

### From WSL:

```bash
# Test if host is reachable
ping <WINDOWS_IP>

# Test if port 8096 is open
curl -v http://<WINDOWS_IP>:8096

# Test alternative: use telnet
telnet <WINDOWS_IP> 8096
```

---

## Expected Output When Working

### From WSL:

```bash
$ curl -I http://<WINDOWS_IP>:8096
HTTP/1.1 302 Found
Location: /web/index.html
...
```

### From bot logs:

```
{"level":"INFO","msg":"Successfully fetched poster image","item_id":"..."}
{"level":"INFO","msg":"Notification sent successfully","chat_id":...}
```

---

## Most Common Fix

**99% of the time, the issue is:**

Jellyfin is listening on `127.0.0.1:8096` (localhost only) instead of `0.0.0.0:8096` (all interfaces).

**Fix:** Go to Jellyfin Dashboard → Networking → Set bind address to `0.0.0.0` or leave empty → Restart Jellyfin

---

## Still Not Working?

1. Check Windows event logs for Jellyfin errors
2. Verify your Windows IP address:
   ```powershell
   ipconfig
   ```
3. Try accessing from Windows browser: http://<WINDOWS_IP>:8096
4. Check Jellyfin logs: `C:\ProgramData\Jellyfin\Server\log\`

---

**After fixing, run from WSL:**
```bash
./verify-setup.sh
```

All checks should be ✅ green!
