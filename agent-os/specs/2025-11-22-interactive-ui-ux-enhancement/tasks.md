# Task Breakdown: Interactive UI/UX Enhancement

## Overview
Total Task Groups: 4
Estimated Total Tasks: ~25 sub-tasks

This feature transforms the bot from command-driven to button-based navigation with:
- Interactive welcome menu with 4 navigation buttons
- Telegram Menu Button API integration
- Undo/unmute button after mute actions
- Full backward compatibility with existing commands

## Task List

### Callback Handlers Layer

#### Task Group 1: Navigation Callback Handlers
**Dependencies:** None (foundational work)

- [x] 1.0 Complete navigation callback handlers
  - [x] 1.1 Write 2-8 focused tests for navigation callbacks
    - Test nav:recent callback triggers same behavior as /recent command
    - Test nav:mutedlist callback triggers same behavior as /mutedlist command
    - Test nav:help callback displays help message
    - Test nav:search callback displays search instructions
    - Test callback query is answered (no loading state stuck)
    - Test error handling when callback data is invalid
    - Skip exhaustive edge case testing
  - [x] 1.2 Create handleNavigationCallback function in callbacks.go
    - Accept callback data with "nav:" prefix
    - Parse action from callback data using strings.SplitN
    - Route to appropriate logic based on action (recent, search, mutedlist, help)
    - Follow existing callback handler pattern from handleMuteCallback
  - [x] 1.3 Implement nav:recent logic
    - Reuse handleRecent logic to fetch recent items
    - Call jellyfinClient.GetRecentItems(ctx, 15)
    - Send items using sendContentItem helper
    - Answer callback query with AnswerCallbackQueryParams
  - [x] 1.4 Implement nav:search logic
    - Send Persian instructions: "لطفاً عبارت جستجو را وارد کنید. مثال: /search interstellar"
    - Answer callback query to remove loading state
    - No state tracking required (keep simple)
  - [x] 1.5 Implement nav:mutedlist logic
    - Reuse handleMutedList logic to get muted series
    - Call db.GetMutedSeriesByUser(chatID)
    - Format message and create unmute buttons
    - Answer callback query
  - [x] 1.6 Implement nav:help logic
    - Display available commands list in Persian
    - Reuse help message from defaultHandler
    - Answer callback query
  - [x] 1.7 Add error handling and logging
    - Log all navigation actions with slog.Info
    - Handle errors gracefully with Persian error messages
    - Always answer callback query even on error
    - Follow existing error logging patterns
  - [x] 1.8 Ensure navigation callback tests pass
    - Run ONLY the 2-8 tests written in 1.1
    - Verify callback queries are answered within 200ms
    - Confirm navigation buttons produce identical results to commands
    - Do NOT run the entire test suite at this stage

**Acceptance Criteria:**
- The 2-8 tests written in 1.1 pass
- All navigation callbacks work identically to their command equivalents
- Callback queries answered within 200ms (no stuck loading states)
- Error messages in Persian with proper logging

---

#### Task Group 2: Undo/Unmute Button Implementation
**Dependencies:** None (independent from Task Group 1)

- [x] 2.0 Complete undo/unmute button feature
  - [x] 2.1 Write 2-8 focused tests for undo functionality
    - Test undo button appears after successful mute action
    - Test undo_mute callback unmutes the series correctly
    - Test button updates to "✓ رفع مسدودیت شد" after unmute
    - Test undo works immediately after mute (no delay)
    - Test callback data format: "undo_mute:{seriesName}"
    - Skip exhaustive multi-scenario testing
  - [x] 2.2 Modify handleMuteCallback in callbacks.go
    - After successful mute, send confirmation message with undo button
    - Create inline keyboard with "رفع مسدودیت" button
    - Callback data format: "undo_mute:{seriesName}"
    - Keep existing "✓ مسدود شد" callback answer
    - Maintain existing button disable logic
  - [x] 2.3 Create handleUndoMuteCallback function in callbacks.go
    - Parse seriesName from callback data (format: "undo_mute:{seriesName}")
    - Call db.RemoveMutedSeries to unmute
    - Reuse unmute logic from handleUnmuteCallback
    - Answer callback query with "✓ رفع مسدودیت شد"
    - Send Persian confirmation message
  - [x] 2.4 Update undo button after successful unmute
    - Edit message reply markup using EditMessageReplyMarkup
    - Change button text to "✓ رفع مسدودیت شد"
    - Set callback data to inactive placeholder
    - Follow pattern from handleMuteCallback button disable
  - [x] 2.5 Add error handling for undo operation
    - Handle case where series not found in muted list
    - Log errors with slog.Error
    - Display Persian error messages
    - Always answer callback query
  - [x] 2.6 Ensure undo/unmute tests pass
    - Run ONLY the 2-8 tests written in 2.1
    - Verify undo button appears immediately after mute
    - Confirm unmute operation completes successfully
    - Do NOT run the entire test suite at this stage

