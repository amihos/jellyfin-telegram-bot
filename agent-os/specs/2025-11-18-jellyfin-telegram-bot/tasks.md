# Task Breakdown: Jellyfin Telegram Bot

## Overview
Total Task Groups: 6
Estimated Total Tasks: ~45 sub-tasks

This is a greenfield project implementing a Telegram bot that monitors a Jellyfin media server and sends Persian-language notifications to subscribed users when new movies or TV episodes are added.

## Task List

### Project Foundation

#### Task Group 1: Project Setup and Architecture
**Dependencies:** None

- [ ] 1.0 Complete project foundation
  - [ ] 1.1 Choose and document tech stack
    - Decision: **Go** with `github.com/go-telegram/bot` library (CHOSEN - Updated based on latest research)
    - Rationale: Zero dependencies, 383 code examples, modern context-based API, built-in webhook support
    - Key advantages: Single binary deployment, excellent concurrency for webhooks, low resource usage (~10-20MB)
    - Library features: Middleware support, clean handler registration, high source reputation
    - Consider: Go's goroutines perfect for handling concurrent webhook processing and notification broadcasting
    - Document choice in: `docs/architecture.md`
  - [ ] 1.2 Initialize project structure
    - Create Go module: `go mod init jellyfin-telegram-bot`
    - Create directory structure: `/cmd/bot`, `/internal/handlers`, `/internal/database`, `/internal/jellyfin`, `/internal/telegram`, `/pkg/models`, `/docs`
    - Set up dependency management with `go.mod` and `go.sum`
    - Create `.env.example` with required environment variables
    - Add `.gitignore` for Go-specific files (vendor/, *.exe, *.test, .env)
  - [ ] 1.3 Configure environment variables
    - TELEGRAM_BOT_TOKEN
    - JELLYFIN_SERVER_URL
    - JELLYFIN_API_KEY
    - WEBHOOK_SECRET (optional, for security)
    - DATABASE_PATH or DATABASE_URL
    - PORT for webhook receiver
  - [ ] 1.4 Set up logging infrastructure
    - Use Go logging library: `log/slog` (Go 1.21+) or `logrus` or `zap`
    - Configure structured logging with timestamps and JSON formatting
    - Log levels: DEBUG, INFO, WARNING, ERROR
    - Separate log files for webhook events and bot commands
    - Rotate logs to prevent disk space issues using `lumberjack` or similar
  - [ ] 1.5 Create documentation structure
    - `README.md` with setup instructions
    - `docs/architecture.md` with system overview
    - `docs/deployment.md` with deployment guide
    - `docs/api-integration.md` for Jellyfin API details

**Acceptance Criteria:**
- Project structure created and documented
- Tech stack chosen and dependencies installable
- Environment configuration template ready
- Logging infrastructure operational
- Basic documentation in place

### Database Layer

#### Task Group 2: Data Models and Storage
**Dependencies:** Task Group 1

- [ ] 2.0 Complete database layer
  - [ ] 2.1 Write 2-8 focused tests for data storage
    - Limit to 2-8 highly focused tests maximum
    - Test only critical behaviors (e.g., subscriber add/remove, duplicate prevention, persistence)
    - Skip exhaustive edge case coverage
  - [ ] 2.2 Choose and initialize database
    - Options: SQLite (simple, file-based) OR PostgreSQL (production-ready)
    - Recommendation: SQLite for MVP, easy to migrate later
    - Use Go database library: `database/sql` with `github.com/mattn/go-sqlite3` driver OR `gorm.io/gorm` (ORM)
    - Create database schema file (SQL migrations)
    - Set up database connection module in `/internal/database`
  - [ ] 2.3 Create Subscriber model/table
    - Fields: id (primary key), chat_id (unique, indexed), username, first_name, subscribed_at (timestamp), is_active (boolean)
    - Indexes: chat_id (unique index)
    - Validation: chat_id must be unique and not null
  - [ ] 2.4 Implement subscriber management functions
    - add_subscriber(chat_id, username, first_name): Add new subscriber
    - remove_subscriber(chat_id): Remove/deactivate subscriber
    - get_all_active_subscribers(): Return list of active chat_ids
    - is_subscribed(chat_id): Check subscription status
  - [ ] 2.5 Create Content Cache model/table (optional but recommended)
    - Fields: id, jellyfin_id (unique), title, type, added_at (timestamp)
    - Purpose: Track already-notified content to prevent duplicates
    - Indexes: jellyfin_id (unique index)
  - [ ] 2.6 Implement content tracking functions
    - is_content_notified(jellyfin_id): Check if already notified
    - mark_content_notified(jellyfin_id, title, type): Record notification
  - [ ] 2.7 Ensure database layer tests pass
    - Run ONLY the 2-8 tests written in 2.1
    - Verify database operations work correctly
    - Do NOT run entire test suite at this stage

