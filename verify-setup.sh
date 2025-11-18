#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "╔════════════════════════════════════════════════════════════╗"
echo "║     Jellyfin Telegram Bot - Setup Verification            ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

# Check .env file exists
echo -e "${BLUE}[1] Checking .env file...${NC}"
if [ -f ".env" ]; then
    echo -e "${GREEN}✅ .env file exists${NC}"
else
    echo -e "${RED}❌ .env file missing${NC}"
    exit 1
fi

# Check Telegram Bot Token
echo ""
echo -e "${BLUE}[2] Checking Telegram Bot Token...${NC}"
BOT_TOKEN=$(grep "TELEGRAM_BOT_TOKEN=" .env | cut -d '=' -f2)
if [ -z "$BOT_TOKEN" ]; then
    echo -e "${RED}❌ Bot token not set${NC}"
elif [ "$BOT_TOKEN" = "YOUR_TELEGRAM_BOT_TOKEN_HERE" ]; then
    echo -e "${RED}❌ Bot token is placeholder value${NC}"
else
    echo -e "${GREEN}✅ Bot token is configured${NC}"
    echo -e "   Token: ${BOT_TOKEN:0:10}...${BOT_TOKEN: -5}"
fi

# Check Jellyfin API Key
echo ""
echo -e "${BLUE}[3] Checking Jellyfin API Key...${NC}"
API_KEY=$(grep "JELLYFIN_API_KEY=" .env | cut -d '=' -f2)
if [ -z "$API_KEY" ]; then
    echo -e "${RED}❌ API key not set${NC}"
    echo -e "${YELLOW}   → Get it from: Jellyfin Dashboard → API Keys → + button${NC}"
elif [ "$API_KEY" = "PASTE_YOUR_API_KEY_HERE" ]; then
    echo -e "${RED}❌ API key is placeholder value${NC}"
    echo -e "${YELLOW}   → Get it from: Jellyfin Dashboard → API Keys → + button${NC}"
else
    echo -e "${GREEN}✅ API key is configured${NC}"
    echo -e "   Key: ${API_KEY:0:8}...${API_KEY: -8}"
fi

# Check Jellyfin Server URL
echo ""
echo -e "${BLUE}[4] Checking Jellyfin Server URL...${NC}"
JELLYFIN_URL=$(grep "JELLYFIN_SERVER_URL=" .env | cut -d '=' -f2)
echo -e "${GREEN}✅ Server URL: $JELLYFIN_URL${NC}"

# Test Jellyfin connectivity
echo ""
echo -e "${BLUE}[5] Testing Jellyfin connectivity...${NC}"
if curl -s -f -m 5 "$JELLYFIN_URL" > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Jellyfin server is reachable from WSL${NC}"
else
    echo -e "${RED}❌ Cannot reach Jellyfin server${NC}"
    echo -e "${YELLOW}   → Check Windows Firewall allows port 8096${NC}"
    echo -e "${YELLOW}   → Verify Jellyfin is running on Windows${NC}"
fi

# Check webhook port
echo ""
echo -e "${BLUE}[6] Checking webhook server...${NC}"
WEBHOOK_PORT=$(grep "PORT=" .env | cut -d '=' -f2)
if curl -s -f "http://localhost:$WEBHOOK_PORT/health" > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Webhook server is responding on port $WEBHOOK_PORT${NC}"
else
    echo -e "${RED}❌ Webhook server not responding${NC}"
    echo -e "${YELLOW}   → Check if bot is running: ps aux | grep jellyfin-bot${NC}"
fi

# Get WSL IP for webhook configuration
echo ""
echo -e "${BLUE}[7] Network Configuration...${NC}"
WSL_IP=$(hostname -I | awk '{print $1}')
echo -e "${GREEN}✅ Your WSL IP: $WSL_IP${NC}"
echo -e "${YELLOW}   → Use this in Jellyfin webhook: http://$WSL_IP:$WEBHOOK_PORT/webhook${NC}"

