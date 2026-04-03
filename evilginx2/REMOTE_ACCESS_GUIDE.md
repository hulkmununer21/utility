# Evilginx2 Enhanced Configuration & Remote Access

## 🚀 Single Command Startup

Now start both CLI and dashboard with a single command:

```bash
./evilginx2 -p ./phishlets -t ./redirectors -web
```

This automatically:
- Starts the terminal CLI interface
- Starts web admin dashboard on port 9000
- Binds to `0.0.0.0` for remote access
- Uses default credentials: `admin` / `admin`

---

## 🔧 Configuration Options

### Basic (Default Settings)

```bash
# Start with web UI on port 9000 (all interfaces)
./evilginx2 -web

# Username: admin
# Password: admin
# Access: http://your-server-ip:9000
```

### Custom Port

```bash
# Use port 8888 instead
./evilginx2 -web -web-port 8888

# Access: http://your-server-ip:8888
```

### Custom Credentials

```bash
# Set custom username and password
./evilginx2 -web -web-user "myuser" -web-pass "MySecurePass123"
```

### Localhost Only (Development)

```bash
# Restrict to localhost only
./evilginx2 -web -web-host 127.0.0.1

# Access: http://localhost:9000
```

### Specific Network Interface

```bash
# Bind to specific IP
./evilginx2 -web -web-host 192.168.1.100

# Access: http://192.168.1.100:9000
```

### HTTPS/TLS (Secure Remote Access)

```bash
# Generate self-signed certificate (if you don't have one)
openssl req -x509 -newkey rsa:4096 -nodes -out cert.pem -keyout key.pem -days 365

# Start with HTTPS
./evilginx2 -web -web-https -web-cert ./cert.pem -web-key ./key.pem

# Access: https://your-server-ip:9000
```

### Full Production Setup

```bash
./evilginx2 \
  -p ./phishlets \
  -t ./redirectors \
  -c ~/.evilginx \
  -web \
  -web-port 9000 \
  -web-host 0.0.0.0 \
  -web-user "admin_user" \
  -web-pass "StrongPassword123" \
  -web-https \
  -web-cert /etc/ssl/certs/evilginx.crt \
  -web-key /etc/ssl/private/evilginx.key \
  -debug
```

---

## 🌐 Remote Access from Another Machine

### Step-by-Step

**1. On Server (running Evilginx2):**
```bash
./evilginx2 -web -web-host 0.0.0.0 -web-user "operator" -web-pass "SecurePass123"
```

This starts the dashboard listening on all network interfaces (`0.0.0.0:9000`)

**2. From Client Machine (your laptop/workstation):**

Option A - Direct Access:
```bash
# Open browser
http://server-ip:9000

# Login with:
Username: operator
Password: SecurePass123
```

Option B - SSH Tunnel (Secure):
```bash
# Create tunnel from local 9000 to remote 9000
ssh -L 9000:localhost:9000 user@server-ip

# Then open browser
http://localhost:9000
```

Option C - VPN/Firewall:
```bash
# Configure firewall to allow port 9000
sudo ufw allow 9000/tcp

# Configure firewall rules (UFW)
sudo ufw allow from 192.168.1.0/24 to any port 9000
```

---

## 🔐 Security Best Practices

### 1. Use HTTPS in Production

```bash
# Generate certificate
sudo openssl req -x509 -nodes -days 365 -newkey rsa:4096 \
  -keyout /etc/ssl/private/evilginx.key \
  -out /etc/ssl/certs/evilginx.crt

# Start with HTTPS
./evilginx2 -web -web-https \
  -web-cert /etc/ssl/certs/evilginx.crt \
  -web-key /etc/ssl/private/evilginx.key
```

### 2. Strong Credentials

```bash
# Use complex password
./evilginx2 -web \
  -web-user "dashboard_admin" \
  -web-pass "$(openssl rand -base64 32)"
```

### 3. Firewall Rules

```bash
# Allow only specific IP ranges
sudo ufw allow from 192.168.1.0/24 to any port 9000

# Deny all others
sudo ufw deny 9000/tcp
```

### 4. SSH Tunnel (Recommended for Remote Access)

```bash
# Server side
./evilginx2 -web -web-host 127.0.0.1

# Client side
ssh -L 9000:127.0.0.1:9000 user@server

# Then access locally at http://localhost:9000
```

### 5. Reverse Proxy with Nginx

```nginx
server {
    listen 443 ssl;
    server_name evilginx.example.com;

    ssl_certificate /etc/ssl/certs/evilginx.crt;
    ssl_certificate_key /etc/ssl/private/evilginx.key;

    location / {
        proxy_pass http://127.0.0.1:9000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

---

## 📊 API Reference

### Authentication

```bash
# Login
curl -X POST http://your-server:9000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'