**Acceptance Criteria:**
- The 2-8 tests written in 2.1 pass
- Database schema created and migrations run successfully
- Subscriber CRUD operations work correctly
- Content tracking prevents duplicate notifications
- Database persists across application restarts

### Webhook Integration

#### Task Group 3: Jellyfin Webhook Receiver
**Dependencies:** Task Groups 1, 2

- [ ] 3.0 Complete webhook receiver
  - [ ] 3.1 Write 2-8 focused tests for webhook handling
    - Limit to 2-8 highly focused tests maximum
    - Test only critical behaviors (e.g., payload parsing, content type filtering, valid vs invalid webhooks)
    - Skip exhaustive testing of all webhook scenarios
  - [ ] 3.2 Set up HTTP server for webhook endpoint
    - Use Go's built-in `net/http` package (no framework needed) OR `gin-gonic/gin` for easier routing
    - Create POST endpoint: `/webhook` or `/jellyfin-webhook`
    - Handle JSON payload parsing using `encoding/json`
    - Add basic error handling with proper HTTP status codes
    - Run server in goroutine to allow concurrent bot polling
  - [ ] 3.3 Implement webhook payload parser
    - Parse Jellyfin webhook JSON structure
    - Extract fields: NotificationType, ItemType, Name, Overview, Year, SeriesName, SeasonNumber, EpisodeNumber, ItemId
    - Reference: Jellyfin webhook plugin documentation for payload structure
  - [ ] 3.4 Filter content types
    - Accept: NotificationType == "ItemAdded"
    - Accept ItemType: "Movie" and "Episode" only
    - Reject: "Series", "Season", "Audio", "Book", etc.
    - Log rejected content types for debugging
  - [ ] 3.5 Extract metadata for notifications
    - Movie: Extract title, overview (description), year, item_id
    - Episode: Extract series name, season number, episode number, episode title, overview, item_id
    - Handle missing fields gracefully with defaults
  - [ ] 3.6 Implement webhook security (optional but recommended)
    - Validate webhook secret token if Jellyfin plugin supports it
    - OR restrict by IP address/network
    - Prevent unauthorized webhook calls
  - [ ] 3.7 Integrate with content tracking
    - Check if content already notified using is_content_notified()
    - Skip duplicate notifications
    - Mark new content as notified
  - [ ] 3.8 Ensure webhook tests pass
    - Run ONLY the 2-8 tests written in 3.1
    - Verify payload parsing works correctly
    - Verify content filtering works
    - Do NOT run entire test suite at this stage

**Acceptance Criteria:**
- The 2-8 tests written in 3.1 pass
- Webhook endpoint receives and parses Jellyfin payloads
- Content type filtering works correctly
- Duplicate content detection prevents repeat notifications
- Error handling prevents crashes from malformed payloads

### Jellyfin API Integration

#### Task Group 4: Jellyfin API Client
**Dependencies:** Task Groups 1, 3

