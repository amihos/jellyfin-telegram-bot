# Running as a Linux System Service

This guide explains how to run the Jellyfin Telegram Bot as a systemd service on Linux. This ensures the bot starts automatically on boot and restarts if it crashes.

## Prerequisites

- Linux system with systemd (most modern distributions)
- Bot binary downloaded or built
- Root or sudo access

## Quick Setup

### 1. Create Bot User

Create a dedicated user to run the bot (for security):

```bash
# Create system user (no home directory, no login)
sudo useradd -r -s /bin/false jellyfin-bot

# Or create regular user with home directory
sudo useradd -m -s /bin/bash jellyfin-bot
```

### 2. Install Bot Binary

```bash
# Create installation directory
sudo mkdir -p /opt/jellyfin-telegram-bot
cd /opt/jellyfin-telegram-bot

# Download the latest release (adjust URL for your architecture)
sudo wget https://github.com/yourusername/jellyfin-telegram-bot/releases/latest/download/jellyfin-telegram-bot-linux-amd64
sudo chmod +x jellyfin-telegram-bot-linux-amd64
sudo mv jellyfin-telegram-bot-linux-amd64 jellyfin-telegram-bot

# Or copy your built binary
sudo cp /path/to/your/jellyfin-telegram-bot /opt/jellyfin-telegram-bot/
sudo chmod +x /opt/jellyfin-telegram-bot/jellyfin-telegram-bot
```

### 3. Create Configuration File

```bash
# Create .env file
sudo nano /opt/jellyfin-telegram-bot/.env
```

Add your configuration:

```env
TELEGRAM_BOT_TOKEN=your_telegram_bot_token
JELLYFIN_SERVER_URL=http://your-jellyfin:8096
JELLYFIN_API_KEY=your_jellyfin_api_key
PORT=8080
DATABASE_PATH=/opt/jellyfin-telegram-bot/data/bot.db
LOG_FILE=/var/log/jellyfin-telegram-bot/bot.log
LOG_LEVEL=INFO
```

Secure the file:

```bash
sudo chmod 600 /opt/jellyfin-telegram-bot/.env
```

### 4. Create Required Directories

```bash
# Create data directory for database
sudo mkdir -p /opt/jellyfin-telegram-bot/data

# Create log directory
sudo mkdir -p /var/log/jellyfin-telegram-bot

# Set ownership
sudo chown -R jellyfin-bot:jellyfin-bot /opt/jellyfin-telegram-bot
sudo chown -R jellyfin-bot:jellyfin-bot /var/log/jellyfin-telegram-bot
```

### 5. Install Systemd Service File

```bash
# Download the service file
sudo curl -o /etc/systemd/system/jellyfin-telegram-bot.service \
  https://raw.githubusercontent.com/yourusername/jellyfin-telegram-bot/main/docs/jellyfin-telegram-bot.service

# Or create it manually
sudo nano /etc/systemd/system/jellyfin-telegram-bot.service
```

Service file content (adjust paths if needed):

```ini
[Unit]
Description=Jellyfin Telegram Bot
Documentation=https://github.com/yourusername/jellyfin-telegram-bot
After=network.target

[Service]
Type=simple
User=jellyfin-bot
Group=jellyfin-bot
WorkingDirectory=/opt/jellyfin-telegram-bot
ExecStart=/opt/jellyfin-telegram-bot/jellyfin-telegram-bot
EnvironmentFile=/opt/jellyfin-telegram-bot/.env
NoNewPrivileges=true
Restart=on-failure
RestartSec=5s
StandardOutput=journal
StandardError=journal
SyslogIdentifier=jellyfin-telegram-bot

[Install]
WantedBy=multi-user.target
```

### 6. Enable and Start Service

```bash
# Reload systemd to recognize new service
sudo systemctl daemon-reload

# Enable service to start on boot
sudo systemctl enable jellyfin-telegram-bot

# Start the service now
sudo systemctl start jellyfin-telegram-bot

# Check status
sudo systemctl status jellyfin-telegram-bot
```

## Managing the Service

### Check Status

```bash
# View service status
sudo systemctl status jellyfin-telegram-bot

# Should show:
# Active: active (running)
```

### View Logs

