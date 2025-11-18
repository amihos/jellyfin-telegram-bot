# Spec Requirements: Jellyfin Telegram Bot

## Initial Description
"i want to create a telegram bot for my jellyfin server so that when i put new content it sends nice messages ot the bot users letting them something was added"

## Requirements Discussion

### First Round Questions

**Q1:** What types of content should the bot notify about?
**Answer:** Movies and TV shows only

**Q2:** What information should be included in the notification messages?
**Answer:** Title, type (movie/TV show), poster image, description, rating

**Q3:** For TV shows, should the bot send one notification per series or one per episode?
**Answer:** Individual notification for every episode

**Q4:** How should users subscribe to notifications?
**Answer:** Users press /start (required by Telegram API)

**Q5:** Should there be any access control or approval process for new subscribers?
**Answer:** Open to anyone, no approval needed

**Q6:** How should the bot detect new content?
**Answer:** Jellyfin webhook plugin (already installed)

**Q7:** What commands should the bot support?
**Answer:** /start, /recent, /search, etc.

**Q8:** Should the bot support Persian/Farsi language?
**Answer:** Yes, all messages in Persian

### Existing Code to Reference

No similar existing features identified for reference. This is a new standalone application.

### Follow-up Questions - Persian/Farsi Implementation

**Follow-up 1:** For Persian language support - should everything be in Persian (commands, responses, content metadata), or should some elements remain in English?
**Answer:**
- Everything communicating messages should be in Persian
- If the title of movie/series is in English, show them in English (preserve original language for titles)
- Show all metadata in original language - English for English, Persian for Persian
- Keep commands in English as per Telegram API suggestion
- All error messages and help text in Persian

**Follow-up 2:** Are there features that should be excluded from the initial version?
**Answer:** Keep it simple - no recommendations, no content requests, no watch tracking

## Visual Assets

### Files Provided:
No visual assets provided.

### Visual Insights:
Not applicable - no visual files were provided.

## Requirements Summary

### Functional Requirements
- Monitor Jellyfin server for new content additions via webhook plugin
- Send formatted notifications to all subscribed Telegram users when new content is added
- Support two content types: Movies and TV Shows
- For TV shows: Send individual notification for each new episode
- Include in each notification:
  - Title (in original language)
  - Content type (Movie/TV Show/Episode)
  - Poster image
  - Description (in original language)
  - Rating
- User subscription via /start command (no approval required)
- Bot commands: /start, /recent, /search
- All user-facing messages, responses, and help text in Persian/Farsi
- Commands remain in English (as per Telegram API best practices)
- Content titles and metadata displayed in their original language

### Reusability Opportunities
Not applicable - this is a new standalone application with no existing similar features in the codebase.

### Scope Boundaries
**In Scope:**
- Webhook-based content detection from Jellyfin
- Telegram bot with /start, /recent, /search commands
- Formatted notifications with poster images
- Persian language interface (messages, responses, help text)
- Movies and TV shows (individual episode notifications)
- User subscription management
- Rating display

**Out of Scope:**
- Content recommendation system
- User content requests
- Watch tracking/progress monitoring
- User approval/access control
- Advanced search filters
- Multi-language support beyond Persian
- Content categories beyond movies and TV shows

### Technical Considerations
- Jellyfin webhook plugin already installed and available
- Telegram Bot API requirements dictate English command structure
- Need to preserve original language for content titles and metadata
- Persian text formatting and RTL (right-to-left) display support required
- Poster images need to be fetched from Jellyfin and sent via Telegram
- Bot must maintain list of subscribed users (persist across restarts)
- Webhook payload parsing from Jellyfin required
