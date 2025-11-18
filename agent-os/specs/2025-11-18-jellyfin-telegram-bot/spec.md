# Specification: Jellyfin Telegram Bot

## Goal
Create a Telegram bot that monitors a Jellyfin media server and sends Persian-language notifications to subscribed users whenever new movies or TV episodes are added, with support for browsing recent content and searching the library.

## User Stories
- As a Jellyfin server user, I want to receive automatic notifications when new content is added so that I know when new movies or episodes are available to watch
- As a viewer, I want to search and browse recently added content through a Persian-language bot interface so that I can discover what's available on the server

## Specific Requirements

**Webhook Integration for Content Detection**
- Receive webhook notifications from Jellyfin's webhook plugin when new content is added
- Parse webhook payload to extract content metadata (title, type, description, rating, poster URL)
- Handle two content types: Movies and TV Shows
- For TV shows, process individual episodes separately (one notification per episode, not per series)
- Extract original-language titles and metadata from webhook data

**User Subscription Management**
- Store subscribed users' chat_ids when they invoke /start command
- Persist subscription list across bot restarts (database or file storage)
- No approval process required - automatic subscription on /start
- Allow users to unsubscribe (optional feature for future enhancement)
- Maintain list of active subscribers for broadcast notifications

**Notification Message Format**
- Send formatted notification to all subscribers when new content detected
- Include poster image fetched from Jellyfin server and sent via Telegram
- Display content title in its original language (English for English content, Persian for Persian content)
- Show content type indicator (Movie/TV Show/Episode) in Persian
- Include description text in original language
- Display rating information
- Format messages for readability with proper Persian RTL (right-to-left) text support

**Start Command (/start)**
- Subscribe user to notifications when they invoke /start
- Save user's chat_id to subscriber list
- Send welcome message in Persian explaining bot features and available commands
- Confirm subscription status to user

**Recent Command (/recent)**
- Display list of recently added content (last 10-20 items)
- Show same metadata as notifications: title, type, poster thumbnail, description, rating
- Format as scrollable list with Persian labels
- Content sorted by date added (newest first)

**Search Command (/search)**
- Accept search query from user (supports both Persian and English search terms)
- Query Jellyfin API to search for movies and TV shows matching the query
- Return search results with title, type, poster, and brief description
- Handle no results found with Persian message
- Limit results to reasonable number (e.g., 10 items) to avoid message flooding

**Persian Language Interface**
- All bot messages, responses, and UI text in Persian/Farsi
- Error messages in Persian
- Help text and command descriptions in Persian
- Commands themselves remain in English (/start, /recent, /search) per Telegram API conventions
- Content titles and metadata preserve original language from Jellyfin

**Jellyfin API Integration**
- Connect to Jellyfin server API for retrieving content details
- Fetch poster images from Jellyfin media library
- Query library for /recent and /search commands
- Handle API authentication with Jellyfin server
- Handle API errors gracefully with user-friendly Persian error messages

## Visual Design
No visual mockups provided - bot will use standard Telegram message formatting with text and images.

## Existing Code to Leverage

**No Existing Code Available**
- This is a new greenfield project with no existing codebase
- No similar features or components found to reuse
- Consider using established Telegram bot libraries for the chosen programming language (e.g., python-telegram-bot for Python, telegraf for Node.js)
- Leverage existing Jellyfin API client libraries if available for the chosen language
- Follow standard patterns for webhook receivers and REST API integrations

## Out of Scope
- Personalized content recommendations based on user preferences or viewing history
- User-initiated content requests or download requests
- Watch status tracking or progress monitoring
- User approval or access control system for subscriptions
- Advanced search filters (genre, year, rating range, etc.)
- Multi-language support beyond Persian (interface) and original content languages
- Content categories beyond movies and TV shows (music, audiobooks, etc.)
- Batch operations or bulk notifications
- User preference settings or customization options
- Integration with other media servers besides Jellyfin