```bash
# View recent logs
sudo journalctl -u jellyfin-telegram-bot -n 50

# Follow logs in real-time
sudo journalctl -u jellyfin-telegram-bot -f

# View logs with timestamps
sudo journalctl -u jellyfin-telegram-bot -f -o short-iso

# View logs for specific time range
sudo journalctl -u jellyfin-telegram-bot --since "1 hour ago"
sudo journalctl -u jellyfin-telegram-bot --since "2024-01-15 10:00:00"

# Also check the log file if LOG_FILE is set
tail -f /var/log/jellyfin-telegram-bot/bot.log
```

### Control Service

```bash
# Start the bot
sudo systemctl start jellyfin-telegram-bot

# Stop the bot
sudo systemctl stop jellyfin-telegram-bot

# Restart the bot
sudo systemctl restart jellyfin-telegram-bot

# Reload configuration (if supported)
sudo systemctl reload jellyfin-telegram-bot

# Enable auto-start on boot
sudo systemctl enable jellyfin-telegram-bot

# Disable auto-start on boot
sudo systemctl disable jellyfin-telegram-bot
```

### Update Configuration

```bash
# Edit configuration
sudo nano /opt/jellyfin-telegram-bot/.env

# Restart to apply changes
sudo systemctl restart jellyfin-telegram-bot

# Verify changes took effect
sudo journalctl -u jellyfin-telegram-bot -n 20
```

### Update Bot Version

```bash
# Stop the service
sudo systemctl stop jellyfin-telegram-bot

# Backup current version
sudo cp /opt/jellyfin-telegram-bot/jellyfin-telegram-bot \
       /opt/jellyfin-telegram-bot/jellyfin-telegram-bot.backup

# Download new version
cd /opt/jellyfin-telegram-bot
sudo wget https://github.com/yourusername/jellyfin-telegram-bot/releases/latest/download/jellyfin-telegram-bot-linux-amd64
sudo chmod +x jellyfin-telegram-bot-linux-amd64
sudo mv jellyfin-telegram-bot-linux-amd64 jellyfin-telegram-bot

# Restart with new version
sudo systemctl start jellyfin-telegram-bot

# Check it's running properly
sudo systemctl status jellyfin-telegram-bot
```

## Troubleshooting

### Service Won't Start

**Check status and logs**:
```bash
sudo systemctl status jellyfin-telegram-bot
sudo journalctl -u jellyfin-telegram-bot -n 50
```

**Common issues**:

1. **Permission denied**:
   ```bash
   # Fix ownership
   sudo chown -R jellyfin-bot:jellyfin-bot /opt/jellyfin-telegram-bot
   sudo chown -R jellyfin-bot:jellyfin-bot /var/log/jellyfin-telegram-bot

   # Fix binary permissions
   sudo chmod +x /opt/jellyfin-telegram-bot/jellyfin-telegram-bot
   ```

2. **Binary not found**:
   ```bash
   # Verify path in service file
   sudo nano /etc/systemd/system/jellyfin-telegram-bot.service

   # Check binary exists
   ls -l /opt/jellyfin-telegram-bot/jellyfin-telegram-bot
   ```

3. **Environment file not found**:
   ```bash
   # Create if missing
   sudo nano /opt/jellyfin-telegram-bot/.env

   # Set permissions
   sudo chmod 600 /opt/jellyfin-telegram-bot/.env
   sudo chown jellyfin-bot:jellyfin-bot /opt/jellyfin-telegram-bot/.env
   ```

4. **Port already in use**:
   ```bash
   # Check what's using the port
   sudo lsof -i :8080

   # Change port in .env
   sudo nano /opt/jellyfin-telegram-bot/.env
   # Change: PORT=8081

   # Restart
   sudo systemctl restart jellyfin-telegram-bot
   ```

### Service Crashes or Restarts

**Check logs for errors**:
```bash
sudo journalctl -u jellyfin-telegram-bot -n 100
```

**Increase restart delay if flapping**:
```bash
sudo nano /etc/systemd/system/jellyfin-telegram-bot.service
```

Change:
```ini
RestartSec=30s
```

Then:
```bash
sudo systemctl daemon-reload
sudo systemctl restart jellyfin-telegram-bot
```

