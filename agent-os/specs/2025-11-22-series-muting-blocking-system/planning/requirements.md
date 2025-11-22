# Spec Requirements: Series Muting/Blocking System

## Initial Description
Allow users to opt-out of notifications for specific TV series they're not interested in. When a user receives a notification about a series episode, they should be able to tell the bot "I don't want to follow this series" so they stop receiving future notifications for new episodes of that series.

This is a user preference system focused on content filtering at the series level.

## Requirements Discussion

### First Round Questions

**Q1:** I assume this should work at the series level (all episodes of a show), not individual season level. Is that correct, or should users be able to mute specific seasons?
**Answer:** Series level only - users can mute/unmute an entire series.

**Q2:** For the user action, I'm thinking an inline button on each notification message (like "دنبال نکردن" / "Unfollow") would be the most convenient. Should we go with that approach, or would you prefer a separate command like `/mute SeriesName`?
**Answer:** Follow-up needed to clarify exact placement.

**Q3:** I assume the muted state should persist permanently until the user explicitly unmutes it. Should there be a way to view and manage their muted series list (like `/mutedlist` command)?
**Answer:** Yes - permanent persistence until unmuted. Yes - include a `/mutedlist` command to view all muted series.

**Q4:** For reversing the action, should the muted series list allow inline unmute buttons, or require a command like `/unmute SeriesName`?
**Answer:** Inline unmute buttons in the `/mutedlist` response for easy management.

**Q5:** I'm assuming we should store this as user preferences in a local SQLite database. Does that align with your current data storage approach?
**Answer:** Yes - SQLite database for user preferences is correct.

**Q6:** Should this only apply to new episode notifications, or also affect other potential series-related messages (like season finale notifications, series cancellation news, etc.)?
**Answer:** Only affects new episode notifications - no other message types exist currently.

**Q7:** When a user mutes a series, should we send a confirmation message (like "✓ You will no longer receive notifications for [Series Name]"), or just silently update the preference?
**Answer:** Send confirmation message with the series name when muted. Also send confirmation when unmuted.

**Q8:** Is there anything that should explicitly NOT be included in this feature (like muting specific genres, actors, or other metadata-based filtering)?
**Answer:** Out of scope: genre/actor/metadata filtering, muting specific seasons, time-based muting schedules.

### Existing Code to Reference

**Similar Features Identified:**
- Feature: Webhook notification system - Path: `internal/handlers/webhook.go`
  - Shows current notification message structure
  - Demonstrates how episode data is processed
  - Contains message formatting logic

- Feature: Telegram message handling - Path: `internal/telegram/notifications.go`
  - Shows how messages are sent to users
  - Demonstrates message formatting patterns

- Feature: Callback handling - Path: `internal/telegram/bot.go`
  - Shows existing inline button handling patterns
  - Demonstrates how to process button clicks

### Follow-up Questions

**Follow-up 1:** For the "دنبال نکردن" button, should we:
- Option A: Integrate it directly into the notification message itself (single message with inline keyboard button below the content)?
- Option B: Send it as a separate follow-up message with the button (keeping notification message clean)?

**Answer:** Option A - Integrate the "دنبال نکردن" button directly into the notification message itself (single message with inline keyboard button below the content). This means we'll modify the existing notification message structure to include the inline button.

## Visual Assets

### Files Provided:
No visual assets provided.

### Visual Insights:
N/A

## Requirements Summary

### Functional Requirements

**Core Functionality:**
- Users can mute specific TV series to stop receiving new episode notifications
- Muting is done via an inline "دنبال نکردن" (Unfollow) button integrated into each notification message
- Muted state persists permanently in SQLite database until explicitly unmuted
- Users can view all muted series via `/mutedlist` command
- Users can unmute series via inline buttons in the `/mutedlist` response
- System sends confirmation messages when series are muted or unmuted

**User Actions Enabled:**
- Click "دنبال نکردن" button on any episode notification to mute that series
- Use `/mutedlist` command to view all currently muted series
- Click unmute button next to any series in the muted list to re-enable notifications
- Receive confirmation feedback for both mute and unmute actions

**Data to be Managed:**
- User ID to series ID mappings in SQLite database
- Series metadata (name, Jellyfin series ID) for display purposes
- Persistent storage of mute preferences per user

### Reusability Opportunities

**Components that exist already:**
- Webhook notification message structure in `internal/handlers/webhook.go`
- Telegram message sending patterns in `internal/telegram/notifications.go`
- Inline button handling logic in `internal/telegram/bot.go`
- SQLite database setup and connection patterns (existing in codebase)

**Backend patterns to follow:**
- Mirror existing callback handling approach for button clicks
- Follow existing message formatting patterns for consistency
- Use established database query patterns

**Similar features to model after:**
- Existing webhook notification flow provides foundation for modification
- Current inline button handling (if any exists) provides pattern for new buttons

### Scope Boundaries

**In Scope:**
- Series-level muting (entire show)
- Inline "دنبال نکردن" button integrated into notification messages
- `/mutedlist` command to view muted series
- Inline unmute buttons in muted list
- SQLite persistence of user preferences
- Confirmation messages for mute/unmute actions
- Filtering new episode notifications based on muted series

**Out of Scope:**
- Season-level muting (muting specific seasons only)
- Genre-based filtering
- Actor-based filtering
- Metadata-based filtering (release year, rating, etc.)
- Time-based muting schedules
- Temporary muting with auto-expiration
- Muting other notification types (only new episode notifications are affected)

### Technical Considerations

**Integration Points:**
- Modify existing notification message structure in `internal/handlers/webhook.go` to include inline button
- Add callback handler for "دنبال نکردن" button clicks
- Implement new `/mutedlist` command handler
- Create database schema for user mute preferences
- Add filtering logic before sending notifications to check mute status

**Existing System Constraints:**
- Must work with current SQLite database setup
- Must integrate with existing Telegram bot message flow
- Must maintain Persian language UI consistency
- Should follow existing code patterns in referenced files

**Technology Stack:**
- Go (existing language)
- SQLite (existing database)
- Telegram Bot API (existing integration)
- Jellyfin webhook data structure (existing format)

**Similar Code Patterns to Follow:**
- Use existing webhook handling patterns from `internal/handlers/webhook.go`
- Follow message formatting conventions from `internal/telegram/notifications.go`
- Mirror callback handling approach from `internal/telegram/bot.go`
- Maintain consistency with existing Persian language UI elements