# Response
{
  "success": true,
  "data": {
    "token": "1234567890",
    "expiry": 1723456789
  }
}
```

### All Endpoints

```bash
# Set token
TOKEN="your-token"

# Get sessions
curl -H "Authorization: Bearer $TOKEN" \
  http://your-server:9000/api/sessions

# Get phishlets
curl -H "Authorization: Bearer $TOKEN" \
  http://your-server:9000/api/phishlets

# Get config
curl -H "Authorization: Bearer $TOKEN" \
  http://your-server:9000/api/config

# Execute command
curl -X POST -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"command":"phishlets"}' \
  http://your-server:9000/api/execute
```

---

## 🎯 Common Scenarios

### Scenario 1: Home Lab Testing
```bash
# Start on localhost only
./evilginx2 -web -web-host 127.0.0.1

# Access: http://localhost:9000
```

### Scenario 2: Team Collaboration
```bash
# Allow team access via HTTPS
./evilginx2 -web \
  -web-user "team_admin" \
  -web-pass "TeamPass123" \
  -web-https \
  -web-cert ./certs/server.crt \
  -web-key ./certs/server.key

# Team members access: https://server-ip:9000
```

### Scenario 3: Remote Management VIA SSH
```bash
# Server starts on localhost only (secure)
./evilginx2 -web -web-host 127.0.0.1

# Admin creates SSH tunnel
ssh -L 9000:127.0.0.1:9000 admin@server

# Admin accesses locally
# http://localhost:9000
```

### Scenario 4: Docker with Remote Access
```bash
# Docker compose with external access
docker-compose up -d

# Services exposed at:
# HTTPS: https://server:9000
# SSH Tunnel available
```

---

## 🔍 Monitoring & Logs

### View Logs

```bash
# Start with debug logging
./evilginx2 -web -debug

# Logs show:
# - Web dashboard startup
# - Login attempts
# - API requests
# - Errors and warnings
```

### Monitor Sessions

```bash
# In CLI terminal, use
sessions

# Or via API
curl -H "Authorization: Bearer $TOKEN" \
  http://server:9000/api/sessions
```

---

## 🚨 Troubleshooting

### Port Already in Use

```bash
# Check what's using port 9000
lsof -i :9000

# Use different port
./evilginx2 -web -web-port 8888
```

### Can't Connect from Remote

```bash
# Verify firewall
sudo ufw status
sudo ufw allow 9000/tcp

# Verify binding
netstat -tulpn | grep 9000

# Check HTTPS cert
ls -la /etc/ssl/certs/
```

### Login Failed

```bash
# Verify credentials
./evilginx2 -web -web-user "admin" -web-pass "admin"

# Check API response
curl -v http://server:9000/api/auth/login
```

### HTTPS Issues

```bash
# Verify certificate
openssl x509 -in cert.pem -text -noout

# Recreate self-signed cert
openssl req -x509 -newkey rsa:4096 -nodes \
  -out cert.pem -keyout key.pem -days 365
```

---

## 💡 Standard Features Added

✅ **Remote Access** - Bind to all interfaces (0.0.0.0)  
✅ **Username + Password Auth** - Dual credential support  
✅ **Custom Port** - Configure any port (default 9000)  
✅ **HTTPS/TLS** - Secure remote connections  
✅ **Single Command** - Start CLI + Web with `-web` flag  
✅ **Token-based Auth** - 24-hour expiry  
✅ **CORS Support** - Cross-origin requests allowed  
✅ **API Documentation** - Self-documenting endpoints  
✅ **Error Handling** - Comprehensive error messages  
✅ **Audit Logging** - All actions logged  

---

## 🎓 Examples

### Example 1: Quick Lab
```bash
./evilginx2 -web
# Access: http://localhost:9000
# Admin: admin / admin
```

### Example 2: Team Collaboration
```bash
./evilginx2 -web \
  -web-user "operator" \
  -web-pass "TeamPass123" \
  -web-port 9000

# Access: http://server-ip:9000
# Admin: operator / TeamPass123
```

### Example 3: Production Deployment
```bash
./evilginx2 \
  -p /opt/phishlets \
  -t /opt/redirectors \
  -web \
  -web-user "prod_admin" \
  -web-pass "ProdSecurePass" \
  -web-https \
  -web-cert /etc/ssl/certs/evilginx.crt \
  -web-key /etc/ssl/private/evilginx.key

# Access: https://evilginx.company.com:9000
```

---

**Version: 2.0 - Enhanced with Remote Access & Security**  
**Status: Production Ready**
