# Evilginx2 Enterprise Edition - Deployment Checklist

Use this checklist to deploy Evilginx2 safely in production with remote access enabled.

---

## 🔍 Pre-Deployment Review

### Security Audit
- [ ] Review all command-line flags and their purposes
- [ ] Understand the authentication flow (username + password)
- [ ] Confirm TLS certificates are valid for your domain
- [ ] Review firewall policy for port 9000
- [ ] Verify SSH keys are in place (for SSH tunnel method)
- [ ] Document any custom passwords or credentials

### System Requirements
- [ ] Go 1.16+ installed (for building)
- [ ] OpenSSL available (for certificate generation if needed)
- [ ] Sufficient disk space for binaries (~20 MB)
- [ ] Network connectivity verified
- [ ] Ports available (9000 or custom port)

### Network Architecture
- [ ] Document server IP address or hostname
- [ ] Confirm DNS records if using domain name
- [ ] Plan firewall rules (UFW, iptables, security groups)
- [ ] Decide on HTTP, HTTPS, or SSH tunnel method
- [ ] Test network path from client to server

---

## 🏗️ Build & Preparation

### Step 1: Clone and Build
```bash
# Clone repository
git clone https://github.com/username/evilginx2.git
cd evilginx2

# Build the project
make clean
make

# Verify binary
ls -lh ./evilginx2
file ./evilginx2
```

**Checklist:**
- [ ] Repository cloned successfully
- [ ] `make` command completed without errors
- [ ] Binary created at `./evilginx2`
- [ ] Binary size is approximately 10-20 MB

### Step 2: Prepare Configuration

**Option A: Simple Deployment**
```bash
# Keep defaults
# Port: 9000
# User: admin
# Password: admin
# Host: 0.0.0.0 (all interfaces)
```

**Option B: Custom Credentials**
```bash
# Record credentials
USERNAME="dashboard_admin"
PASSWORD="Generate strong password"
PORT="9000"
HOSTNAME="your-server.com"
```

**Checklist:**
- [ ] Choose username (not "admin")
- [ ] Generate strong password (minimum 16 chars, mixed case + numbers + symbols)
- [ ] Decide on port (9000 recommended)
- [ ] Hostname/IP address documented

### Step 3: TLS Certificate Setup

**For Self-Signed Certificate:**
```bash
# Generate certificate
openssl req -x509 -newkey rsa:4096 -nodes -days 365 \
  -out evilginx.crt -keyout evilginx.key \
  -subj "/CN=your-server.com"

# Move to secure location
sudo mkdir -p /etc/evilginx2
sudo mv evilginx.crt /etc/evilginx2/
sudo mv evilginx.key /etc/evilginx2/
sudo chmod 600 /etc/evilginx2/evilginx.key
```

**For Let's Encrypt Certificate:**
```bash
# (If using domain name)
# Use existing certificate management
# Typically at: /etc/letsencrypt/live/your-domain/

# Use existing certificate paths:
# -web-cert /etc/letsencrypt/live/your-domain/fullchain.pem
# -web-key /etc/letsencrypt/live/your-domain/privkey.pem
```

**Checklist:**
- [ ] Certificate file exists and is readable
- [ ] Private key file exists and is only readable by owner
- [ ] Certificate is valid (check expiry date)
- [ ] Paths are correct (absolute paths recommended)

---

## 🌐 Deployment Methods

### Method 1: HTTP (Local Network / SSH Tunnel)

**Recommended for:** Development, home lab, secure networks with SSH tunnel

```bash
# Start server on localhost only
./evilginx2 -web -web-host 127.0.0.1 -web-user admin -web-pass MySecretPass

# On client machine, create SSH tunnel
ssh -L 9000:127.0.0.1:9000 admin@server-ip

# Access dashboard
# Browser: http://localhost:9000
```

**Checklist:**
- [ ] SSH access configured
- [ ] Using non-default port if possible
- [ ] Strong password set
- [ ] Firewall allows SSH (port 22)
- [ ] SSH keys configured (password auth disabled preferred)

**Deployment Commands:**
```bash
# Start server
./evilginx2 \
  -p ./phishlets \
  -web \
  -web-host 127.0.0.1 \
  -web-port 9000 \
  -web-user admin \
  -web-pass "YourSecurePassword123"
```

### Method 2: HTTPS (Direct Remote Access)

