# Specification: Interactive UI/UX Enhancement

## Goal
Transform the bot from a command-driven interface to an intuitive, button-based user experience by adding interactive welcome menus, Telegram menu button integration, and improved mute/unmute workflows for easier navigation and reduced cognitive load.

## User Stories
- As a bot user, I want to access main features through buttons instead of typing commands so that I can navigate the bot more intuitively
- As a user who accidentally mutes a series, I want an immediate undo button so that I can quickly recover without navigating to the muted list

## Specific Requirements

**Interactive Welcome Menu on /start**
- Add inline keyboard with 4 buttons arranged in 2x2 grid layout
- Button labels: "تازه‌ها" (Recent Content), "جستجو" (Search), "سریال‌های مسدود شده" (Muted List), "راهنما" (Help)
- Each button triggers the same behavior as typing the corresponding command
- Maintain existing welcome message text above the buttons
- Use callback data format: "nav:recent", "nav:search", "nav:mutedlist", "nav:help"

**Navigation Button Callback Handlers**
- Create unified navigation callback handler with prefix "nav:"
- Recent button executes same logic as /recent command handler
- Search button sends instructions: "لطفاً عبارت جستجو را وارد کنید. مثال: /search interstellar"
- Muted List button executes same logic as /mutedlist command handler
- Help button displays available commands list
- All handlers must answer callback query to remove loading state

**Telegram Menu Button Integration**
- Implement SetMyCommands API to register bot commands
- Register commands: /start, /recent, /search, /mutedlist
- Set Persian descriptions for each command
- Initialize during bot startup in NewBot function

**Undo/Unmute Button After Mute Action**
- After successful mute, send confirmation message with inline "رفع مسدودیت" (Unmute) button
- Callback data format: "undo_mute:{seriesName}"
- Button triggers unmute operation for the just-muted series
- Reuse existing unmute logic from handleUnmuteCallback
- Update button to "✓ رفع مسدودیت شد" after successful unmute

**Persian Language Labels**
- All button labels must use specified Persian text
- Consistent RTL text direction handling
- Follow existing Persian text patterns in codebase
- Error messages and confirmations in Persian

**Backward Compatibility**
- Keep all existing command handlers active (/start, /recent, /search, /mutedlist)
- Both button clicks and typed commands must work identically
- No changes to existing command handler logic
- Maintain existing database operations and API calls

**Code Organization**
- Add navigation callback handlers to callbacks.go
- Update bot.go NewBot function to register new callback handlers
- Update bot.go NewBot function to call SetMyCommands
- Modify handleStart in handlers.go to include inline keyboard
- Modify handleMuteCallback in callbacks.go to add undo button

**Error Handling**
- Answer all callback queries even on error
- Use existing error logging patterns with slog
- Graceful fallback if keyboard creation fails
- Maintain existing error message patterns

**Testing Considerations**
- Navigation buttons must produce identical results to commands
- Undo button must work immediately after mute action
- Menu button should appear in Telegram UI
- All Persian text renders correctly

**Performance Requirements**
- No additional database queries beyond existing handlers
- Callback handlers should respond within 200ms
- Button interactions should feel instant to users

## Existing Code to Leverage

**Inline Keyboard Creation Pattern from handleMutedList**
- Uses botModels.InlineKeyboardMarkup with InlineKeyboard array
- Creates button rows with Text and CallbackData fields
- Sends via SendMessageWithKeyboard helper method
- Pattern found at internal/telegram/handlers.go lines 188-214

**Callback Handler Structure from callbacks.go**
- Extracts callback data using strings.SplitN for prefix parsing
- Answers callback query with AnswerCallbackQueryParams
- Updates message using EditMessageReplyMarkup or EditMessageText
- Error handling with slog.Error logging and user-friendly Persian messages

**Mute Button Implementation from notifications.go**
- createMuteButton function creates single-button keyboard
- shouldShowMuteButton validates content type and series name
- Callback data format: "mute:{seriesName}"
- Integration with SendPhotoBytesWithKeyboard and SendMessageWithKeyboard

**Command Handler Registration from bot.go**
- Uses bot.WithMessageTextHandler for command routing
- MatchTypeExact for exact command matches
- MatchTypePrefix for commands with parameters
- bot.WithCallbackQueryDataHandler with MatchTypePrefix for callbacks

**Existing Helper Methods from bot.go**
- SendMessage for plain text messages
- SendMessageWithKeyboard for messages with inline keyboards
- SendPhotoBytesWithKeyboard for images with keyboards
- All methods handle context and chatID consistently

## Out of Scope
- Custom user preferences UI for personalization settings
- Notification settings UI for controlling when/how notifications are sent
- Admin controls or dashboard features
- Action buttons on individual content items beyond mute/unmute
- Interactive search with state tracking or multi-step search wizards
- Pagination for search results or recent content
- Voice command support
- Multi-language support (only Persian required)
- Analytics or usage tracking dashboard
- Scheduled notifications or reminder features