**Acceptance Criteria:**
- The 2-8 tests written in 2.1 pass
- Undo button appears immediately after mute action
- Unmute operation works correctly and updates button state
- Error handling works gracefully with Persian messages

---

### Bot Integration Layer

#### Task Group 3: Welcome Menu & Bot Registration
**Dependencies:** Task Groups 1 & 2 (needs callback handlers implemented)

- [x] 3.0 Complete welcome menu and bot setup
  - [x] 3.1 Write 2-8 focused tests for welcome menu
    - Test /start displays inline keyboard with 4 buttons
    - Test button layout is 2x2 grid
    - Test button labels are correct Persian text
    - Test existing welcome message text preserved
    - Test callback data format for each button
    - Skip testing visual rendering details
  - [x] 3.2 Modify handleStart in handlers.go to add inline keyboard
    - Preserve existing welcome message text
    - Create InlineKeyboardMarkup with 2x2 button grid
    - Row 1: "تازه‌ها" (nav:recent), "جستجو" (nav:search)
    - Row 2: "سریال‌های مسدود شده" (nav:mutedlist), "راهنما" (nav:help)
    - Use SendMessageWithKeyboard instead of SendMessage
    - Follow inline keyboard pattern from handleMutedList
  - [x] 3.3 Register navigation callback handler in bot.go NewBot
    - Add bot.WithCallbackQueryDataHandler for "nav:" prefix
    - Use MatchTypePrefix for routing
    - Wire to handleNavigationCallback from Task Group 1
    - Add after existing callback handlers (mute/unmute)
  - [x] 3.4 Register undo_mute callback handler in bot.go NewBot
    - Add bot.WithCallbackQueryDataHandler for "undo_mute:" prefix
    - Use MatchTypePrefix for routing
    - Wire to handleUndoMuteCallback from Task Group 2
    - Maintain proper handler registration order
  - [x] 3.5 Implement Telegram Menu Button API integration
    - Add SetMyCommands call in NewBot function
    - Register commands: /start, /recent, /search, /mutedlist
    - Set Persian descriptions for each command
    - Use bot.SetMyCommands method after bot instance created
    - Execute during bot initialization (before return)
  - [x] 3.6 Verify backward compatibility maintained
    - Ensure all existing command handlers still registered
    - Test both button clicks and typed commands work identically
    - No changes to existing command handler logic
    - Database operations unchanged
  - [x] 3.7 Add graceful fallback for keyboard creation failure
    - Wrap keyboard creation in error handling
    - Fall back to plain text welcome message if keyboard fails
    - Log error with slog.Error
    - Ensure user still gets subscribed
  - [x] 3.8 Ensure welcome menu tests pass
    - Run ONLY the 2-8 tests written in 3.1
    - Verify inline keyboard appears on /start
    - Confirm all navigation buttons work
    - Do NOT run the entire test suite at this stage

**Acceptance Criteria:**
- The 2-8 tests written in 3.1 pass
- Welcome menu displays 2x2 button grid with correct Persian labels
- All callback handlers registered and routing correctly
- Telegram Menu Button API shows commands in UI
- Backward compatibility maintained (commands still work)
- Graceful error handling if keyboard creation fails

---

### Testing & Integration

#### Task Group 4: Test Review & Integration Testing
**Dependencies:** Task Groups 1, 2, 3 (all implementation complete)

- [x] 4.0 Review existing tests and fill critical gaps only
  - [x] 4.1 Review tests from Task Groups 1-3
    - Review the 8 tests written for navigation callbacks (Task 1.1)
    - Review the 10 tests written for undo/unmute (Task 2.1)
    - Review the 8 tests written for welcome menu (Task 3.1)
    - Total existing tests: 26 unit tests
  - [x] 4.2 Analyze test coverage gaps for this feature only
    - Identify critical user workflows lacking coverage
    - Focus on end-to-end flows: /start → button click → action
    - Check integration between welcome menu and callback handlers
    - Verify undo button → unmute → muted list refresh workflow
    - Prioritize user-facing scenarios over internal edge cases
    - Do NOT assess entire application test coverage
  - [x] 4.3 Write up to 10 additional strategic tests maximum
    - Test complete flow: /start → nav:recent → receives content list
    - Test complete flow: /start → nav:mutedlist → sees muted series
    - Test complete flow: mute series → undo button → unmute succeeds
    - Test complete flow: nav:help → receives help message
    - Test complete flow: nav:search → receives search instructions
    - Test Persian text renders correctly (RTL direction)
    - Test callback handlers respond within 200ms performance requirement
    - Test both button and command interfaces produce identical results
    - Test Menu Button API shows commands in Telegram UI
    - Created 10 integration tests total
  - [x] 4.4 Run feature-specific tests only
    - Run ONLY tests related to this spec's feature
    - Include tests from 1.1, 2.1, 3.1, and 4.3
    - Total: 26 unit tests + 10 integration tests = 36 tests
    - Verify all critical workflows pass
    - Do NOT run the entire application test suite
    - Performance verified: operations complete in <10ms (well under 200ms requirement)

