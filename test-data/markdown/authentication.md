# Authentication & Authorization

## Overview

DataFlow Platform uses multiple authentication methods to secure access to your account, data, and API resources. This guide covers all supported authentication mechanisms and best practices.

## Authentication Methods

### 1. API Keys

API keys are the primary method for authenticating API requests.

#### Creating an API Key

1. Navigate to **Settings** → **API Keys**
2. Click **Generate New Key**
3. Provide a descriptive name (e.g., "Production ETL Service")
4. Select permissions scope
5. Copy the key immediately (it won't be shown again)

#### Using API Keys

Include the API key in the `Authorization` header:

```bash
curl -H "Authorization: Bearer dfp_live_1a2b3c4d5e6f7g8h9i0j" \
  https://api.dataflow-platform.io/v1/pipelines
```

#### Security Best Practices

- Never commit API keys to version control
- Use environment variables to store keys
- Rotate keys every 90 days
- Use separate keys for development and production
- Revoke unused keys immediately

### 2. OAuth 2.0

OAuth is recommended for third-party applications that need to access user data.

#### Supported Flows

- **Authorization Code Flow**: For web applications
- **Client Credentials Flow**: For server-to-server communication
- **Refresh Token Flow**: For long-lived access

#### Authorization Code Example

**Step 1: Redirect user to authorization URL**

```
https://auth.dataflow-platform.io/authorize?
  client_id=YOUR_CLIENT_ID&
  redirect_uri=https://yourapp.com/callback&
  response_type=code&
  scope=pipelines:read pipelines:write connections:read
```

**Step 2: Exchange code for access token**

```bash
curl -X POST https://auth.dataflow-platform.io/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=authorization_code" \
  -d "code=AUTH_CODE" \
  -d "client_id=YOUR_CLIENT_ID" \
  -d "client_secret=YOUR_CLIENT_SECRET" \
  -d "redirect_uri=https://yourapp.com/callback"
```

**Response:**

```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "rt_1a2b3c4d5e6f",
  "scope": "pipelines:read pipelines:write connections:read"
}
```

### 3. Service Accounts

Service accounts are designed for automated systems and CI/CD pipelines.

#### Creating a Service Account

```bash
curl -X POST https://api.dataflow-platform.io/v1/service-accounts \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "CI/CD Pipeline Account",
    "permissions": ["pipelines:deploy", "connections:test"]
  }'
```

## Authorization & Permissions

### Permission Scopes

| Scope | Description |
|-------|-------------|
| `pipelines:read` | View pipeline configurations and runs |
| `pipelines:write` | Create, update, and delete pipelines |
| `pipelines:execute` | Trigger pipeline executions |
| `connections:read` | View connection configurations |
| `connections:write` | Create and modify connections |
| `connections:test` | Test connection credentials |
| `users:read` | View user information |
| `users:write` | Manage users and permissions |
| `billing:read` | View billing information |
| `billing:write` | Update billing settings |

### Role-Based Access Control (RBAC)

DataFlow Platform supports three default roles:

#### Admin
- Full access to all resources
- Manage users and permissions
- Access billing information
- Delete organization

#### Developer
- Create and manage pipelines
- Create and test connections
- View execution logs
- Cannot manage users or billing

#### Viewer
- Read-only access to pipelines
- View connection metadata (not credentials)
- View execution history
- Cannot modify any resources

### Custom Roles

Create custom roles with specific permission combinations:

```json
{
  "name": "Data Analyst",
  "permissions": [
    "pipelines:read",
    "pipelines:execute",
    "connections:read"
  ],
  "description": "Can view and run pipelines but not modify them"
}
```

## Token Management

### Token Expiration

| Token Type | Default Expiration |
|------------|-------------------|
| API Key | Never (manual revocation required) |
| OAuth Access Token | 1 hour |
| OAuth Refresh Token | 30 days |
| Service Account Token | 90 days |

### Refreshing Tokens

Use the refresh token to obtain a new access token:

```bash
curl -X POST https://auth.dataflow-platform.io/token \
  -d "grant_type=refresh_token" \
  -d "refresh_token=rt_1a2b3c4d5e6f" \
  -d "client_id=YOUR_CLIENT_ID" \
  -d "client_secret=YOUR_CLIENT_SECRET"
```

## Security Features

### IP Allowlisting

Restrict API access to specific IP addresses:

1. Go to **Settings** → **Security**
2. Enable **IP Allowlisting**
3. Add allowed IP ranges (CIDR notation supported)
4. Save changes

Example:
```
192.168.1.0/24
10.0.0.1/32
```

### Audit Logging

All authentication attempts and API calls are logged:

- Login/logout events
- API key creation and revocation
- Permission changes
- Failed authentication attempts

Access audit logs at **Settings** → **Audit Logs**

## Troubleshooting

### 401 Unauthorized

- Verify API key is correct and not expired
- Check that key has required permissions
- Ensure Authorization header is properly formatted

### 403 Forbidden

- User or API key lacks necessary permissions
- Resource may belong to different organization
- IP address may be blocked

### Token Expired

- Use refresh token to get new access token
- For API keys, generate a new key if rotation is needed

## Related Documentation

For more information on securing your data pipelines, see:
- [Security Best Practices](https://docs.stripe.com/security/guide)
- [Compliance & Certifications](https://aws.amazon.com/compliance/)
