# Spec Requirements: Interactive UI/UX Enhancement

## Initial Description

Enhance the bot's user experience by adding interactive buttons and menus to make it easier to use. Current state: Users must type commands like /start, /recent, /search, /mutedlist. Proposed improvements include:

1. Inline keyboard buttons for main navigation (welcome menu with buttons)
2. Menu button integration (Telegram's native menu button feature)
3. Enhanced welcome/start message with inline buttons
4. Possibly add action buttons to content items

The goal is to make the bot more intuitive and user-friendly by reducing the need to remember and type commands, while maintaining a clean interface.

**Project Context:**
- Go-based Telegram bot using github.com/go-telegram/bot library
- Already has working commands and inline buttons for mute/unmute features
- Production bot with active users
- Persian (Farsi) UI language

## Requirements Discussion

### First Round Questions

**Q1:** Should the main menu prioritize one-click access to Recent Content, Search, and Muted List through inline keyboard buttons on /start?
**Answer:** Confirmed - Add interactive welcome menu on /start with inline keyboard buttons for Recent Content, Search, Muted List, and Help.

**Q2:** Should we implement Telegram's native Menu Button API to show available commands?
**Answer:** Confirmed - Implement using Telegram's native Menu Button API to show available commands.

**Q3:** When users interact with navigation buttons, should they work exactly like typing the commands manually?
**Answer:** Confirmed - Buttons should work exactly like typing commands manually. No visual differences needed.

**Q4:** For search functionality, should the button trigger an interactive state or simply show instructions?
**Answer:** Keep it SIMPLE - Search button should show instructions explaining how to use /search command (not interactive state tracking).

**Q5:** Should we add an inline "Undo/Unmute" button immediately after a user mutes a series for quick recovery?
**Answer:** Option C selected - Add inline "Undo/Unmute" button immediately after user mutes a series. This allows quick recovery if they clicked by mistake.

**Q6:** What are the preferred Persian button labels?
**Answer:**
- Recent Content = "تازه‌ها" (not "محتوای اخیر")
- Search = "جستجو" (confirmed)
- Muted List = "سریال‌های مسدود شده" (confirmed)
- Help = "راهنما" (confirmed)

**Q7:** Should we maintain backward compatibility with command-based interface?
**Answer:** YES - Keep both button and command interfaces active. Users can use either method.

**Q8:** What features are explicitly OUT OF SCOPE?
**Answer:** The following are OUT OF SCOPE for this spec:
- Custom user preferences UI
- Notification settings UI
- Admin controls
- Any other advanced features not mentioned

### Existing Code to Reference

**Similar Features Identified:**
- Feature: Current mute/unmute button implementation - Path: `internal/telegram/callbacks.go`
- Feature: Inline keyboard usage - Path: `internal/telegram/handlers.go` (handleMutedList function)
- Feature: Existing command handlers - Path: `internal/telegram/handlers.go`

**Components to potentially reuse:**
- Inline keyboard creation patterns from existing mute/unmute buttons
- Callback handling structure from callbacks.go
- Command handler patterns from handlers.go

## Visual Assets

### Files Provided:
No visual assets provided.

### Visual Insights:
N/A - No visual files were found in the visuals folder.

## Requirements Summary

### Functional Requirements

**Main Menu Enhancement:**
- Interactive welcome menu on /start command
- Inline keyboard with four buttons: Recent Content ("تازه‌ها"), Search ("جستجو"), Muted List ("سریال‌های مسدود شده"), and Help ("راهنما")
- Buttons should trigger the same behavior as typing the corresponding commands

**Telegram Menu Button Integration:**
- Implement Telegram's native Menu Button API
- Display available commands in the menu button
- Standard Telegram UX pattern

**Search Button Behavior:**
- Simple approach: Show instructions on how to use /search command
- No interactive state tracking required
- Keeps implementation straightforward

**Mute Action Enhancement:**
- Add inline "Undo/Unmute" button immediately after user mutes a series
- Allows quick recovery from accidental mutes
- Single-click unmute functionality

**Language Requirements:**
- All button labels must be in Persian (Farsi)
- Use specific translations provided by user

**Backward Compatibility:**
- Maintain existing command-based interface
- Users can choose between buttons or typing commands
- Both methods should work identically

### Reusability Opportunities

**Existing Patterns to Follow:**
- Mute/unmute button implementation in `internal/telegram/callbacks.go`
- Inline keyboard creation in `internal/telegram/handlers.go` (handleMutedList)
- Command handler structure in `internal/telegram/handlers.go`

**Backend Patterns:**
- Callback handling mechanism already in place
- Command routing system already established
- Button interaction patterns already working in production

### Scope Boundaries

**In Scope:**
- Interactive welcome menu with inline keyboard buttons on /start
- Telegram Menu Button API integration
- Four main navigation buttons (Recent Content, Search, Muted List, Help)
- Undo/Unmute button after mute actions
- Persian language labels as specified
- Backward compatibility with existing commands

**Out of Scope:**
- Custom user preferences UI
- Notification settings UI
- Admin controls
- Action buttons on individual content items (beyond mute/unmute)
- Any advanced features not explicitly mentioned
- Interactive search state tracking

### Technical Considerations

**Integration Points:**
- Uses github.com/go-telegram/bot library
- Must integrate with existing callback handling in callbacks.go
- Must integrate with existing command handlers in handlers.go
- Production bot with active users - changes must be non-breaking

**Existing System Constraints:**
- Persian (Farsi) UI language requirement
- Must maintain current mute/unmute functionality
- Must preserve all existing command functionality

**Technology Preferences:**
- Follow existing Go code patterns in the codebase
- Use established inline keyboard patterns
- Leverage Telegram's native Menu Button API

**Similar Code Patterns to Follow:**
- Inline keyboard button creation from handleMutedList
- Callback data structure and handling from callbacks.go
- Command handler registration and routing patterns
