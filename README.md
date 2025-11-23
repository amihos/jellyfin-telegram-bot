# Jellyfin Telegram Bot

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.22%2B-blue)](https://go.dev/)
[![Test](https://github.com/yourusername/jellyfin-telegram-bot/workflows/Test/badge.svg)](https://github.com/yourusername/jellyfin-telegram-bot/actions/workflows/test.yml)
[![Build](https://github.com/yourusername/jellyfin-telegram-bot/workflows/Build/badge.svg)](https://github.com/yourusername/jellyfin-telegram-bot/actions/workflows/build.yml)
[![Docker](https://github.com/yourusername/jellyfin-telegram-bot/workflows/Docker/badge.svg)](https://github.com/yourusername/jellyfin-telegram-bot/actions/workflows/docker.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/jellyfin-telegram-bot)](https://goreportcard.com/report/github.com/yourusername/jellyfin-telegram-bot)

A Telegram bot that sends instant notifications when new movies, TV shows, or episodes are added to your Jellyfin media server. Get notified in your preferred language (English or Persian) with beautiful media posters and detailed information.

## Features

- **Real-time Notifications**: Get instant alerts when new content is added to your Jellyfin server
- **Multi-language Support**: Full interface in English and Persian (Farsi), with automatic language detection
- **Beautiful Media Cards**: Notifications include poster images, ratings, genres, and descriptions
- **Browse Recent Content**: View recently added media with the `/recent` command
- **Search Your Library**: Find movies and TV shows instantly with `/search`
- **Smart Mute Controls**: Mute notifications for specific TV series while continuing to receive others
- **Interactive UI**: Inline keyboard navigation for browsing content
- **Simple Subscription**: Just send `/start` to subscribe to notifications
- **Lightweight & Fast**: Single binary deployment with minimal resource usage (< 50MB RAM)
- **Docker Support**: Easy deployment with Docker or docker-compose

## Quick Start

Get your bot running in under 10 minutes!

### Prerequisites

- A Telegram account
- A running Jellyfin server (with admin access)
- One of the following:
  - Docker installed (easiest option)
  - Go 1.22+ installed (to build from source)
  - Or download a pre-built binary for your platform

### Step 1: Create Your Telegram Bot

1. Open Telegram and message [@BotFather](https://t.me/BotFather)
2. Send the command `/newbot`
3. Follow the prompts to choose a name and username for your bot
4. Copy the bot token (looks like `123456789:ABCdefGHIjklMNOpqrsTUVwxyz`)
5. Save this token - you'll need it in Step 3

### Step 2: Get Your Jellyfin API Key

1. Log in to your Jellyfin server web interface
2. Click the menu icon (☰) → **Dashboard**
3. In the Advanced section, click **API Keys**
4. Click the **+** button to create a new key
5. Give it a name (e.g., "Telegram Bot")
6. Copy the API key that appears
7. Save this key - you'll need it in Step 3

### Step 3: Choose Your Installation Method

#### Option A: Docker (Recommended)

```bash
# Create a directory for the bot
mkdir jellyfin-telegram-bot
cd jellyfin-telegram-bot

# Download the example docker-compose file
curl -O https://raw.githubusercontent.com/yourusername/jellyfin-telegram-bot/main/docker-compose.example.yml

# Rename it to docker-compose.yml
mv docker-compose.example.yml docker-compose.yml

# Create .env file with your credentials
cat > .env << EOF
TELEGRAM_BOT_TOKEN=your_bot_token_from_step1
JELLYFIN_SERVER_URL=http://your-jellyfin-server:8096
JELLYFIN_API_KEY=your_api_key_from_step2
PORT=8080
EOF

# Start the bot
docker-compose up -d

# View logs to confirm it's running
docker-compose logs -f
```

#### Option B: Pre-built Binary

Download the latest release for your platform:

**Linux (amd64):**
```bash
wget https://github.com/yourusername/jellyfin-telegram-bot/releases/latest/download/jellyfin-telegram-bot-linux-amd64
chmod +x jellyfin-telegram-bot-linux-amd64
mv jellyfin-telegram-bot-linux-amd64 jellyfin-telegram-bot
```

**Linux (arm64 - Raspberry Pi, etc.):**
```bash
wget https://github.com/yourusername/jellyfin-telegram-bot/releases/latest/download/jellyfin-telegram-bot-linux-arm64
chmod +x jellyfin-telegram-bot-linux-arm64
mv jellyfin-telegram-bot-linux-arm64 jellyfin-telegram-bot
```

**Windows (amd64):**
Download from: https://github.com/yourusername/jellyfin-telegram-bot/releases/latest/download/jellyfin-telegram-bot-windows-amd64.exe

**macOS (Intel):**
```bash
wget https://github.com/yourusername/jellyfin-telegram-bot/releases/latest/download/jellyfin-telegram-bot-darwin-amd64
chmod +x jellyfin-telegram-bot-darwin-amd64
mv jellyfin-telegram-bot-darwin-amd64 jellyfin-telegram-bot
```

**macOS (Apple Silicon):**
```bash
wget https://github.com/yourusername/jellyfin-telegram-bot/releases/latest/download/jellyfin-telegram-bot-darwin-arm64
chmod +x jellyfin-telegram-bot-darwin-arm64
mv jellyfin-telegram-bot-darwin-arm64 jellyfin-telegram-bot
```

Then create a `.env` file:
```bash
cat > .env << EOF
TELEGRAM_BOT_TOKEN=your_bot_token_from_step1
JELLYFIN_SERVER_URL=http://your-jellyfin-server:8096
JELLYFIN_API_KEY=your_api_key_from_step2
PORT=8080
EOF
```

And run the bot:
```bash
# Linux/macOS
./jellyfin-telegram-bot

# Windows (in Command Prompt or PowerShell)
jellyfin-telegram-bot.exe
```

#### Option C: Build from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/jellyfin-telegram-bot.git
cd jellyfin-telegram-bot

# Download dependencies
go mod download

# Build the binary
go build -o jellyfin-telegram-bot cmd/bot/main.go

# Create .env file
cp .env.example .env

# Edit .env with your credentials
nano .env  # or use your preferred editor

# Run the bot
./jellyfin-telegram-bot
```

### Step 4: Configure Jellyfin Webhook

1. In your Jellyfin dashboard, go to **Plugins** → **Catalog**
2. Search for "Webhook" and install it
3. Restart Jellyfin server
4. Go to **Dashboard** → **Plugins** → **Webhook**
5. Click **Add Generic Destination**
6. Configure the webhook:
   - **Webhook Name**: Telegram Bot
   - **Webhook URL**: `http://your-bot-server-ip:8080/webhook`
     - If bot runs on same machine as Jellyfin: `http://localhost:8080/webhook`
     - If bot runs on different machine: `http://192.168.1.x:8080/webhook`
   - **Notification Type**: Check **Item Added**
   - **Item Type**: Check **Movies** and **Episodes**
   - **Send All Properties**: Enable (recommended)
7. Click **Save**

### Step 5: Subscribe and Test

1. Open Telegram and find your bot (search for the username you created)
2. Send `/start` to subscribe to notifications
3. The bot will respond in your language (based on your Telegram settings)
4. Add a new movie or episode to Jellyfin to test
5. You should receive a notification within seconds!

**Congratulations!** Your bot is now running.

## Installation

### Docker Installation

The easiest way to run the bot is with Docker. See [docker-compose.example.yml](docker-compose.example.yml) for a complete example.

**Run with Docker:**
```bash
docker run -d \
  --name jellyfin-telegram-bot \
  -e TELEGRAM_BOT_TOKEN=your_token \
  -e JELLYFIN_SERVER_URL=http://your-jellyfin:8096 \
  -e JELLYFIN_API_KEY=your_api_key \
  -e PORT=8080 \
  -v ./data:/app/data \
  -v ./logs:/app/logs \
  -p 8080:8080 \
  --restart unless-stopped \
  ghcr.io/yourusername/jellyfin-telegram-bot:latest
```

**Run with docker-compose:**
```bash
# Download docker-compose.example.yml and rename to docker-compose.yml
# Edit it with your credentials
docker-compose up -d
```

**View logs:**
```bash
docker-compose logs -f
```

**Update to latest version:**
```bash
docker-compose pull
docker-compose up -d
```

For more details, see [docs/deployment.md](docs/deployment.md).

### Binary Installation

Pre-built binaries are available for:
- Linux (amd64, arm64)
- Windows (amd64)
- macOS (amd64, arm64 / Apple Silicon)

Download from the [Releases](https://github.com/yourusername/jellyfin-telegram-bot/releases) page.

#### Running as a System Service (Linux)

For production deployments, you can run the bot as a systemd service. See [docs/linux-service.md](docs/linux-service.md) for a complete guide.

**Quick setup:**
```bash
# Download the systemd service file
sudo curl -o /etc/systemd/system/jellyfin-telegram-bot.service \
  https://raw.githubusercontent.com/yourusername/jellyfin-telegram-bot/main/docs/jellyfin-telegram-bot.service

# Edit the service file to set your paths
sudo nano /etc/systemd/system/jellyfin-telegram-bot.service

# Reload systemd
sudo systemctl daemon-reload

# Enable and start the service
sudo systemctl enable jellyfin-telegram-bot
sudo systemctl start jellyfin-telegram-bot

# Check status
sudo systemctl status jellyfin-telegram-bot

# View logs
sudo journalctl -u jellyfin-telegram-bot -f
```

#### Running on Windows

See [docs/windows-service.md](docs/windows-service.md) for instructions on running as a Windows service.

### Building from Source

**Requirements:**
- Go 1.22 or higher
- Git

**Build steps:**
```bash
# Clone the repository
git clone https://github.com/yourusername/jellyfin-telegram-bot.git
cd jellyfin-telegram-bot

# Download dependencies
go mod download

# Build
go build -o jellyfin-telegram-bot cmd/bot/main.go

# Or build for a specific platform
GOOS=linux GOARCH=amd64 go build -o jellyfin-telegram-bot-linux-amd64 cmd/bot/main.go
```

## Configuration

The bot is configured using environment variables. You can set them in a `.env` file or pass them directly.

### Required Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `TELEGRAM_BOT_TOKEN` | Your bot token from @BotFather | `123456789:ABCdef...` |
| `JELLYFIN_SERVER_URL` | URL of your Jellyfin server | `http://localhost:8096` |
| `JELLYFIN_API_KEY` | Jellyfin API key for authentication | `a1b2c3d4e5f6...` |

### Optional Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Port for webhook server | `8080` |
| `WEBHOOK_SECRET` | Secret for webhook validation | (none) |
| `DATABASE_PATH` | Path to SQLite database | `./bot.db` |
| `LOG_LEVEL` | Log verbosity (DEBUG, INFO, WARN, ERROR) | `INFO` |
| `LOG_FILE` | Path to log file | `./logs/bot.log` |

For a complete reference of all configuration options, see [docs/configuration.md](docs/configuration.md).

### Example .env File

```env
# Required
TELEGRAM_BOT_TOKEN=123456789:ABCdefGHIjklMNOpqrsTUVwxyz-123456789
JELLYFIN_SERVER_URL=http://192.168.1.100:8096
JELLYFIN_API_KEY=a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6

# Optional
PORT=8080
WEBHOOK_SECRET=my-secure-secret-token
DATABASE_PATH=./data/bot.db
LOG_LEVEL=INFO
LOG_FILE=./logs/bot.log
```

## Usage

### Available Commands

Send these commands to your bot in Telegram:

- `/start` - Subscribe to notifications
- `/language` - Change bot language (English/Persian)
- `/recent` - View recently added content (last 15 items)
- `/search <query>` - Search for movies or TV shows
- `/help` - Show help message with all available commands

### Notification Features

When new content is added to Jellyfin, subscribers receive a message with:
- **Poster image** (if available)
- **Title** and year
- **Type** (Movie, Episode, Series)
- **Rating** (e.g., ⭐ 8.5/10)
- **Genres** (e.g., Action, Drama, Thriller)
- **Description** (plot summary)
- **Interactive buttons** to mute notifications for specific series

### Browsing Content

Use `/recent` to see the latest additions:
- Navigate with ◀️ Previous / Next ▶️ buttons
- View 3 items per page
- See poster images and full details
- Mute notifications for series you're not interested in

## Supported Languages

The bot currently supports:
- **English** (en)
- **Persian/Farsi** (fa) - فارسی

The bot automatically detects your language from Telegram settings. You can change it anytime with the `/language` command.

### Adding More Languages

We welcome translations! To add a new language:

1. Copy `locales/active.en.toml` to `locales/active.{language_code}.toml`
2. Translate all message strings
3. Test your translations
4. Submit a pull request

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed instructions.

## Architecture

The bot is built with Go and uses:
- **Telegram Bot API** via `github.com/go-telegram/bot`
- **SQLite database** with GORM for data persistence
- **Jellyfin API** for fetching media information
- **go-i18n** for internationalization
- **Structured logging** with Go's built-in slog

Project structure:
```
jellyfin-telegram-bot/
├── cmd/bot/              # Application entry point
├── internal/
│   ├── config/           # Configuration and logging
│   ├── database/         # Database layer (GORM models)
│   ├── handlers/         # HTTP webhook handlers
│   ├── telegram/         # Telegram bot logic
│   ├── jellyfin/         # Jellyfin API client
│   └── i18n/             # Internationalization
├── locales/              # Translation files
├── docs/                 # Documentation
└── test/                 # Integration tests
```

For detailed architecture documentation, see [docs/architecture.md](docs/architecture.md).

## Troubleshooting

### Bot Not Responding

**Check if the bot is running:**
```bash
# Docker
docker-compose ps

# Systemd
sudo systemctl status jellyfin-telegram-bot

# Manual
ps aux | grep jellyfin-telegram-bot
```

**Check the logs:**
```bash
# Docker
docker-compose logs -f

# Systemd
sudo journalctl -u jellyfin-telegram-bot -f

# Manual (if LOG_FILE is set)
tail -f ./logs/bot.log
```

**Common issues:**
- Invalid bot token: Check `TELEGRAM_BOT_TOKEN` in .env
- Bot not started: Start the bot service
- Network issues: Ensure bot can reach Telegram API

### Webhook Not Working

**Test the webhook endpoint:**
```bash
curl http://localhost:8080/health
# Should return: {"status":"healthy"}
```

**Check Jellyfin webhook configuration:**
- Webhook URL must be accessible from Jellyfin server
- If bot runs on different machine, use IP address not localhost
- Port must match `PORT` in .env
- Firewall must allow incoming connections on that port

**Test sending a webhook:**
```bash
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -d '{"NotificationType":"ItemAdded","ItemType":"Movie","Name":"Test Movie"}'
```

### Not Receiving Notifications

**Check subscription:**
- Send `/start` to the bot to ensure you're subscribed
- Check if you've muted the specific series (if it's a TV show)

**Check Jellyfin webhook:**
- Go to Jellyfin Dashboard → Plugins → Webhook
- Verify webhook is enabled
- Check "Item Added" notification type is selected
- Check "Movies" and "Episodes" item types are selected

**Enable debug logging:**
Edit .env and set:
```env
LOG_LEVEL=DEBUG
```
Restart the bot and check logs for detailed information.

### Language Not Changing

**Check translation files exist:**
```bash
ls locales/
# Should show: active.en.toml, active.fa.toml
```

**Reset language preference:**
Send `/language` and select your preferred language again.

**Check database:**
```bash
# If using SQLite
sqlite3 bot.db "SELECT chat_id, language_code FROM subscribers;"
```

For more troubleshooting help, see [docs/troubleshooting.md](docs/troubleshooting.md) or [open an issue](https://github.com/yourusername/jellyfin-telegram-bot/issues).

## Contributing

We welcome contributions! Whether you're fixing bugs, adding features, improving documentation, or translating to new languages, your help is appreciated.

### How to Contribute

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`go test ./...`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Development Setup

```bash
# Clone your fork
git clone https://github.com/yourusername/jellyfin-telegram-bot.git
cd jellyfin-telegram-bot

# Install dependencies
go mod download

# Create .env for testing
cp .env.example .env
# Edit .env with test credentials

# Run tests
go test ./...

# Run the bot
go run cmd/bot/main.go
```

### Code Quality

Before submitting a PR:
- Run `go fmt ./...` to format code
- Run `go test ./...` to ensure tests pass
- Run `golangci-lint run` to check for issues (if you have it installed)
- Update documentation if you've changed functionality
- Add tests for new features

For detailed guidelines, see [CONTRIBUTING.md](CONTRIBUTING.md).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Copyright (c) 2025 Hossein Amirkhalili

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

## Acknowledgments

- Built with [go-telegram/bot](https://github.com/go-telegram/bot)
- Internationalization with [go-i18n](https://github.com/nicksnyder/go-i18n)
- Inspired by the [Jellyfin](https://jellyfin.org) community

## Support

- **Issues**: [GitHub Issues](https://github.com/yourusername/jellyfin-telegram-bot/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/jellyfin-telegram-bot/discussions)
- **Jellyfin Community**: [jellyfin.org](https://jellyfin.org)

---

Made with ❤️ for the Jellyfin community
