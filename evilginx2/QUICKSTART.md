# Evilginx2 - Quick Start Guide

## ⚡ 30-Second Startup

```bash
# Build
make

# Run with web dashboard
./evilginx2 -web
```

**That's it!** 🎉

### Access the Dashboard
- **URL:** http://localhost:9000
- **Username:** admin
- **Password:** admin

---

## 🌐 Accessing from Another Machine

**Option 1: Direct Access**
```bash
# On server
./evilginx2 -web

# On your laptop
# Open browser: http://server-ip:9000
```

**Option 2: SSH Tunnel (More Secure)**
```bash
# On your laptop
ssh -L 9000:localhost:9000 user@server

# Then open browser: http://localhost:9000
```

**Option 3: Custom Credentials**
```bash
# On server
./evilginx2 -web -web-user "myuser" -web-pass "mypass"

# On your laptop
# Browser: http://server-ip:9000
# Login: myuser / mypass
```

---

## 🔒 Secure Remote Access (HTTPS)

```bash
# Generate certificate (one-time)
openssl req -x509 -newkey rsa:4096 -nodes \
  -out cert.pem -keyout key.pem -days 365

# Run with HTTPS
./evilginx2 -web -web-https -web-cert cert.pem -web-key key.pem

# Access: https://server-ip:9000
```

---

## 📋 Common Commands

| Task | Command |
|------|---------|
| **Quick start** | `./evilginx2 -web` |
| **Custom port** | `./evilginx2 -web -web-port 8888` |
| **Custom user/pass** | `./evilginx2 -web -web-user admin -web-pass secret` |
| **With HTTPS** | `./evilginx2 -web -web-https -web-cert cert.pem -web-key key.pem` |
| **Localhost only** | `./evilginx2 -web -web-host 127.0.0.1` |
| **All options** | `./evilginx2 -web -web-host 0.0.0.0 -web-port 9000 -web-user admin` |

---

## ✅ What Works Now

✅ Start CLI + Web dashboard with **one command**  
✅ Access from **any machine** on your network  
✅ **Username + password** login  
✅ **HTTPS/TLS** support  
✅ **Port 9000** as standard  
✅ Single `-web` flag enables everything  

---

## 📚 More Details

See [REMOTE_ACCESS_GUIDE.md](REMOTE_ACCESS_GUIDE.md) for:
- Security best practices
- Firewall configuration
- Nginx reverse proxy setup
- API reference
- Troubleshooting

---

**That's all you need! Go build something awesome! 🚀**
