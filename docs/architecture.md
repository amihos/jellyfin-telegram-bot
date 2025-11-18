# Architecture: Jellyfin Telegram Bot

## Tech Stack Decision

### Language: Go 1.21+

**Chosen Library:** `github.com/go-telegram/bot`

### Rationale

#### Why Go?
- **Single Binary Deployment:** Compile once, deploy anywhere - no runtime dependencies needed
- **Excellent Concurrency:** Goroutines are perfect for handling webhooks and simultaneous notification broadcasting to multiple users
- **Low Resource Usage:** ~10-20MB RAM footprint (vs 50-100MB+ for Python/Node.js)
- **Built-in HTTP Server:** No need for external web framework
- **Fast Compilation:** Quick build times for rapid iteration
- **Production Ready:** Designed for long-running services and 24/7 operation
- **Strong Typing:** Catches bugs at compile time, reducing runtime errors

#### Why `github.com/go-telegram/bot`?
- **Zero Dependencies:** Even better deployment - no transitive dependencies to manage or update
- **383 Code Examples:** Extensive documentation with real-world usage examples
- **High Source Reputation:** Well-maintained and trusted library
- **Modern Context-Based API:** Clean design using Go's context patterns for cancellation and timeout handling
- **Built-in Webhook Support:** Native support for both polling and webhook modes
- **Middleware System:** Extensible architecture for adding cross-cutting concerns
- **Clean Handler Registration:** Simple API using `RegisterHandler` for command routing

### Complete Technology Stack

#### Core Dependencies
- **Language:** Go 1.21+ (latest stable recommended for `log/slog` support)
- **Telegram Bot Library:** `github.com/go-telegram/bot`
- **Database:** SQLite with `gorm.io/gorm` ORM
  - Driver: `gorm.io/driver/sqlite`
  - Migration support with `db.AutoMigrate()`
- **Configuration:** `github.com/joho/godotenv` for .env file support
- **Logging:** `log/slog` (Go 1.21+) for structured logging
- **Log Rotation:** `gopkg.in/natefinch/lumberjack.v2` for log file management

#### HTTP & API
- **Jellyfin API Client:** Custom HTTP client using Go's `net/http` package
- **Webhook Server:** Built-in `net/http` package (no framework needed)
- **JSON Parsing:** Built-in `encoding/json` package

#### Testing
- **Test Framework:** Go's built-in `testing` package
- **Assertions:** `github.com/stretchr/testify` for readable test assertions

## System Architecture

### High-Level Components

```
┌─────────────────────────────────────────────────────────────┐
│                     Jellyfin Server                          │
│  (Webhook Plugin configured to send notifications)          │
└──────────────────┬──────────────────────────────────────────┘
                   │ HTTP POST (webhook)
                   ▼
┌─────────────────────────────────────────────────────────────┐
│                  Jellyfin Telegram Bot                       │
│                                                              │
│  ┌────────────────────────────────────────────────────┐    │
│  │         Webhook Receiver (HTTP Server)              │    │
│  │  - Receives webhook notifications                   │    │
│  │  - Parses JSON payload                              │    │
│  │  - Filters content types (Movies/Episodes)          │    │
│  └──────┬─────────────────────────────────────────────┘    │
│         │                                                    │
│         ▼                                                    │
│  ┌────────────────────────────────────────────────────┐    │
│  │         Content Tracking Layer                      │    │
│  │  - Checks for duplicate notifications               │    │
│  │  - Marks content as notified                        │    │
│  └──────┬─────────────────────────────────────────────┘    │
│         │                                                    │
│         ▼                                                    │
│  ┌────────────────────────────────────────────────────┐    │
│  │         Jellyfin API Client                         │    │
│  │  - Fetches poster images                            │    │
│  │  - Retrieves content metadata                       │    │
│  │  - Handles /recent and /search queries              │    │
│  └──────┬─────────────────────────────────────────────┘    │
│         │                                                    │
│         ▼                                                    │
│  ┌────────────────────────────────────────────────────┐    │
│  │      Notification Broadcaster                       │    │
│  │  - Formats messages in Persian                      │    │
│  │  - Broadcasts to all subscribers                    │    │
│  │  - Handles rate limiting                            │    │
│  └──────┬─────────────────────────────────────────────┘    │
│         │                                                    │
│         │                                                    │
│  ┌────────────────────────────────────────────────────┐    │
│  │         Telegram Bot Handler                        │    │
│  │  - /start command (subscription)                    │    │
│  │  - /recent command (view recent content)            │    │
│  │  - /search command (search library)                 │    │
│  └──────┬─────────────────────────────────────────────┘    │
│         │                                                    │
│         │                                                    │
│  ┌────────────────────────────────────────────────────┐    │
│  │         Database Layer (SQLite + GORM)              │    │
│  │  - Subscriber management                            │    │
│  │  - Content cache (duplicate prevention)             │    │
│  └────────────────────────────────────────────────────┘    │
└──────────────────┬──────────────────────────────────────────┘
                   │ Telegram API calls
                   ▼
┌─────────────────────────────────────────────────────────────┐
│              Telegram API (Subscribed Users)                 │
└─────────────────────────────────────────────────────────────┘
```