- [ ] 4.0 Complete Jellyfin API integration
  - [ ] 4.1 Write 2-8 focused tests for Jellyfin API client
    - Limit to 2-8 highly focused tests maximum
    - Test only critical behaviors (e.g., authentication, image fetching, search query)
    - Use mock responses to avoid requiring live Jellyfin server
    - Skip exhaustive API endpoint coverage
  - [ ] 4.2 Set up Jellyfin API client
    - Create custom HTTP client using Go's `net/http` package
    - OR use community library if available (search for Jellyfin Go client)
    - Configure authentication with API key in request headers
    - Handle connection errors and timeouts gracefully using `context.WithTimeout`
    - Create reusable client struct in `/internal/jellyfin`
  - [ ] 4.3 Implement image fetching function
    - Function: get_poster_image(item_id) -> image_bytes or image_url
    - Fetch primary poster image for movies/episodes
    - Use Jellyfin API endpoint: `/Items/{itemId}/Images/Primary`
    - Handle missing images with placeholder or skip image
    - Return image in format suitable for Telegram (bytes or URL)
  - [ ] 4.4 Implement recent content query
    - Function: get_recent_items(limit=20) -> list of content items
    - Query Jellyfin API for recently added movies and episodes
    - Sort by DateCreated descending
    - Filter to Movies and Episodes only
    - Return structured data: title, type, item_id, overview, rating, year
  - [ ] 4.5 Implement search function
    - Function: search_content(query, limit=10) -> list of content items
    - Query Jellyfin API with search term
    - Support both Persian and English queries
    - Search in: Name, Overview fields
    - Filter to Movies and Episodes only
    - Return same structured data as get_recent_items()
  - [ ] 4.6 Extract rating information
    - Parse CommunityRating (e.g., IMDB rating) from API responses
    - Parse OfficialRating (e.g., PG-13, TV-MA) if available
    - Format for display: "Rating: 8.5/10" or "Rating: N/A"
  - [ ] 4.7 Handle API errors gracefully
    - Network errors: Retry with exponential backoff
    - 401 Unauthorized: Log error about invalid API key
    - 404 Not Found: Handle missing items gracefully
    - Timeout: Set reasonable timeout (10-30 seconds)
  - [ ] 4.8 Ensure Jellyfin API tests pass
    - Run ONLY the 2-8 tests written in 4.1
    - Verify authentication works
    - Verify image fetching and content queries work
    - Do NOT run entire test suite at this stage

**Acceptance Criteria:**
- The 2-8 tests written in 4.1 pass
- Jellyfin API client authenticates successfully
- Poster images can be fetched and are Telegram-compatible
- Recent content query returns correct data
- Search functionality works with Persian and English
- Errors handled without crashing bot

### Telegram Bot Implementation

#### Task Group 5: Bot Commands and Notification System
**Dependencies:** Task Groups 1, 2, 3, 4

