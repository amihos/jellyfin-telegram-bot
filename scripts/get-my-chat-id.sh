#!/bin/bash

# Helper script to get your Telegram chat ID
# Just run this script and send any message to your bot

echo "🔍 Getting your Telegram Chat ID..."
echo ""
echo "Instructions:"
echo "1. Make sure your bot is running"
echo "2. Send ANY message to your bot in Telegram"
echo "3. Your chat ID will appear in the bot logs"
echo ""
echo "To find your chat ID, look in the bot logs for lines like:"
echo '  "chatID": 123456789'
echo ""
echo "Alternatively, you can:"
echo "1. Send /start to your bot"
echo "2. Check the database: sqlite3 bot.db 'SELECT chat_id FROM subscribers;'"
echo ""

# If database exists, show current subscribers
if [ -f "bot.db" ]; then
    echo "📊 Current subscribers in database:"
    sqlite3 bot.db "SELECT chat_id FROM subscribers;" 2>/dev/null || echo "  (unable to read database)"
    echo ""
fi

echo "💡 Once you have your chat ID, add it to .env:"
echo "   ENABLE_BETA_FEATURES=true"
echo "   TESTER_CHAT_IDS=YOUR_CHAT_ID_HERE"
echo ""
