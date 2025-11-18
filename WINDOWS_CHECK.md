# Quick Windows Jellyfin Check

## Run These Commands in Windows PowerShell

### 1. Check if Jellyfin is running and what IP it's bound to:

```powershell
netstat -ano | findstr :8096
```

**What you should see:**
```
TCP    0.0.0.0:8096           0.0.0.0:0              LISTENING       <PID>
```

**If you see this instead (PROBLEM):**
```
TCP    127.0.0.1:8096         0.0.0.0:0              LISTENING       <PID>
```
This means Jellyfin is bound to localhost only!

---

### 2. Check your Windows IP addresses:

```powershell
ipconfig | findstr IPv4
```

Look for an IP starting with `10.`, `192.168.`, or `172.`

---

### 3. Test if you can access Jellyfin locally:

```powershell
curl http://localhost:8096
```

Should return HTML or redirect response.

---

### 4. Create/verify firewall rule:

```powershell
# Run PowerShell as Administrator first!
New-NetFirewallRule -DisplayName "Jellyfin WSL" -Direction Inbound -Protocol TCP -LocalPort 8096 -Action Allow -Profile Any
```

---

## If Jellyfin is bound to 127.0.0.1 (localhost only):

### Fix in Jellyfin Dashboard:

1. Open: http://localhost:8096
2. Login to Jellyfin
3. Go to: **Dashboard** (top right)
4. Click: **Networking** (left sidebar)
5. Look for: **"Bind to local network address"** or **"Server address settings"**
6. Options:
   - **Option A**: Set to `0.0.0.0` (bind to all interfaces)
   - **Option B**: Leave it empty/blank
   - **Option C**: Set to your specific Windows IP (from ipconfig)
7. Click **Save**
8. **Restart Jellyfin** (Dashboard → restart button OR restart Windows service)

### Alternative: Edit config file

1. Stop Jellyfin service
2. Edit: `C:\ProgramData\Jellyfin\Server\config\network.xml`
3. Find `<PublicHttpsPort>` section
4. Make sure you have:
   ```xml
   <PublicPort>8096</PublicPort>
   <EnableRemoteAccess>true</EnableRemoteAccess>
   ```
5. Save and restart Jellyfin

---

## After Making Changes:

### Test from Windows PowerShell:

```powershell
# Should work now with network IP
curl http://YOUR-WINDOWS-IP:8096
```

### Then test from WSL:

```bash
curl -I http://10.255.255.254:8096
```

---

## Full Diagnostic Script

I've created `diagnose-windows.ps1` - copy it to Windows and run:

```powershell
# In Windows PowerShell as Administrator
cd C:\path\to\script
.\diagnose-windows.ps1
```

This will check everything and tell you exactly what's wrong.

---

## Quick Summary

**Most likely issue:** Jellyfin is listening on `127.0.0.1:8096` instead of `0.0.0.0:8096`

**Quick fix:**
1. Jellyfin Dashboard → Networking
2. Set bind address to `0.0.0.0` or leave empty
3. Restart Jellyfin
4. Test from WSL: `curl http://10.255.255.254:8096`

✅ When working, you'll see HTTP response instead of "Connection refused"
