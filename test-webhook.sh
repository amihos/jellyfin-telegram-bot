#!/bin/bash

# Test webhook notifications without Jellyfin
# This sends fake webhook payloads to test the bot

WEBHOOK_URL="http://localhost:8080/webhook"
SECRET="optional_webhook_secret_for_security"

echo "=========================================="
echo "  Jellyfin Webhook Notification Tester"
echo "=========================================="
echo ""

# Check if webhook server is running
echo "Checking if webhook server is running..."
if curl -s -f "$WEBHOOK_URL" > /dev/null 2>&1 || curl -s "http://localhost:8080/health" > /dev/null 2>&1; then
    echo "✓ Webhook server is running"
else
    echo "✗ Webhook server is NOT running"
    echo ""
    echo "Please start the bot first:"
    echo "  ./jellyfin-bot"
    exit 1
fi

echo ""
echo "Select notification type to test:"
echo "  1) Movie notification"
echo "  2) TV Episode notification"
echo "  3) Both (movie + episode)"
echo ""
read -p "Enter choice (1-3): " choice

send_movie_notification() {
    echo ""
    echo "Sending movie notification..."
    echo "--------------------------------------"

    RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$WEBHOOK_URL" \
        -H "Content-Type: application/json" \
        -H "X-Webhook-Secret: $SECRET" \
        -d '{
            "NotificationType": "ItemAdded",
            "ItemType": "Movie",
            "ItemName": "Interstellar",
            "ItemId": "test-movie-12345",
            "Year": 2014,
            "Overview": "A team of explorers travel through a wormhole in space in an attempt to ensure humanity'\''s survival.",
            "ServerName": "Test Jellyfin Server"
        }')

    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    BODY=$(echo "$RESPONSE" | sed '$d')

    if [ "$HTTP_CODE" = "200" ]; then
        echo "✓ Movie notification sent successfully (HTTP $HTTP_CODE)"
        echo "Check your Telegram bot for the notification!"
    else
        echo "✗ Failed to send notification (HTTP $HTTP_CODE)"
        echo "Response: $BODY"
    fi
}

send_episode_notification() {
    echo ""
    echo "Sending episode notification..."
    echo "--------------------------------------"

    RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$WEBHOOK_URL" \
        -H "Content-Type: application/json" \
        -H "X-Webhook-Secret: $SECRET" \
        -d '{
            "NotificationType": "ItemAdded",
            "ItemType": "Episode",
            "ItemName": "The One Where It All Begins",
            "ItemId": "test-episode-67890",
            "SeriesName": "Breaking Bad",
            "SeasonNumber": 1,
            "EpisodeNumber": 1,
            "Year": 2008,
            "Overview": "Walter White, a chemistry teacher, discovers that he has cancer and decides to get into the meth-making business.",
            "ServerName": "Test Jellyfin Server"
        }')

    HTTP_CODE=$(echo "$RESPONSE" | tail -n1)
    BODY=$(echo "$RESPONSE" | sed '$d')

    if [ "$HTTP_CODE" = "200" ]; then
        echo "✓ Episode notification sent successfully (HTTP $HTTP_CODE)"
        echo "Check your Telegram bot for the notification!"
    else
        echo "✗ Failed to send notification (HTTP $HTTP_CODE)"
        echo "Response: $BODY"
    fi
}

# Process user choice
case $choice in
    1)
        send_movie_notification
        ;;
    2)
        send_episode_notification
        ;;
    3)
        send_movie_notification
        sleep 1
        send_episode_notification
        ;;
    *)
        echo "Invalid choice"
        exit 1
        ;;
esac

echo ""
echo "=========================================="
echo "Testing complete!"
echo ""
echo "If you didn't receive notifications:"
echo "  1. Make sure you sent /start to the bot"
echo "  2. Check bot logs: tail -f logs/bot.log"
echo "  3. Verify TELEGRAM_BOT_TOKEN is correct"
echo "=========================================="
