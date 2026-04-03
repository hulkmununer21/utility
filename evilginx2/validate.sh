#!/bin/bash
# Validation script to verify Web UI integration

echo "================================"
echo "Evilginx2 Web UI Validator"
echo "================================"
echo ""

ERRORS=0

# Check backend files
echo "[*] Checking backend files..."
files_to_check=(
    "core/api_service.go"
    "core/http_api.go"
    "main.go"
)

for file in "${files_to_check[@]}"; do
    if [ -f "$file" ]; then
        echo "[+] $file ✓"
    else
        echo "[-] $file ✗"
        ((ERRORS++))
    fi
done

echo ""
echo "[*] Checking frontend structure..."
frontend_files=(
    "frontend/package.json"
    "frontend/vite.config.js"
    "frontend/index.html"
    "frontend/src/App.jsx"
    "frontend/src/index.jsx"
    "frontend/src/api/client.js"
    "frontend/src/pages/LoginPage.jsx"
    "frontend/src/pages/SessionsPage.jsx"
    "frontend/src/pages/PhishletsPage.jsx"
    "frontend/src/pages/ConsolePage.jsx"
    "frontend/src/pages/ConfigPage.jsx"
)

for file in "${frontend_files[@]}"; do
    if [ -f "$file" ]; then
        echo "[+] $file ✓"
    else
        echo "[-] $file ✗"
        ((ERRORS++))
    fi
done

echo ""
echo "[*] Checking documentation..."
docs=(
    "WEB_UI_README.md"
    "WEB_UI_COMPLETE_GUIDE.md"
    "INTEGRATION_SUMMARY.md"
    "QUICK_REFERENCE.md"
)

for doc in "${docs[@]}"; do
    if [ -f "$doc" ]; then
        echo "[+] $doc ✓"
    else
        echo "[-] $doc ✗"
        ((ERRORS++))
    fi
done

echo ""
echo "[*] Checking deployment files..."
deploy_files=(
    "setup.sh"
    "Dockerfile"
    "docker-compose.yml"
    ".env.example"
)

for file in "${deploy_files[@]}"; do
    if [ -f "$file" ]; then
        echo "[+] $file ✓"
    else
        echo "[-] $file ✗"
        ((ERRORS++))
    fi
done

echo ""
echo "[*] Checking CSS files..."
css_files=(
    "frontend/src/styles/app.css"
    "frontend/src/styles/auth.css"
    "frontend/src/styles/sessions.css"
    "frontend/src/styles/phishlets.css"
    "frontend/src/styles/console.css"
    "frontend/src/styles/config.css"
    "frontend/src/index.css"
)

for file in "${css_files[@]}"; do
    if [ -f "$file" ]; then
        echo "[+] $file ✓"
    else
        echo "[-] $file ✗"
        ((ERRORS++))
    fi
done

echo ""
echo "================================"
if [ $ERRORS -eq 0 ]; then
    echo "[+] All files present! ✓"
    echo ""
    echo "Next steps:"
    echo "1. chmod +x build_frontend.sh setup.sh validate.sh"
    echo "2. ./build_frontend.sh"
    echo "3. go build -o evilginx2 main.go"
    echo "4. ./evilginx2 -api -api-pwd 'your_password'"
    echo ""
else
    echo "[-] $ERRORS file(s) missing!"
    echo ""
fi
echo "================================"