**Recommended for:** Cloud servers, DMZ, multi-team access

```bash
# Start server with HTTPS
./evilginx2 -web \
  -web-https \
  -web-cert /etc/evilginx2/evilginx.crt \
  -web-key /etc/evilginx2/evilginx.key \
  -web-user admin \
  -web-pass "YourSecurePassword"

# Access from remote
# Browser: https://your-server-ip:9000
```

**Security Configuration:**
```bash
# 1. Firewall - Allow HTTPS traffic
sudo ufw allow 9000/tcp

# 2. Firewall - Restrict to specific IPs (optional)
sudo ufw allow from 192.168.1.0/24 to any port 9000

# 3. Firewall - Deny all other traffic
sudo ufw default deny incoming

# 4. Verify rules
sudo ufw status
```

**Checklist:**
- [ ] TLS certificate generated and valid
- [ ] Private key permissions correct (600)
- [ ] Firewall rules applied (`sudo ufw status` shows rules)
- [ ] Port 9000 is open and accessible from client
- [ ] Certificate matches hostname if using domain
- [ ] Tested HTTPS connection with curl

### Method 3: Reverse Proxy (Nginx/Apache)

**Recommended for:** Production, multiple services, load balancing

```nginx
# /etc/nginx/sites-available/evilginx2
server {
    listen 443 ssl http2;
    server_name evilginx.company.com;

    ssl_certificate /etc/ssl/certs/evilginx.crt;
    ssl_certificate_key /etc/ssl/private/evilginx.key;
    
    # Security headers
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    add_header Strict-Transport-Security "max-age=31536000" always;

    location / {
        proxy_pass http://127.0.0.1:9000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# HTTP redirect
server {
    listen 80;
    server_name evilginx.company.com;
    return 301 https://$server_name$request_uri;
}
```

**Setup Steps:**
```bash
# 1. Create configuration
sudo nano /etc/nginx/sites-available/evilginx2

# 2. Enable site
sudo ln -s /etc/nginx/sites-available/evilginx2 /etc/nginx/sites-enabled/

# 3. Test configuration
sudo nginx -t

# 4. Reload Nginx
sudo systemctl reload nginx

# 5. Start Evilginx2 on localhost
./evilginx2 -web -web-host 127.0.0.1
```

**Checklist:**
- [ ] Nginx installed and running
- [ ] Configuration file created and tested
- [ ] SSL certificates configured in Nginx
- [ ] DNS record points to server
- [ ] Firewall allows ports 80/443
- [ ] Evilginx2 listens on `127.0.0.1:9000`
- [ ] Tested access via https://your-domain.com

---

## 🚀 Production Startup

### Full Production Command

```bash
./evilginx2 \
  -p /opt/phishlets \
  -t /opt/redirectors \
  -c ~/.evilginx \
  -web \
  -web-port 9000 \
  -web-host 0.0.0.0 \
  -web-user prod_admin \
  -web-pass "ProductionPassword123!@#" \
  -web-https \
  -web-cert /etc/evilginx2/evilginx.crt \
  -web-key /etc/evilginx2/evilginx.key \
  -debug
```

**Checklist:**
- [ ] All paths exist and are readable
- [ ] Username is business-appropriate
- [ ] Password is strong and documented securely
- [ ] Certificates are in place
- [ ] Output shows no errors
- [ ] Web dashboard accessible
- [ ] CLI terminal working

### Systemd Service (Optional)

```ini
# /etc/systemd/system/evilginx2.service
[Unit]
Description=Evilginx2 Phishing Engine
After=network.target

[Service]
Type=simple
User=evilginx
WorkingDirectory=/opt/evilginx2
ExecStart=/opt/evilginx2/evilginx2 \
    -p /opt/phishlets \
    -web \
    -web-port 9000 \
    -web-user admin \
    -web-pass YourPassword \
    -web-https \
    -web-cert /etc/evilginx2/evilginx.crt \
    -web-key /etc/evilginx2/evilginx.key

Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
```

**Enable Service:**
```bash
sudo systemctl daemon-reload
sudo systemctl enable evilginx2
sudo systemctl start evilginx2
sudo systemctl status evilginx2
```

**Checklist:**
- [ ] Service file created in `/etc/systemd/system/`
- [ ] All paths are absolute
- [ ] Permissions are correct
- [ ] Service starts without errors
- [ ] Service restarts on failure
- [ ] Service survives system reboot

