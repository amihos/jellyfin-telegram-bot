# Webhook Implementation Documentation

## Overview

The webhook receiver is implemented in `/internal/handlers/webhook.go` and provides a secure HTTP endpoint for receiving Jellyfin webhook notifications when new content is added to the media server.

## Architecture

### Components

1. **WebhookHandler** - Main handler struct that processes incoming webhooks
2. **ContentTracker Interface** - Abstraction for database operations
3. **HTTP Server** - Built-in Go `net/http` server running on configurable port

### Flow

```
Jellyfin Server
    |
    | POST /webhook (JSON payload)
    v
WebhookHandler.HandleWebhook()
    |
    ├─> Validate request method (POST only)
    ├─> Validate webhook secret (if configured)
    ├─> Parse JSON payload
    ├─> Validate content type (ItemAdded, Movie/Episode only)
    ├─> Check for duplicates (IsContentNotified)
    ├─> Extract metadata
    ├─> Mark as notified (MarkContentNotified)
    └─> Log for debugging (actual notification in Task Group 5)
```

## API Specification

### Endpoint

- **URL**: `/webhook`
- **Method**: `POST`
- **Content-Type**: `application/json`

### Request Headers

| Header | Required | Description |
|--------|----------|-------------|
| `Content-Type` | Yes | Must be `application/json` |
| `X-Webhook-Secret` | Conditional | Required if `WEBHOOK_SECRET` environment variable is set |

### Request Body

The webhook expects a JSON payload matching the `JellyfinWebhook` structure:

```json
{
  "NotificationType": "ItemAdded",
  "Timestamp": "2024-01-15T10:30:00Z",
  "ServerId": "abc123",
  "ServerName": "My Jellyfin Server",
  "ServerUrl": "http://localhost:8096",
  "ServerVersion": "10.8.0",
  "ItemId": "movie123",
  "ItemName": "Interstellar",
  "ItemType": "Movie",
  "Year": 2014,
  "Overview": "A team of explorers travel through a wormhole in space.",
  "ItemPath": "/media/movies/Interstellar.mkv",
  "UserName": "admin",
  "UserId": "user123"
}
```

For TV episodes, additional fields are included:

```json
{
  "NotificationType": "ItemAdded",
  "ItemId": "episode456",
  "ItemName": "Pilot",
  "ItemType": "Episode",
  "SeriesName": "Breaking Bad",
  "SeasonNumber": 1,
  "EpisodeNumber": 1,
  "Overview": "A high school chemistry teacher turned meth cook."
}
```

### Response Codes

| Code | Description |
|------|-------------|
| 200 | Webhook processed successfully (or ignored due to filtering) |
| 400 | Invalid JSON payload |
| 401 | Invalid or missing webhook secret |
| 405 | Method not allowed (only POST accepted) |
| 500 | Internal server error (database failure) |

## Content Filtering

The webhook handler automatically filters content based on:

1. **Notification Type**: Only `ItemAdded` notifications are processed
2. **Item Type**: Only `Movie` and `Episode` items are processed
3. **Duplicates**: Content already in the database is skipped

### Rejected Content Types

The following item types are ignored:
- Series
- Season
- Audio
- Book
- Photo
- MusicAlbum
- Any other non-Movie/Episode type

## Security

### Webhook Secret

To prevent unauthorized webhook calls, configure a secret token:

1. Set `WEBHOOK_SECRET` environment variable
2. Configure Jellyfin webhook plugin to send the same secret in `X-Webhook-Secret` header
3. Requests without matching secret will receive 401 Unauthorized response

**Example Configuration**:
```bash
WEBHOOK_SECRET=your-random-secret-token-here
```

### Best Practices

1. **Use HTTPS**: Deploy behind a reverse proxy with HTTPS in production
2. **Firewall Rules**: Restrict webhook endpoint to Jellyfin server IP only
3. **Strong Secret**: Use a randomly generated secret (32+ characters)
4. **Monitor Logs**: Watch for unauthorized access attempts

## Usage

### Starting the Webhook Server

```go
package main

import (
    "jellyfin-telegram-bot/internal/database"
    "jellyfin-telegram-bot/internal/handlers"
    "log/slog"
    "os"
)

func main() {
    // Initialize database
    db, err := database.NewDB("./bot.db")
    if err != nil {
        slog.Error("Failed to initialize database", "error", err)
        os.Exit(1)
    }
    defer db.Close()

    // Create webhook handler
    webhookSecret := os.Getenv("WEBHOOK_SECRET")
    webhookHandler := handlers.NewWebhookHandler(db, webhookSecret)

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    if err := handlers.StartWebhookServer(port, webhookHandler); err != nil {
        slog.Error("Webhook server error", "error", err)
        os.Exit(1)
    }
}
```

### Testing the Webhook

Use `curl` to test the webhook endpoint:

```bash
# Test movie webhook
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Secret: your-secret-here" \
  -d '{
    "NotificationType": "ItemAdded",
    "ItemType": "Movie",
    "ItemId": "test123",
    "ItemName": "Test Movie",
    "Year": 2024,
    "Overview": "A test movie"
  }'

# Test episode webhook
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Secret: your-secret-here" \
  -d '{
    "NotificationType": "ItemAdded",
    "ItemType": "Episode",
    "ItemId": "test456",
    "ItemName": "Pilot",
    "SeriesName": "Test Series",
    "SeasonNumber": 1,
    "EpisodeNumber": 1,
    "Overview": "A test episode"
  }'
```

## Metadata Extraction

The webhook handler extracts and validates the following metadata:

### For Movies
- **Title**: `ItemName` (defaults to "Unknown")
- **Overview**: `Overview` (defaults to "No description available")
- **Year**: `Year`
- **ItemID**: `ItemId`

### For Episodes
- **Title**: `ItemName` (episode title)
- **Series Name**: `SeriesName` (defaults to "Unknown Series")
- **Season Number**: `SeasonNumber`
- **Episode Number**: `EpisodeNumber`
- **Overview**: `Overview` (defaults to "No description available")
- **ItemID**: `ItemId`

All missing fields are handled gracefully with sensible defaults.

## Logging

The webhook handler logs at different levels:

- **INFO**: Successful webhook processing, new content detected
- **WARN**: Invalid webhook secret attempts
- **DEBUG**: Filtered/ignored content
- **ERROR**: JSON parsing errors, database errors

Example log output:
```
INFO: Received webhook notification_type=ItemAdded item_type=Movie item_id=movie123 item_name=Interstellar
INFO: New content ready for notification item_id=movie123 type=Movie title=Interstellar year=2014
INFO: Content marked as notified item_id=movie123 item_name=Interstellar
```

## Testing

The webhook implementation includes 8 comprehensive tests:

1. **TestWebhookHandler_ValidMoviePayload** - Valid movie processing
2. **TestWebhookHandler_ValidEpisodePayload** - Valid episode processing
3. **TestWebhookHandler_FilterInvalidContentType** - Content type filtering
4. **TestWebhookHandler_DuplicateContent** - Duplicate detection
5. **TestWebhookHandler_InvalidJSON** - Malformed payload handling
6. **TestWebhookHandler_WrongNotificationType** - Notification type filtering
7. **TestWebhookHandler_WithSecret** - Webhook security validation
8. **Additional edge case tests** - Various error scenarios

Run tests with:
```bash
go test -v ./internal/handlers/
```

## Jellyfin Configuration

### Step 1: Install Webhook Plugin

1. Open Jellyfin admin dashboard
2. Navigate to: Dashboard > Plugins > Catalog
3. Search for "Webhook"
4. Install the webhook plugin
5. Restart Jellyfin server

### Step 2: Configure Webhook

1. Navigate to: Dashboard > Plugins > Webhook
2. Add new webhook destination
3. Configure:
   - **Webhook URL**: `http://your-server:8080/webhook`
   - **Notification Type**: Check "Item Added"
   - **Item Type Filter**: Select "Movies" and "Episodes"
   - **Request Header** (if using secret):
     - Name: `X-Webhook-Secret`
     - Value: `your-secret-token`

### Step 3: Test Configuration

1. Add a new movie or episode to Jellyfin
2. Check bot logs for webhook receipt
3. Verify content is marked in database

## Integration with Task Group 5

The webhook handler is designed to integrate seamlessly with the notification broadcaster (Task Group 5):

```go
// Current implementation (Task Group 3)
slog.Info("New content ready for notification", "item_id", payload.ItemID)
// TODO: Task Group 5 - Send notification to subscribers here

// Future implementation (Task Group 5)
metadata := h.extractMetadata(&payload)
if err := h.notifier.BroadcastNotification(metadata); err != nil {
    slog.Error("Failed to broadcast notification", "error", err)
}
```

## Performance Considerations

- **Asynchronous Processing**: Consider processing webhooks in goroutines for high-volume servers
- **Rate Limiting**: Jellyfin may send multiple webhooks rapidly during library scans
- **Database Connection Pool**: GORM handles connection pooling automatically
- **Memory Usage**: Minimal - each webhook payload is ~1-2KB

## Troubleshooting

### Common Issues

1. **401 Unauthorized**
   - Check `WEBHOOK_SECRET` matches between bot and Jellyfin
   - Verify `X-Webhook-Secret` header is being sent

2. **Content Not Processed**
   - Check logs for content type filtering
   - Verify `NotificationType` is "ItemAdded"
   - Ensure `ItemType` is "Movie" or "Episode"

3. **Duplicate Notifications**
   - Database content tracking should prevent this
   - Check `content_caches` table for existing entries

4. **Server Not Starting**
   - Check port is not already in use
   - Verify `PORT` environment variable is set correctly

## Future Enhancements

Potential improvements for future versions:

1. **Webhook Queue**: Process webhooks asynchronously via message queue
2. **Batch Processing**: Handle multiple webhooks in batches
3. **Retry Logic**: Retry failed database operations
4. **Metrics**: Track webhook processing rates and errors
5. **IP Whitelisting**: Additional security layer beyond secret token