### High Resource Usage

**Check resource usage**:
```bash
systemctl status jellyfin-telegram-bot
```

**Set resource limits**:
```bash
sudo nano /etc/systemd/system/jellyfin-telegram-bot.service
```

Add under `[Service]`:
```ini
# Limit memory to 128MB
MemoryMax=128M

# Limit CPU to 50%
CPUQuota=50%

# Limit number of file descriptors
LimitNOFILE=1024
```

Then:
```bash
sudo systemctl daemon-reload
sudo systemctl restart jellyfin-telegram-bot
```

## Advanced Configuration

### Running Multiple Instances

To run multiple bot instances (e.g., for different Jellyfin servers):

**Instance 1**:
```bash
# Create separate directories
sudo mkdir -p /opt/jellyfin-bot-1
# Configure as above with different port (8081)

# Create service file
sudo nano /etc/systemd/system/jellyfin-telegram-bot-1.service
# Adjust paths and WorkingDirectory

sudo systemctl enable jellyfin-telegram-bot-1
sudo systemctl start jellyfin-telegram-bot-1
```

**Instance 2**:
```bash
sudo mkdir -p /opt/jellyfin-bot-2
# Configure with different port (8082)

sudo nano /etc/systemd/system/jellyfin-telegram-bot-2.service

sudo systemctl enable jellyfin-telegram-bot-2
sudo systemctl start jellyfin-telegram-bot-2
```

### Security Hardening

Add these settings to the service file for increased security:

```ini
[Service]
# Prevent privilege escalation
NoNewPrivileges=true

# Restrict access to system
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true

# Allow writes only to specific directories
ReadWritePaths=/opt/jellyfin-telegram-bot/data
ReadWritePaths=/var/log/jellyfin-telegram-bot

# Network restrictions
RestrictAddressFamilies=AF_INET AF_INET6

# System call restrictions
SystemCallFilter=@system-service
SystemCallErrorNumber=EPERM
```

### Log Rotation

If not using systemd journal exclusively:

```bash
# Create logrotate config
sudo nano /etc/logrotate.d/jellyfin-telegram-bot
```

Content:
```
/var/log/jellyfin-telegram-bot/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 0644 jellyfin-bot jellyfin-bot
    postrotate
        systemctl reload jellyfin-telegram-bot > /dev/null 2>&1 || true
    endscript
}
```

### Monitoring with Systemd

**Monitor service health**:
```bash
# Check if service is running
systemctl is-active jellyfin-telegram-bot

# Check if service is enabled
systemctl is-enabled jellyfin-telegram-bot

# View failure count
systemctl show jellyfin-telegram-bot -p NRestarts
```

**Email notifications on failure**:

Install systemd-mail:
```bash
sudo apt install systemd-mail
```

Edit service:
```ini
[Service]
OnFailure=status-email@%n.service
```

## Uninstalling

To completely remove the bot:

```bash
# Stop and disable service
sudo systemctl stop jellyfin-telegram-bot
sudo systemctl disable jellyfin-telegram-bot

# Remove service file
sudo rm /etc/systemd/system/jellyfin-telegram-bot.service

# Reload systemd
sudo systemctl daemon-reload

# Remove bot files (CAUTION: This deletes your data!)
sudo rm -rf /opt/jellyfin-telegram-bot
sudo rm -rf /var/log/jellyfin-telegram-bot

# Remove user
sudo userdel jellyfin-bot
```

## Best Practices

1. **Always use a dedicated user** - Don't run as root
2. **Secure your .env file** - chmod 600, readable only by bot user
3. **Monitor logs regularly** - Check for errors
4. **Backup your database** - Regularly backup bot.db
5. **Test updates** - Test new versions before deploying
6. **Use resource limits** - Prevent runaway resource usage
7. **Enable on boot** - Use `systemctl enable`
8. **Document changes** - Keep notes on customizations

## Support

For issues with systemd service:
- Check logs: `sudo journalctl -u jellyfin-telegram-bot -n 100`
- Review [Troubleshooting Guide](troubleshooting.md)
- Ask on [GitHub Discussions](https://github.com/yourusername/jellyfin-telegram-bot/discussions)
