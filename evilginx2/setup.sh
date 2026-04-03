#!/bin/bash
# Quick deployment script for Evilginx2 with Web UI

set -e

echo "================================"
echo "Evilginx2 Web UI Setup Script"
echo "================================"
echo ""

# Check prerequisites
echo "[*] Checking prerequisites..."

if ! command -v go &> /dev/null; then
    echo "[-] Go is not installed"
    exit 1
fi

if ! command -v node &> /dev/null; then
    echo "[-] Node.js is not installed. Installing..."
    curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
    sudo apt-get install -y nodejs
fi

echo "[+] Prerequisites OK"
echo ""

# Get configuration
echo "[*] Configuration"
read -p "Enter API password (default: admin): " API_PWD
API_PWD=${API_PWD:-admin}

read -p "Enter API port (default: 8080): " API_PORT
API_PORT=${API_PORT:-8080}

read -p "Enter API bind address (default: 127.0.0.1): " API_ADDR
API_ADDR=${API_ADDR:-127.0.0.1}

echo ""
echo "[*] Building frontend..."
chmod +x build_frontend.sh
./build_frontend.sh
echo "[+] Frontend build complete"
echo ""

echo "[*] Building Go backend..."
go build -o evilginx2 main.go
echo "[+] Backend build complete"
echo ""

echo "[*] Verifying configuration directories..."
mkdir -p phishlets redirectors
echo "[+] Directories ready"
echo ""

echo "================================"
echo "Setup Complete!"
echo "================================"
echo ""
echo "Start the application with:"
echo ""
echo "  ./evilginx2 -p ./phishlets -t ./redirectors \\"
echo "    -api -api-addr $API_ADDR:$API_PORT -api-pwd '$API_PWD'"
echo ""
echo "Then open your browser to:"
echo "  http://$API_ADDR:$API_PORT"
echo ""
echo "Login password: $API_PWD"
echo ""