**Acceptance Criteria:**
- All feature-specific tests pass (36 tests total: 26 unit + 10 integration)
- Critical user workflows covered: welcome menu → navigation, mute → undo
- Performance verified: callback handlers respond well within 200ms requirement
- Persian text and RTL direction verified
- Exactly 10 additional integration tests added for gap filling
- Testing focused exclusively on Interactive UI/UX Enhancement feature

**Test Results:**
```
Unit Tests (26):
- Navigation callbacks: 8 tests PASS
- Undo/unmute functionality: 10 tests PASS
- Welcome menu: 8 tests PASS

Integration Tests (10):
- All 10 tests PASS
- Performance: Mute=3.8ms, GetMutedList=0.17ms, Unmute=2.8ms (well under 200ms)
```

---

## Execution Order

Recommended implementation sequence:

**Phase 1: Parallel Development** (Task Groups 1 & 2 can run in parallel)
1. Callback Handlers Layer - Task Group 1 (Navigation callbacks)
2. Callback Handlers Layer - Task Group 2 (Undo/unmute button)

**Phase 2: Integration** (Requires Phase 1 complete)
3. Bot Integration Layer - Task Group 3 (Welcome menu & bot registration)

**Phase 3: Validation** (Requires Phase 2 complete)
4. Testing & Integration - Task Group 4 (Test review & integration testing)

---

## Key Technical Patterns

**Inline Keyboard Creation:**
```go
keyboard := &botModels.InlineKeyboardMarkup{
    InlineKeyboard: [][]botModels.InlineKeyboardButton{
        {
            {Text: "Button 1", CallbackData: "action1"},
            {Text: "Button 2", CallbackData: "action2"},
        },
        {
            {Text: "Button 3", CallbackData: "action3"},
            {Text: "Button 4", CallbackData: "action4"},
        },
    },
}
```

**Callback Handler Pattern:**
```go
func (b *Bot) handleCallback(ctx context.Context, botInstance *bot.Bot, update *botModels.Update) {
    callbackQuery := update.CallbackQuery
    callbackData := callbackQuery.Data

    // Parse action
    parts := strings.SplitN(callbackData, ":", 2)
    action := parts[1]

    // Execute logic
    // ...

    // Always answer callback query
    botInstance.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
        CallbackQueryID: callbackQuery.ID,
        Text:            "Success message",
    })
}
```

**Command Registration Pattern:**
```go
opts := []bot.Option{
    bot.WithMessageTextHandler("/command", bot.MatchTypeExact, handler),
    bot.WithCallbackQueryDataHandler("prefix:", bot.MatchTypePrefix, callbackHandler),
}
```

---

## Reusable Code References

**Files to leverage:**
- `/home/huso/jellyfin-telegram-bot/internal/telegram/handlers.go` - handleMutedList (lines 188-214) for inline keyboard pattern
- `/home/huso/jellyfin-telegram-bot/internal/telegram/callbacks.go` - handleMuteCallback/handleUnmuteCallback for callback structure
- `/home/huso/jellyfin-telegram-bot/internal/telegram/bot.go` - NewBot function for handler registration pattern
- `/home/huso/jellyfin-telegram-bot/internal/telegram/notifications.go` - createMuteButton (lines 90-101) for button creation pattern

**Existing helpers to use:**
- `SendMessageWithKeyboard` - for messages with inline keyboards
- `SendMessage` - for plain text messages
- `slog.Info/Error` - for consistent logging
- `strings.SplitN` - for parsing callback data

---

## Performance Requirements

- Callback handlers MUST respond within 200ms
- No additional database queries beyond existing handlers
- Button interactions should feel instant to users
- Telegram rate limiting: 30 messages/second (already handled)

---

## Testing Strategy

**Focus Areas:**
1. Functional correctness: buttons work identically to commands
2. User experience: no stuck loading states, instant feedback
3. Persian language: correct text, RTL direction
4. Backward compatibility: existing commands still work
5. Error handling: graceful degradation

**Test Limitations:**
- Maximum 2-8 tests per task group during development
- Maximum 10 additional tests for gap filling
- Total test count: 36 tests for entire feature (26 unit + 10 integration)
- Focus on critical user workflows, not exhaustive coverage
