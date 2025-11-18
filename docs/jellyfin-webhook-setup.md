# Jellyfin Webhook Configuration Guide

This guide provides step-by-step instructions for configuring Jellyfin to send webhooks to the Telegram bot.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Install Webhook Plugin](#install-webhook-plugin)
3. [Configure Webhook](#configure-webhook)
4. [Test Webhook](#test-webhook)
5. [Troubleshooting](#troubleshooting)

## Prerequisites

- Jellyfin server running (version 10.8.0 or later recommended)
- Jellyfin administrator access
- Telegram bot deployed and running
- Bot server webhook URL accessible from Jellyfin server

## Install Webhook Plugin

### Step 1: Access Plugin Catalog

1. Log into Jellyfin as administrator
2. Navigate to **Dashboard** (Settings icon → Dashboard)
3. Click on **Plugins** in the left sidebar
4. Click on **Catalog** tab

### Step 2: Install Webhook Plugin

1. Search for "Webhook" in the catalog
2. Click on **Webhook** plugin
3. Click **Install** button
4. Wait for installation to complete
5. Restart Jellyfin server when prompted

### Step 3: Verify Installation

1. Go back to **Plugins** section
2. Click on **My Plugins** tab
3. Verify "Webhook" is listed and enabled
4. Note the plugin version

## Configure Webhook

### Step 1: Access Webhook Settings

1. In Jellyfin Dashboard, go to **Plugins**
2. Click on **My Plugins** tab
3. Click on **Webhook** plugin
4. This opens the webhook configuration page

### Step 2: Add Generic Destination

1. Click **Add Generic Destination** button
2. A configuration form will appear

### Step 3: Basic Configuration

Fill in the following fields:

**Webhook Name**: `Telegram Bot` (or any descriptive name)

**Webhook Url**:
```
http://your-bot-server-ip:8080/webhook
```

Examples:
- Same server: `http://localhost:8080/webhook`
- Different server: `http://192.168.1.100:8080/webhook`
- With domain: `http://bot.example.com:8080/webhook`
- With HTTPS (if using reverse proxy): `https://bot.example.com/webhook`

**Request Method**: `POST` (default)

**Content Type**: `application/json` (default)

### Step 4: Notification Configuration

**Notification Type** (Select these):
- [x] Item Added

**Do NOT select**:
- [ ] Item Updated
- [ ] Item Removed
- [ ] Playback Start
- [ ] Playback Stop
- [ ] User Created
- [ ] Authentication Success/Failure
- [ ] etc.

We only want "Item Added" notifications.

### Step 5: Item Type Filter

Under "Item Type", select:
- [x] Movie
- [x] Episode

**Do NOT select**:
- [ ] Series (notifications are per episode, not series)
- [ ] Season (notifications are per episode)
- [ ] Audio
- [ ] Book
- [ ] Photo
- [ ] etc.

### Step 6: Template Configuration

**Send All Properties**: `Enabled` (toggle ON)

This ensures all metadata fields are included in the webhook payload.

**Custom Template**: Leave empty (use default)

### Step 7: Add Authentication Header

Click **Add Header** button:

**Header Name**: `X-Webhook-Secret`

**Header Value**: `your-webhook-secret-from-env`

This must match the `WEBHOOK_SECRET` value in your bot's `.env` file.

Example:
```
X-Webhook-Secret: my-super-secret-webhook-token-12345
```

### Step 8: Save Configuration

1. Review all settings
2. Click **Save** button at the bottom
3. Webhook is now configured

## Test Webhook

### Test 1: Add Content to Jellyfin

1. Add a new movie or TV episode to your Jellyfin library
2. Wait for Jellyfin to scan and add the item (usually automatic)
3. Check bot logs for webhook reception:
   ```bash
   sudo journalctl -u jellyfin-bot -f
   ```

Expected log output:
```
INFO  Received webhook notification_type=ItemAdded item_type=Movie item_name="Test Movie"
INFO  New content ready for notification
INFO  Content marked as notified
INFO  Notification broadcast initiated
INFO  Broadcast completed success=X failures=0 blocked=0
```

4. Check Telegram for notification message

### Test 2: Manual Webhook Test

Send a test webhook manually:

```bash
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Secret: your-secret-here" \
  -d '{
    "NotificationType": "ItemAdded",
    "ItemType": "Movie",
    "Name": "Test Movie",
    "Overview": "This is a test notification",
    "Year": 2024,
    "ItemId": "test-123-456"
  }'
```

Expected response: `HTTP 200 OK`

### Test 3: Check Jellyfin Webhook Logs

Jellyfin logs webhook attempts:

1. Go to Jellyfin Dashboard → **Logs**
2. Look for webhook-related messages
3. Check for successful POST requests to webhook URL
4. Verify no error messages

## Webhook Payload Structure

The bot expects this JSON structure from Jellyfin:

### Movie Example

```json
{
  "NotificationType": "ItemAdded",
  "ItemType": "Movie",
  "Name": "Interstellar",
  "Overview": "A team of explorers travel through a wormhole in space...",
  "Year": 2014,
  "ItemId": "abc123def456",
  "ServerId": "server-id",
  "ServerUrl": "http://localhost:8096"
}
```

### Episode Example

```json
{
  "NotificationType": "ItemAdded",
  "ItemType": "Episode",
  "Name": "Pilot",
  "SeriesName": "Breaking Bad",
  "SeasonNumber": 1,
  "EpisodeNumber": 1,
  "Overview": "High school chemistry teacher Walter White's life is suddenly...",
  "Year": 2008,
  "ItemId": "episode-123-456",
  "ServerId": "server-id",
  "ServerUrl": "http://localhost:8096"
}
```

## Troubleshooting

### Webhook Not Sent

**Issue**: Jellyfin doesn't send webhooks when content is added.

**Solutions**:
1. Verify webhook plugin is installed and enabled
2. Check webhook configuration is saved
3. Restart Jellyfin server:
   ```bash
   sudo systemctl restart jellyfin
   ```
4. Check Jellyfin logs for errors
5. Verify "Item Added" event is selected
6. Ensure Movie/Episode types are selected

### Webhook Fails (401 Unauthorized)

**Issue**: Bot logs show "Unauthorized" error.

**Cause**: Webhook secret mismatch.

**Solutions**:
1. Verify webhook secret in Jellyfin matches bot's `.env`
2. Check header name is exactly: `X-Webhook-Secret`
3. No extra spaces in header value
4. Restart bot after changing `.env`:
   ```bash
   sudo systemctl restart jellyfin-bot
   ```

### Webhook Fails (Connection Refused)

**Issue**: Jellyfin cannot reach webhook URL.

**Solutions**:
1. Verify bot is running:
   ```bash
   sudo systemctl status jellyfin-bot
   ```
2. Test webhook endpoint manually from Jellyfin server:
   ```bash
   curl http://bot-server-ip:8080/health
   ```
3. Check firewall allows port 8080:
   ```bash
   sudo ufw status
   sudo ufw allow 8080/tcp
   ```
4. Verify bot's webhook server is listening:
   ```bash
   sudo netstat -tlnp | grep 8080
   ```
5. Try using IP address instead of hostname in webhook URL

### Webhook Succeeds But No Notification

**Issue**: Bot receives webhook but doesn't send notification.

**Solutions**:
1. Check bot logs for errors:
   ```bash
   sudo journalctl -u jellyfin-bot -n 100
   ```
2. Verify bot has active subscribers:
   ```bash
   # Subscribe by sending /start to bot in Telegram
   ```
3. Check database for subscribers:
   ```bash
   sqlite3 /opt/jellyfin-bot/bot.db "SELECT * FROM subscribers;"
   ```
4. Verify content hasn't been notified already (duplicate check):
   ```bash
   sqlite3 /opt/jellyfin-bot/bot.db "SELECT * FROM content_cache WHERE jellyfin_id='item-id';"
   ```
5. Check Telegram bot token is valid

### Duplicate Notifications

**Issue**: Receiving multiple notifications for the same content.

**Solutions**:
1. Check if multiple webhooks are configured in Jellyfin
2. Verify content tracking in database:
   ```bash
   sqlite3 /opt/jellyfin-bot/bot.db "SELECT * FROM content_cache ORDER BY created_at DESC LIMIT 10;"
   ```
3. Ensure only one bot instance is running:
   ```bash
   ps aux | grep jellyfin-bot
   ```

### Wrong Content Type Notified

**Issue**: Receiving notifications for Series, Seasons, Audio, etc.

**Cause**: Wrong item types selected in webhook config.

**Solutions**:
1. Go to Jellyfin → Plugins → Webhook
2. Edit webhook configuration
3. Ensure ONLY Movie and Episode are selected
4. Uncheck all other item types
5. Save configuration

## Advanced Configuration

### Multiple Webhooks

You can configure multiple webhook destinations:

1. One for Telegram bot
2. One for logging/monitoring service
3. One for other integrations

Each webhook is independent and can have different:
- URLs
- Events
- Item types
- Headers

### Webhook Retry Configuration

Jellyfin webhook plugin doesn't have built-in retry mechanism. If webhook fails:
- It won't be retried automatically
- Content won't be notified
- Recommendation: Implement health checks and monitoring

### Filtering by Library

To send webhooks only for specific libraries:

Unfortunately, Jellyfin webhook plugin doesn't support library filtering directly. Workarounds:
1. Run multiple bot instances (one per library)
2. Implement library filtering in bot code (future enhancement)
3. Use separate Jellyfin servers for different libraries

### HTTPS Webhooks

For production deployments with HTTPS:

1. Set up reverse proxy (nginx/apache) with SSL certificate
2. Configure reverse proxy to forward `/webhook` to bot
3. Update Jellyfin webhook URL to HTTPS endpoint:
   ```
   https://bot.example.com/webhook
   ```

Nginx example:
```nginx
server {
    listen 443 ssl;
    server_name bot.example.com;

    ssl_certificate /etc/letsencrypt/live/bot.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/bot.example.com/privkey.pem;

    location /webhook {
        proxy_pass http://localhost:8080/webhook;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Security Best Practices

1. **Always use webhook secret** for authentication
2. **Use HTTPS** in production environments
3. **Restrict webhook endpoint** to Jellyfin server IP only (firewall rules)
4. **Rotate webhook secret** periodically
5. **Monitor webhook logs** for suspicious activity
6. **Use reverse proxy** to hide internal bot server

## Monitoring

Monitor webhook health:

1. **Jellyfin side**: Check Dashboard → Logs for webhook POST requests
2. **Bot side**: Monitor logs for webhook reception:
   ```bash
   sudo journalctl -u jellyfin-bot | grep "Received webhook"
   ```
3. **Health check**: Regular health endpoint tests:
   ```bash
   */5 * * * * curl -f http://localhost:8080/health
   ```

## Reference

### Jellyfin Webhook Plugin
- Official repository: https://github.com/jellyfin/jellyfin-plugin-webhook
- Documentation: https://github.com/jellyfin/jellyfin-plugin-webhook/wiki

### Webhook Events Reference
- `ItemAdded`: New content added to library (USED)
- `ItemUpdated`: Content metadata updated (NOT USED)
- `ItemRemoved`: Content removed from library (NOT USED)
- `PlaybackStart`: User starts playback (NOT USED)
- `PlaybackStop`: User stops playback (NOT USED)
- Others: See plugin documentation

### Item Types Reference
- `Movie`: Feature films (USED)
- `Episode`: TV show episodes (USED)
- `Series`: TV show series container (NOT USED - notified per episode)
- `Season`: TV show season container (NOT USED - notified per episode)
- `Audio`: Music files (NOT USED)
- `Book`: Ebooks (NOT USED)
- Others: See Jellyfin documentation

## Support

If you encounter issues:

1. Check bot logs: `sudo journalctl -u jellyfin-bot -f`
2. Check Jellyfin logs: Dashboard → Logs
3. Test webhook manually with curl
4. Verify network connectivity between Jellyfin and bot
5. Review this guide's troubleshooting section
6. Check bot health endpoint: `curl http://localhost:8080/health`
