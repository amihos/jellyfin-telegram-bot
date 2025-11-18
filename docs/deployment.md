# Deployment Guide

This guide covers deploying the Jellyfin Telegram Bot to a production server.

## Server Requirements

### Minimum Requirements
- **OS**: Linux (Ubuntu 20.04+, Debian 11+, CentOS 8+, etc.)
- **RAM**: 128 MB (typically uses 10-20 MB)
- **Disk**: 50 MB (excluding logs and database)
- **CPU**: 1 core (minimal usage)
- **Go Version**: 1.21+ (only needed for building, not for running)

### Network Requirements
- Outbound HTTPS access to Telegram API (api.telegram.org)
- Outbound HTTP/HTTPS access to Jellyfin server
- Inbound HTTP access on configured webhook port (default: 8080)

## Deployment Methods

### Method 1: Single Binary Deployment (Recommended)

This is the simplest deployment method, leveraging Go's single binary compilation.

#### 1. Build on Development Machine

```bash
# Clone repository
git clone <repository-url>
cd jellyfin-telegram-bot

# Install dependencies
go mod download

# Build for Linux (if building on different OS)
GOOS=linux GOARCH=amd64 go build -o jellyfin-bot cmd/bot/main.go

# Or build for current system
go build -o jellyfin-bot cmd/bot/main.go
```

#### 2. Copy to Server

```bash
# Copy binary to server
scp jellyfin-bot user@your-server:/opt/jellyfin-bot/

# Copy environment template
scp .env.example user@your-server:/opt/jellyfin-bot/.env
```

#### 3. Configure on Server

```bash
# SSH to server
ssh user@your-server

# Create application directory
sudo mkdir -p /opt/jellyfin-bot
sudo chown $USER:$USER /opt/jellyfin-bot
cd /opt/jellyfin-bot

# Make binary executable
chmod +x jellyfin-bot

# Configure environment
nano .env
# Fill in your credentials (see Configuration section)
```

#### 4. Set Up systemd Service

Create systemd service file:

```bash
sudo nano /etc/systemd/system/jellyfin-bot.service
```

Add the following content:

```ini
[Unit]
Description=Jellyfin Telegram Bot
After=network.target

[Service]
Type=simple
User=your-username
WorkingDirectory=/opt/jellyfin-bot
ExecStart=/opt/jellyfin-bot/jellyfin-bot
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=jellyfin-bot

# Environment file
EnvironmentFile=/opt/jellyfin-bot/.env

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/jellyfin-bot

[Install]
WantedBy=multi-user.target
```

#### 5. Start and Enable Service

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable service to start on boot
sudo systemctl enable jellyfin-bot

# Start service
sudo systemctl start jellyfin-bot

# Check status
sudo systemctl status jellyfin-bot

# View logs
sudo journalctl -u jellyfin-bot -f
```

### Method 2: Docker Deployment (Coming Soon)

Docker deployment will be added in a future update.

## Configuration

### Environment Variables

Create a `.env` file in the application directory:

```env
# Required
TELEGRAM_BOT_TOKEN=1234567890:ABCdefGHIjklMNOpqrsTUVwxyz
JELLYFIN_SERVER_URL=http://192.168.1.100:8096
JELLYFIN_API_KEY=abc123def456ghi789jkl

# Optional
WEBHOOK_SECRET=my-secure-secret-token
PORT=8080
DATABASE_PATH=/opt/jellyfin-bot/bot.db
LOG_LEVEL=INFO
LOG_FILE=/opt/jellyfin-bot/logs/bot.log
```

### Security Considerations

1. **File Permissions**:
   ```bash
   chmod 600 /opt/jellyfin-bot/.env  # Restrict .env access
   chmod 755 /opt/jellyfin-bot/jellyfin-bot  # Binary executable
   ```

2. **Webhook Secret**: Set `WEBHOOK_SECRET` in `.env` and configure the same in Jellyfin webhook plugin

3. **Firewall**: Only expose webhook port if Jellyfin is on different server:
   ```bash
   # If using UFW
   sudo ufw allow 8080/tcp
   ```

4. **User Permissions**: Run as non-root user (configured in systemd service)

## Database Management

### Backup

```bash
# Stop bot
sudo systemctl stop jellyfin-bot

# Backup database
cp /opt/jellyfin-bot/bot.db /opt/jellyfin-bot/bot.db.backup

