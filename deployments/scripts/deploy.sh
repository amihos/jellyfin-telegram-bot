#!/bin/bash
# Deployment script for Jellyfin Telegram Bot

set -e  # Exit on error

echo "==================================="
echo "Jellyfin Telegram Bot Deployment"
echo "==================================="

# Configuration
APP_NAME="jellyfin-bot"
INSTALL_DIR="/opt/${APP_NAME}"
SERVICE_FILE="deployments/systemd/${APP_NAME}.service"
BINARY_NAME="${APP_NAME}"
USER="${APP_NAME}"
GROUP="${APP_NAME}"

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo "ERROR: This script must be run as root (use sudo)"
    exit 1
fi

# Step 1: Build the binary
echo ""
echo "Step 1: Building binary..."
if ! command -v go &> /dev/null; then
    echo "ERROR: Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

go build -o ${BINARY_NAME} cmd/bot/main.go
if [ ! -f "${BINARY_NAME}" ]; then
    echo "ERROR: Build failed"
    exit 1
fi
echo "  Binary built successfully: ${BINARY_NAME}"

# Step 2: Create user and group (if they don't exist)
echo ""
echo "Step 2: Creating system user..."
if ! id -u ${USER} > /dev/null 2>&1; then
    useradd --system --no-create-home --shell /bin/false ${USER}
    echo "  User '${USER}' created"
else
    echo "  User '${USER}' already exists"
fi

# Step 3: Create installation directory
echo ""
echo "Step 3: Setting up installation directory..."
mkdir -p ${INSTALL_DIR}
chown ${USER}:${GROUP} ${INSTALL_DIR}
chmod 755 ${INSTALL_DIR}
echo "  Installation directory created: ${INSTALL_DIR}"

# Step 4: Copy binary to installation directory
echo ""
echo "Step 4: Installing binary..."
cp ${BINARY_NAME} ${INSTALL_DIR}/
chown ${USER}:${GROUP} ${INSTALL_DIR}/${BINARY_NAME}
chmod 755 ${INSTALL_DIR}/${BINARY_NAME}
echo "  Binary installed to ${INSTALL_DIR}/${BINARY_NAME}"

# Step 5: Copy .env.example if .env doesn't exist
echo ""
echo "Step 5: Setting up configuration..."
if [ ! -f "${INSTALL_DIR}/.env" ]; then
    if [ -f ".env.example" ]; then
        cp .env.example ${INSTALL_DIR}/.env
        chown ${USER}:${GROUP} ${INSTALL_DIR}/.env
        chmod 600 ${INSTALL_DIR}/.env
        echo "  Configuration template copied to ${INSTALL_DIR}/.env"
        echo "  WARNING: Please edit ${INSTALL_DIR}/.env with your actual credentials!"
    else
        echo "  WARNING: .env.example not found. Please create ${INSTALL_DIR}/.env manually."
    fi
else
    echo "  Configuration file already exists: ${INSTALL_DIR}/.env"
fi

# Step 6: Install systemd service
echo ""
echo "Step 6: Installing systemd service..."
if [ -f "${SERVICE_FILE}" ]; then
    cp ${SERVICE_FILE} /etc/systemd/system/${APP_NAME}.service

    # Update service file with correct paths and environment file
    sed -i "s|WorkingDirectory=.*|WorkingDirectory=${INSTALL_DIR}|g" /etc/systemd/system/${APP_NAME}.service
    sed -i "s|ExecStart=.*|ExecStart=${INSTALL_DIR}/${BINARY_NAME}|g" /etc/systemd/system/${APP_NAME}.service
    sed -i "s|# EnvironmentFile=.*|EnvironmentFile=${INSTALL_DIR}/.env|g" /etc/systemd/system/${APP_NAME}.service

    # Reload systemd
    systemctl daemon-reload
    echo "  Systemd service installed: /etc/systemd/system/${APP_NAME}.service"
else
    echo "  ERROR: Service file not found: ${SERVICE_FILE}"
    exit 1
fi

# Step 7: Initialize database (if needed)
echo ""
echo "Step 7: Database setup..."
echo "  Database will be created automatically on first run at ${INSTALL_DIR}/bot.db"

# Step 8: Enable and start service
echo ""
echo "Step 8: Starting service..."
read -p "Do you want to start the service now? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    systemctl enable ${APP_NAME}
    systemctl start ${APP_NAME}
    echo "  Service enabled and started"

    # Show status
    echo ""
    echo "Service status:"
    systemctl status ${APP_NAME} --no-pager
else
    echo "  Service not started. You can start it later with:"
    echo "    sudo systemctl start ${APP_NAME}"
    echo "  And enable it to start on boot with:"
    echo "    sudo systemctl enable ${APP_NAME}"
fi

echo ""
echo "==================================="
echo "Deployment Complete!"
echo "==================================="
echo ""
echo "Important next steps:"
echo "1. Edit ${INSTALL_DIR}/.env with your actual credentials"
echo "2. Restart the service: sudo systemctl restart ${APP_NAME}"
echo "3. Check logs: sudo journalctl -u ${APP_NAME} -f"
echo "4. Test health endpoint: curl http://localhost:8080/health"
echo "5. Configure Jellyfin webhook: http://your-server:8080/webhook"
echo ""
