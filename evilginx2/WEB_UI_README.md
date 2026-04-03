# Web UI Integration Guide

## Quick Start

### 1. Build the Frontend
```bash
chmod +x build_frontend.sh
./build_frontend.sh
```

This runs:
- `npm install` - Installs React and dependencies
- `npm run build` - Builds the React app into `frontend/dist`

### 2. Start Evilginx2 with Web API

```bash
./evilginx2 -api -api-pwd "your_secure_password"
```

**Flags:**
- `-api` - Enable web API server
- `-api-addr "127.0.0.1:8080"` - Set API address (default: 127.0.0.1:8080)
- `-api-pwd "password"` - Set API password (default: admin)

### 3. Access the Web UI
Open your browser and go to:
```
http://localhost:8080
```

Login with the password you set.

---

## API Endpoints

### Authentication
- `POST /api/auth/login` - Login with password

### Sessions
- `GET /api/sessions` - List all sessions
- `GET /api/sessions/{id}` - Get specific session
- `PUT /api/sessions/{id}` - Update session field
- `DELETE /api/sessions/{id}` - Delete session
- `DELETE /api/sessions` - Clear all sessions

### Phishlets
- `GET /api/phishlets` - List all phishlets
- `POST /api/phishlets/{name}/enable` - Enable phishlet
- `POST /api/phishlets/{name}/disable` - Disable phishlet

### Lures
- `GET /api/lures` - List all lures
- `POST /api/lures` - Create new lure
- `DELETE /api/lures/{id}` - Delete lure

### Configuration
- `GET /api/config` - Get current config
- `PUT /api/config` - Update config

### Command Execution
- `POST /api/execute` - Execute any CLI command

#### Example: Execute Command
```bash
curl -X POST http://localhost:8080/api/execute \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"command": "phishlets"}'
```

---

## Frontend Features

- **Sessions Dashboard** - View, update, and delete captured sessions
- **Phishlets Manager** - Enable/disable phishlets
- **Lures Manager** - Create and manage lures
- **Configuration Editor** - Modify settings
- **Command Console** - Execute CLI commands directly
- **Real-time Updates** - Auto-refresh data
- **Dark Theme** - Easy on the eyes

---

## Security Notes

1. **Change default password** - Always set `-api-pwd` to a strong password
2. **Only expose locally** - By default API listens on 127.0.0.1 only
3. **Use HTTPS** (recommended) - Run behind a reverse proxy with SSL
4. **Token expiry** - Tokens expire after 24 hours
5. **Audit logging** - All operations are logged

---

## Development

### Frontend Development Server
```bash
cd frontend
npm install
npm run dev
```

This starts a dev server with hot reload on `http://localhost:5173`

The API proxy is configured in `vite.config.js` to forward `/api` calls to `http://localhost:8080`

### Building for Production
```bash
cd frontend
npm run build
```

Output files are in `frontend/dist/` and served by the Go backend.

---

## Troubleshooting

### Port already in use
```bash
./evilginx2 -api -api-addr "127.0.0.1:9000"
```

### Can't find frontend files
Make sure you ran `build_frontend.sh` first to build the React app.

### Session data not persisting
Check database file: `~/.evilginx/data.db`

### API errors
Check logs by running with debug flag:
```bash
./evilginx2 -api -debug
```
