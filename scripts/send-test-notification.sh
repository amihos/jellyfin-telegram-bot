#!/bin/bash

# Mock Webhook Sender for Testing
# Usage: ./scripts/send-test-notification.sh [type] [series_name] [season] [episode]
#
# Examples:
#   ./scripts/send-test-notification.sh episode "Breaking Bad" 5 3
#   ./scripts/send-test-notification.sh movie "The Matrix"

set -e

# Configuration
WEBHOOK_URL="${WEBHOOK_URL:-http://localhost:8080/webhook}"
WEBHOOK_SECRET="${WEBHOOK_SECRET:-}"

# Parse arguments
TYPE="${1:-episode}"
SERIES_NAME="${2:-Test Series}"
SEASON="${3:-1}"
EPISODE="${4:-1}"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}📤 Sending test notification...${NC}"
echo -e "${YELLOW}Type:${NC} $TYPE"

if [ "$TYPE" == "episode" ]; then
    echo -e "${YELLOW}Series:${NC} $SERIES_NAME"
    echo -e "${YELLOW}Season:${NC} $SEASON"
    echo -e "${YELLOW}Episode:${NC} $EPISODE"

    # Generate episode webhook payload
    PAYLOAD=$(cat <<EOF
{
  "NotificationType": "ItemAdded",
  "ItemType": "Episode",
  "ItemId": "test-$(date +%s)",
  "ItemName": "$SERIES_NAME - S${SEASON}E${EPISODE}",
  "SeriesName": "$SERIES_NAME",
  "SeasonNumber": "$SEASON",
  "EpisodeNumber": "$EPISODE",
  "Year": "2024",
  "Overview": "This is a test notification for episode $EPISODE of season $SEASON.",
  "Provider_tmdb": "12345",
  "Video_0_Title": "1080p H.264",
  "Video_0_Type": "Video",
  "Video_0_Codec": "h264"
}
EOF
)
else
    echo -e "${YELLOW}Movie:${NC} $SERIES_NAME"

    # Generate movie webhook payload
    PAYLOAD=$(cat <<EOF
{
  "NotificationType": "ItemAdded",
  "ItemType": "Movie",
  "ItemId": "test-movie-$(date +%s)",
  "ItemName": "$SERIES_NAME",
  "Year": "2024",
  "Overview": "This is a test notification for movie $SERIES_NAME.",
  "Provider_tmdb": "67890",
  "Video_0_Title": "1080p H.264",
  "Video_0_Type": "Video",
  "Video_0_Codec": "h264"
}
EOF
)
fi

echo -e "\n${BLUE}Payload:${NC}"
echo "$PAYLOAD" | jq '.' 2>/dev/null || echo "$PAYLOAD"

echo -e "\n${BLUE}Sending to:${NC} $WEBHOOK_URL"

# Send webhook
RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" \
    -X POST \
    -H "Content-Type: application/json" \
    $([ -n "$WEBHOOK_SECRET" ] && echo "-H \"X-Webhook-Secret: $WEBHOOK_SECRET\"") \
    -d "$PAYLOAD" \
    "$WEBHOOK_URL")

HTTP_CODE=$(echo "$RESPONSE" | grep HTTP_CODE | cut -d':' -f2)
BODY=$(echo "$RESPONSE" | sed '/HTTP_CODE/d')

echo -e "\n${BLUE}Response:${NC}"
if [ "$HTTP_CODE" == "200" ] || [ "$HTTP_CODE" == "202" ]; then
    echo -e "${GREEN}✅ Success! (HTTP $HTTP_CODE)${NC}"
    [ -n "$BODY" ] && echo "$BODY"
else
    echo -e "${YELLOW}⚠️  HTTP $HTTP_CODE${NC}"
    [ -n "$BODY" ] && echo "$BODY"
fi

echo -e "\n${BLUE}💡 Tips:${NC}"
echo "  - Check your Telegram for the notification"
echo "  - If testing mute button, click it and send another notification"
echo "  - Check bot logs for detailed information"
echo ""
