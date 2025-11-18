#!/bin/bash

# Simple monitoring script for the Jellyfin Telegram Bot

echo "╔════════════════════════════════════════════════════════════╗"
echo "║     Jellyfin Telegram Bot - Live Monitor                  ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

# Check if bot is running
if pgrep -x "jellyfin-bot" > /dev/null; then
    echo "✅ Bot Status: RUNNING"
    echo "   PID: $(pgrep -x jellyfin-bot)"
else
    echo "❌ Bot Status: NOT RUNNING"
    echo ""
    echo "Start it with: ./jellyfin-bot"
    exit 1
fi

# Check subscribers
if [ -f "bot.db" ]; then
    SUBS=$(sqlite3 bot.db "SELECT COUNT(*) FROM subscribers WHERE is_active = 1;" 2>/dev/null || echo "0")
    echo "👥 Active Subscribers: $SUBS"

    if [ "$SUBS" -gt 0 ]; then
        echo ""
        echo "📋 Subscriber List:"
        sqlite3 bot.db "SELECT chat_id, username, first_name, created_at FROM subscribers WHERE is_active = 1;" 2>/dev/null | \
        while IFS='|' read -r chat_id username first_name created_at; do
            echo "   • $first_name (@$username) - Chat ID: $chat_id"
        done
    fi
else
    echo "⚠️  Database not found"
fi

# Check webhook endpoint
echo ""
echo -n "🌐 Webhook Health: "
if curl -s -f "http://localhost:8080/health" > /dev/null 2>&1; then
    echo "✅ OK"
else
    echo "❌ Not responding"
fi

# Show recent activity
echo ""
echo "📊 Recent Activity (last 10 lines):"
echo "════════════════════════════════════════════════════════════"
tail -10 logs/bot.log 2>/dev/null | while read line; do
    echo "$line" | python3 -c "import sys, json; [print(f\"[{d['level']}] {d['msg']}\") for d in [json.loads(line) for line in sys.stdin]]" 2>/dev/null || echo "$line"
done

echo ""
echo "════════════════════════════════════════════════════════════"
echo ""
echo "Commands:"
echo "  • Watch logs: tail -f logs/bot.log"
echo "  • Test webhook: ./test-webhook.sh"
echo "  • Stop bot: pkill jellyfin-bot"
echo ""