---

## ✅ Testing & Validation

### Step 1: Server Testing

```bash
# Check if service is running
ps aux | grep evilginx2

# Check if port is listening
netstat -tulpn | grep 9000
# or
lsof -i :9000

# Test with curl (HTTP)
curl -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}' \
  http://localhost:9000/api/auth/login

# Test with curl (HTTPS - ignore cert for self-signed)
curl -k -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}' \
  https://localhost:9000/api/auth/login
```

**Checklist:**
- [ ] Process running (`ps aux` shows evilginx2)
- [ ] Port listening (`netstat` shows port 9000)
- [ ] API responds to login request
- [ ] Response includes token
- [ ] No error messages in output

### Step 2: Remote Access Testing

```bash
# From another machine, test connectivity
curl -v http://server-ip:9000/api/auth/login

# If HTTPS
curl -k -v https://server-ip:9000/api/auth/login

# Test login in browser
# Navigate to: http://server-ip:9000
# Or: https://server-ip:9000
```

**Checklist:**
- [ ] Server responds from remote machine
- [ ] Login page loads
- [ ] Can submit username/password form
- [ ] Receives token on successful login
- [ ] Dashboard loads after login

### Step 3: Security Validation

```bash
# Verify firewall rules
sudo ufw status

# Check SSL certificate
openssl x509 -in /etc/evilginx2/evilginx.crt -text -noout

# Monitor access logs
tail -f ~/.evilginx/logs/* 2>/dev/null || echo "No logs yet"

# Test invalid credentials
curl -d '{"username":"admin","password":"wrongpass"}' \
  http://localhost:9000/api/auth/login
```

**Checklist:**
- [ ] Firewall configured correctly
- [ ] SSL certificate is valid
- [ ] Logs are being created
- [ ] Invalid login attempts are rejected
- [ ] No sensitive data in logs

---

## 📊 Post-Deployment

### Monitoring

```bash
# Monitor disk usage
df -h ~/.evilginx

# Monitor processes
top -p $(pidof evilginx2)

# Monitor connections
netstat -tulpn | grep 9000
```

**Checklist:**
- [ ] Disk usage reasonable
- [ ] Process memory stable
- [ ] Network connections normal
- [ ] No error messages

### Maintenance

- [ ] Backup phishlets directory weekly
- [ ] Backup sessions database weekly
- [ ] Rotate TLS certificates annually
- [ ] Update Evilginx2 when new versions available
- [ ] Review and archive logs monthly

### Documentation

- [ ] Record final deployment command
- [ ] Document custom credentials (store securely)
- [ ] Document TLS certificate paths
- [ ] Document firewall rules
- [ ] Document backup procedures

---

## 🆘 Troubleshooting

| Problem | Solution |
|---------|----------|
| Port 9000 already in use | Use `-web-port 8888` or kill existing process |
| Can't connect from remote | Check firewall: `sudo ufw allow 9000/tcp` |
| HTTPS certificate error | Regenerate or check path: `openssl x509 -in cert.pem -text` |
| Login fails | Verify username/password flags match |
| Service won't start | Check permissions: `ls -l /opt/evilginx2/evilginx2` |
| High memory usage | Reduce number of open phishlets or restart service |

---

## ✨ Final Checklist

- [ ] Build successful
- [ ] Credentials configured and noted
- [ ] TLS certificates in place
- [ ] Firewall rules applied
- [ ] Deployment method chosen (HTTP/HTTPS/Proxy)
- [ ] Server started without errors
- [ ] Local tests passed
- [ ] Remote tests passed
- [ ] Security validation passed
- [ ] Service configured (if using systemd)
- [ ] Monitoring in place
- [ ] Backups configured
- [ ] Documentation updated

---

## 🎉 Deployment Complete!

Your Evilginx2 enterprise deployment is ready for production use. Users can now:

✅ Access the dashboard remotely  
✅ Log in with username + password  
✅ Manage phishing campaigns securely  
✅ Monitor captured sessions  

**Next Steps:**
1. Load your phishlets
2. Configure redirectors
3. Monitor incoming sessions
4. Review captured data

**For Support:** See [REMOTE_ACCESS_GUIDE.md](REMOTE_ACCESS_GUIDE.md) for configuration options and examples.

---

*Last Updated: 2024*  
*Version: 2.0 Enterprise Edition*