# Check bot status
echo ""
echo -e "${BLUE}[8] Checking bot process...${NC}"
if pgrep -x "jellyfin-bot" > /dev/null; then
    PID=$(pgrep -x "jellyfin-bot")
    echo -e "${GREEN}✅ Bot is running (PID: $PID)${NC}"

    # Check for errors in logs
    if [ -f "logs/bot.log" ]; then
        RECENT_ERRORS=$(tail -20 logs/bot.log | grep -i "error\|unauthorized" | wc -l)
        if [ "$RECENT_ERRORS" -gt 0 ]; then
            echo -e "${YELLOW}   ⚠️  Found $RECENT_ERRORS recent errors in logs${NC}"
            echo -e "${YELLOW}   → Run: tail -20 logs/bot.log | grep -i error${NC}"
        else
            echo -e "${GREEN}   No recent errors detected${NC}"
        fi
    fi
else
    echo -e "${RED}❌ Bot is not running${NC}"
    echo -e "${YELLOW}   → Start it: ./jellyfin-bot${NC}"
fi

# Check database
echo ""
echo -e "${BLUE}[9] Checking database...${NC}"
if [ -f "bot.db" ]; then
    SUBSCRIBER_COUNT=$(sqlite3 bot.db "SELECT COUNT(*) FROM subscribers WHERE is_active = 1;" 2>/dev/null || echo "0")
    echo -e "${GREEN}✅ Database exists${NC}"
    echo -e "   Active subscribers: $SUBSCRIBER_COUNT"
else
    echo -e "${YELLOW}⚠️  Database not created yet${NC}"
    echo -e "   (Will be created when bot starts)"
fi

# Summary
echo ""
echo "╔════════════════════════════════════════════════════════════╗"
echo "║                    NEXT STEPS                              ║"
echo "╚════════════════════════════════════════════════════════════╝"
echo ""

STEPS_NEEDED=0

if [ "$API_KEY" = "PASTE_YOUR_API_KEY_HERE" ]; then
    STEPS_NEEDED=$((STEPS_NEEDED + 1))
    echo -e "${YELLOW}[ ] Step 1: Get Jellyfin API Key${NC}"
    echo "    1. Open http://10.255.255.254:8096 in browser"
    echo "    2. Go to Dashboard → API Keys"
    echo "    3. Click '+' button, name it 'Telegram Bot'"
    echo "    4. Copy the generated key"
    echo "    5. Edit .env and replace PASTE_YOUR_API_KEY_HERE with your key"
    echo ""
fi

WEBHOOK_STATUS=$(curl -s -f "$JELLYFIN_URL/System/Info" 2>/dev/null)
if [ -z "$WEBHOOK_STATUS" ]; then
    STEPS_NEEDED=$((STEPS_NEEDED + 1))
    echo -e "${YELLOW}[ ] Step 2: Configure Jellyfin Webhook${NC}"
    echo "    1. In Jellyfin: Dashboard → Plugins → Webhook"
    echo "    2. Click 'Add Generic Destination'"
    echo "    3. Set URL: http://$WSL_IP:$WEBHOOK_PORT/webhook"
    echo "    4. Check 'Item Added' notification"
    echo "    5. Check 'Movies' and 'Episodes' item types"
    echo "    6. Enable the webhook"
    echo ""
fi

if pgrep -x "jellyfin-bot" > /dev/null && grep -q "error\|unauthorized" logs/bot.log 2>/dev/null; then
    STEPS_NEEDED=$((STEPS_NEEDED + 1))
    echo -e "${YELLOW}[ ] Step 3: Restart Bot${NC}"
    echo "    After updating .env, restart the bot:"
    echo "    pkill jellyfin-bot && sleep 2 && nohup ./jellyfin-bot > logs/bot.log 2>&1 &"
    echo ""
fi

if [ $STEPS_NEEDED -eq 0 ]; then
    echo -e "${GREEN}✅ All configuration complete!${NC}"
    echo ""
    echo "Test the setup:"
    echo "  1. Send /start to your bot in Telegram"
    echo "  2. Add a movie/episode to Jellyfin"
    echo "  3. Check your Telegram for notification"
else
    echo -e "${YELLOW}📝 $STEPS_NEEDED configuration step(s) remaining${NC}"
fi

echo ""
echo "────────────────────────────────────────────────────────────"
echo "For detailed instructions, see: WSL_SETUP.md"
echo ""
