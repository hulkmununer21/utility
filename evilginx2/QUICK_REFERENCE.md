# Evilginx2 Web UI - Quick Reference

## 📋 Installation (5 minutes)

```bash
# 1. Clone and navigate
cd /workspaces/utility/evilginx2

# 2. Build frontend
chmod +x build_frontend.sh
./build_frontend.sh

# 3. Build backend
go build -o evilginx2 main.go

# 4. Run with API
./evilginx2 -p ./phishlets -t ./redirectors -api -api-pwd "MyPassword"

# 5. Open browser
# http://localhost:8080
```

---

## 🚀 Quick Commands

### Start with API
```bash
./evilginx2 -api
```

### Start with custom password
```bash
./evilginx2 -api -api-pwd "SecurePassword123"
```

### Start with custom port
```bash
./evilginx2 -api -api-addr "127.0.0.1:9000"
```

### Full configuration
```bash
./evilginx2 \
  -p ./phishlets \
  -t ./redirectors \
  -c ~/.evilginx \
  -api \
  -api-addr 127.0.0.1:8080 \
  -api-pwd "Password123" \
  -debug
```

### Docker
```bash
docker-compose up -d
```

---

## 🌐 Web UI Access

**URL:** `http://localhost:8080`  
**Default Password:** `admin`

### Pages
- 📊 **Sessions** - View/edit captured credentials
- 🎣 **Phishlets** - Enable/disable phishing targets
- 🔗 **Lures** - Create/manage lure links
- ⚙️ **Config** - Edit settings
- 💻 **Console** - Execute commands

---

## 🔑 API Authentication

Get token:
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"password":"admin"}'
```

Use in requests:
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/sessions
```

---

## 📡 Common API Calls

### Sessions
```bash
# List all
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/sessions

# Delete all
curl -X DELETE -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/sessions

# Get specific
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/sessions/1
```

### Phishlets
```bash
# List
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/phishlets

# Enable
curl -X POST -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/phishlets/google/enable

# Disable
curl -X POST -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/phishlets/google/disable
```

### Commands
```bash
# Execute command
curl -X POST -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"command":"phishlets"}' \
  http://localhost:8080/api/execute
```

---

## 🛠️ Troubleshooting

### Port in use
```bash
lsof -i :8080
kill -9 <PID>
# Or use different port: -api-addr "127.0.0.1:9000"
```

### Frontend not loading
```bash
rm -rf frontend/dist
./build_frontend.sh
```

### API not responding
```bash
# Check logs
./evilginx2 -api -debug

# Verify port
netstat -tulpn | grep 8080
```

### Permission denied
```bash
chmod +x setup.sh
chmod +x build_frontend.sh
```

---

## 📁 Directory Structure

```
evilginx2/
├── core/
│   ├── api_service.go        ← Service layer
│   ├── http_api.go           ← REST API
│   └── ...
├── frontend/
│   ├── src/
│   │   ├── pages/            ← UI pages
│   │   ├── styles/           ← CSS files
│   │   └── api/client.js     ← HTTP client
│   ├── dist/                 ← Built files (after build)
│   └── package.json
├── main.go                   ← Entry point (modified)
├── build_frontend.sh         ← Build script
├── setup.sh                  ← Setup wizard
├── Dockerfile
├── docker-compose.yml
└── WEB_UI_*.md               ← Documentation
```

---

## 🔐 Security Checklist

- [ ] Change default password with `-api-pwd`
- [ ] Run on localhost only (default)
- [ ] Use HTTPS in production
- [ ] Enable firewall rules
- [ ] Use strong passwords
- [ ] Regular database backups

---

## 📊 Default Configuration

| Setting | Default | Override |
|---------|---------|----------|
| API Port | 8080 | `-api-addr` |
| Bind Address | 127.0.0.1 | `-api-addr` |
| Password | admin | `-api-pwd` |
| Config Dir | ~/.evilginx | `-c` |
| Phishlets | ./phishlets | `-p` |
| Redirectors | ./redirectors | `-t` |

---

## 📚 Documentation Files

- **WEB_UI_README.md** - Quick start
- **WEB_UI_COMPLETE_GUIDE.md** - Comprehensive guide with examples
- **INTEGRATION_SUMMARY.md** - Technical details
- **QUICK_REFERENCE.md** - This file

---

## 🧪 Development Mode

Frontend development server (auto hot reload):
```bash
cd frontend
npm install
npm run dev
```

Access at: `http://localhost:5173`

---

## 🎯 Common Workflows

### Setup New Campaign
1. Open Web UI
2. Go to Phishlets → Enable target
3. Go to Lures → Create new lure
4. Share lure link
5. Monitor Sessions in real-time

### Check Captures
1. Login to Web UI
2. Go to Sessions
3. Review captured credentials
4. Export or delete as needed

### Modify Configuration
1. Go to Config page
2. Edit domain/IP as needed
3. Changes apply immediately

### Execute Commands
1. Go to Console page
2. Type command (e.g., `phishlets`, `sessions delete 1`)
3. View real-time output

---

## 🆘 Common Commands via Console

```
phishlets
phishlets enable <name>
phishlets disable <name>

sessions
sessions delete <id>
sessions clear

lures
lures delete <id>

config
config domain <domain.com>
config ipv4 external <ip>

help
help <command>
```

---

## 📞 Support

1. Check debug: `./evilginx2 -api -debug`
2. Review docs: WEB_UI_COMPLETE_GUIDE.md
3. Check API response error messages
4. Verify file permissions with `ls -la`

---

## ✅ Verification Steps

After setup:

- [ ] Frontend builds: `npm run build` in frontend/
- [ ] Backend compiles: `go build -o evilginx2 main.go`
- [ ] API starts: `./evilginx2 -api`
- [ ] Web accessible: `http://localhost:8080`
- [ ] API responds: `curl http://localhost:8080/api/config`
- [ ] Auth works: Login in web UI

---

## 🔗 Useful Links

- **Go Docs:** https://golang.org/doc
- **React Docs:** https://react.dev
- **Gorilla Mux:** https://github.com/gorilla/mux
- **Axios Docs:** https://axios-http.com
- **Vite Docs:** https://vitejs.dev

---

**Version:** 1.0  
**Last Updated:** 2024  
**Status:** Production Ready
