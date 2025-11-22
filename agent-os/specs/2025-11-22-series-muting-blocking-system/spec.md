# Specification: Series Muting/Blocking System

## Goal
Allow users to opt-out of notifications for specific TV series by clicking an inline button on episode notifications, with persistent storage in SQLite and the ability to view and manage their muted series list through a dedicated command.

## User Stories
- As a bot subscriber, I want to click a "دنبال نکردن" button on episode notifications so that I can stop receiving future notifications for that series
- As a user who has muted series, I want to use a /mutedlist command so that I can view all my muted series and unmute them with inline buttons

## Specific Requirements

**Inline Mute Button on Episode Notifications**
- Add inline keyboard with "دنبال نکردن" button to all episode notifications sent via BroadcastNotification
- Button should be integrated directly into notification message structure (not a separate follow-up message)
- Button callback data should encode series ID and series name for processing
- Only appear on Episode type notifications, not Movie notifications
- Button should disappear or be disabled after being clicked to prevent duplicate mute actions

**Series Mute Database Schema**
- Create new GORM model MutedSeries with fields: ChatID (int64), SeriesID (string), SeriesName (string), timestamps
- Use composite unique index on (ChatID, SeriesID) to prevent duplicate mutes
- Store Jellyfin series ID from webhook payload SeriesName field mapping
- Follow existing database patterns in models/subscriber.go and database/db.go
- Auto-migrate schema in database/db.go NewDB function alongside existing models

**Callback Handler for Mute Button**
- Register callback query handler in telegram/bot.go using bot.RegisterHandler with HandlerTypeCallbackQuery
- Extract series ID and name from callback data payload
- Insert mute preference into database using new database method
- Send Persian confirmation message to user: "شما دیگر اعلان‌های [Series Name] را دریافت نخواهید کرد"
- Answer callback query to remove loading state from button
- Update original message to show button as disabled or with different text

**Filter Notifications Based on Mute Preferences**
- Modify BroadcastNotification in telegram/notifications.go before subscriber loop
- Query database for users who have muted the series (using SeriesName from NotificationContent)
- Exclude muted users from subscriber list for this specific notification
- Keep existing rate limiting and error handling logic intact
- Log filtered user count in broadcast statistics

**/mutedlist Command Handler**
- Register new command handler in telegram/bot.go for /mutedlist
- Query database for all muted series for the requesting user's chat ID
- Format response message in Persian with list of muted series names
- Include inline "رفع مسدودیت" (Unmute) button next to each series in the list
- Handle empty list case with message: "شما هیچ سریالی را مسدود نکرده‌اید"
- Encode series ID in callback data for unmute action

**Callback Handler for Unmute Button**
- Register separate callback query handler for unmute actions
- Extract series ID from callback data and delete from database using chat ID
- Send Persian confirmation message: "[Series Name] از لیست مسدودی‌ها حذف شد"
- Refresh the /mutedlist message to remove the unmuted series from display
- Answer callback query to provide user feedback

**Database Operations Layer**
- Create database/muted_series.go with methods: AddMutedSeries, RemoveMutedSeries, GetMutedSeriesByUser, IsSeriesMuted
- Follow error handling patterns from database/subscriber.go
- Return GORM errors for not found cases
- Use transactions for consistency where appropriate
- Add logging using slog for all database operations

**Series ID Extraction from Webhook**
- Webhook payload contains SeriesName field for episodes but not explicit SeriesID
- Use SeriesName as unique identifier for muting (store as SeriesID in database)
- Ensure consistency between notification filtering and mute storage by using same SeriesName field
- Handle edge case where SeriesName might be empty or "Unknown Series" by preventing mute action

**Confirmation Message Formatting**
- Mute confirmation: "✓ شما دیگر اعلان‌های [Series Name] را دریافت نخواهید کرد"
- Unmute confirmation: "✓ [Series Name] از لیست مسدودی‌ها حذف شد"
- Use Persian right-to-left text formatting
- Include checkmark emoji for visual confirmation
- Keep messages concise and consistent with existing Persian UI patterns

**Help Command Update**
- Add /mutedlist to help message in defaultHandler function in telegram/bot.go
- Add description: "مشاهده سریال‌های مسدود شده" next to /mutedlist command
- Update /start welcome message to include new command
- Maintain Persian language consistency

## Visual Design
No visual assets provided.

## Existing Code to Leverage

**internal/telegram/notifications.go - BroadcastNotification function**
- Shows subscriber iteration pattern for broadcasting notifications
- Demonstrates rate limiting with time.Sleep(35 * time.Millisecond)
- Contains error handling for blocked users and logging patterns
- Use as template for adding mute filtering before subscriber loop

**internal/database/subscriber.go - Database operation patterns**
- Shows GORM query patterns: Where(), FirstOrCreate(), Update()
- Demonstrates error wrapping with fmt.Errorf
- Uses gorm.ErrRecordNotFound for not found cases
- Follow same structure for muted_series.go database methods

**pkg/models/subscriber.go - GORM model structure**
- Shows model definition with gorm tags and JSON tags
- Demonstrates uniqueIndex usage for ChatID field
- Shows TableName() method pattern for custom table names
- Replicate structure for MutedSeries model with composite index

**internal/telegram/bot.go - Handler registration patterns**
- Shows bot.RegisterHandler usage for command handlers
- Demonstrates HandlerTypeMessageText with MatchTypeExact and MatchTypePrefix
- Use same pattern to register HandlerTypeCallbackQuery for inline button callbacks
- Follow existing logging and error handling patterns

**internal/handlers/webhook.go - NotificationContent structure**
- Contains SeriesName, SeasonNumber, EpisodeNumber fields for episodes
- Shows content type differentiation between Movie and Episode
- Use SeriesName field as series identifier for mute matching
- Leverage existing structure without modification

## Out of Scope
- Season-level muting (only entire series can be muted)
- Genre-based filtering or actor-based filtering
- Metadata-based filtering such as release year or rating
- Time-based muting schedules or temporary muting with auto-expiration
- Muting notification types other than new episode notifications
- Muting movies (only TV series/episodes)
- Batch mute/unmute operations
- Search or filter functionality within muted list
- Export or import of muted preferences
- Series recommendation system based on muted preferences
