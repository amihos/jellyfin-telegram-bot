---
name: Bug Report
about: Report a bug or issue with the Jellyfin Telegram Bot
title: '[BUG] '
labels: bug
assignees: ''
---

## Description

A clear and concise description of what the bug is.

## Steps to Reproduce

Steps to reproduce the behavior:

1. Go to '...'
2. Send command '...'
3. Configure '...'
4. See error

## Expected Behavior

A clear and concise description of what you expected to happen.

## Actual Behavior

A clear and concise description of what actually happened.

## Environment

Please complete the following information:

- **OS**: [e.g., Ubuntu 22.04, Windows 11, macOS 14.0]
- **Bot Version**: [e.g., v1.0.0 - run `jellyfin-telegram-bot --version` or check Docker image tag]
- **Go Version** (if building from source): [e.g., 1.23.0]
- **Installation Method**: [Docker, Pre-built binary, Built from source]
- **Jellyfin Version**: [e.g., 10.8.13]
- **Telegram Language**: [e.g., English, Persian]

## Logs

Please provide relevant log output. To enable debug logging, set `LOG_LEVEL=DEBUG` in your `.env` file and restart the bot.

```
Paste your log output here
```

## Configuration

Please share relevant parts of your configuration (remember to remove sensitive information like tokens and API keys):

```env
# Example - DO NOT include actual tokens
TELEGRAM_BOT_TOKEN=<redacted>
JELLYFIN_SERVER_URL=http://localhost:8096
PORT=8080
LOG_LEVEL=INFO
```

## Additional Context

Add any other context about the problem here:

- Screenshots (if applicable)
- Error messages
- Network configuration (if webhook-related)
- Any recent changes to your setup

## Possible Solution

If you have an idea of what might be causing the issue or how to fix it, please share it here.
