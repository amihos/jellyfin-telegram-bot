# Running as a Windows Service

This guide explains how to run the Jellyfin Telegram Bot as a Windows service, so it starts automatically when your computer boots and runs in the background.

## Table of Contents

- [Method 1: Using NSSM (Recommended)](#method-1-using-nssm-recommended)
- [Method 2: Using Task Scheduler](#method-2-using-task-scheduler)
- [Method 3: Using sc.exe](#method-3-using-scexe)
- [Managing the Service](#managing-the-service)
- [Troubleshooting](#troubleshooting)

## Prerequisites

- Windows 10/11 or Windows Server 2016+
- Administrator access
- Bot binary downloaded for Windows

## Method 1: Using NSSM (Recommended)

NSSM (Non-Sucking Service Manager) is the easiest way to run programs as Windows services.

### 1. Download NSSM

1. Download NSSM from [nssm.cc](https://nssm.cc/download)
2. Extract the ZIP file
3. Copy `nssm.exe` from the `win64` folder to a permanent location (e.g., `C:\Program Files\nssm\nssm.exe`)

Or use Chocolatey:
```powershell
choco install nssm
```

### 2. Download Bot Binary

1. Download the latest Windows release from [GitHub Releases](https://github.com/yourusername/jellyfin-telegram-bot/releases)
2. Create a folder: `C:\Program Files\JellyfinTelegramBot\`
3. Copy `jellyfin-telegram-bot-windows-amd64.exe` to this folder
4. Rename to `jellyfin-telegram-bot.exe`

### 3. Create .env File

Create `C:\Program Files\JellyfinTelegramBot\.env` with your configuration:

```env
TELEGRAM_BOT_TOKEN=your_telegram_bot_token
JELLYFIN_SERVER_URL=http://your-jellyfin:8096
JELLYFIN_API_KEY=your_jellyfin_api_key
PORT=8080
DATABASE_PATH=C:\Program Files\JellyfinTelegramBot\data\bot.db
LOG_FILE=C:\Program Files\JellyfinTelegramBot\logs\bot.log
LOG_LEVEL=INFO
```

**Important**: Use full Windows paths with backslashes (`C:\path\to\file`)

### 4. Create Required Folders

```powershell
# Create directories
New-Item -Path "C:\Program Files\JellyfinTelegramBot\data" -ItemType Directory -Force
New-Item -Path "C:\Program Files\JellyfinTelegramBot\logs" -ItemType Directory -Force
```

### 5. Install as Service with NSSM

Open PowerShell or Command Prompt **as Administrator**:

```powershell
# Navigate to NSSM directory
cd "C:\Program Files\nssm"

# Install service
.\nssm.exe install JellyfinTelegramBot "C:\Program Files\JellyfinTelegramBot\jellyfin-telegram-bot.exe"
```

This opens the NSSM configuration window:

**Application Tab**:
- Path: `C:\Program Files\JellyfinTelegramBot\jellyfin-telegram-bot.exe`
- Startup directory: `C:\Program Files\JellyfinTelegramBot`
- Arguments: (leave empty)

**Details Tab**:
- Display name: `Jellyfin Telegram Bot`
- Description: `Telegram bot for Jellyfin media server notifications`
- Startup type: `Automatic`

**Log on Tab**:
- Select: `Local System account` (or use a specific user if needed)

**Environment Tab**:
- Click "Choose file" and select `C:\Program Files\JellyfinTelegramBot\.env`

**I/O Tab**:
- Output (stdout): `C:\Program Files\JellyfinTelegramBot\logs\stdout.log`
- Error (stderr): `C:\Program Files\JellyfinTelegramBot\logs\stderr.log`

Click **Install service**.

### 6. Start the Service

```powershell
# Start service
.\nssm.exe start JellyfinTelegramBot

# Or use Windows Services
services.msc
# Find "Jellyfin Telegram Bot", right-click, Start
```

### 7. Verify Service is Running

```powershell
# Check service status
.\nssm.exe status JellyfinTelegramBot

# Or check Windows Services
services.msc

# Test webhook endpoint
curl http://localhost:8080/health
```

### Managing NSSM Service

```powershell
# Start service
nssm start JellyfinTelegramBot

# Stop service
nssm stop JellyfinTelegramBot

# Restart service
nssm restart JellyfinTelegramBot

# Check status
nssm status JellyfinTelegramBot

# Edit service configuration
nssm edit JellyfinTelegramBot

# Remove service
nssm remove JellyfinTelegramBot confirm
```

## Method 2: Using Task Scheduler

Task Scheduler can run programs at startup without creating a formal Windows service.

### 1. Prepare Bot Files

1. Create folder: `C:\JellyfinTelegramBot\`
2. Copy bot executable to this folder
3. Create `.env` file with configuration (use full Windows paths)
4. Create `data` and `logs` folders

### 2. Create Batch Script

Create `C:\JellyfinTelegramBot\start-bot.bat`:

```batch
@echo off
cd /d "C:\JellyfinTelegramBot"
jellyfin-telegram-bot.exe
```

### 3. Create Scheduled Task

1. Open **Task Scheduler** (search in Start menu)
2. Click **Create Basic Task**
3. Name: `Jellyfin Telegram Bot`
4. Trigger: **When the computer starts**
5. Action: **Start a program**
6. Program: `C:\JellyfinTelegramBot\start-bot.bat`
7. Finish the wizard

### 4. Configure Advanced Settings

1. Right-click the task → **Properties**
2. **General** tab:
   - Check: **Run whether user is logged on or not**
   - Check: **Run with highest privileges**
3. **Triggers** tab:
   - Edit trigger
   - Check: **Enabled**
4. **Settings** tab:
   - Check: **Allow task to be run on demand**
   - Uncheck: **Stop the task if it runs longer than**
5. Click **OK**

### 5. Test the Task

1. Right-click task → **Run**
2. Check if bot starts
3. Open browser: `http://localhost:8080/health`

### Managing Task Scheduler Task

```powershell
# Start task
schtasks /run /tn "Jellyfin Telegram Bot"

# Stop task (kills the process)
taskkill /F /IM jellyfin-telegram-bot.exe

# Disable task
schtasks /change /tn "Jellyfin Telegram Bot" /disable

# Enable task
schtasks /change /tn "Jellyfin Telegram Bot" /enable

# Delete task
schtasks /delete /tn "Jellyfin Telegram Bot" /f
```

## Method 3: Using sc.exe

This method uses Windows' built-in service control manager. **Note**: Most Go programs don't natively support Windows service APIs, so NSSM is recommended instead.

If your bot is compiled with Windows service support:

```powershell
# Create service
sc.exe create JellyfinTelegramBot binPath= "C:\Program Files\JellyfinTelegramBot\jellyfin-telegram-bot.exe" start= auto

# Start service
sc.exe start JellyfinTelegramBot

# Stop service
sc.exe stop JellyfinTelegramBot

# Delete service
sc.exe delete JellyfinTelegramBot
```

## Managing the Service

### Using Windows Services Manager

1. Press `Win + R`
2. Type `services.msc` and press Enter
3. Find "Jellyfin Telegram Bot" in the list
4. Right-click for options:
   - **Start** - Start the service
   - **Stop** - Stop the service
   - **Restart** - Restart the service
   - **Properties** - Configure service settings

### Using PowerShell

```powershell
# Get service status
Get-Service -Name JellyfinTelegramBot

# Start service
Start-Service -Name JellyfinTelegramBot

# Stop service
Stop-Service -Name JellyfinTelegramBot

# Restart service
Restart-Service -Name JellyfinTelegramBot

# Check if running
Get-Service -Name JellyfinTelegramBot | Select-Object Status
```

### Using Command Prompt (Admin)

```batch
# Start service
net start JellyfinTelegramBot

# Stop service
net stop JellyfinTelegramBot

# Check status
sc query JellyfinTelegramBot
```

### View Logs

**NSSM Logs**:
```powershell
# View stdout log
notepad "C:\Program Files\JellyfinTelegramBot\logs\stdout.log"

# View stderr log
notepad "C:\Program Files\JellyfinTelegramBot\logs\stderr.log"

# View application log (if LOG_FILE is set)
notepad "C:\Program Files\JellyfinTelegramBot\logs\bot.log"

# Or use PowerShell to tail logs
Get-Content "C:\Program Files\JellyfinTelegramBot\logs\bot.log" -Wait -Tail 50
```

**Event Viewer**:
1. Open Event Viewer (`eventvwr.msc`)
2. Go to **Windows Logs** → **Application**
3. Look for events from "JellyfinTelegramBot"

## Troubleshooting

### Service Won't Start

**Check Event Viewer**:
1. `eventvwr.msc` → Windows Logs → Application
2. Look for errors from JellyfinTelegramBot

**Common Issues**:

1. **Path issues**:
   - Use full paths in .env (e.g., `C:\path\to\file.db`)
   - Use backslashes (`\`) not forward slashes (`/`)
   - Quote paths with spaces: `"C:\Program Files\..."`

2. **Permission issues**:
   - Ensure service user has access to folders
   - Try running as Administrator
   - Check folder permissions in Windows Explorer

3. **Port already in use**:
   ```powershell
   # Check what's using port 8080
   netstat -ano | findstr :8080

   # Kill the process (if needed)
   taskkill /PID <pid> /F

   # Or change port in .env
   PORT=8081
   ```

4. **Missing .env file**:
   ```powershell
   # Verify file exists
   Test-Path "C:\Program Files\JellyfinTelegramBot\.env"

   # View contents
   Get-Content "C:\Program Files\JellyfinTelegramBot\.env"
   ```

### Service Crashes Immediately

**Check logs**:
```powershell
# View error log
Get-Content "C:\Program Files\JellyfinTelegramBot\logs\stderr.log"

# View application log
Get-Content "C:\Program Files\JellyfinTelegramBot\logs\bot.log" -Tail 50
```

**Test manually first**:
```powershell
# Run bot manually to see errors
cd "C:\Program Files\JellyfinTelegramBot"
.\jellyfin-telegram-bot.exe
```

If it works manually but not as service, it's likely a permissions or path issue.

### Can't Connect to Jellyfin

**From Windows Command Prompt**:
```powershell
# Test connectivity
curl http://your-jellyfin:8096/System/Info

# If curl not available, use PowerShell
Invoke-WebRequest -Uri "http://your-jellyfin:8096/System/Info"
```

**Firewall**:
```powershell
# Check if Windows Firewall is blocking
# Open Windows Firewall settings
firewall.cpl

# Add inbound rule for port 8080
New-NetFirewallRule -DisplayName "Jellyfin Telegram Bot" -Direction Inbound -LocalPort 8080 -Protocol TCP -Action Allow
```

### Update Bot Version

```powershell
# Stop service
Stop-Service -Name JellyfinTelegramBot

# Backup current version
Copy-Item "C:\Program Files\JellyfinTelegramBot\jellyfin-telegram-bot.exe" `
          "C:\Program Files\JellyfinTelegramBot\jellyfin-telegram-bot.exe.backup"

# Download new version and replace

# Start service
Start-Service -Name JellyfinTelegramBot

# Verify
Get-Service -Name JellyfinTelegramBot | Select-Object Status
```

### Remove Service Completely

**NSSM**:
```powershell
# Stop service
nssm stop JellyfinTelegramBot

# Remove service
nssm remove JellyfinTelegramBot confirm

# Delete files (CAUTION: This deletes your data!)
Remove-Item -Path "C:\Program Files\JellyfinTelegramBot" -Recurse -Force
```

**Task Scheduler**:
```powershell
# Delete task
schtasks /delete /tn "Jellyfin Telegram Bot" /f

# Delete files
Remove-Item -Path "C:\JellyfinTelegramBot" -Recurse -Force
```

## Best Practices

1. **Use NSSM** - Simplest and most reliable method
2. **Use full paths** - Always use complete Windows paths in .env
3. **Set up logging** - Configure LOG_FILE to track issues
4. **Regular backups** - Backup bot.db database file
5. **Test manually first** - Run bot.exe manually before creating service
6. **Monitor logs** - Check logs regularly for errors
7. **Secure .env** - Keep bot token and API key secure

## Windows-Specific Tips

### Running as Specific User

If you need the service to run as a specific Windows user:

**NSSM**:
1. `nssm edit JellyfinTelegramBot`
2. **Log on** tab
3. Select **This account**
4. Enter username and password
5. Click **Install service**

**Services Manager**:
1. `services.msc`
2. Right-click service → Properties
3. **Log On** tab
4. Select **This account**
5. Browse for user or enter: `DOMAIN\username`
6. Enter password

### Startup Delay

To delay service start (wait for network):

**NSSM**:
```powershell
# Set 30 second startup delay
nssm set JellyfinTelegramBot AppStopMethodSkip 0
nssm set JellyfinTelegramBot AppThrottle 30000
```

**Services Manager**:
1. Service Properties → **Dependencies** tab
2. Add dependency on "Network Connection"

### Resource Limits

Windows doesn't have built-in resource limits for services, but you can:

1. **Set process priority**:
   ```powershell
   # NSSM
   nssm set JellyfinTelegramBot AppPriority BELOW_NORMAL_PRIORITY_CLASS
   ```

2. **Use Process Lasso** - Third-party tool for resource management

## Support

For Windows-specific issues:
- Check logs in `C:\Program Files\JellyfinTelegramBot\logs\`
- Check Event Viewer: `eventvwr.msc`
- Review [Troubleshooting Guide](troubleshooting.md)
- Ask on [GitHub Discussions](https://github.com/yourusername/jellyfin-telegram-bot/discussions)

## Additional Resources

- [NSSM Documentation](https://nssm.cc/)
- [Windows Task Scheduler Guide](https://docs.microsoft.com/en-us/windows/win32/taskschd/task-scheduler-start-page)
- [Windows Services Overview](https://docs.microsoft.com/en-us/dotnet/framework/windows-services/)