### Directory Structure

```
jellyfin-telegram-bot/
├── cmd/
│   └── bot/
│       └── main.go              # Application entry point
├── internal/
│   ├── handlers/
│   │   ├── commands.go          # Telegram command handlers
│   │   └── webhook.go           # Webhook receiver handlers
│   ├── database/
│   │   ├── db.go                # Database connection & setup
│   │   └── repository.go        # Data access layer
│   ├── jellyfin/
│   │   ├── client.go            # Jellyfin API client
│   │   ├── images.go            # Image fetching
│   │   └── search.go            # Content search & recent queries
│   ├── telegram/
│   │   ├── bot.go               # Bot initialization
│   │   ├── formatter.go         # Message formatting (Persian)
│   │   └── broadcaster.go       # Notification broadcasting
│   └── config/
│       └── config.go            # Configuration management
├── pkg/
│   └── models/
│       ├── subscriber.go        # Subscriber data model
│       ├── content.go           # Content cache model
│       └── webhook.go           # Webhook payload models
├── docs/
│   ├── architecture.md          # This file
│   ├── deployment.md            # Deployment guide
│   └── api-integration.md       # Jellyfin API details
├── .env.example                 # Environment variables template
├── .gitignore                   # Git ignore patterns
├── go.mod                       # Go module definition
├── go.sum                       # Dependency checksums
└── README.md                    # Setup instructions
```

## Data Flow

### Notification Flow (Webhook → Users)
1. Jellyfin server adds new content (movie/episode)
2. Jellyfin webhook plugin sends HTTP POST to bot's `/webhook` endpoint
3. Webhook handler parses JSON payload
4. Content tracking layer checks if already notified
5. If new content:
   - Jellyfin API client fetches poster image
   - Notification formatter creates Persian message
   - Broadcaster sends to all active subscribers
   - Content marked as notified in database

### User Command Flow (/start, /recent, /search)
1. User sends command to Telegram bot
2. Bot handler receives and routes command
3. For `/start`:
   - Add user to subscribers database
   - Send Persian welcome message
4. For `/recent`:
   - Query Jellyfin API for recent items
   - Fetch poster images
   - Format and send to user
5. For `/search <query>`:
   - Query Jellyfin API with search term
   - Fetch poster images for results
   - Format and send to user

## Concurrency Model

Go's goroutines enable efficient concurrent processing:

1. **Webhook Processing:** Each webhook handled in separate goroutine
2. **Notification Broadcasting:** Parallel message sending to subscribers (with rate limiting)
3. **Graceful Shutdown:** Context-based cancellation ensures clean termination

## Error Handling Strategy

- **Database Errors:** Log and retry with exponential backoff
- **Jellyfin API Errors:**
  - Network errors: Retry with timeout
  - 401 Unauthorized: Log critical error about invalid API key
  - 404 Not Found: Handle gracefully, skip missing content
- **Telegram API Errors:**
  - User blocked bot: Mark subscriber as inactive
  - Rate limit: Queue and retry
  - Network timeout: Retry with backoff
- **Webhook Errors:** Log malformed payloads, return HTTP 400

## Security Considerations

- **Webhook Security:** Optional secret token validation
- **Environment Variables:** Sensitive data (tokens, API keys) stored in `.env` file, never committed
- **Database:** Local SQLite file with restricted permissions
- **API Authentication:** Jellyfin API key stored securely

## Scalability

While designed as a single-instance application, the architecture supports:
- **Horizontal Scaling:** Multiple bot instances with load-balanced webhook endpoint
- **Database Migration:** Easy switch from SQLite to PostgreSQL/MySQL
- **Caching:** Future addition of Redis for content metadata caching

## Performance Targets

- **Webhook Response Time:** < 100ms (acknowledge receipt)
- **Notification Delivery:** All subscribers notified within 10 seconds
- **Command Response Time:** < 2 seconds for /recent and /search
- **Memory Footprint:** < 50MB under normal operation
- **Uptime:** 99.9% (designed for 24/7 operation)
