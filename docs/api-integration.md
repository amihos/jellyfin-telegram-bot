# Jellyfin API Integration

This document provides details about integrating with the Jellyfin API for the Telegram bot.

## Overview

The bot uses the Jellyfin REST API to:
1. Fetch poster images for content
2. Retrieve recently added items
3. Search the media library
4. Extract metadata (ratings, descriptions, etc.)

## Authentication

### API Key Authentication

Jellyfin uses API key authentication for programmatic access:

```
X-MediaBrowser-Token: your_api_key_here
```

All API requests must include this header.

### Getting an API Key

1. Log in to Jellyfin web interface
2. Navigate to Dashboard → API Keys
3. Click "+" to create new key
4. Name it (e.g., "Telegram Bot")
5. Copy the generated key

## Base URL Structure

```
{JELLYFIN_SERVER_URL}/[endpoint]
```

Example:
```
http://192.168.1.100:8096/Users/{userId}/Items
```

## Key API Endpoints

### 1. Get Items (Recently Added)

**Endpoint**: `GET /Users/{userId}/Items`

**Purpose**: Retrieve recently added movies and episodes

**Query Parameters**:
- `SortBy=DateCreated` - Sort by creation date
- `SortOrder=Descending` - Newest first
- `IncludeItemTypes=Movie,Episode` - Filter content types
- `Recursive=true` - Search all libraries
- `Limit=20` - Number of results
- `Fields=Overview,CommunityRating,OfficialRating,PremiereDate,ProviderIds` - Metadata to include

**Example Request**:
```
GET /Users/{userId}/Items?SortBy=DateCreated&SortOrder=Descending&IncludeItemTypes=Movie,Episode&Recursive=true&Limit=20&Fields=Overview,CommunityRating,OfficialRating
```

**Example Response**:
```json
{
  "Items": [
    {
      "Name": "Interstellar",
      "Id": "abc123",
      "Type": "Movie",
      "Overview": "A team of explorers travel through a wormhole...",
      "CommunityRating": 8.6,
      "OfficialRating": "PG-13",
      "ProductionYear": 2014,
      "PremiereDate": "2014-11-07T00:00:00.0000000Z",
      "ImageTags": {
        "Primary": "def456"
      }
    }
  ],
  "TotalRecordCount": 150
}
```

### 2. Search Items

**Endpoint**: `GET /Users/{userId}/Items`

**Purpose**: Search for movies and episodes by name

**Query Parameters**:
- `SearchTerm={query}` - Search query
- `IncludeItemTypes=Movie,Episode` - Filter types
- `Recursive=true` - Search all libraries
- `Limit=10` - Maximum results
- `Fields=Overview,CommunityRating,OfficialRating` - Metadata

**Example Request**:
```
GET /Users/{userId}/Items?SearchTerm=interstellar&IncludeItemTypes=Movie,Episode&Recursive=true&Limit=10&Fields=Overview,CommunityRating
```

**Note**: Jellyfin search supports Unicode, so Persian queries work:
```
GET /Users/{userId}/Items?SearchTerm=میان‌ستاره‌ای&...
```

### 3. Get Item Images

**Endpoint**: `GET /Items/{itemId}/Images/Primary`

**Purpose**: Fetch poster/thumbnail image for content

**Query Parameters**:
- `maxHeight=600` - Resize image (optional)
- `maxWidth=400` - Resize image (optional)
- `quality=90` - JPEG quality (optional)

**Example Request**:
```
GET /Items/abc123/Images/Primary?maxHeight=600&quality=90
```

**Response**: Binary image data (JPEG/PNG)

**Headers to Send**:
```
X-MediaBrowser-Token: your_api_key
```

### 4. Get User ID

**Endpoint**: `GET /Users/Public`

**Purpose**: Get list of public users (to find userId for queries)

**Example Request**:
```
GET /Users/Public
```

**Example Response**:
```json
[
  {
    "Name": "admin",
    "ServerId": "server123",
    "Id": "user456",
    "HasPassword": true,
    "HasConfiguredPassword": true
  }
]
```

**Note**: Typically use the first user's ID for system-wide queries.

## Webhook Payload Structure

When Jellyfin sends webhook notifications, the payload structure is:

```json
{
  "NotificationType": "ItemAdded",
  "Timestamp": "2024-01-15T10:30:00Z",
  "ServerId": "server123",
  "ServerName": "My Jellyfin Server",
  "ServerUrl": "http://192.168.1.100:8096",
  "ServerVersion": "10.8.13",
  "ItemId": "abc123def456",
  "ItemName": "Interstellar",
  "ItemType": "Movie",
  "Year": 2014,
  "Overview": "A team of explorers travel through a wormhole in space...",
  "ItemPath": "/media/movies/Interstellar (2014)/Interstellar.mkv",
  "UserName": "admin",
  "UserId": "user456"
}
```

### For Episodes:

```json
{
  "NotificationType": "ItemAdded",
  "ItemType": "Episode",
  "ItemName": "Pilot",
  "SeriesName": "Breaking Bad",
  "SeasonNumber": 1,
  "EpisodeNumber": 1,
  "Overview": "Walter White, a chemistry teacher...",
  "ItemId": "episode789",
  ...
}
```

### Webhook Fields to Extract

**For Movies**:
- `ItemId` - For fetching images
- `ItemName` - Title
- `ItemType` - "Movie"
- `Year` - Release year
- `Overview` - Description

**For Episodes**:
- `ItemId` - For fetching images
- `ItemName` - Episode title
- `ItemType` - "Episode"
- `SeriesName` - TV show name
- `SeasonNumber` - Season number
- `EpisodeNumber` - Episode number
- `Overview` - Episode description

## Error Handling

### Common Status Codes

- **200 OK**: Success
- **401 Unauthorized**: Invalid or missing API key
- **404 Not Found**: Item/resource not found
- **500 Internal Server Error**: Jellyfin server error
- **503 Service Unavailable**: Jellyfin temporarily unavailable

### Best Practices

1. **Retry Logic**: Implement exponential backoff for 5xx errors
2. **Timeouts**: Set reasonable timeouts (10-30 seconds)
3. **Rate Limiting**: Respect Jellyfin server load
4. **Error Logging**: Log all API errors for debugging
5. **Graceful Degradation**: Continue operation if image fetch fails

## Example Go Implementation

### API Client Structure

```go
type JellyfinClient struct {
    serverURL string
    apiKey    string
    client    *http.Client
}

func NewJellyfinClient(serverURL, apiKey string) *JellyfinClient {
    return &JellyfinClient{
        serverURL: serverURL,
        apiKey:    apiKey,
        client: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}
```

### Making Authenticated Requests

```go
func (c *JellyfinClient) doRequest(ctx context.Context, method, path string) (*http.Response, error) {
    req, err := http.NewRequestWithContext(ctx, method, c.serverURL+path, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Set("X-MediaBrowser-Token", c.apiKey)
    req.Header.Set("Accept", "application/json")

    return c.client.Do(req)
}
```

### Fetching Image

```go
func (c *JellyfinClient) GetPosterImage(ctx context.Context, itemID string) ([]byte, error) {
    path := fmt.Sprintf("/Items/%s/Images/Primary?maxHeight=600&quality=90", itemID)

    resp, err := c.doRequest(ctx, "GET", path)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to fetch image: status %d", resp.StatusCode)
    }

    return io.ReadAll(resp.Body)
}
```

## Testing with curl

### Test Authentication

```bash
curl -H "X-MediaBrowser-Token: your_api_key" \
     http://your-server:8096/System/Info
```

### Test Recent Items

```bash
curl -H "X-MediaBrowser-Token: your_api_key" \
     "http://your-server:8096/Users/user_id/Items?SortBy=DateCreated&SortOrder=Descending&IncludeItemTypes=Movie,Episode&Limit=5&Recursive=true"
```

### Test Image Fetch

```bash
curl -H "X-MediaBrowser-Token: your_api_key" \
     http://your-server:8096/Items/item_id/Images/Primary \
     --output poster.jpg
```

### Test Search

```bash
curl -H "X-MediaBrowser-Token: your_api_key" \
     "http://your-server:8096/Users/user_id/Items?SearchTerm=test&IncludeItemTypes=Movie,Episode&Limit=10"
```

## Rate Limiting Considerations

Jellyfin doesn't enforce strict rate limits, but:

1. **Be Considerate**: Don't hammer the API
2. **Cache When Possible**: Store metadata to reduce requests
3. **Batch Operations**: Group requests when feasible
4. **Monitor Load**: Check Jellyfin server resources

## References

- [Jellyfin API Documentation](https://api.jellyfin.org/)
- [Jellyfin Webhook Plugin](https://github.com/jellyfin/jellyfin-plugin-webhook)
- [Jellyfin OpenAPI Spec](https://api.jellyfin.org/openapi/api.html)

## Troubleshooting

### Authentication Failures

**Problem**: Getting 401 Unauthorized

**Solutions**:
- Verify API key is correct
- Check API key is not expired/revoked
- Ensure header name is exact: `X-MediaBrowser-Token`

### Image Not Found

**Problem**: Getting 404 for image requests

**Solutions**:
- Check if item has `ImageTags.Primary` in metadata
- Try alternative image types: `Backdrop`, `Thumb`
- Some items may not have images

### Search Not Working with Persian

**Problem**: Persian search returns no results

**Solutions**:
- Ensure content has Persian titles/metadata in Jellyfin
- Try searching original (English) titles
- Verify Jellyfin metadata providers include Persian sources

### Slow Response Times

**Problem**: API requests taking too long

**Solutions**:
- Check Jellyfin server load
- Reduce image quality/size parameters
- Implement caching for frequently accessed data
- Consider using Jellyfin CDN features if available
