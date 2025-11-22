# Manual Testing Checklist: Series Muting/Blocking System

## Test Date: [To be filled in during manual testing]
## Tester: [To be filled in]
## Environment: Development

---

## Prerequisites
- [ ] Jellyfin server is running and accessible
- [ ] Telegram bot is configured and running
- [ ] At least one test user is subscribed to the bot
- [ ] Test content (TV series with multiple episodes) is available in Jellyfin

---

## Test Workflow 1: Basic Mute Functionality

### 1.1 Subscribe to Bot
- [ ] Send `/start` command to bot
- [ ] Verify welcome message is received
- [ ] Verify welcome message includes `/mutedlist` in the commands list

### 1.2 Receive Episode Notification
- [ ] Trigger episode webhook from Jellyfin (add new episode)
- [ ] Verify notification is received in Telegram
- [ ] Verify notification includes:
  - [ ] Episode information (series name, season, episode number)
  - [ ] "دنبال نکردن" (Unfollow/Mute) button

### 1.3 Mute Series
- [ ] Click the "دنبال نکردن" button
- [ ] Verify callback response appears immediately
- [ ] Verify confirmation message is received: "✓ شما دیگر اعلان‌های [Series Name] را دریافت نخواهید کرد"
- [ ] Verify button text changes to "✓ مسدود شده" (Muted)
- [ ] Verify button becomes non-clickable

### 1.4 Verify Mute is Active
- [ ] Trigger another episode notification for the same series
- [ ] Verify muted user does NOT receive the notification
- [ ] If possible, verify non-muted users still receive the notification

---

## Test Workflow 2: View and Manage Muted List

### 2.1 View Muted List
- [ ] Send `/mutedlist` command to bot
- [ ] Verify response shows muted series list with header: "سریال‌های مسدود شده:"
- [ ] Verify each muted series has:
  - [ ] Series name displayed correctly
  - [ ] "رفع مسدودیت: [Series Name]" unmute button

### 2.2 Empty List Case
- [ ] (If no series are muted) Send `/mutedlist`
- [ ] Verify message: "شما هیچ سریالی را مسدود نکرده‌اید"

---

## Test Workflow 3: Unmute Functionality

### 3.1 Unmute Series
- [ ] From `/mutedlist`, click "رفع مسدودیت" button for a muted series
- [ ] Verify callback response appears
- [ ] Verify confirmation message: "✓ [Series Name] از لیست مسدودی‌ها حذف شد"
- [ ] Verify the `/mutedlist` message refreshes automatically
- [ ] Verify unmuted series is removed from the list

### 3.2 Verify Unmute Restores Notifications
- [ ] Trigger episode notification for the unmuted series
- [ ] Verify user now receives the notification
- [ ] Verify notification includes the "دنبال نکردن" button again

---

## Test Workflow 4: Edge Cases

### 4.1 Movie Notifications (Not Affected by Muting)
- [ ] Mute a TV series
- [ ] Trigger a movie notification
- [ ] Verify movie notification is still received
- [ ] Verify movie notification does NOT have mute button

### 4.2 Multiple Series Management
- [ ] Mute 3-5 different series
- [ ] Send `/mutedlist` command
- [ ] Verify all muted series appear in the list
- [ ] Unmute one series from the middle of the list
- [ ] Verify other series remain muted
- [ ] Verify unmuted series no longer appears in list

### 4.3 Special Characters in Series Names
- [ ] Test with series containing:
  - [ ] Persian/Arabic characters
  - [ ] Special characters (colons, ampersands, dashes)
  - [ ] Long series names
- [ ] Verify all operations (mute/unmute/list) work correctly

### 4.4 Duplicate Mute Attempts
- [ ] Mute a series
- [ ] Receive another episode notification for the same series
- [ ] Try clicking the mute button again (if it appears)
- [ ] Verify graceful handling (no error, no duplicate)

### 4.5 Invalid Series Names
- [ ] (If possible) Trigger episode with empty series name
- [ ] Verify notification is received
- [ ] Verify no mute button appears

---

## Test Workflow 5: Help and Documentation

### 5.1 Help Command
- [ ] Send `/help` command (or any unrecognized command for default handler)
- [ ] Verify help message includes `/mutedlist` command
- [ ] Verify description: "مشاهده سریال‌های مسدود شده"

### 5.2 Start Command
- [ ] Send `/start` command
- [ ] Verify welcome message includes `/mutedlist` in commands list

---

## Test Workflow 6: Multi-User Scenarios

### 6.1 Independent Muting (Requires 2+ Test Users)
- [ ] User 1: Mute series "Breaking Bad"
- [ ] User 2: Don't mute any series
- [ ] Trigger "Breaking Bad" episode notification
- [ ] Verify User 1 does not receive notification
- [ ] Verify User 2 receives notification

### 6.2 Same Series, Different Users
- [ ] User 1: Mute series "Breaking Bad"
- [ ] User 2: Mute series "Breaking Bad"
- [ ] User 1: Unmute "Breaking Bad"
- [ ] Trigger "Breaking Bad" episode notification
- [ ] Verify User 1 receives notification (unmuted)
- [ ] Verify User 2 does not receive notification (still muted)

---

## Persian/RTL Text Verification

### UI Text Correctness
- [ ] All button text displays correctly in Persian
- [ ] All confirmation messages display correctly in Persian
- [ ] `/mutedlist` output formats correctly (right-to-left)
- [ ] No text truncation or display issues
- [ ] Checkmark emoji (✓) appears correctly

---

## Performance and Reliability

### Database Persistence
- [ ] Mute a series
- [ ] Restart the bot
- [ ] Send `/mutedlist`
- [ ] Verify muted series persists after restart

### Response Time
- [ ] Mute button click responds within 1-2 seconds
- [ ] Unmute button click responds within 1-2 seconds
- [ ] `/mutedlist` command responds within 1-2 seconds

---

## Issues Found

| Issue # | Description | Severity | Steps to Reproduce |
|---------|-------------|----------|-------------------|
| 1       |             |          |                   |
| 2       |             |          |                   |
| 3       |             |          |                   |

**Severity Levels:**
- **Critical**: Feature is broken/unusable
- **High**: Major functionality issue
- **Medium**: Minor functionality issue
- **Low**: Cosmetic or minor inconvenience

---

## Test Summary

**Total Tests Executed:** _____ / 49

**Passed:** _____

**Failed:** _____

**Blocked:** _____

**Notes:**
[Add any additional observations, comments, or recommendations]

---

## Sign-Off

**Tester Signature:** ___________________

**Date:** ___________________

**Approved for Production:** [ ] Yes [ ] No

**Additional Comments:**
