#!/bin/bash

# Color codes for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "============================================"
echo "   Jellyfin Telegram Bot - Setup Tester"
echo "============================================"
echo ""

# Check if binary exists
echo -n "Checking bot binary... "
if [ -f "./jellyfin-bot" ]; then
    echo -e "${GREEN}✓ Found${NC}"
else
    echo -e "${RED}✗ Not found${NC}"
    echo "Run: go build -o jellyfin-bot cmd/bot/main.go"
    exit 1
fi

# Check if .env exists
echo -n "Checking .env file... "
if [ -f ".env" ]; then
    echo -e "${GREEN}✓ Found${NC}"

    # Check for required variables
    echo ""
    echo "Checking environment variables:"

    if grep -q "TELEGRAM_BOT_TOKEN=" .env && ! grep -q "TELEGRAM_BOT_TOKEN=YOUR" .env; then
        echo -e "  TELEGRAM_BOT_TOKEN: ${GREEN}✓ Set${NC}"
    else
        echo -e "  TELEGRAM_BOT_TOKEN: ${RED}✗ Not set${NC}"
    fi

    if grep -q "JELLYFIN_SERVER_URL=" .env && ! grep -q "JELLYFIN_SERVER_URL=YOUR" .env; then
        echo -e "  JELLYFIN_SERVER_URL: ${GREEN}✓ Set${NC}"
    else
        echo -e "  JELLYFIN_SERVER_URL: ${RED}✗ Not set${NC}"
    fi

    if grep -q "JELLYFIN_API_KEY=" .env && ! grep -q "JELLYFIN_API_KEY=YOUR" .env; then
        echo -e "  JELLYFIN_API_KEY: ${GREEN}✓ Set${NC}"
    else
        echo -e "  JELLYFIN_API_KEY: ${RED}✗ Not set${NC}"
    fi

else
    echo -e "${RED}✗ Not found${NC}"
    echo ""
    echo "Creating .env template..."
    cp .env.example .env 2>/dev/null || cat > .env <<'EOF'
# Telegram Configuration
TELEGRAM_BOT_TOKEN=YOUR_BOT_TOKEN_HERE

# Jellyfin Configuration
JELLYFIN_SERVER_URL=http://YOUR_JELLYFIN_SERVER:8096
JELLYFIN_API_KEY=YOUR_API_KEY_HERE

# Webhook Configuration
WEBHOOK_PORT=8080
WEBHOOK_SECRET=optional-secret-key

# Database Configuration
DATABASE_PATH=./jellyfin_bot.db

# Logging Configuration
LOG_LEVEL=INFO
LOG_FILE=./logs/bot.log
EOF
    echo -e "${YELLOW}Created .env file - please edit it with your credentials${NC}"
    exit 1
fi

# Check if logs directory exists
echo ""
echo -n "Checking logs directory... "
if [ -d "logs" ]; then
    echo -e "${GREEN}✓ Found${NC}"
else
    echo -e "${YELLOW}Creating...${NC}"
    mkdir -p logs
    echo -e "${GREEN}✓ Created${NC}"
fi

# Test if bot can start (dry run)
echo ""
echo "Testing bot startup (5 second test)..."
echo "----------------------------------------"

timeout 5s ./jellyfin-bot > /tmp/bot_test.log 2>&1 &
BOT_PID=$!
sleep 2

if ps -p $BOT_PID > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Bot started successfully!${NC}"
    kill $BOT_PID 2>/dev/null
    wait $BOT_PID 2>/dev/null

    # Show last few lines of output
    echo ""
    echo "Bot output:"
    tail -5 /tmp/bot_test.log
else
    echo -e "${RED}✗ Bot failed to start${NC}"
    echo ""
    echo "Error output:"
    cat /tmp/bot_test.log
    exit 1
fi

# Summary
echo ""
echo "============================================"
echo "          Setup Status Summary"
echo "============================================"
echo ""

ALL_GOOD=true

# Check all requirements
if [ -f "./jellyfin-bot" ]; then
    echo -e "Binary:      ${GREEN}✓ Ready${NC}"
else
    echo -e "Binary:      ${RED}✗ Missing${NC}"
    ALL_GOOD=false
fi

if [ -f ".env" ]; then
    if grep -q "YOUR" .env; then
        echo -e ".env file:   ${YELLOW}⚠ Needs configuration${NC}"
        ALL_GOOD=false
    else
        echo -e ".env file:   ${GREEN}✓ Configured${NC}"
    fi
else
    echo -e ".env file:   ${RED}✗ Missing${NC}"
    ALL_GOOD=false
fi

echo ""

if [ "$ALL_GOOD" = true ]; then
    echo -e "${GREEN}Everything looks good! Ready to run.${NC}"
    echo ""
    echo "To start the bot:"
    echo "  ./jellyfin-bot"
    echo ""
    echo "To run in background:"
    echo "  nohup ./jellyfin-bot > logs/bot.log 2>&1 &"
    echo ""
    echo "To view logs:"
    echo "  tail -f logs/bot.log"
else
    echo -e "${YELLOW}Please complete the setup:${NC}"
    echo "  1. Edit .env with your credentials"
    echo "  2. Get Telegram bot token from @BotFather"
    echo "  3. Get Jellyfin API key from Dashboard → API Keys"
    echo ""
    echo "See SETUP_GUIDE.md for detailed instructions"
fi

echo ""
echo "============================================"

# Cleanup
rm -f /tmp/bot_test.log