# Start bot
sudo systemctl start jellyfin-bot
```

### Restore

```bash
# Stop bot
sudo systemctl stop jellyfin-bot

# Restore database
cp /opt/jellyfin-bot/bot.db.backup /opt/jellyfin-bot/bot.db

# Start bot
sudo systemctl start jellyfin-bot
```

### Automated Backups

Add to crontab for daily backups:

```bash
crontab -e
```

Add line:
```
0 2 * * * cp /opt/jellyfin-bot/bot.db /opt/jellyfin-bot/backups/bot-$(date +\%Y\%m\%d).db
```

## Log Management

### Viewing Logs

```bash
# Real-time logs
sudo journalctl -u jellyfin-bot -f

# Last 100 lines
sudo journalctl -u jellyfin-bot -n 100

# Logs from specific date
sudo journalctl -u jellyfin-bot --since "2024-01-15"
```

### Log Rotation

Logs are automatically rotated by lumberjack library:
- Max size: 10 MB per file
- Max backups: 5 files
- Max age: 30 days
- Compression: Enabled

Configuration in `.env`:
```env
LOG_FILE=/opt/jellyfin-bot/logs/bot.log
```

## Monitoring

### Health Check

Check if bot is running:

```bash
# Service status
sudo systemctl status jellyfin-bot

# Process check
ps aux | grep jellyfin-bot

# Port check (webhook server)
sudo netstat -tlnp | grep 8080
```

### Performance Monitoring

```bash
# Memory usage
ps aux | grep jellyfin-bot | awk '{print $6}'

# CPU usage
top -p $(pgrep jellyfin-bot)
```

## Updating

### Update Binary

```bash
# Build new version on development machine
go build -o jellyfin-bot cmd/bot/main.go

# Stop service on server
sudo systemctl stop jellyfin-bot

# Copy new binary
scp jellyfin-bot user@your-server:/opt/jellyfin-bot/

# Start service
sudo systemctl start jellyfin-bot

# Verify
sudo systemctl status jellyfin-bot
```

### Database Migrations

Future versions may include database migrations. Follow release notes for migration instructions.

## Troubleshooting

### Bot Won't Start

1. Check service status:
   ```bash
   sudo systemctl status jellyfin-bot
   ```

2. Check logs:
   ```bash
   sudo journalctl -u jellyfin-bot -n 50
   ```

3. Verify configuration:
   ```bash
   cat /opt/jellyfin-bot/.env
   ```

4. Test binary manually:
   ```bash
   cd /opt/jellyfin-bot
   ./jellyfin-bot
   ```

### Webhooks Not Working

1. Test webhook endpoint:
   ```bash
   curl -X POST http://localhost:8080/webhook
   ```

2. Check firewall:
   ```bash
   sudo ufw status
   ```

3. Verify Jellyfin can reach webhook URL

4. Check webhook logs in bot output

### High Memory Usage

Normal memory usage: 10-20 MB. If higher:

1. Check for memory leaks in logs
2. Restart service:
   ```bash
   sudo systemctl restart jellyfin-bot
   ```
3. Review recent activity (large broadcasts, etc.)

## Reverse Proxy Setup (Optional)

If you want to use HTTPS for webhooks or hide the webhook port:

### Nginx Example

```nginx
server {
    listen 443 ssl;
    server_name bot.yourdomain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location /webhook {
        proxy_pass http://localhost:8080/webhook;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Update Jellyfin webhook URL to: `https://bot.yourdomain.com/webhook`

## Production Checklist

- [ ] Binary built and tested
- [ ] Environment variables configured
- [ ] Database directory writable
- [ ] Logs directory created
- [ ] systemd service installed and enabled
- [ ] Firewall rules configured
- [ ] Webhook secret configured in both bot and Jellyfin
- [ ] Bot tested with `/start` command
- [ ] Webhook tested by adding content to Jellyfin
- [ ] Backup strategy implemented
- [ ] Monitoring in place

## Performance Optimization

For large installations (100+ users):

1. **Database**: Consider migrating from SQLite to PostgreSQL
2. **Caching**: Implement Redis for content metadata caching
3. **Rate Limiting**: Adjust Telegram broadcast delays if needed
4. **Load Balancing**: Run multiple instances behind load balancer

These optimizations are not needed for typical installations (<100 users).
