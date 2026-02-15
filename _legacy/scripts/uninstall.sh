#!/bin/bash

# Pryx Clean Uninstaller
# Removes all Pryx binaries, services, and data.

echo "⚠️  Pryx Clean Uninstaller"
echo "This will permanently delete all Pryx data including sessions, configuration, and keys."
read -p "Are you sure you want to proceed? (y/N) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Aborted."
    exit 1
fi

# 1. Stop and remove services
echo "Stopping services..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    launchctl unload ~/Library/LaunchAgents/com.pryx.runtime.plist 2>/dev/null
    rm ~/Library/LaunchAgents/com.pryx.runtime.plist 2>/dev/null
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    sudo systemctl stop pryx.service 2>/dev/null
    sudo systemctl disable pryx.service 2>/dev/null
    sudo rm /etc/systemd/system/pryx.service 2>/dev/null
    sudo systemctl daemon-reload 2>/dev/null
fi

# 2. Remove binaries
echo "Removing binaries..."
rm -f /usr/local/bin/pryx 2>/dev/null
rm -f /usr/local/bin/pryx-core 2>/dev/null
rm -f ~/.local/bin/pryx 2>/dev/null
rm -f ~/.local/bin/pryx-core 2>/dev/null

# 3. Remove data directory
echo "Removing data directory (~/.pryx)..."
rm -rf ~/.pryx

echo "✓ Pryx has been cleanly uninstalled."
