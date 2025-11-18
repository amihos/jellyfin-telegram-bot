# Jellyfin Telegram Bot

A Telegram bot that monitors a Jellyfin media server and sends Persian-language notifications to subscribed users when new movies or TV episodes are added.

## Features

- **Automatic Notifications**: Receive notifications when new movies or episodes are added to Jellyfin
- **Persian Interface**: All bot messages and UI in Persian/Farsi
- **Browse Recent Content**: View recently added media with `/recent` command
- **Search Library**: Search for movies and TV shows with `/search` command
- **Simple Subscription**: Just send `/start` to subscribe to notifications

## Requirements

- Go 1.21 or higher
- Jellyfin server with Webhook plugin installed and configured
- Telegram Bot Token (obtained from [@BotFather](https://t.me/BotFather))
- Jellyfin API Key

## Quick Start

### 1. Clone and Build

```bash
git clone <repository-url>
cd jellyfin-telegram-bot
go mod download
go build -o jellyfin-bot cmd/bot/main.go
```

### 2. Configure Environment

Copy the example environment file and fill in your credentials:

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
TELEGRAM_BOT_TOKEN=your_telegram_bot_token_here
JELLYFIN_SERVER_URL=http://your-jellyfin-server:8096
JELLYFIN_API_KEY=your_jellyfin_api_key_here
WEBHOOK_SECRET=optional_webhook_secret
PORT=8080
DATABASE_PATH=./bot.db
LOG_LEVEL=INFO
LOG_FILE=./logs/bot.log
```

### 3. Run the Bot

```bash
./jellyfin-bot
```

Or run directly:

```bash
go run cmd/bot/main.go
```

## Configuration

### Getting a Telegram Bot Token

1. Message [@BotFather](https://t.me/BotFather) on Telegram
2. Send `/newbot` command
3. Follow the instructions to create your bot
4. Copy the token provided by BotFather

### Getting Jellyfin API Key

1. Log in to your Jellyfin server
2. Go to Dashboard → API Keys
3. Click "+" to create a new API key
4. Give it a name (e.g., "Telegram Bot")
5. Copy the generated API key

### Configuring Jellyfin Webhook Plugin

1. Install the Webhook plugin from Jellyfin's plugin catalog
2. Go to Dashboard → Plugins → Webhook
3. Add a new webhook destination:
   - **Webhook Name**: Telegram Bot
   - **Webhook URL**: `http://your-server-ip:8080/webhook`
   - **Notification Type**: Select "Item Added"
   - **Item Type**: Select "Movies" and "Episodes"
   - **User Filter**: (optional) Select specific users
   - **Send All Properties**: Enabled (recommended)
4. Save the configuration

## Available Commands

- `/start` - Subscribe to notifications
- `/recent` - View recently added content (last 15 items)
- `/search <query>` - Search for movies or TV shows

## Project Structure

```
jellyfin-telegram-bot/
├── cmd/bot/              # Application entry point
├── internal/
│   ├── handlers/         # Request handlers (commands, webhooks)
│   ├── database/         # Database layer
│   ├── jellyfin/         # Jellyfin API client
│   ├── telegram/         # Telegram bot logic
│   └── config/           # Configuration management
├── pkg/models/           # Data models
├── docs/                 # Documentation
└── logs/                 # Log files (auto-created)
```

## Documentation

- [Architecture](docs/architecture.md) - System design and tech stack
- [Deployment](docs/deployment.md) - Production deployment guide
- [API Integration](docs/api-integration.md) - Jellyfin API details

## Development

### Building from Source

```bash
go build -o jellyfin-bot cmd/bot/main.go
```

### Running Tests

```bash
go test ./...
```

### Code Style

This project follows standard Go conventions. Format code with:

```bash
go fmt ./...
```

## Troubleshooting

### Bot not receiving webhooks

- Ensure the bot is running and the webhook endpoint is accessible
- Check firewall rules allow incoming connections on the configured port
- Verify the webhook URL in Jellyfin matches your server address
- Check bot logs for webhook-related errors

### Notifications not sending

- Verify users have subscribed with `/start` command
- Check Telegram API token is valid
- Review bot logs for Telegram API errors
- Ensure bot has not been blocked by users

### Database errors

- Check database file permissions
- Ensure database directory is writable
- Review logs for specific database errors

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License

[Add your license here]

## Support

For issues and questions, please open an issue on the project repository.