- [ ] 5.0 Complete Telegram bot implementation
  - [ ] 5.1 Write 2-8 focused tests for bot commands
    - Limit to 2-8 highly focused tests maximum
    - Test only critical behaviors (e.g., /start subscribes user, /recent returns content, /search handles query)
    - Use Telegram bot testing utilities or mocks
    - Skip exhaustive command permutation testing
  - [ ] 5.2 Initialize Telegram bot
    - Use library: `github.com/go-telegram/bot` (zero dependencies, modern API)
    - Install: `go get github.com/go-telegram/bot`
    - Initialize with context: `bot.New(token, bot.WithDefaultHandler(handler))`
    - Use `RegisterHandler` for specific commands (e.g., /start, /recent, /search)
    - Configure bot token from environment variables
    - Start bot with context for graceful shutdown: `b.Start(ctx)`
    - Supports both polling (recommended) and webhook modes
  - [ ] 5.3 Implement /start command
    - Handler function for /start
    - Call add_subscriber() with user's chat_id, username, first_name
    - Send welcome message in Persian:
      ```
      سلام! به ربات اطلاع‌رسانی جلیفین خوش آمدید.

      شما از این پس اطلاعیه‌های محتوای جدید را دریافت خواهید کرد.

      دستورات موجود:
      /start - عضویت در ربات
      /recent - مشاهده محتوای اخیر
      /search - جستجوی محتوا
      ```
    - Confirm subscription status
  - [ ] 5.4 Implement /recent command
    - Handler function for /recent
    - Call get_recent_items(limit=15)
    - For each item, send formatted message with:
      - Poster image (using get_poster_image())
      - Title in original language
      - Type in Persian: "فیلم" (Movie) or "سریال - فصل X قسمت Y" (Episode)
      - Description in original language
      - Rating: "امتیاز: 8.5/10"
    - Handle empty results: "محتوای اخیری یافت نشد"
    - Send multiple messages (one per item) or use album/carousel if supported
  - [ ] 5.5 Implement /search command
    - Handler function for /search expecting argument
    - Handle two cases:
      - No argument: Send help message "لطفاً عبارت جستجو را وارد کنید. مثال: /search interstellar"
      - With argument: Call search_content(query, limit=10)
    - Format results same as /recent
    - Handle no results: "نتیجه‌ای برای '{query}' یافت نشد"
    - Support Persian and English search terms
  - [ ] 5.6 Create notification message formatter
    - Function: format_notification(content_item) -> formatted_text
    - Persian type labels:
      - Movie: "🎬 فیلم جدید"
      - Episode: "📺 قسمت جدید"
    - Template:
      ```
      🎬 فیلم جدید

      نام: {title}
      توضیحات: {overview}
      امتیاز: {rating}
      ```
    - For episodes:
      ```
      📺 قسمت جدید

      سریال: {series_name}
      فصل {season_number} - قسمت {episode_number}
      نام قسمت: {episode_title}
      توضیحات: {overview}
      امتیاز: {rating}
      ```
    - Handle missing fields gracefully
    - Ensure proper RTL formatting for Persian text
  - [ ] 5.7 Implement broadcast notification function
    - Function: broadcast_notification(content_item)
    - Get all active subscribers using get_all_active_subscribers()
    - For each subscriber:
      - Send poster image (if available)
      - Send formatted notification message
      - Handle errors per-user (don't stop broadcast if one fails)
      - Log failed deliveries
    - Track broadcast completion
    - Handle rate limiting from Telegram API (max 30 messages/second)
  - [ ] 5.8 Integrate webhook with notifications
    - When webhook receives new content:
      - Extract metadata using webhook parser (Task 3.3)
      - Enrich with Jellyfin API data if needed (Task 4)
      - Call broadcast_notification()
    - Complete flow: Webhook -> Parse -> Fetch Image -> Broadcast
  - [ ] 5.9 Add error handling and user feedback
    - Handle bot blocked by user: Remove from subscribers or mark inactive
    - Handle network errors: Retry with backoff
    - Unknown commands: Send help message in Persian
    - Rate limit errors: Queue and retry
  - [ ] 5.10 Ensure bot command tests pass
    - Run ONLY the 2-8 tests written in 5.1
    - Verify /start, /recent, /search work correctly
    - Verify notifications format properly
    - Do NOT run entire test suite at this stage

**Acceptance Criteria:**
- The 2-8 tests written in 5.1 pass
- /start command subscribes users and sends welcome message
- /recent command displays recently added content with images
- /search command finds content with Persian and English queries
- Notifications broadcast to all subscribers when webhook triggered
- Persian text displays correctly with RTL support
- Error handling prevents bot crashes

### Testing and Deployment

#### Task Group 6: Integration Testing and Deployment Preparation
**Dependencies:** Task Groups 1-5

- [ ] 6.0 Complete testing and deployment preparation
  - [ ] 6.1 Review existing tests from Task Groups 1-5
    - Review the 2-8 tests written by database-engineer (Task 2.1)
    - Review the 2-8 tests written by webhook-engineer (Task 3.1)
    - Review the 2-8 tests written by api-engineer (Task 4.1)
    - Review the 2-8 tests written by bot-engineer (Task 5.1)
    - Total existing tests: approximately 8-32 tests
  - [ ] 6.2 Analyze test coverage gaps for THIS feature only
    - Identify critical end-to-end workflows lacking coverage:
      - Complete notification flow: Webhook -> Parse -> Jellyfin API -> Broadcast
      - User subscribes -> receives notification workflow
      - Search with Persian characters
      - Image fetching and delivery to Telegram
    - Focus ONLY on gaps related to Jellyfin Telegram Bot requirements
    - Do NOT assess entire application test coverage
    - Prioritize integration tests over additional unit tests
  - [ ] 6.3 Write up to 10 additional strategic tests maximum
    - Add maximum of 10 new integration/end-to-end tests
    - Focus on:
      - Complete webhook-to-notification pipeline (3-4 tests)
      - Persian text handling and RTL formatting (1-2 tests)
      - Error scenarios (API failures, blocked users) (2-3 tests)
      - Multi-user broadcast (1 test)
    - Do NOT write comprehensive edge case coverage
    - Skip load testing, stress testing, and performance tests for MVP
  - [ ] 6.4 Run feature-specific tests only
    - Run ONLY tests related to Jellyfin Telegram Bot (tests from 2.1, 3.1, 4.1, 5.1, and 6.3)
    - Expected total: approximately 18-42 tests maximum
    - Do NOT run entire application test suite (no suite exists yet - this is new project)
    - Verify all critical workflows pass
  - [ ] 6.5 Create deployment documentation
    - Document server requirements (Go version 1.21+ recommended, no other dependencies)
    - Document build process: `go build -o jellyfin-bot cmd/bot/main.go`
    - Document single binary deployment (just copy executable to server)
    - Document environment variable configuration
    - Document Jellyfin webhook plugin setup
    - Document Telegram bot creation process (BotFather)
    - Document database initialization steps
  - [ ] 6.6 Set up process management
    - Use systemd (Linux) for process management (recommended for Go binaries)
    - Create systemd service file: `jellyfin-bot.service`
    - Configure auto-restart on failure with `Restart=always`
    - Configure log rotation with journald or separate log file
    - Alternative: Run binary directly in screen/tmux for simple deployments
  - [ ] 6.7 Create deployment script
    - Script to install dependencies
    - Script to set up database
    - Script to run migrations
    - Script to start bot service
    - Script to verify deployment (health check)
  - [ ] 6.8 Document Jellyfin webhook configuration
    - Webhook URL format: `http://your-server:port/webhook`
    - Webhook events to enable: "Item Added"
    - Recommended: Add webhook secret for security
    - Test webhook delivery from Jellyfin admin panel
  - [ ] 6.9 Create monitoring and health check
    - Implement /health endpoint for bot status
    - Log metrics: notifications sent, active subscribers, errors
    - Set up alerts for critical errors (optional)
  - [ ] 6.10 Perform end-to-end testing
    - Test complete flow: Add content to Jellyfin -> Receive notification in Telegram
    - Test all commands with real Telegram bot
    - Test Persian text display on mobile devices
    - Test image delivery
    - Verify with multiple subscribers

**Acceptance Criteria:**
- All feature-specific tests pass (approximately 18-42 tests total)
- Critical end-to-end workflows validated
- No more than 10 additional tests added when filling testing gaps
- Deployment documentation complete and accurate
- Deployment script tested and functional
- Bot can be deployed to production server
- Health monitoring in place
- End-to-end testing confirms all requirements met

## Execution Order

Recommended implementation sequence:

1. **Project Foundation** (Task Group 1) - Set up project structure, tech stack, logging
2. **Database Layer** (Task Group 2) - Create subscriber management and content tracking
3. **Jellyfin API Integration** (Task Group 4) - Build API client for images and content queries
4. **Webhook Integration** (Task Group 3) - Implement webhook receiver and parser
5. **Telegram Bot Implementation** (Task Group 5) - Build bot commands and notification system
6. **Testing and Deployment** (Task Group 6) - Integration testing and deployment preparation

Note: Task Groups 3 and 4 can be developed in parallel after Task Group 2 is complete, as they have no direct dependencies on each other.

## Technical Stack (CHOSEN - Updated with Latest Research)

### Go Stack ✅
- **Language:** Go 1.21+ (latest stable recommended)
- **Telegram Bot Library:** `github.com/go-telegram/bot` (UPDATED CHOICE)
  - Zero dependencies (even better deployment!)
  - 383 documented code examples (vs 28 for alternatives)
  - High source reputation
  - Modern context-based API design
  - Built-in webhook and polling support
  - Middleware system for extensibility
  - Clean handler registration with `RegisterHandler`
- **Jellyfin API:** Custom HTTP client using `net/http` package (or community library if available)
- **Web Server:** Built-in `net/http` package (no framework needed) OR `github.com/gin-gonic/gin` for easier routing
- **Database:** SQLite with `github.com/mattn/go-sqlite3` driver (or `gorm.io/driver/sqlite` with GORM)
  - Recommended: Use GORM ORM for cleaner code
  - Simple setup: `gorm.Open(sqlite.Open("bot.db"), &gorm.Config{})`
  - Auto-migration support with `db.AutoMigrate(&Model{})`
- **ORM:** `gorm.io/gorm` (RECOMMENDED - well-documented, 381 code examples)
- **Configuration:** `github.com/joho/godotenv` for .env file support
- **Logging:** `log/slog` (Go 1.21+), `github.com/sirupsen/logrus`, or `go.uber.org/zap`
- **Testing:** Go's built-in `testing` package with `github.com/stretchr/testify` for assertions

### Why Go for This Project
- **Single Binary Deployment:** Compile once, deploy anywhere - no runtime dependencies
- **Zero-Dependency Bot Library:** No transitive dependencies to manage or update
- **Excellent Concurrency:** Goroutines perfect for handling webhooks + simultaneous notification broadcasting
- **Low Resource Usage:** ~10-20MB RAM footprint (vs 50-100MB+ for Python)
- **Built-in HTTP Server:** No need for external web framework
- **Fast Compilation:** Quick build times for rapid iteration
- **Production Ready:** Designed for long-running services and 24/7 operation
- **Strong Typing:** Catches bugs at compile time
- **Modern Patterns:** Context-based cancellation and graceful shutdown

## Modern Implementation Patterns (From Latest Documentation)

### Bot Initialization Pattern
Based on `github.com/go-telegram/bot` best practices:

```go
package main

import (
    "context"
    "os"
    "os/signal"

    "github.com/go-telegram/bot"
    "github.com/go-telegram/bot/models"
)

func main() {
    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
    defer cancel()

    opts := []bot.Option{
        bot.WithDefaultHandler(defaultHandler),
    }

    b, err := bot.New(os.Getenv("TELEGRAM_BOT_TOKEN"), opts...)
    if err != nil {
        panic(err)
    }

    // Register command handlers
    b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler)
    b.RegisterHandler(bot.HandlerTypeMessageText, "/recent", bot.MatchTypeExact, recentHandler)
    b.RegisterHandler(bot.HandlerTypeMessageText, "/search", bot.MatchTypePrefix, searchHandler)

    b.Start(ctx)
}
```

### GORM Database Setup Pattern
Based on GORM 2.0+ best practices:

```go
import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type Subscriber struct {
    gorm.Model
    ChatID    int64  `gorm:"uniqueIndex;not null"`
    Username  string
    FirstName string
    IsActive  bool `gorm:"default:true"`
}

type ContentCache struct {
    gorm.Model
    JellyfinID string `gorm:"uniqueIndex;not null"`
    Title      string
    Type       string
}

func initDB() (*gorm.DB, error) {
    db, err := gorm.Open(sqlite.Open("bot.db"), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    // Auto-migrate schema
    db.AutoMigrate(&Subscriber{}, &ContentCache{})

    return db, nil
}
```

### Webhook Handler Pattern
For concurrent webhook processing:

```go
func webhookHandler(w http.ResponseWriter, r *http.Request) {
    var payload JellyfinWebhook
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, "Invalid payload", http.StatusBadRequest)
        return
    }

    // Process webhook asynchronously
    go func() {
        ctx := context.Background()
        processNewContent(ctx, payload)
    }()

    w.WriteHeader(http.StatusOK)
}
```

### Telegram API Best Practices
From official Telegram Bot API documentation:

- **Webhook Setup:** Use `setWebhook` with proper HTTPS URL and secret token for security
- **File Uploads:** Photos < 5MB, other files < 20MB can use HTTP URLs
- **sendPhoto:** Include `caption` for image descriptions in Persian
- **Rate Limiting:** Max 30 messages/second to avoid API throttling
- **Graceful Shutdown:** Use context cancellation for clean bot termination

## Persian Language Considerations

- All user-facing messages must be in Persian/Farsi
- Commands remain in English: /start, /recent, /search
- Content titles and descriptions preserve original language from Jellyfin
- Ensure proper RTL (right-to-left) text rendering in Telegram
- Test Persian text display on both iOS and Android Telegram clients
- Use Persian numerals (۰-۹) or Western numerals (0-9) consistently - Western recommended for ratings
- Error messages in Persian

## Key Integration Points

1. **Jellyfin Webhook Plugin** -> Bot Webhook Receiver (Task 3)
2. **Bot Webhook Receiver** -> Jellyfin API Client (Task 4) for image fetching
3. **Bot Webhook Receiver** -> Database (Task 2) for duplicate checking
4. **Bot Webhook Receiver** -> Telegram Bot (Task 5) for broadcasting
5. **Telegram Bot Commands** -> Jellyfin API Client (Task 4) for /recent and /search
6. **Telegram Bot Commands** -> Database (Task 2) for subscriber management

## Success Metrics

- Bot successfully receives webhooks from Jellyfin
- All subscribers receive notifications within 10 seconds of content addition
- /start command successfully subscribes users
- /recent command displays last 15 items with images
- /search command returns relevant results for Persian and English queries
- Zero crashes during normal operation
- Persian text displays correctly with proper RTL formatting
- All critical tests pass (18-42 tests)
